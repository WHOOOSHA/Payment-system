package database

import (
	"database/sql"
	"fmt"
	"log"
	"payment_system/configs"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(config *configs.DatabaseConfig) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBname, config.SSLmode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	log.Println("Successful connection to the database")
}
