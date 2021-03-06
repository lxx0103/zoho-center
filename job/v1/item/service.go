package item

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

type itemService struct {
}

func NewItemService() ItemService {
	return &itemService{}
}

type ItemService interface {
	GetItemList(string, int) (bool, *[]string, error)
	UpdateItem(string, string) error
}

func (s itemService) GetItemList(token string, page int) (bool, *[]string, error) {
	var res []string
	url := config.ReadConfig("zoho.item_uri")
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
	var itemList ZohoItemList
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		return false, nil, err
	}
	if itemList.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(itemList.Code) + itemList.Message
		return false, nil, errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return false, nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	for i := 0; i < len(itemList.Items); i++ {
		lastModifiedTime, err := repo.GetZohoUpdated(itemList.Items[i].ItemID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				*lastModifiedTime, _ = time.Parse(time.RFC3339, "2000-01-01 00:00:00+0800")
			} else {
				return false, nil, err
			}
		}
		newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(itemList.Items[i].LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
		if lastModifiedTime.Before(newModifiedTime) {
			res = append(res, itemList.Items[i].ItemID)
		}
	}
	return itemList.Page.HasMorePage, &res, nil
	// return false, nil
}

func (s itemService) UpdateItem(token string, itemID string) error {
	duration := time.Duration(3) * time.Second
	time.Sleep(duration)
	url := config.ReadConfig("zoho.item_uri") + itemID
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
	var itemDetail ZohoItemDetail
	err = json.Unmarshal([]byte(body), &itemDetail)
	if err != nil {
		return err
	}
	if itemDetail.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(itemDetail.Code) + itemDetail.Message
		return errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)

	_, err = repo.GetZohoUpdated(itemDetail.Item.ItemID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = repo.AddItem(itemDetail.Item)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		err = repo.UpdateItem(itemDetail.Item)
		if err != nil {
			return err
		}
	}

	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(itemDetail.Item)
	err = rabbit.Publish("ItemUpdated", msg)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
