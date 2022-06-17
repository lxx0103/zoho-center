package salesorder

import (
	"fmt"
	"zoho-center/job/v1/auth"
)

func GetSalesorderList() {
	token, err := auth.GetCode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	salesorderService := NewSalesorderService()
	hasMorePage, new, err := salesorderService.GetSalesorderList(token, 1)
	fmt.Println(hasMorePage)
	fmt.Println(err)

	for i := 0; i < len(*new); i++ {
		token, err := auth.GetCode()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = salesorderService.UpdateSalesorder(token, (*new)[i])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
