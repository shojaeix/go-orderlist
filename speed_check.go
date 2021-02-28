package main

import (
	"errors"
	"fmt"
	"go-order-list/olist"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func main() {
	add_samples_and_show_performance_details()
}

func add_a_simple() {
	bidsList := olist.OrderList{}

	_, err := bidsList.AddOrder(olist.Order{
		Price:  1499,
		Volume: 5,
	})

	if err != nil {
		panic(err)
	}
}

func add_some_samples_and_show_developer_parameters(){
	bidsList := olist.OrderList{}

	var count = 10000
	var repeat int64 = 3

	var t int64
	for t = 1; t <= repeat; t++ {
		for i := 1; i <= count; i++ {
			_, _ = bidsList.AddOrder(olist.Order{
				Price: randBetween(1500, 1600),
				//Volume: (randBetween(0, 12000000)),
				Volume: 1,
			})

		}
	}

	println(bidsList.GetLowestAndHighestPrices())
	bidsList.SetSort("ASC")
	edgeOrder, err := bidsList.GetEdgeOrder()
	println("Edge order Price: ", edgeOrder.Price)

	if errors.Is(err, olist.OrderNotFoundError{}) {
		println("Order not found")
	} else if err != nil {
		println(err.Error())
	}
	//id, _ := bidsList.AddOrder(olist.Order{
	//	Price:  1498,
	//	Volume: 1,
	//})
	//
	//bidsList.AddOrder(olist.Order{
	//	Price:  1499,
	//	Volume: 1,
	//})

	////time.Sleep(time.Second * 2)
	//row, volume := bidsList.GetRowAndAheadVolume(id)
	//fmt.Printf("Id: %d, Row: %d, Volume: %d\n", id, row, volume)
	//bidsList.PrintAll(false, false)
}

func add_samples_and_show_performance_details() {

	baseMem := getUsedMemoryByByte()
	printer := message.NewPrinter(language.English)

	bidsList := olist.OrderList{}

	var count = 1000000
	var repeat int64 = 30

	var totalSamples = repeat * int64(count)
	printer.Println("Config:")
	printer.Printf("\tTotal samples to create = %d million\n", (int64(count) * repeat)/1000000)
	castedCount, _ := strconv.ParseUint(fmt.Sprint(count), 10, 64)
	printer.Printf("\nCreating samples ...\n")
	var estimates = make([]int64, repeat)
	var t int64

	// initiator
	_, _ = bidsList.AddOrder(olist.Order{Price: 499, Volume: 5})

	start := time.Now().UnixNano()
	for t = 1; t <= repeat; t++ {
		startTime := time.Now().UnixNano()
		for i := 1; i <= count; i++ {
			randomNumber := int32(time.Now().UnixNano() % 1000) // this is necessary to use time for random prices, because generating new random numbers is a heavy operations for this test

			newOrder := olist.Order{
				Price:  (randomNumber) + 1,
				Volume: (randomNumber / 2) + 1,
			}
			_, err := bidsList.AddOrder(newOrder)
			if err != nil {
				panic(err)
			}
		}
		endTime := time.Now().UnixNano()
		if endTime-startTime > 0 {
			estimates = append(estimates, (endTime-startTime)/int64(count))
		}
	}

	tookedTime := time.Now().UnixNano() - start
	tookedSec := tookedTime / 1000000000
	secFloat := float64(tookedSec) + float64(tookedTime%1000000000)/1000000000
	//printer.Printf("time %d sec %d sec float %f", tookedTime, tookedSec, secFloat)

	printer.Printf("\tAdded %d Orders per sec\n", int(float64(totalSamples)/secFloat))

	printer.Println("\nTime:")
	printer.Printf("\tTooked about %f seconds totally\n", secFloat)

	var sum int64 = 0
	for _, value := range estimates {
		sum += value
	}
	printer.Println("\tExact average time:", tookedTime/totalSamples, "nanoseconds per each Order")

	printer.Println("\nOrder book memory usage:")
	printer.Printf("\tUsed %d byte for each order. \n\tTotal used memory: %d(%dmb) \n", (getUsedMemoryByByte()-baseMem)/castedCount, (getUsedMemoryByByte() - baseMem), bToMb((getUsedMemoryByByte() - baseMem)))

	printer.Println("\nData:")
	bidsList.PrintAll(false, false)
}

func getUsedMemoryByByte() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func randBetween(min int32, max int32) int32 {
	rand.Seed(time.Now().UnixNano())
	result := rand.Int31n(max-min) + min
	if result < min || result > max {
		panic("Wrong random number!")
	}
	return result

}
