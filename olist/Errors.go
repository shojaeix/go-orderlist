package olist

type OrderNotFoundError struct {
	Message string
}

func (err *OrderNotFoundError) Error() string {
	return err.Message
}
