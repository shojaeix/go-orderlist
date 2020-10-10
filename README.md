# About
This package will help you to keep your bids/asks list.
 
This is dynamic-resizing mechanism and n(1) time complexity in most of operations.


# How to use
After importing the package, you should create a new `OrderList`.

Then you can use below functions to work on your list
- add/delete/update order
- get order's row and ahead volume
- get unique ID for each order
- get lowest/highest price and order
- get total number and volume of orders


## Examples

#### import the package
````go 
import "github.com/shojaeix/go-order-list/order-list
````

#### Create new list
````go
bidsList := olist.OrderList{}
bidsList.setSort("lowest") // lower bid must come first
````
#### Create new order
````go
newOrder := olist.Order{
   "price": 4302.23,
   "volume": 20,
}
````
#### Add order to list
````go
id, err := bidsList.AddOrder(newOrder) // the newOrder.ID will overwrite
````

#### Get order row and volume ahead
````go
row, volume, err := bidsList.GetRowAndAheadVolume(id)
````
#### Get edge order
````go
edgeOrder := bidsList.GetEdgeOrder()
````

#### Get lowest and highest price
````go
lowPrice, highPrice := bidsList.GetLowestAndHighestPrice()
````

#### Get total number and volume of orders
````go
totalOrders, totalVolume := bidsList.GetTotals()
````

#### Update order
````go
bidsList.UpdateOrder(orderId, newPrice, newVolume)
````

#### Delete order
````go
bidsList.DeleteOrder(orderId)
````
 
