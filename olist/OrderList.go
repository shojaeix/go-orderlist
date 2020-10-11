package olist

import (
	"errors"
	"fmt"
)

type OrderList struct {
	plusPointer                 *[]*[]Order
	negativePointer             *[]*[]Order
	indexDispute                int // price + indexDispute = index
	greatestIndex               int
	smallestIndex				int
	newOrdersChannel            *chan Order
	newOrdersFlagUpdaterChannel *chan Order
	ordersById                  []*Order
	ordersLastId                uint
}

type orderArray []Order

func (ol *OrderList) Run() {
	// go process new orders
	// go update flags after adding new order
}
func (ol *OrderList) AddOrder(order Order) (uint, error) {
	// validate values
	if !ol.validateOrderValues(&order) {
		return 0, errors.New("Order's values are invalid.")
	}

	// store order and set an ID
	ol.accociateIdToOrder(&order)

	// send order to NewOrdersChannel
	go ol.pushOrderToArray(&order)
	
	// return ID
	return order.id, nil
}

func (ol *OrderList) validateOrderValues (order *Order) bool {
	return order.Price > 0 && order.Volume > 0;
}

func (ol *OrderList) accociateIdToOrder(order *Order){
	// associate an ID  to order
	ol.ordersById = append(ol.ordersById, order)
	ol.ordersLastId++
	order.id = ol.ordersLastId
}
func (ol *OrderList) initiateArrays(indexDispute int){

	// make ordersById array
	if ol.ordersById == nil {
		ol.ordersById = make([]*Order, 0, 100000)
		ol.ordersLastId = 0
	}

	// positive list
	if ol.plusPointer == nil {
		list := make([]*[]Order, 100)
		ol.plusPointer = &list
		ol.indexDispute = indexDispute
		ol.greatestIndex = 99
	}

	// negative list
	if ol.negativePointer == nil {
		list := make([]*[]Order, 100)
		ol.negativePointer = &list
		ol.smallestIndex = 99
	}
}
func (ol *OrderList) pushOrderToArray(order *Order){

	ol.initiateArrays(order.Price)


	//
	index := order.Price - ol.indexDispute
	const incrementStep = 1000


	if index >= 0 { // positive

		// increase list size if necessary
		if index > ol.greatestIndex {
			newList := make([]*[]Order, index+incrementStep)
			ol.greatestIndex = index + incrementStep - 1
			copy(newList, *ol.plusPointer)
			ol.plusPointer = &newList
		}

		list := *ol.plusPointer
		// fill list[index] if necessary
		if list[index] == nil {
			var orderArr []Order = make([]Order, incrementStep)
			list[index] = &orderArr
		}

		// append order
		indexOrderList := append(*list[index], *order)
		list[index] = &indexOrderList
	} else { // negative

		absoluteIndex := index * -1

		// increase list size if necessary
		if absoluteIndex > ol.smallestIndex {
			newList := make([]*[]Order, absoluteIndex+incrementStep)
			ol.smallestIndex = absoluteIndex + incrementStep - 1
			copy(newList, *ol.negativePointer)
			ol.negativePointer = &newList
		}

		list := *ol.negativePointer
		// fill list[index] if necessary
		if list[absoluteIndex - 1] == nil {
			var orderArr = make([]Order, incrementStep)
			list[absoluteIndex - 1] = &orderArr
		}

		// append order
		indexOrderList := append(*list[absoluteIndex - 1], *order)
		list[absoluteIndex - 1] = &indexOrderList
	}

}

func (ol *OrderList) DeleteOrder(id uint) bool {

	if id > ol.ordersLastId {
		println("id is bigger than last id")
		return false
	}

	if ol.ordersById[id-1] == nil {
		fmt.Printf("Ids length: %d\n", len(ol.ordersById))
		panic("id is nil")
		return false
	}

	order := *ol.ordersById[id-1]
	println("plusPointer: ", order.id)

	// search in main list
	index := order.Price - ol.indexDispute

	orderArr := (*ol.plusPointer)[index]

	for key, value := range *orderArr {
		if value.id == order.id {
			println((*orderArr)[key].id)
			(*orderArr)[key] = Order{}
			ol.ordersById[id-1] = nil
			return true
		}
	}

	return false
}

func (ol *OrderList) PrintAll(printOrders bool) {

	fmt.Printf("Ids length: %d, capacity: %d\n", len(ol.ordersById), cap(ol.ordersById))
	if ol.plusPointer != nil {


		lenLvl2 := 0
		capLvl2 := 0

			for index, orderArr := range *ol.plusPointer {
				if printOrders {
					fmt.Printf("Index: %d, The price: %d\n", index, index+ol.indexDispute)
				}
				if orderArr != nil {
					lenLvl2 += len(*orderArr)
					capLvl2 += cap(*orderArr)
					if printOrders {
						for _, orderPointer := range *orderArr {
							if orderPointer.id != 0 {
								fmt.Printf("Order details: Id: %d, price: %d, volume: %d \n", orderPointer.id, orderPointer.Price, orderPointer.Volume)
							}
						}
					}
				}

		}
		fmt.Printf("Positive Length: %d, Capacity %d Length lvl 2: %d Capacity lvl 2: %d\n", len(*ol.plusPointer), cap(*ol.plusPointer), lenLvl2, capLvl2)

	}

	if ol.negativePointer != nil {

		lenLvl2 := 0
		capLvl2 := 0

		for index, orderArr := range *ol.negativePointer {
			if printOrders {
				fmt.Printf("Index: %d, The price: %d\n", index, index+ol.indexDispute-1)
			}
			if orderArr != nil {
				lenLvl2 += len(*orderArr)
				capLvl2 += cap(*orderArr)
				if printOrders {
					for _, orderPointer := range *orderArr {
						if orderPointer.id != 0 {
							fmt.Printf("Negative Order details: Id: %d, price: %d, volume: %d \n", orderPointer.id, orderPointer.Price, orderPointer.Volume)
						}
					}
				}
			}
		}

		fmt.Printf("Negative Length: %d, Capacity %d  Length lvl 2: %d Capacity lvl 2: %d\n", len(*ol.negativePointer), cap(*ol.negativePointer), lenLvl2, capLvl2)
	}
}

