package orderprocessor

type Order struct {
	orderID     uint
	userID      uint
	totalAmount uint
	items       []string
}
