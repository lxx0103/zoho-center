package purchaseorder

import "time"

type Purchaseorder struct {
	ID                   int64     `db:"id" json:"id"`
	ZohoID               string    `db:"zoho_id" json:"zoho_id"`
	PurchaseorderNumber  string    `db:"purchaseorder_number" json:"purchaseorder_number"`
	Date                 string    `db:"date" json:"date"`
	ExpectedDeliveryDate string    `db:"expected_delivery_date" json:"expected_delivery_date"`
	VendorID             string    `db:"vendor_id" json:"vendor_id"`
	VendorName           string    `db:"vendor_name" json:"vendor_name"`
	OrderStatus          string    `db:"order_status" json:"order_status"`
	ReceivedStatus       string    `db:"received_status" json:"received_status"`
	BillStatus           string    `db:"bill_status" json:"bill_status"`
	SubTotal             float64   `db:"sub_total" json:"sub_total"`
	TaxTotal             float64   `db:"tax_total" json:"tax_total"`
	Total                float64   `db:"total" json:"total"`
	Rate                 float64   `db:"rate" json:"rate"`
	Status               string    `db:"status" json:"status"`
	ZohoUpdated          time.Time `db:"zoho_updated" json:"zoho_updated"`
	Created              time.Time `db:"created" json:"created"`
	CreatedBy            string    `db:"created_by" json:"created_by"`
	Updated              time.Time `db:"updated" json:"updated"`
	UpdatedBy            string    `db:"updated_by" json:"updated_by"`
}

type PurchaseorderItem struct {
	ID                int64     `db:"id" json:"id"`
	ZohoID            string    `db:"zoho_id" json:"zoho_id"`
	PurchaseorderID   string    `db:"purchaseorder_id" json:"purchaseorder_id"`
	ItemID            string    `db:"item_id" json:"item_id"`
	SKU               string    `db:"sku" json:"sku"`
	ItemName          string    `db:"item_name" json:"item_name"`
	Rate              string    `db:"rate" json:"rate"`
	Quantity          string    `db:"quantity" json:"quantity"`
	Discount          string    `db:"discount" json:"discount"`
	ItemTotal         string    `db:"item_total" json:"item_total"`
	TaxTotal          string    `db:"tax_total" json:"tax_total"`
	QuantityReceived  string    `db:"quantity_received" json:"quantity_received"`
	QuantityCancelled string    `db:"quantity_cancelled" json:"quantity_cancelled"`
	Quantitybilled    string    `db:"quantity_billed" json:"quantity_billed"`
	Created           time.Time `db:"created" json:"created"`
	CreatedBy         string    `db:"created_by" json:"created_by"`
	Updated           time.Time `db:"updated" json:"updated"`
	UpdatedBy         string    `db:"updated_by" json:"updated_by"`
}
