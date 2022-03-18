package cost

import "fmt"

func Calculate() {
	// fmt.Println("STEP2: INSERT ALL THE OPENING STOCK")
	// fmt.Println("STEP3: INSERT ALL THE BILLS")
	// fmt.Println("STEP4: INSERT ALL THE ADJUSTMENTS GREATER THAN 0")
	// fmt.Println("STEP5: INSERT ALL INVOICE WHICH STATUS ARE RIGHT")
	// fmt.Println("STEP6: INSERT ALL THE ADJUSTMENTS LESS THAN 0")
	costService := NewCostService()
	err := costService.Calculate()
	fmt.Println("AAAA")
	fmt.Println(err)
}
