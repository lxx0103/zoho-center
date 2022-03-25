package cost

import (
	"database/sql"
)

type costRepository struct {
	tx *sql.Tx
}

func NewCostRepository(transaction *sql.Tx) CostRepository {
	return &costRepository{
		tx: transaction,
	}
}

type CostRepository interface {
	GetItems() (*[]string, error)
	ClearCost(string) error
	InsertOpeningStock(string) error
	InsertBillItems(string) error
	InsertCreditnote(string) error
	InsertAdjustmentCost(string) error
	InsertAdjustmentLog(string) error
	InsertInvoiceLog(string) error
	GetLogs() (*[]ItemLog, error)
	GetFirstCost() (*ItemCost, error)
	UpdateCost(int64, float64) error
	UpdateLog(int64, float64) error
	UpdateInvoiceItemCost() error
	UpdateInvoiceCost() error
	UpdateInvoiceCredit() error
}

func (r *costRepository) GetItems() (*[]string, error) {
	var res []string
	// rows, err := r.tx.Query(`SELECT item_id FROM items WHERE 1 = 1 AND status = ? AND item_id = '8581000000120257'`, "active")
	rows, err := r.tx.Query(`SELECT item_id FROM items `)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var itemID string
		err = rows.Scan(&itemID)
		if err != nil {
			return nil, err
		}
		res = append(res, itemID)
	}
	return &res, nil
}

func (r *costRepository) ClearCost(itemID string) error {
	_, err := r.tx.Exec(`Truncate table item_costs`)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`Truncate table item_logs`)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`Update invoice_items set cost = 0, cost_row = '' WHERE item_id = ?`, itemID)
	return err
}

func (r *costRepository) InsertOpeningStock(itemID string) error {
	_, err := r.tx.Exec(`
	INSERT into item_costs 
	(reference_id, type, date, item_id, quantity, rate, balance) 
	SELECT item_id, 'OpenStock', '2020-01-01', item_id, initial_stock, initial_stock_rate, initial_stock 
	FROM items 
	WHERE item_id = ?`, itemID)
	return err
}

func (r *costRepository) InsertBillItems(itemID string) error {
	_, err := r.tx.Exec(`
	INSERT into item_costs 
	(reference_id, type, date, item_id, quantity, rate, balance) 
	SELECT 
	b.bill_id, 
	'Bill', 
	b.date, 
	bi.item_id, 
	bi.quantity, 
	bi.rate, 
	bi.quantity 
	FROM bill_items bi
	LEFT JOIN bills b
	ON b.bill_id  = bi.bill_id 
	WHERE bi.item_id = ?`, itemID)
	return err
}

func (r *costRepository) InsertCreditnote(itemID string) error {
	_, err := r.tx.Exec(`	
	INSERT into item_costs 
	(reference_id, type, date, item_id, quantity, rate, balance) 
	SELECT 
	c.creditnote_id, 
	'Creditnote', 
	c.date, 
	ci.item_id, 
	ci.quantity, 
	i.purchase_rate , 
	ci.quantity 
	FROM creditnote_items ci
	LEFT JOIN creditnotes c 
	ON c.creditnote_id  = ci.creditnote_id 
	LEFT JOIN items i
	ON i.item_id = ci.item_id 
	WHERE ci.item_id = ?
	`, itemID)
	return err
}

func (r *costRepository) InsertAdjustmentCost(itemID string) error {
	_, err := r.tx.Exec(`	
	INSERT into item_costs 
	(reference_id, type, date, item_id, quantity, rate, balance) 
	SELECT 
	a.inventory_adjustment_id , 
	'Adjustment', 
	a.date, 
	ai.item_id, 
	ai.quantity_adjusted , 
	ai.price, 
	ai.quantity_adjusted 
	FROM adjustment_items ai
	LEFT JOIN adjustments a 
	ON a.inventory_adjustment_id  = ai.inventory_adjustment_id 
	WHERE ai.item_id = ?
	AND a.status = 'adjusted'
	AND ai.quantity_adjusted > 0`, itemID)
	return err
}

