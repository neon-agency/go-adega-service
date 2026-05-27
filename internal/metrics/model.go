package metrics

type Overview struct {
	EntriesCents     int `json:"entries_cents"`
	SalesCents       int `json:"sales_cents"`
	CostOfGoodsCents int `json:"cost_of_goods_cents"`
	GrossProfitCents int `json:"gross_profit_cents"`
	NetProfitCents   int `json:"net_profit_cents"`
	OrdersCount      int `json:"orders_count"`
}
