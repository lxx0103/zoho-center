package salesorder

import "time"

type Salesorder struct {
	ID                   int64     `db:"id" json:"id"`
	ZohoID               string    `db:"zoho_id" json:"zoho_id"`
	SalesorderNumber     string    `db:"salesorder_number" json:"salesorder_number"`
	Date                 string    `db:"date" json:"date"`
	ExpectedShipmentDate string    `db:"expected_shipment_date" json:"expected_shipment_date"`
	CustomerID           string    `db:"customer_id" json:"customer_id"`
	CustomerName         string    `db:"customer_name" json:"customer_name"`
	OrderStatus          string    `db:"order_status" json:"order_status"`
	InvoicedStatus       string    `db:"invoiced_status" json:"invoiced_status"`
	PaidStatus           string    `db:"paid_status" json:"paid_status"`
	ShippedStatus        string    `db:"shipped_status" json:"shipped_status"`
	Source               string    `db:"source" json:"source"`
	SalespersonID        string    `db:"salesperson_id" json:"salesperson_id"`
	SalespersonName      string    `db:"salesperson_name" json:"salesperson_name"`
	ShippingCharge       float64   `db:"shipping_charge" json:"shipping_charge"`
	SubTotal             float64   `db:"sub_total" json:"sub_total"`
	DiscountTotal        float64   `db:"discount_total" json:"discount_total"`
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

type SalesorderItem struct {
	ID              int64     `db:"id" json:"id"`
	ZohoID          string    `db:"zoho_id" json:"zoho_id"`
	SalesorderID    string    `db:"salesorder_id" json:"salesorder_id"`
	ItemID          string    `db:"item_id" json:"item_id"`
	SKU             string    `db:"sku" json:"sku"`
	Rate            string    `db:"rate" json:"rate"`
	Quantity        string    `db:"quantity" json:"quantity"`
	Discount        string    `db:"discount" json:"discount"`
	ItemTotal       string    `db:"item_total" json:"item_total"`
	TaxTotal        string    `db:"tax_total" json:"tax_total"`
	QuantityPacked  string    `db:"quantity_packed" json:"quantity_packed"`
	QuantityShipped string    `db:"quantity_shipped" json:"quantity_shipped"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}
