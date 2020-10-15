package olist

type Order struct {
	id     uint64
	Price  int32 // TODO convert to float
	Volume int32
}
