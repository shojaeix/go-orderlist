package olist

type Order struct {
	id     uint64
	Price  int32 // TODO convert to float
	Volume int32
}

// Positive number, or zero if id has been not set
func (o *Order) GetId() uint64 {
	return o.id
}