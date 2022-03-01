package item

import "time"

type Item struct {
	ID           int64     `db:"id" json:"id"`
	ZohoID       string    `db:"zoho_id" json:"zoho_id"`
	Name         string    `db:"name" json:"name"`
	SKU          string    `db:"sku" json:"sku"`
	Status       string    `db:"status" json:"status"`
	Um           string    `db:"um" json:"um"`
	Desc         string    `db:"desc" json:"desc"`
	InitialStock float64   `db:"initial_stock" json:"initial_stock"`
	InitialRate  float64   `db:"initial_rate" json:"initial_rate"`
	PurchaseRate float64   `db:"purchase_rate" json:"purchase_rate"`
	SalesRate    float64   `db:"sales_rate" json:"sales_rate"`
	VendorID     string    `db:"vendor_id" json:"vendor_id"`
	Source       string    `db:"source" json:"source"`
	ZohoUpdated  time.Time `db:"zoho_updated" json:"zoho_updated"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	Updated      time.Time `db:"updated" json:"updated"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}
