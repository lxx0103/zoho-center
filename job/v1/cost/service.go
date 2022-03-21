package cost

import (
	"errors"
	"fmt"
	"zoho-center/core/database"
)

type costService struct {
}

func NewCostService() CostService {
	return &costService{}
}

type CostService interface {
	Calculate() error
}

func (s costService) Calculate() error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误:" + err.Error()
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewCostRepository(tx)
	items, err := repo.GetItems()
	if err != nil {
		msg := "获取商品错误:" + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	for i := 0; i < len(*items); i++ {
		err = itemCalculate((*items)[i])
		if err != nil {
			return err
		}
	}
	tx2, err := db.Begin()
	if err != nil {
		msg := "事务开启错误:" + err.Error()
		return errors.New(msg)
	}
	defer tx2.Rollback()
	repo2 := NewCostRepository(tx2)
	err = repo2.UpdateInvoiceCost()
	if err != nil {
		msg := "更新Invoice成本错误:" + err.Error()
		return errors.New(msg)
	}
	tx2.Commit()
	return nil
	// return false, nil
}

func itemCalculate(itemID string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误:" + err.Error()
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewCostRepository(tx)
	err = repo.ClearCost(itemID)
	if err != nil {
		msg := itemID + "清空数据错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertOpeningStock(itemID)
	if err != nil {
		msg := itemID + "初始库存错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertBillItems(itemID)
	if err != nil {
		msg := itemID + "BILL错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertCreditnote(itemID)
	if err != nil {
		msg := itemID + "Creditnote错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertAdjustmentCost(itemID)
	if err != nil {
		msg := itemID + "AdjustmentCost错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertAdjustmentLog(itemID)
	if err != nil {
		msg := itemID + "AdjustmentLog错误:" + err.Error()
		return errors.New(msg)
	}
	err = repo.InsertInvoiceLog(itemID)
	if err != nil {
		msg := itemID + "Invoice错误:" + err.Error()
		return errors.New(msg)
	}
	logs, err := repo.GetLogs()
	if err != nil {
		msg := "获取LOGS错误:" + err.Error()
		return errors.New(msg)
	}
	for j := 0; j < len(*logs); j++ {
		// for j := 0; j < 5; j++ {
		fmt.Println((*logs)[j])
		var remainQty int
		var nowTotalCost float64
		remainQty = (*logs)[j].Quantity
		nowTotalCost = 0
		for remainQty > 0 {
			firstCost, err := repo.GetFirstCost()
			fmt.Println(firstCost)
			if err != nil {
				msg := "获取COST错误:" + err.Error()
				return errors.New(msg)
			}
			if firstCost.Balance >= (*logs)[j].Quantity {
				err = repo.UpdateCost(firstCost.ID, remainQty)
				if err != nil {
					msg := "扣除Balance错误:" + err.Error()
					return errors.New(msg)
				}
				nowTotalCost = nowTotalCost + firstCost.Rate*float64(remainQty)
				remainQty = 0
			} else {
				err = repo.UpdateCost(firstCost.ID, firstCost.Balance)
				if err != nil {
					msg := "扣除Balance错误:" + err.Error()
					return errors.New(msg)
				}
				nowTotalCost = nowTotalCost + firstCost.Rate*float64(firstCost.Balance)
				remainQty = remainQty - firstCost.Balance
			}
		}
		fmt.Println(nowTotalCost)
		err := repo.UpdateLog((*logs)[j].ID, nowTotalCost)
		if err != nil {
			msg := "更新LOG错误:" + err.Error()
			return errors.New(msg)
		}
	}
	err = repo.UpdateInvoiceItemCost()
	if err != nil {
		msg := "更新InvoiceItem成本错误:" + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}