func (r *costRepository) InsertAdjustmentLog(itemID string) error {
	_, err := r.tx.Exec(`	
	INSERT into item_logs 
	(reference_id, type, date, item_id, quantity, rate) 
	SELECT 
	a.inventory_adjustment_id , 
	'Adjustment', 
	a.date, 
	ai.item_id, 
	-ai.quantity_adjusted , 
	ai.price
	FROM adjustment_items ai
	LEFT JOIN adjustments a 
	ON a.inventory_adjustment_id  = ai.inventory_adjustment_id 
	WHERE ai.item_id = ?
	AND a.status = 'adjusted'
	AND ai.quantity_adjusted < 0
	ORDER BY a.date ASC`, itemID)
	return err
}

func (r *costRepository) InsertInvoiceLog(itemID string) error {
	_, err := r.tx.Exec(`	
	INSERT into item_logs 
	(reference_id, type, date, item_id, quantity, rate) 
	SELECT 
	i.invoice_id , 
	'Invoice', 
	i.date, 
	ii.item_id, 
	ii.quantity,
	ii.rate 
	FROM invoice_items ii
	LEFT JOIN invoices i 
	ON i.invoice_id = ii.invoice_id  
	WHERE ii.item_id = ?
	AND i.status not in ('void', 'draft')
	ORDER BY i.date ASC`, itemID)
	return err
}

func (r *costRepository) GetLogs() (*[]ItemLog, error) {
	var res []ItemLog
	rows, err := r.tx.Query(`SELECT log_id, quantity FROM item_logs WHERE status = 0 ORDER BY date ASC`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var log ItemLog
		err = rows.Scan(&log.ID, &log.Quantity)
		if err != nil {
			return nil, err
		}
		res = append(res, log)
	}
	return &res, nil
}

func (r *costRepository) GetFirstCost() (*ItemCost, error) {
	var res ItemCost
	err := r.tx.QueryRow(`SELECT cost_id, balance, rate FROM item_costs WHERE balance > 0  ORDER BY date ASC`).Scan(&res.ID, &res.Balance, &res.Rate)
	return &res, err
}

func (r *costRepository) UpdateCost(costID int64, quantity float64) error {
	_, err := r.tx.Exec(`	
	UPDATE item_costs SET 
	balance = balance - ? 
	WHERE cost_id = ?`, quantity, costID)
	return err
}

func (r *costRepository) UpdateLog(logID int64, total float64) error {
	_, err := r.tx.Exec(`	
	UPDATE item_logs SET 
	status = 1,
	cost = ?/quantity
	WHERE log_id = ?`, total, logID)
	return err
}

func (r *costRepository) UpdateInvoiceItemCost() error {
	_, err := r.tx.Exec(`	
	UPDATE item_logs il 
	LEFT JOIN invoice_items ii
	ON il.reference_id = ii.invoice_id 
	AND il.item_id = ii.item_id 
	SET ii.cost = round(il.cost, 2)
	WHERE il.type = "invoice"`)
	return err
}

func (r *costRepository) UpdateInvoiceCost() error {
	_, err := r.tx.Exec(`	
	UPDATE invoices i ,
	(SELECT invoice_id, SUM(cost*quantity) AS cost FROM invoice_items GROUP BY invoice_id) s
	SET i.cost = round(s.cost, 2) 
	WHERE i.invoice_id = s.invoice_id`)
	return err
}

func (r *costRepository) UpdateInvoiceCredit() error {
	_, err := r.tx.Exec(`
		UPDATE invoices i 
		INNER JOIN 
		(
			SELECT invoice_id ,sum(total) as credits  FROM creditnotes 
			WHERE invoice_id != ''
			GROUP by invoice_id
		) t  
		ON i.invoice_id  = t.invoice_id 
		SET i.cost =  (i.cost * (i.total -t.credits)/i.total)
	`)
	return err
}
