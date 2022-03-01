package item

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type itemRepository struct {
	tx *sql.Tx
}

func NewItemRepository(transaction *sql.Tx) ItemRepository {
	return &itemRepository{
		tx: transaction,
	}
}

type ItemRepository interface {
	GetZohoUpdated(string) (*time.Time, error)
	UpdateItem(ItemUpdate) error
	AddItem(ItemUpdate) error
}

func (r *itemRepository) GetZohoUpdated(zohoID string) (*time.Time, error) {
	var res time.Time
	row := r.tx.QueryRow(`SELECT zoho_updated FROM items WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return &res, err
}

func (r *itemRepository) UpdateItem(info ItemUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		Update items SET
		name = ?,
		sku = ?,
		status = ?,
		um = ?,
		description = ?,
		rate = ?,
		initial_stock = ?,
		initial_rate = ?,
		purchase_rate = ?,
		sales_rate = ?,
		stock_on_hand = ?,
		available_stock = ?,
		actual_available_stock = ?,
		vendor_id = ?,
		source = ?,
		zoho_updated = ?,
		updated = ?,
		updated_by = ?
		WHERE zoho_id = ?
	`, info.Name, info.Sku, info.Status, info.Unit, info.Description, info.Rate, info.InitialStock, info.InitialStockRate, info.PurchaseRate, info.SalesRate, info.StockOnHand, info.AvailableStock, info.ActualAvailableStock, info.VendorID, info.Source, newModifiedTime, time.Now(), "SYSTEM", info.ItemID)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *itemRepository) AddItem(info ItemUpdate) error {
	newModifiedTime, _ := time.Parse(time.RFC3339, strings.Replace(strings.Replace(info.LastModifiedTime, " ", "T", 1), "+0800", "+08:00", 1))
	_, err := r.tx.Exec(`
		INSERT INTO items ( 
			zoho_id,
			name,
			sku, 
			status, 
			um, 
			description, 
			rate, 
			initial_stock, 
			initial_rate, 
			purchase_rate, 
			sales_rate,
			stock_on_hand, 
			available_stock, 
			actual_available_stock, 
			vendor_id, 
			source, 
			zoho_updated, 
			created,
			created_by,
			updated, 
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, info.ItemID, info.Name, info.Sku, info.Status, info.Unit, info.Description, info.Rate, info.InitialStock, info.InitialStockRate, info.PurchaseRate, info.SalesRate, info.StockOnHand, info.AvailableStock, info.ActualAvailableStock, info.VendorID, info.Source, newModifiedTime, time.Now(), "SYSTEM", time.Now(), "SYSTEM")
	if err != nil {
		msg := "新建失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}
