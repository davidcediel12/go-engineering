package retrybackoff

import (
	"fmt"
	"log"
	"math/rand/v2"
)

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	PerformOperation() error
}

type ServiceImp struct{}

func (s *ServiceImp) PerformOperation() error {
	if rand.IntN(15) == 0 {
		log.Printf("Succeed c:")
		return nil
	}
	return fmt.Errorf("Oops :c")
}
