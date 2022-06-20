package purchaseorder

type ZohoPurchaseorderList struct {
	Code           int                 `json:"code"`
	Message        string              `json:"message"`
	Purchaseorders []ZohoPurchaseorder `json:"purchaseorders"`
	Page           ZohoPage            `json:"page_context"`
}

type ZohoPurchaseorder struct {
	PurchaseorderID  string `json:"purchaseorder_id"`
	LastModifiedTime string `json:"last_modified_time"`
}

type ZohoPage struct {
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	HasMorePage bool `json:"has_more_page"`
}

type ZohoPurchaseorderDetail struct {
	Code          int                 `json:"code"`
	Message       string              `json:"message"`
	Purchaseorder PurchaseorderUpdate `json:"purchaseorder"`
}

type PurchaseorderUpdate struct {
	PurchaseorderID      string                    `json:"purchaseorder_id"`
	PurchaseorderNumber  string                    `json:"purchaseorder_number"`
	Date                 string                    `json:"date"`
	ExpectedDeliveryDate string                    `json:"delivery_date"`
	VendorID             string                    `json:"vendor_id"`
	VendorName           string                    `json:"vendor_name"`
	OrderStatus          string                    `json:"order_status"`
	ReceivedStatus       string                    `json:"received_status"`
	BilledStatus         string                    `json:"billed_status"`
	Status               string                    `json:"status"`
	SubTotal             float64                   `json:"sub_total"`
	TaxTotal             float64                   `json:"tax_total"`
	Total                float64                   `json:"total"`
	LastModifiedTime     string                    `json:"last_modified_time"`
	Items                []PurchaseorderItemUpdate `json:"line_items"`
}

type PurchaseorderItemUpdate struct {
	PurchaseorderItemID string  `json:"line_item_id"`
	ItemID              string  `json:"item_id"`
	SKU                 string  `json:"sku"`
	ItemName            string  `json:"name"`
	Rate                float64 `json:"rate"`
	Quantity            float64 `json:"quantity"`
	DiscountTotal       float64 `json:"discount"`
	ItemTotal           float64 `json:"item_total"`
	QuantityReceived    float64 `json:"quantity_received"`
	QuantityCancelled   float64 `json:"quantity_cancelled"`
	Quantitybilled      float64 `json:"quantity_billed"`
}
