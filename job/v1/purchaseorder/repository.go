package purchaseorder

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type purchaseorderRepository struct {
	tx *sql.Tx
}

func NewPurchaseorderRepository(transaction *sql.Tx) PurchaseorderRepository {
	return &purchaseorderRepository{
		tx: transaction,
	}
}

type PurchaseorderRepository interface {
	GetZohoUpdated(string) (*time.Time, error)
	UpdatePurchaseorder(PurchaseorderUpdate) error
	AddPurchaseorder(PurchaseorderUpdate) error
	AddPurchaseorderItem(PurchaseorderUpdate, PurchaseorderItemUpdate) error
	UpdatePurchaseorderItem(PurchaseorderItemUpdate) error
	GetPurchaseorderItem(string) (int, error)
	GetPurchaseorderItemID(string, string) (string, error)
}

func (r *purchaseorderRepository) GetZohoUpdated(zohoID string) (*time.Time, error) {
	var res time.Time
	row := r.tx.QueryRow(`SELECT zoho_updated FROM purchaseorders WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return &res, err
}

func (r *purchaseorderRepository) UpdatePurchaseorder(info PurchaseorderUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		UPDATE purchaseorders SET 
		purchaseorder_number = ?, 
		date = ?, 
		expected_delivery_date = ?, 
		vendor_id = ?, 
		vendor_name = ?, 
		order_status = ?, 
		received_status = ?, 
		billed_status = ?, 
		sub_total = ?, 
		tax_total = ?, 
		total = ?, 
		status = ?, 
		zoho_updated = ?, 
		updated = ?, 
		updated_by = ?
		WHERE zoho_id = ?
	`, info.PurchaseorderNumber, info.Date, info.ExpectedDeliveryDate, info.VendorID, info.VendorName, info.OrderStatus, info.ReceivedStatus, info.BilledStatus, info.SubTotal, info.TaxTotal, info.Total, info.Status, newModifiedTime, time.Now(), "SYSTEM", info.PurchaseorderID)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *purchaseorderRepository) AddPurchaseorder(info PurchaseorderUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		INSERT INTO purchaseorders (
			zoho_id, 
			purchaseorder_number, 
			date, 
			expected_delivery_date, 
			vendor_id, 
			vendor_name, 
			order_status, 
			received_status, 
			billed_status, 
			sub_total, 
			tax_total, 
			total, 
			status, 
			zoho_updated, 
			created, 
			created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, info.PurchaseorderID, info.PurchaseorderNumber, info.Date, info.ExpectedDeliveryDate, info.VendorID, info.VendorName, info.OrderStatus, info.ReceivedStatus, info.BilledStatus, info.SubTotal, info.TaxTotal, info.Total, info.Status, newModifiedTime, time.Now(), "SYSTEM")
	if err != nil {
		msg := "新建失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *purchaseorderRepository) AddPurchaseorderItem(info PurchaseorderUpdate, item PurchaseorderItemUpdate) error {
	_, err := r.tx.Exec(`
		INSERT INTO purchaseorder_items (
			zoho_id, 
			purchaseorder_id, 
			item_id, 
			sku, 
			item_name, 
			rate, 
			quantity, 
			discount, 
			item_total, 
			tax_total, 
			quantity_received, 
			quantity_cancelled, 
			quantity_billed, 
			created, 
			created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, item.PurchaseorderItemID, info.PurchaseorderID, item.ItemID, item.SKU, item.ItemName, item.Rate, item.Quantity, item.DiscountTotal, item.ItemTotal, 0, item.QuantityReceived, item.QuantityCancelled, item.Quantitybilled, time.Now(), "SYSTEM")
	if err != nil {
		msg := "新建失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *purchaseorderRepository) GetPurchaseorderItem(zohoID string) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM purchaseorder_items WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return res, err
}

func (r *purchaseorderRepository) UpdatePurchaseorderItem(item PurchaseorderItemUpdate) error {
	_, err := r.tx.Exec(`
		UPDATE purchaseorder_items SET
		rate = ?, 
		quantity = ?, 
		discount = ?, 
		item_total = ?, 
		updated = ?, 
		updated_by = ?
		WHERE zoho_id = ?
		`, item.Rate, item.Quantity, item.DiscountTotal, item.ItemTotal, time.Now(), "SYSTEM", item.PurchaseorderItemID)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *purchaseorderRepository) GetPurchaseorderItemID(poID, itemID string) (string, error) {
	var res string
	row := r.tx.QueryRow(`SELECT zoho_id FROM purchaseorder_items WHERE purchaseorder_id = ? AND item_id = ? LIMIT 1`, poID, itemID)
	err := row.Scan(&res)
	return res, err
}
