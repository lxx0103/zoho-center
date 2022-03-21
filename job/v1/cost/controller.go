package cost

import "fmt"

func Calculate() {
	costService := NewCostService()
	err := costService.Calculate()
	fmt.Println(err)
}
