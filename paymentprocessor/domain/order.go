package domain

type Order struct {
	ID         uint
	Amount     int64
	PaidAmount int64
}

func NewOrder(amount int64) Order {
	return Order{
		Amount: amount,
	}
}
