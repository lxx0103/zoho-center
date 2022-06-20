package purchaseorder

import (
	"zoho-center/core/queue"
)

func Subscribe(conn *queue.Conn) {
	// conn.StartConsumer("NewReceiveCreated", "NewReceiveCreated", AddTransaction)
}

// func AddTransaction(d amqp.Delivery) bool {
// 	if d.Body == nil {
// 		return false
// 	}
// 	var newReceiveCreated NewReceiveCreated
// 	err := json.Unmarshal(d.Body, &newReceiveCreated)
// 	if err != nil {
// 		fmt.Println("1")
// 		fmt.Println(err)
// 		return false
// 	}
// 	db := database.InitMySQL()
// 	repo := NewInventoryRepository(db)
// 	po, err := repo.GetPurchaseOrderByID(newReceiveCreated.POID)
// 	if err != nil {
// 		fmt.Println("11")
// 		fmt.Println(err)
// 		return false
// 	}
// 	var transation TransactionNew
// 	transation.POID = newReceiveCreated.POID
// 	transation.PONumber = po.PONumber
// 	transation.ItemName = newReceiveCreated.ItemName
// 	transation.SKU = newReceiveCreated.SKU
// 	transation.Quantity = newReceiveCreated.Quantity
// 	transation.ShelfCode = newReceiveCreated.ShelfCode
// 	transation.ShelfLocation = newReceiveCreated.ShelfLocation
// 	transation.LocationCode = newReceiveCreated.LocationCode
// 	transation.LocationLevel = newReceiveCreated.LocationLevel
// 	transation.User = newReceiveCreated.User
// 	err = repo.CreateTransaction(transation)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return true
// }
