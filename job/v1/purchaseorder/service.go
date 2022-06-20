package purchaseorder

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

type purchaseorderService struct {
}

func NewPurchaseorderService() PurchaseorderService {
	return &purchaseorderService{}
}

type PurchaseorderService interface {
	GetPurchaseorderList(string, int) (bool, *[]string, error)
	UpdatePurchaseorder(string, string) error
}

func (s purchaseorderService) GetPurchaseorderList(token string, page int) (bool, *[]string, error) {
	var res []string
	url := config.ReadConfig("zoho.purchaseorder_uri")
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
	var purchaseorderList ZohoPurchaseorderList
	err = json.Unmarshal(body, &purchaseorderList)
	if err != nil {
		return false, nil, err
	}
	if purchaseorderList.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(purchaseorderList.Code) + purchaseorderList.Message
		return false, nil, errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return false, nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	for i := 0; i < len(purchaseorderList.Purchaseorders); i++ {
		lastModifiedTime, err := repo.GetZohoUpdated(purchaseorderList.Purchaseorders[i].PurchaseorderID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				*lastModifiedTime, _ = time.Parse(time.RFC3339, "2000-01-01 00:00:00+0800")
			} else {
				return false, nil, err
			}
		}
		newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(purchaseorderList.Purchaseorders[i].LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
		if lastModifiedTime.Before(newModifiedTime) {
			res = append(res, purchaseorderList.Purchaseorders[i].PurchaseorderID)
		}
	}
	return purchaseorderList.Page.HasMorePage, &res, nil
	// return false, nil
}

func (s purchaseorderService) UpdatePurchaseorder(token string, purchaseorderID string) error {
	purchaseorderID = "8581000024488308"
	duration := time.Duration(3) * time.Second
	time.Sleep(duration)
	url := config.ReadConfig("zoho.purchaseorder_uri") + purchaseorderID
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
	var purchaseorderDetail ZohoPurchaseorderDetail
	err = json.Unmarshal([]byte(body), &purchaseorderDetail)
	if err != nil {
		return err
	}
	if purchaseorderDetail.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(purchaseorderDetail.Code) + purchaseorderDetail.Message
		return errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)

	_, err = repo.GetZohoUpdated(purchaseorderDetail.Purchaseorder.PurchaseorderID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = repo.AddPurchaseorder(purchaseorderDetail.Purchaseorder)
			if err != nil {
				return err
			}
			for _, purchaseorderItem := range purchaseorderDetail.Purchaseorder.Items {
				err = repo.AddPurchaseorderItem(purchaseorderDetail.Purchaseorder, purchaseorderItem)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		err = repo.UpdatePurchaseorder(purchaseorderDetail.Purchaseorder)
		if err != nil {
			return err
		}
		for _, purchaseorderItem := range purchaseorderDetail.Purchaseorder.Items {
			count, err := repo.GetPurchaseorderItem(purchaseorderItem.PurchaseorderItemID)
			if err != nil {
				if count == 0 {
					err = repo.AddPurchaseorderItem(purchaseorderDetail.Purchaseorder, purchaseorderItem)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				err = repo.UpdatePurchaseorderItem(purchaseorderItem)
				if err != nil {
					return err
				}
			}
		}
	}

	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(purchaseorderDetail.Purchaseorder)
	err = rabbit.Publish("PurchaseOrderUpdated", msg)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
