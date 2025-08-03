package seed

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/lib/pq"

	"payment_system/internal/database"
)

func generateRandomAddress() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func isUniqueViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}

func insertUniqueWallet() error {
	for {
		addr, err := generateRandomAddress()
		if err != nil {
			return err
		}

		_, err = database.DB.Exec("INSERT INTO wallets (addr, balance) VALUES ($1, 100.00)", addr)
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			return err
		}

		return nil
	}
}

/* Вероятность получить два одинаковых адреса ничтожно мала, но не равна нулю
поэтому адреса проверяются на уровне базы данных. Если адрес занят, адрес генерируется снова. */

func SeedWallets() error {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("SeedWallets: wallets already exist, skipping seeding")
		return nil
	}

	for i := 0; i < 10; i++ {
		err := insertUniqueWallet()
		if err != nil {
			return err
		}
	}

	log.Println("SeedWallets: successfully created 10 wallets with balance 100.00")
	return nil
}
