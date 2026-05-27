package reports

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, period string) (Data, error) {
	since := sinceFor(period)
	data := Data{}

	series, err := r.series(ctx, since)
	if err != nil {
		return Data{}, err
	}
	data.Series = series

	payments, err := r.payments(ctx, since)
	if err != nil {
		return Data{}, err
	}
	data.Payments = payments

	top, err := r.top(ctx, since)
	if err != nil {
		return Data{}, err
	}
	data.Top = top

	byCategory, err := r.byCategory(ctx, since)
	if err != nil {
		return Data{}, err
	}
	data.ByCategory = byCategory

	return data, nil
}

func (r *Repository) series(ctx context.Context, since time.Time) ([]SeriesPoint, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT to_char(date_trunc('day', created_at), 'DD/MM') AS label,
		       COALESCE(SUM(total_cents), 0), COUNT(*)
		FROM orders
		WHERE status <> 'canceled' AND created_at >= $1
		GROUP BY date_trunc('day', created_at)
		ORDER BY date_trunc('day', created_at)`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]SeriesPoint, 0)
	for rows.Next() {
		var p SeriesPoint
		var cents int
		if err := rows.Scan(&p.Label, &cents, &p.Orders); err != nil {
			return nil, err
		}
		p.Value = centsToBRL(cents)
		p.Target = p.Value
		if p.Target < 100 {
			p.Target = 100
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) payments(ctx context.Context, since time.Time) ([]PaymentSlice, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT payment_method::text, COUNT(*)
		FROM orders
		WHERE status <> 'canceled' AND created_at >= $1
		GROUP BY payment_method::text`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{}
	total := 0
	for rows.Next() {
		var method string
		var count int
		if err := rows.Scan(&method, &count); err != nil {
			return nil, err
		}
		counts[method] += count
		total += count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if total == 0 {
		return []PaymentSlice{}, nil
	}
	return []PaymentSlice{
		{Label: "PIX", Pct: pct(counts["pix"], total), Color: "var(--wine)"},
		{Label: "Cartão", Pct: pct(counts["credit_card"]+counts["debit_card"]+counts["card"], total), Color: "var(--gold)"},
		{Label: "Dinheiro", Pct: pct(counts["cash"], total), Color: "var(--ink-3)"},
	}, nil
}

func (r *Repository) top(ctx context.Context, since time.Time) ([]TopProduct, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(p.id::text, ''), oi.product_name, COALESCE(p.image_url, ''), oi.unit_price_cents,
		       SUM(oi.quantity), SUM(oi.total_cents)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.status <> 'canceled' AND o.created_at >= $1
		GROUP BY p.id, oi.product_name, p.image_url, oi.unit_price_cents
		ORDER BY SUM(oi.quantity) DESC
		LIMIT 5`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]TopProduct, 0)
	for rows.Next() {
		var p TopProduct
		var price, revenue int
		if err := rows.Scan(&p.ID, &p.Name, &p.Swatch, &price, &p.Sold, &revenue); err != nil {
			return nil, err
		}
		p.Price = centsToBRL(price)
		p.Revenue = centsToBRL(revenue)
		if strings.TrimSpace(p.Swatch) == "" {
			p.Swatch = "#6E1F2C"
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) byCategory(ctx context.Context, since time.Time) ([]CategoryReport, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(p.category, 'sem-categoria'), COALESCE(p.category, 'Sem categoria'),
		       SUM(oi.quantity), SUM(oi.total_cents)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.status <> 'canceled' AND o.created_at >= $1
		GROUP BY p.category
		ORDER BY SUM(oi.total_cents) DESC`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CategoryReport, 0)
	for rows.Next() {
		var c CategoryReport
		var revenue int
		if err := rows.Scan(&c.ID, &c.Label, &c.Sold, &revenue); err != nil {
			return nil, err
		}
		c.Revenue = centsToBRL(revenue)
		out = append(out, c)
	}
	return out, rows.Err()
}

func sinceFor(period string) time.Time {
	now := time.Now()
	switch period {
	case "today":
		y, m, d := now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	case "7d":
		return now.AddDate(0, 0, -7)
	case "90d":
		return now.AddDate(0, 0, -90)
	default:
		return now.AddDate(0, 0, -30)
	}
}

func centsToBRL(cents int) float64 {
	return float64(cents) / 100
}

func pct(part, total int) int {
	if total == 0 {
		return 0
	}
	return int(float64(part)/float64(total)*100 + 0.5)
}
