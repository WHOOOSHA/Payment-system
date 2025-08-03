package main

import (
	"log"
	"net/http"
	"payment_system/configs"
	"payment_system/internal/database"
	"payment_system/internal/router"
	"payment_system/internal/seed"
)

func main() {
	DBconfig, err := configs.LoadDatabaseConfig()
	if err != nil {
		log.Fatal("Cannot read database config: ", err)
	}

	//Создание необходимых таблиц
	database.RunMigrations(DBconfig)
	database.ConnectDB(DBconfig)

	//При первом запуске создаётся 10 кошельков с случайными адресами и 100.0 у.е. на счету
	err = seed.SeedWallets()
	if err != nil {
		log.Fatal("Faileds to seed wallets: ", err)
	}

	httpConfig, err := configs.LoadHTTPConfig()
	if err != nil {
		log.Fatal("Cannot read http config:", err)
	}

	/*Для вашего удобства проверка подписи отключена
	  чтобы включить сделать useAuth равным true
	  как получить подпись см. в Readme.me*/
	useAuth := false
	r := router.NewRouter(httpConfig.Secret, useAuth)

	log.Println("Server started")
	err = http.ListenAndServe(httpConfig.Host+":"+httpConfig.Port, r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
