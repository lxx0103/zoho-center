package purchaseorder

import (
	"fmt"
	"zoho-center/job/v1/auth"
)

func GetPurchaseorderList() {
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	purchaseorderService := NewPurchaseorderService()
	hasMorePage, new, err := purchaseorderService.GetPurchaseorderList(token, 1)
	fmt.Println(hasMorePage)
	fmt.Println(err)

	for i := 0; i < len(*new); i++ {
		token, err := auth.GetCode()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = purchaseorderService.UpdatePurchaseorder(token, (*new)[i])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		break
	}
}
