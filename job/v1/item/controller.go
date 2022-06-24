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
	hasMorePage, new, err := itemService.GetItemList(token, 1)
	fmt.Println(hasMorePage)
	fmt.Println(err)

	for i := 0; i < len(*new); i++ {
		token, err := auth.GetCode()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = itemService.UpdateItem(token, (*new)[i])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func GetItemListTest() {
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	itemService := NewItemService()
	itemService.UpdateItem(token, "8581000023858747")
	itemService.UpdateItem(token, "8581000023858760")
	itemService.UpdateItem(token, "8581000023858773")
}
