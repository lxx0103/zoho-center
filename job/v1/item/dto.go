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
