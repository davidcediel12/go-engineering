package main

import (
	"context"
	"flag"

	repo "github.com/davidcediel12/go-engineering/paymentprocessor/infrastructure/gorm"
	"github.com/davidcediel12/go-engineering/paymentprocessor/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	generateData := flag.Bool("generate", true, "Randomness for the delay")
	records := flag.Int("records", 1000, "Number of payment processes to generate")
	flag.Parse()

	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=America/Bogota"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	paymentProcessorRepo := repo.NewPaymentProcessRepo(db)
	orderRepo := repo.NewOrderRepo(db)
	txManager := repo.NewTxManager(db)

	if *generateData {
		paymentProcessCreator := service.NewProcessCreator(paymentProcessorRepo, orderRepo, txManager)
		paymentProcessCreator.CreateRandomPaymentProcesses(context.Background(), *records)
	}
}
