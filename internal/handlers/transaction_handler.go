package handlers

import (
	"encoding/json"
	"net/http"
	"payment_system/internal/database"
	"strconv"
	"time"
)

type SendRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

/*Если включена проверка подписи, тело данного запроса не должно
содержать пробелов, табуляции и переносов строк, чтобы не испортить
подпись.
Пример: {"from":"92d962ac5ae2c733e550883baf32012e5747ce8d","to":"ddff98c499af435cb739e7a33b5501e43f731771","amount":15}
(первые кошельки генерируются случайно и данный запрос не сработает (вероятность невероятно мала))
Если проверка подписи отключена, то данное требование не актуально
*/

func Send(w http.ResponseWriter, r *http.Request) {
	var req SendRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	/*в случае если обработка перевода прервется не дойдя до конца,
	никаких изменений в базе данных быть не должно. Для этого используются транзакции.*/
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction begin failed", http.StatusInternalServerError)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
	}()

	var fromID, toID int
	var fromBalance float64

	err = tx.QueryRow(`SELECT id, balance FROM wallets WHERE addr = $1 FOR UPDATE`, req.From).Scan(&fromID, &fromBalance)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Sender wallet not found", http.StatusNotFound)
		return
	}

	if fromBalance < req.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	err = tx.QueryRow(`SELECT id FROM wallets WHERE addr = $1 FOR UPDATE`, req.To).Scan(&toID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Recipient wallet not found", http.StatusNotFound)
		return
	}

	_, err = tx.Exec(`UPDATE wallets SET balance = balance - $1 WHERE id = $2`, req.Amount, fromID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update sender balance", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, req.Amount, toID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update recipient balance", http.StatusInternalServerError)
		return
	}

	/*СУБД автоматически заполнит поле created_at*/
	_, err = tx.Exec(`INSERT INTO transfers (id_from, id_to, amount) VALUES ($1, $2, $3)`, fromID, toID, req.Amount)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to record transfer", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type TransferInfo struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func GetLastTransactions(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			f.addr as from_addr, 
			t.addr as to_addr, 
			tr.amount, 
			tr.created_at
		FROM transfers tr
		JOIN wallets f ON f.id = tr.id_from
		JOIN wallets t ON t.id = tr.id_to
		ORDER BY tr.created_at DESC
		LIMIT $1
	`

	rows, err := database.DB.Query(query, count)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transfers []TransferInfo
	for rows.Next() {
		var t TransferInfo
		err := rows.Scan(&t.From, &t.To, &t.Amount, &t.CreatedAt)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		transfers = append(transfers, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transfers)
}
