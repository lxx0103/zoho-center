package cost

type ItemLog struct {
	ID       int64 `db:"log_id"`
	Quantity int   `db:"quantity"`
}

type ItemCost struct {
	ID      int64   `db:"cost_id"`
	Balance int     `db:"balance"`
	Rate    float64 `db:"rate"`
}
