package cost

type ItemLog struct {
	ID       int64   `db:"log_id"`
	Quantity float64 `db:"quantity"`
}

type ItemCost struct {
	ID      int64   `db:"cost_id"`
	Balance float64 `db:"balance"`
	Rate    float64 `db:"rate"`
}
