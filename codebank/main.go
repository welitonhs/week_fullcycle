package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/welitonhs/fccodebank/domain"
	"github.com/welitonhs/fccodebank/domain/infrasctructure/repository"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Weliton"
	cc.ExpirationYear = 2022
	cc.ExpirationMonth = 8
	cc.CVV = 123
	cc.Limit = 1200
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

// func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
// 	transactionRepository := repository.NewTransactionRepositoryDb(db)
// 	useCase := usecase.NewUseCaseTransaction(transactionRepository)
// 	return useCase
// }

func setupDb() *sql.DB {
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"host.docker.internal",
		"5432",
		"postgres",
		"root",
		"codebank",
	)

	db, err := sql.Open("postgres", pgsqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}
