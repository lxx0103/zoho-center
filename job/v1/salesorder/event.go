package salesorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"zoho-center/core/config"
	"zoho-center/core/database"
	"zoho-center/core/queue"
	"zoho-center/job/v1/auth"

	"github.com/streadway/amqp"
)

type NewPackedToZoho struct {
	SOID     string `json:"so_id"`
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("NewPackedToZoho", "NewPackedToZoho", SyncPackedToZoho)
}

func SyncPackedToZoho(d amqp.Delivery) bool {
	fmt.Println("aaa")
	if d.Body == nil {
		return false
	}
	fmt.Println("aaa")
	var newPackedToZoho NewPackedToZoho
	err := json.Unmarshal(d.Body, &newPackedToZoho)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
		return false
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	lineID, err := repo.GetSalesorderItemID(newPackedToZoho.SOID, newPackedToZoho.SKU)
	fmt.Println(lineID)
	if err != nil {
		fmt.Println("11")
		fmt.Println(err)
		return false
	}
	// packageNumber := "WMS-PA-" + time.Now().Format("20060102150405")
	packageDate := time.Now().Format("2006-01-02")
	jsonReq := []byte(`{  "date" : "` + packageDate + `", "line_items" : [ { "so_line_item_id" : "` + lineID + `", "quantity": ` + strconv.FormatInt(newPackedToZoho.Quantity, 10) + ` } ] }`)
	package_uri := config.ReadConfig("zoho.package_uri") + "?salesorder_id=" + newPackedToZoho.SOID
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err)
		return false
	}
	authHeader := "Zoho-oauthtoken " + token
	req, err := http.NewRequest("POST", package_uri, bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("body", string(body))
	// var purchaseorderDetail ZohoSalesorderDetail
	// err = json.Unmarshal([]byte(body), &purchaseorderDetail)
	// if err != nil {
	// 	return err
	// }
	// if purchaseorderDetail.Code != 0 {
	// 	msg := "状态码错误:" + fmt.Sprint(purchaseorderDetail.Code) + purchaseorderDetail.Message
	// 	return errors.New(msg)
	// }
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
