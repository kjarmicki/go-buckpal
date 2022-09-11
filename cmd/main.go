package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	account_adapter_in_web "github.com/kjarmicki/go-buckpal/pkg/account/adapter/in/web"
	account_adapter_out_lock "github.com/kjarmicki/go-buckpal/pkg/account/adapter/out/lock"
	account_adapter_out_persistence "github.com/kjarmicki/go-buckpal/pkg/account/adapter/out/persistence"
	account_application_service "github.com/kjarmicki/go-buckpal/pkg/account/application/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := connectToDb()
	activityWindowRepository := account_adapter_out_persistence.NewActivityWindowGormMysqlRepository(db)
	accountRepository := account_adapter_out_persistence.NewAccountGormMySqlRepository(db, activityWindowRepository)
	accountPersistenceAdapter := account_adapter_out_persistence.NewAccountPersistenceAdapter(accountRepository, activityWindowRepository)
	sendMoneyService := account_application_service.NewSendMoneyService(accountPersistenceAdapter, &account_adapter_out_lock.AccountNoopLockAdapter{}, accountPersistenceAdapter)
	sendMoneyController := account_adapter_in_web.NewSendMoneyController(sendMoneyService)

	router := mux.NewRouter()
	sendMoneyController.AttachToRouter(router)

	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
}

func connectToDb() *gorm.DB {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		panic("Database DSN not found in env variable DB_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return db
}
