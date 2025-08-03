package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"payment_system/internal/database"

	"github.com/go-chi/chi/v5"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

	var balance float64
	err := database.DB.QueryRow(`SELECT balance FROM wallets WHERE addr = $1`, address).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Address string  `json:"address"`
		Balance float64 `json:"balance"`
	}{
		Address: address,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
