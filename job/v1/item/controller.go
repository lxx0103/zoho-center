package item

import (
	"fmt"
	"zoho-center/job/v1/auth"
)

func GetItemList() {
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	itemService := NewItemService()
	hasMorePage, err := itemService.GetItemList(token, 1)
	fmt.Println(hasMorePage)
	fmt.Println(err)
}
