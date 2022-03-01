package item

type ZohoItemList struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Items   []ZohoItem `json:"items"`
	Page    ZohoPage   `json:"page_context"`
}

type ZohoItem struct {
	ItemID           string `json:"item_id"`
	LastModifiedTime string `json:"last_modified_time"`
}

type ZohoPage struct {
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	HasMorePage bool `json:"has_more_page"`
}

type ZohoItemDetail struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Item    ItemUpdate `json:"item"`
}

type ItemUpdate struct {
	ItemID               string  `json:"item_id"`
	Name                 string  `json:"name"`
	Sku                  string  `json:"sku"`
	Status               string  `json:"status"`
	Unit                 string  `json:"unit"`
	Description          string  `json:"description"`
	Rate                 float64 `json:"rate"`
	InitialStock         float64 `json:"initial_stock"`
	InitialStockRate     float64 `json:"initial_stock_rate"`
	PurchaseRate         float64 `json:"purchase_rate"`
	SalesRate            float64 `json:"sales_rate"`
	StockOnHand          float64 `json:"stock_on_hand"`
	AvailableStock       float64 `json:"available_stock"`
	ActualAvailableStock float64 `json:"actual_available_stock"`
	VendorID             string  `json:"vendor_id"`
	Source               string  `json:"source"`
	LastModifiedTime     string  `json:"last_modified_time"`
}
