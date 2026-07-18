package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	repo "github.com/davidcediel12/go-engineering/paymentprocessor/infrastructure/gorm"
	"github.com/davidcediel12/go-engineering/paymentprocessor/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// startPostgres brings up the docker-compose Postgres instance and blocks
// until it reports healthy, using the compose file's healthcheck.
func startPostgres() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("unable to determine caller information")
	}
	composeDir := filepath.Dir(filepath.Dir(filename)) // paymentprocessor/

	cmd := exec.Command("docker", "compose", "up", "-d", "--wait")
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	generateData := flag.Bool("generate", true, "Randomness for the delay")
	records := flag.Int("records", 1000, "Number of payment processes to generate")
	flag.Parse()

	if err := startPostgres(); err != nil {
		panic(fmt.Sprintf("failed to start postgres: %v", err))
	}

	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5433 sslmode=disable TimeZone=America/Bogota"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(&repo.Order{}, &repo.PaymentProcess{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
	paymentProcessorRepo := repo.NewPaymentProcessRepo(db)
	orderRepo := repo.NewOrderRepo(db)
	txManager := repo.NewTxManager(db)

	if *generateData {
		paymentProcessCreator := service.NewProcessCreator(paymentProcessorRepo, orderRepo, txManager)
		paymentProcessCreator.CreateRandomPaymentProcesses(context.Background(), *records)
	}
}
