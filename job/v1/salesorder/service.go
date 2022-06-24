package salesorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"zoho-center/core/config"
	"zoho-center/core/database"
	"zoho-center/core/queue"
)

type salesorderService struct {
}

func NewSalesorderService() SalesorderService {
	return &salesorderService{}
}

type SalesorderService interface {
	GetSalesorderList(string, int) (bool, *[]string, error)
	UpdateSalesorder(string, string) error
}

func (s salesorderService) GetSalesorderList(token string, page int) (bool, *[]string, error) {
	var res []string
	url := config.ReadConfig("zoho.salesorder_uri")
	auth := "Zoho-oauthtoken " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", auth)
	q := req.URL.Query()
	q.Add("page", fmt.Sprint(page))
	q.Add("sort_column", "last_modified_time")
	q.Add("sort_order", "D")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}
	var salesorderList ZohoSalesorderList
	err = json.Unmarshal(body, &salesorderList)
	if err != nil {
		return false, nil, err
	}
	if salesorderList.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(salesorderList.Code) + salesorderList.Message
		return false, nil, errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return false, nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	for i := 0; i < len(salesorderList.Salesorders); i++ {
		lastModifiedTime, err := repo.GetZohoUpdated(salesorderList.Salesorders[i].SalesorderID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				*lastModifiedTime, _ = time.Parse(time.RFC3339, "2000-01-01 00:00:00+0800")
			} else {
				return false, nil, err
			}
		}
		newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(salesorderList.Salesorders[i].LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
		if lastModifiedTime.Before(newModifiedTime) {
			res = append(res, salesorderList.Salesorders[i].SalesorderID)
		}
	}
	return salesorderList.Page.HasMorePage, &res, nil
	// return false, nil
}

func (s salesorderService) UpdateSalesorder(token string, salesorderID string) error {
	duration := time.Duration(3) * time.Second
	time.Sleep(duration)
	url := config.ReadConfig("zoho.salesorder_uri") + salesorderID
	auth := "Zoho-oauthtoken " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", auth)
	// q := req.URL.Query() // Get a copy of the query values.
	// //q.Add("page", fmt.Sprint(page)) // Add a new value to the set.
	// req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var salesorderDetail ZohoSalesorderDetail
	err = json.Unmarshal([]byte(body), &salesorderDetail)
	if err != nil {
		return err
	}
	if salesorderDetail.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(salesorderDetail.Code) + salesorderDetail.Message
		return errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)

	_, err = repo.GetZohoUpdated(salesorderDetail.Salesorder.SalesorderID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = repo.AddSalesorder(salesorderDetail.Salesorder)
			if err != nil {
				return err
			}
			for _, salesorderItem := range salesorderDetail.Salesorder.Items {
				err = repo.AddSalesorderItem(salesorderDetail.Salesorder, salesorderItem)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		err = repo.UpdateSalesorder(salesorderDetail.Salesorder)
		if err != nil {
			return err
		}
		for _, salesorderItem := range salesorderDetail.Salesorder.Items {
			count, err := repo.GetSalesorderItem(salesorderItem.SalesorderItemID)
			if err != nil {
				if count == 0 {
					err = repo.AddSalesorderItem(salesorderDetail.Salesorder, salesorderItem)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				err = repo.UpdateSalesorderItem(salesorderItem)
				if err != nil {
					return err
				}
			}
		}
	}

	rabbit, _ := queue.GetConn()
	salesorderDetail.Salesorder.Status = strings.ToUpper(salesorderDetail.Salesorder.Status)
	msg, _ := json.Marshal(salesorderDetail.Salesorder)
	err = rabbit.Publish("SalesOrderUpdated", msg)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
