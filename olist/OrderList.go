package olist

import (
	"errors"
	"fmt"
)

const listsIncrementStep = 10000


type OrderList struct {
	positiveList                *[]*[]Order
	negativeList                *[]*[]Order
	indexDispute                int // price - indexDispute = index
	greatestIndex               int
	smallestIndex               int
	newOrdersChannel            *chan *Order
	newOrdersFlagUpdaterChannel *chan *Order
	ordersById                  []*Order
	ordersLastId                uint
	running                     bool
}

func (ol *OrderList) run() {

	ol.initiate()

	// go process new orders
	go ol.processNewOrders()

	// go update flags after adding new order

	// set a boolean flag after run
	ol.running = true
}
func (ol *OrderList) processNewOrders(){
	// get new orders from channel
	for order := range *ol.newOrdersChannel {
		// process
		ol.pushOrderToArray(order)
	}
}

// Add an order to list and associate an unique uint ID
// return id, error
func (ol *OrderList) AddOrder(order Order) (uint, error) {
	// run
	if !ol.running {
		ol.run()
	}

	// validate the order's values
	if !ol.validateOrderValues(&order) {
		return 0, errors.New("the order's values are invalid")
	}

	// store order and set an ID
	ol.associateIdToOrder(&order)

	// send order to NewOrdersChannel to process
	*ol.newOrdersChannel <- &order

	// return ID
	return order.id, nil
}

// check order's price & volume to be positive number
func (ol *OrderList) validateOrderValues (order *Order) bool {
	return order.Price > 0 && order.Volume > 0;
}

func (ol *OrderList) associateIdToOrder(order *Order){
	// associate an ID  to order
	ol.ordersById = append(ol.ordersById, order)
	ol.ordersLastId++
	order.id = ol.ordersLastId
}

func (ol *OrderList) initiate(){

	// make new orders channel
	newOrdersChannel := make(chan *Order, 10000)
	ol.newOrdersChannel = &newOrdersChannel


	// make new orders' flags channel
	newOrdersFlagUpdaterChannel := make(chan *Order, 10000)
	ol.newOrdersFlagUpdaterChannel = &newOrdersFlagUpdaterChannel

	// make ordersById array
	if ol.ordersById == nil {
		ol.ordersById = make([]*Order, 0, 100000)
		ol.ordersLastId = 0
	}

	// set initial value for dispute
	ol.indexDispute = -1

	// make the positive list
	if ol.positiveList == nil {
		list := make([]*[]Order, 100)
		ol.positiveList = &list
		ol.greatestIndex = 99
	}

	// make the negative list
	if ol.negativeList == nil {
		list := make([]*[]Order, 100)
		ol.negativeList = &list
		ol.smallestIndex = 99
	}
}
func (ol *OrderList) pushOrderToArray(order *Order){

	// set dispute if it has not been set before
	if ol.indexDispute == -1 {
		ol.indexDispute = order.Price
	}

	//
	index := order.Price - ol.indexDispute

	if index >= 0 { // positive indices

		// increase positiveList size if necessary
		if index > ol.greatestIndex {
			newList := make([]*[]Order, index+listsIncrementStep)
			ol.greatestIndex = index + listsIncrementStep - 1
			copy(newList, *ol.positiveList)
			ol.positiveList = &newList
		}

		positiveList := *ol.positiveList

		// fill positiveList[index] if necessary
		if positiveList[index] == nil {
			var ordersArr = make([]Order, listsIncrementStep)
			positiveList[index] = &ordersArr
		}

		// append order
		indexOrderList := append(*positiveList[index], *order)
		positiveList[index] = &indexOrderList

	} else { // negative indices

		// absolute the index, for saving in array
		absoluteIndex := index * -1

		// increase list size if necessary
		if absoluteIndex > ol.smallestIndex {
			newList := make([]*[]Order, absoluteIndex+listsIncrementStep)
			ol.smallestIndex = absoluteIndex + listsIncrementStep - 1
			copy(newList, *ol.negativeList)
			ol.negativeList = &newList
		}

		list := *ol.negativeList
		// fill list[index] if necessary
		if list[absoluteIndex - 1] == nil {
			var orderArr = make([]Order, listsIncrementStep)
			list[absoluteIndex - 1] = &orderArr
		}

		// append order
		indexOrderList := append(*list[absoluteIndex - 1], *order)
		list[absoluteIndex - 1] = &indexOrderList
	}

}

func (ol *OrderList) DeleteOrder(id uint) bool {

	if id > ol.ordersLastId { // couldn't be exist
		return false
	}

	if ol.ordersById[id-1] == nil { // not exists
		return false
	}

	// get order info
	order := *ol.ordersById[id-1]

	// search in main list
	index := order.Price - ol.indexDispute

	orderArr := (*ol.positiveList)[index]

	for key, value := range *orderArr {
		if value.id == order.id {
			println((*orderArr)[key].id)
			(*orderArr)[key] = Order{}
			ol.ordersById[id-1] = nil
			return true
		}
	}

	panic("order not found in orders list")
	return false
}

func (ol *OrderList) PrintAll(printOrders bool) {

	fmt.Printf("Ids length: %d, capacity: %d\n", len(ol.ordersById), cap(ol.ordersById))
	if ol.positiveList != nil {


		lenLvl2 := 0
		capLvl2 := 0

			for index, orderArr := range *ol.positiveList {
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
		fmt.Printf("Positive Length: %d, Capacity %d Length lvl 2: %d Capacity lvl 2: %d\n", len(*ol.positiveList), cap(*ol.positiveList), lenLvl2, capLvl2)

	}

	if ol.negativeList != nil {

		lenLvl2 := 0
		capLvl2 := 0

		for index, orderArr := range *ol.negativeList {
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

		fmt.Printf("Negative Length: %d, Capacity %d  Length lvl 2: %d Capacity lvl 2: %d\n", len(*ol.negativeList), cap(*ol.negativeList), lenLvl2, capLvl2)
	}
}

