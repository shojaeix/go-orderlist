# About
This package will help you to keep your bids/asks list with minimum memory usage and fast speed.
 
This package supports dynamic-resizing mechanism and have n(1) time complexity in most operations.


# How to use
After importing the package, you should create a new variable with type `OrderList`.

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

#### Create a new list
````go
bidsList := olist.OrderList{}
````
#### Set list's sort
````go
bidsList.setSort("lowest") // lower bid must come first
````
#### Create a new order
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
 
## TODO
- Return callback chan in result of the AddOrder() func
- Test
- Implement the GetRowAndAheadVolume() func
- Implement the setSort() functionality