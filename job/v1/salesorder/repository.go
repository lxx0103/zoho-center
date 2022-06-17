package salesorder

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type salesorderRepository struct {
	tx *sql.Tx
}

func NewSalesorderRepository(transaction *sql.Tx) SalesorderRepository {
	return &salesorderRepository{
		tx: transaction,
	}
}

type SalesorderRepository interface {
	GetZohoUpdated(string) (*time.Time, error)
	UpdateSalesorder(SalesorderUpdate) error
	AddSalesorder(SalesorderUpdate) error
	AddSalesorderItem(SalesorderUpdate, SalesorderItemUpdate) error
	UpdateSalesorderItem(SalesorderItemUpdate) error
	GetSalesorderItem(string) (int, error)
}

func (r *salesorderRepository) GetZohoUpdated(zohoID string) (*time.Time, error) {
	var res time.Time
	row := r.tx.QueryRow(`SELECT zoho_updated FROM salesorders WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return &res, err
}

func (r *salesorderRepository) UpdateSalesorder(info SalesorderUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		UPDATE salesorders SET 
		salesorder_number = ?, 
		date = ?, 
		expected_shipment_date = ?, 
		customer_id = ?, 
		customer_name = ?, 
		order_status = ?, 
		invoiced_status = ?, 
		paid_status = ?, 
		shipped_status = ?, 
		source = ?, 
		salesperson_id = ?, 
		salesperson_name = ?, 
		shipping_charge = ?, 
		sub_total = ?, 
		discount_total = ?, 
		tax_total = ?, 
		total = ?, 
		status = ?, 
		zoho_updated = ?, 
		updated = ?, 
		updated_by = ?
		WHERE zoho_id = ?
	`, info.SalesorderNumber, info.Date, info.ExpectedShipmentDate, info.CustomerID, info.CustomerName, info.OrderStatus, info.InvoicedStatus, info.PaidStatus, info.ShippedStatus, info.Source, info.SalespersonID, info.SalespersonName, info.ShippingCharge, info.SubTotal, info.DiscountTotal, info.TaxTotal, info.Total, info.Status, newModifiedTime, time.Now(), "SYSTEM", info.SalesorderID)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *salesorderRepository) AddSalesorder(info SalesorderUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		INSERT INTO salesorders (
			zoho_id, 
			salesorder_number, 
			date, 
			expected_shipment_date, 
			customer_id, 
			customer_name, 
			order_status, 
			invoiced_status, 
			paid_status, 
			shipped_status, 
			source, 
			salesperson_id, 
			salesperson_name, 
			shipping_charge, 
			sub_total, 
			discount_total, 
			tax_total, 
			total, 
			status, 
			zoho_updated, 
			created, 
			created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, info.SalesorderID, info.SalesorderNumber, info.Date, info.ExpectedShipmentDate, info.CustomerID, info.CustomerName, info.OrderStatus, info.InvoicedStatus, info.PaidStatus, info.ShippedStatus, info.Source, info.SalespersonID, info.SalespersonName, info.ShippingCharge, info.SubTotal, info.DiscountTotal, info.TaxTotal, info.Total, info.Status, newModifiedTime, time.Now(), "SYSTEM")
	if err != nil {
		msg := "新建失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *salesorderRepository) AddSalesorderItem(info SalesorderUpdate, item SalesorderItemUpdate) error {
	_, err := r.tx.Exec(`
		INSERT INTO salesorder_items (
			zoho_id, 
			salesorder_id, 
			item_id, 
			sku, 
			item_name, 
			rate, 
			quantity, 
			discount, 
			item_total, 
			tax_total, 
			quantity_packed, 
			quantity_shipped, 
			created, 
			created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, item.SalesorderItemID, info.SalesorderID, item.ItemID, item.SKU, item.ItemName, item.Rate, item.Quantity, item.DiscountTotal, item.ItemTotal, 0, item.QuantityPacked, item.QuantityShipped, time.Now(), "SYSTEM")
	if err != nil {
		msg := "新建失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *salesorderRepository) GetSalesorderItem(zohoID string) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM salesorder_items WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return res, err
}

func (r *salesorderRepository) UpdateSalesorderItem(item SalesorderItemUpdate) error {
	_, err := r.tx.Exec(`
		UPDATE salesorder_items SET
		rate = ?, 
		quantity = ?, 
		discount = ?, 
		item_total = ?, 
		updated = ?, 
		updated_by = ?
		WHERE zoho_id = ?
		`, item.Rate, item.Quantity, item.DiscountTotal, item.ItemTotal, time.Now(), "SYSTEM", item.SalesorderItemID)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}
