package olist

import (
	"errors"
	"fmt"
	"unsafe"
)

const listsIncrementStep = 10000


type OrderList struct {
	positiveList                *[]*[]Order
	negativeList                *[]*[]Order
	indexDispute                int32 // only positive.  price - indexDispute = index
	greatestIndex               int32
	smallestIndex               int32
	newOrdersChannel            *chan *Order
	processingNewOrders         bool
	newOrdersFlagUpdaterChannel *chan *Order
	ordersById                  []*Order
	ordersLastId                uint64
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

func (ol *OrderList) processNewOrders() {
	// get new orders from channel
	for order := range *ol.newOrdersChannel {
		// process
		ol.pushOrderToArray(order)
	}
}

// Add an order to list and associate an unique uint ID
// return id, error
func (ol *OrderList) AddOrder(order Order) (id uint64, err error) {
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
func (ol *OrderList) validateOrderValues(order *Order) bool {
	return order.Price > 0 && order.Volume > 0
}

func (ol *OrderList) associateIdToOrder(order *Order) {
	// associate an ID  to order
	ol.ordersById = append(ol.ordersById, order)
	ol.ordersLastId++
	order.id = ol.ordersLastId
}

func (ol *OrderList) initiate() {

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

func (ol *OrderList) pushOrderToArray(order *Order) {
	// set dispute if it's the first order
	if order.id == 1 {
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
		absoluteIndex := -index

		// increase list size if necessary
		if absoluteIndex > ol.smallestIndex {
			newList := make([]*[]Order, absoluteIndex+listsIncrementStep)
			ol.smallestIndex = absoluteIndex + listsIncrementStep - 1
			copy(newList, *ol.negativeList)
			ol.negativeList = &newList
		}

		list := *ol.negativeList

		// fill list[index] if necessary
		if list[absoluteIndex-1] == nil {
			var orderArr = make([]Order, listsIncrementStep)
			list[absoluteIndex-1] = &orderArr
		}

		// append order
		indexOrderList := append(*list[absoluteIndex-1], *order)
		list[absoluteIndex-1] = &indexOrderList
	}

}

func (ol *OrderList) DeleteOrder(id uint64) bool {

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

	return false
}

func (ol *OrderList) PrintAll(printOrders bool, printIndices bool) {

	fmt.Printf("Ids length: %d, capacity: %d, size: %d\n", len(ol.ordersById), cap(ol.ordersById), unsafe.Sizeof(ol.ordersById))

	totalCap := 0

	if ol.positiveList != nil {

		lenLvl2 := 0
		capLvl2 := 0
		var sizeLvl2 uintptr = 0
		var orderArr *[]Order
		for i := len(*ol.positiveList) - 1; i >= 0; i-- {
			orderArr = (*ol.positiveList)[i]
			if printOrders || printIndices {
				fmt.Printf("Index: %d, The price: %d\n", i, int64(i)+int64(ol.indexDispute))
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
		sizeLvl2 += (unsafe.Sizeof([61440000]Order{})) / 1024 / 1024
		fmt.Printf("Positive Length: %d, Capacity %d Length lvl 2: %d Capacity lvl 2: %d, size: %d\n", len(*ol.positiveList), cap(*ol.positiveList), lenLvl2, capLvl2, sizeLvl2)
		totalCap += capLvl2
	}

	if ol.negativeList != nil {

		lenLvl2 := 0
		capLvl2 := 0

		var index int64 = 0
		var orderArr *[]Order
		var key int
		for key, orderArr = range *ol.negativeList {
			index = int64(key)
			if printOrders || printIndices {
				fmt.Printf("Index: %d, The price: %d\n", index, (index)+int64(ol.indexDispute)-1)
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
		totalCap += capLvl2
	}
	println("Total cap: ", totalCap)
}

func (ol *OrderList) GetRowAndAheadVolume(id uint64) (uint64, uint64) {
	if !ol.running || id == 0 || ol.ordersLastId < id || ol.ordersById[id-1] == nil {
		return 0, 0
	}

	// find the order in lists
	order := ol.ordersById[id-1]
	index := order.Price - ol.indexDispute
	var row uint64 = 0
	var volume uint64 = 0

	if index >= 0 {
		if index > ol.greatestIndex {
			return 0, 0
		}

		// count all items in positiveList from index 0 until (index of the order - 1)
		var i int32 = 0
		for ; i < index; i++ {
			if (*ol.positiveList)[i] != nil {
				for _, itemPointer := range *(*ol.positiveList)[i] {
					if itemPointer.id > 0 {
						//print(itemPointer.Price, "   ")
						row++
						volume += uint64(itemPointer.Volume)
					}
				}
			}
		}

		// include items with the same price
		if (*ol.positiveList)[index] == nil {
			return 0, 0
		}
		for _, itemPointer := range *(*ol.positiveList)[index] {
			if itemPointer.id == id {
				break
			}
			row++
			volume += uint64(itemPointer.Volume)
		}
		return row, volume
	} else {
		absoluteIndex := -index
		if absoluteIndex > ol.smallestIndex {
			return 0, 0
		}
		for i := absoluteIndex + 1 - 1; i <= ol.smallestIndex; i++ {
			if (*ol.negativeList)[i] != nil {
				for _, itemPointer := range *(*ol.negativeList)[i] {
					if itemPointer.id > 0 {
						row++
						volume += uint64(itemPointer.Volume)
					}
				}
			}
		}

		// include items with the same price
		if (*ol.negativeList)[absoluteIndex-1] == nil {
			return 0, 0
		}
		for _, itemPointer := range *(*ol.negativeList)[absoluteIndex-1] {
			if itemPointer.id > 0 {
				row++
				volume += uint64(itemPointer.Volume)
			}

			if itemPointer.id == id {
				break
			}
		}
	}

	return row, volume
}
