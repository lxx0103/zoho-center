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
	fmt.Println(token)
}
