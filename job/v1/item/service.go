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
)

type itemService struct {
}

func NewItemService() ItemService {
	return &itemService{}
}

type ItemService interface {
	GetItemList(string, int) (bool, error)
}

func (s itemService) GetItemList(token string, page int) (bool, error) {
	url := config.ReadConfig("zoho.item_uri")
	auth := "Zoho-oauthtoken " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
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
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var itemList ZohoItemList
	err = json.Unmarshal(body, &itemList)
	if err != nil {
		return false, err
	}
	if itemList.Code != 0 {
		msg := "状态码错误:" + fmt.Sprint(itemList.Code)
		return false, errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return false, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	for i := 0; i < len(itemList.Items); i++ {
		lastModifiedTime, err := repo.GetZohoUpdated(itemList.Items[i].ItemID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				*lastModifiedTime, _ = time.Parse(time.RFC3339, "2000-01-01 00:00:00+0800")
			} else {
				return false, err
			}
		}
		newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(itemList.Items[i].LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
		if lastModifiedTime.Before(newModifiedTime) {
			fmt.Println("Item has been changed:")
			fmt.Println(itemList.Items[i].ItemID)
		}
	}
	// return itemList.Page.HasMorePage, nil
	return false, nil
}
