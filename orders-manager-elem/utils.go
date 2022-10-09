package managerElem

import (
	"sync/atomic"
	"time"
)

// here i have to send to dinning_hall each order
func ConstructResponseOrder(clientOrder ClientOrder) ClientOrderResponse {
	var clientOrderResponse ClientOrderResponse
	//creating already null struct values
	var ordersListRespone = make([]OrderResponse, len(clientOrderResponse.Orders))

	atomic.AddInt32(&orderIdInc, 1)
	clientOrderResponse.OrderId = orderIdInc

	for i, order := range clientOrder.Orders {

		//first of all perfom order post

		ordersListRespone[i].RestaurantId = order.RestaurantId

		address, _ := RestaurantsAddress.Load(order.RestaurantId)
		ordersListRespone[i].RestaurantAddress = address.(string)

		ordersListRespone[i].OrderId = i + 1
		ordersListRespone[i].EstimatedWaitingTime = 0
		ordersListRespone[i].CreatedTime = order.CreatedTime
		ordersListRespone[i].RegisteredTime = time.Now()
	}

	clientOrderResponse.Orders = ordersListRespone

	return clientOrderResponse

}
