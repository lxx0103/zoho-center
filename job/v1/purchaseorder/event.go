package purchaseorder

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

type NewReceiveToZoho struct {
	POID     string `json:"po_id"`
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("NewReceiveToZoho", "NewReceiveToZoho", SyncReceiveToZoho)
}

func SyncReceiveToZoho(d amqp.Delivery) bool {
	fmt.Println("aaa")
	if d.Body == nil {
		return false
	}
	fmt.Println("aaa")
	var newReceiveToZoho NewReceiveToZoho
	err := json.Unmarshal(d.Body, &newReceiveToZoho)
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
	repo := NewPurchaseorderRepository(tx)
	lineID, err := repo.GetPurchaseorderItemID(newReceiveToZoho.POID, newReceiveToZoho.SKU)
	fmt.Println(lineID)
	if err != nil {
		fmt.Println("11")
		fmt.Println(err)
		return false
	}
	receiveNumber := "WMS-PR-" + time.Now().Format("20060102150405")
	receiveDate := time.Now().Format("2006-01-02")
	jsonReq := []byte(`{ "receive_number" : "` + receiveNumber + `", "date" : "` + receiveDate + `", "line_items" : [ { "line_item_id" : "` + lineID + `", "quantity": ` + strconv.FormatInt(newReceiveToZoho.Quantity, 10) + ` } ] }`)
	receive_uri := config.ReadConfig("zoho.receive_uri") + "?purchaseorder_id=" + newReceiveToZoho.POID
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err)
		return false
	}
	authHeader := "Zoho-oauthtoken " + token
	req, err := http.NewRequest("POST", receive_uri, bytes.NewBuffer(jsonReq))
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
	// var purchaseorderDetail ZohoPurchaseorderDetail
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
