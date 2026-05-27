package reports

type SeriesPoint struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Orders int     `json:"orders"`
	Target float64 `json:"target"`
}

type PaymentSlice struct {
	Label string `json:"label"`
	Pct   int    `json:"pct"`
	Color string `json:"color"`
}

type TopProduct struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Swatch  string  `json:"swatch"`
	Price   float64 `json:"price"`
	Sold    int     `json:"sold"`
	Revenue float64 `json:"revenue"`
}

type CategoryReport struct {
	ID      string  `json:"id"`
	Label   string  `json:"label"`
	Sold    int     `json:"sold"`
	Revenue float64 `json:"revenue"`
}

type Data struct {
	Series     []SeriesPoint    `json:"series"`
	Payments   []PaymentSlice   `json:"payments"`
	Top        []TopProduct     `json:"top"`
	ByCategory []CategoryReport `json:"byCategory"`
}
