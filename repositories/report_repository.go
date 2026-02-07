package repositories

import "database/sql"

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) TodaySummary() (int, int, string, int, error) {

	var totalRevenue int
	var totalTransaksi int
	var productName string
	var qty int

	r.db.QueryRow(`
	SELECT COALESCE(SUM(total_amount),0),
	       COUNT(*)
	FROM transactions
	WHERE DATE(created_at)=CURRENT_DATE
	`).Scan(&totalRevenue, &totalTransaksi)

	r.db.QueryRow(`
	SELECT p.name, SUM(td.quantity) as qty
	FROM transaction_details td
	JOIN products p ON p.id = td.product_id
	JOIN transactions t ON t.id = td.transaction_id
	WHERE DATE(t.created_at)=CURRENT_DATE
	GROUP BY p.name
	ORDER BY qty DESC
	LIMIT 1
	`).Scan(&productName, &qty)

	return totalRevenue, totalTransaksi, productName, qty, nil
}

func (r *ReportRepository) SummaryByDate(startDate, endDate string) (int, int, string, int, error) {

	var totalRevenue int
	var totalTransaksi int
	var productName string
	var qty int

	// total revenue & transaksi
	err := r.db.QueryRow(`
	SELECT
		COALESCE(SUM(total_amount),0),
		COUNT(*)
	FROM transactions
	WHERE DATE(created_at) BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)

	if err != nil {
		return 0, 0, "", 0, err
	}

	// produk terlaris
	_ = r.db.QueryRow(`
	SELECT p.name, SUM(td.quantity) as qty
	FROM transaction_details td
	JOIN products p ON p.id = td.product_id
	JOIN transactions t ON t.id = td.transaction_id
	WHERE DATE(t.created_at) BETWEEN $1 AND $2
	GROUP BY p.name
	ORDER BY qty DESC
	LIMIT 1
	`, startDate, endDate).Scan(&productName, &qty)

	return totalRevenue, totalTransaksi, productName, qty, nil
}
