package handlers

import (
	"encoding/json"
	"kasir-api/repositories"
	"net/http"
)

type ReportHandler struct {
	repo *repositories.ReportRepository
}

func NewReportHandler(repo *repositories.ReportRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

func (h *ReportHandler) Today(w http.ResponseWriter, r *http.Request) {

	totalRevenue, totalTransaksi, name, qty, _ :=
		h.repo.TodaySummary()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_revenue": totalRevenue,
		"total_transaksi": totalTransaksi,
		"produk_terlaris": map[string]interface{}{
			"nama": name,
			"qty_terjual": qty,
		},
	})
}

func (h *ReportHandler) Report(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var totalRevenue, totalTransaksi, qty int
	var name string
	var err error

	if startDate != "" && endDate != "" {
		totalRevenue, totalTransaksi, name, qty, err =
			h.repo.SummaryByDate(startDate, endDate)
	} else {
		totalRevenue, totalTransaksi, name, qty, err =
			h.repo.TodaySummary()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_revenue": totalRevenue,
		"total_transaksi": totalTransaksi,
		"produk_terlaris": map[string]interface{}{
			"nama": name,
			"qty_terjual": qty,
		},
	})
}
