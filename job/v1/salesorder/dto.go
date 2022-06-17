package salesorder

type ZohoSalesorderList struct {
	Code        int              `json:"code"`
	Message     string           `json:"message"`
	Salesorders []ZohoSalesorder `json:"salesorders"`
	Page        ZohoPage         `json:"page_context"`
}

type ZohoSalesorder struct {
	SalesorderID     string `json:"salesorder_id"`
	LastModifiedTime string `json:"last_modified_time"`
}

type ZohoPage struct {
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	HasMorePage bool `json:"has_more_page"`
}

type ZohoSalesorderDetail struct {
	Code       int              `json:"code"`
	Message    string           `json:"message"`
	Salesorder SalesorderUpdate `json:"salesorder"`
}

type SalesorderUpdate struct {
	SalesorderID         string                 `json:"salesorder_id"`
	SalesorderNumber     string                 `json:"salesorder_number"`
	Date                 string                 `json:"date"`
	ExpectedShipmentDate string                 `json:"shipment_date"`
	CustomerID           string                 `json:"customer_id"`
	CustomerName         string                 `json:"customer_name"`
	OrderStatus          string                 `json:"order_status"`
	InvoicedStatus       string                 `json:"invoiced_status"`
	PaidStatus           string                 `json:"paid_status"`
	ShippedStatus        string                 `json:"shipped_status"`
	Status               string                 `json:"status"`
	Source               string                 `json:"source"`
	SalespersonID        string                 `json:"salesperson_id"`
	SalespersonName      string                 `json:"salesperson_name"`
	ShippingCharge       float64                `json:"shipping_charge"`
	SubTotal             float64                `json:"sub_total"`
	DiscountTotal        float64                `json:"discount_total"`
	TaxTotal             float64                `json:"tax_total"`
	Total                float64                `json:"total"`
	LastModifiedTime     string                 `json:"last_modified_time"`
	Items                []SalesorderItemUpdate `json:"line_items"`
}

type SalesorderItemUpdate struct {
	SalesorderItemID string  `json:"line_item_id"`
	ItemID           string  `json:"item_id"`
	SKU              string  `json:"sku"`
	ItemName         string  `json:"name"`
	Rate             float64 `json:"rate"`
	Quantity         float64 `json:"quantity"`
	DiscountTotal    float64 `json:"discount"`
	ItemTotal        float64 `json:"item_total"`
	QuantityPacked   float64 `json:"quantity_packed"`
	QuantityShipped  float64 `json:"quantity_shipped"`
}
