package managerElem

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func RegisterRestaurant(registerInfo RestaurantRegister) {
	var restaurInfo RestaurantInfo
	RestaurList.RestaurantsNum++

	restaurInfo.Name = registerInfo.Name
	restaurInfo.MenuItems = registerInfo.MenuItems
	restaurInfo.Menu = registerInfo.Menu
	restaurInfo.Rating = registerInfo.Rating

	RestaurantsAddress.Store(registerInfo.Id, registerInfo.Address)
	RestaurList.RestaurantsData = append(RestaurList.RestaurantsData, restaurInfo)

}

// here i have to send to dinning_hall each order, after i send  to dinning_hall, i receive the orderId,  which is sent together with estimatedTime, as a response to Client
func ConstructResponseOrder(clientOrder ClientOrder) ClientOrderResponse {
	var clientOrderResponse ClientOrderResponse
	//log.Fatal(len(clientOrderResponse.Orders))
	//creating already null struct values
	var ordersListResponse = make([]OrderResponse, len(clientOrder.Orders))

	atomic.AddInt32(&orderIdInc, 1)
	clientOrderResponse.OrderId = orderIdInc

	for i, order := range clientOrder.Orders {

		//first of all perfom order post request to dinning hall in order to receive important information
		orderId, estimatedTime, registeredTime := SendOrderToRestaurant(order)

		ordersListResponse[i].RestaurantId = order.RestaurantId

		address, _ := RestaurantsAddress.Load(order.RestaurantId)
		ordersListResponse[i].RestaurantAddress = address.(string)

		//the order is get from the received info from dinning hall
		ordersListResponse[i].OrderId = orderId
		ordersListResponse[i].EstimatedWaitingTime = estimatedTime
		ordersListResponse[i].CreatedTime = order.CreatedTime
		//NO, IT IS FROM THE RESPONSE OF THE DINNING HALL ?
		ordersListResponse[i].RegisteredTime = registeredTime
	}

	clientOrderResponse.Orders = ordersListResponse

	return clientOrderResponse

}

func SendOrderToRestaurant(order ReceivedFromClientOrder) (orderId int, estimatedTime int, registeredTime time.Time) {
	var orderForRestaurant OrderForRestaurant

	restaurantAddress, ok := RestaurantsAddress.Load(order.RestaurantId)
	if ok == false {
		log.Fatal("Could not find address of restaurant ", order.RestaurantId)
	}

	orderForRestaurant.Priority = order.Priority
	orderForRestaurant.MaxWait = order.MaxWait
	orderForRestaurant.CreatedTime = order.CreatedTime
	orderForRestaurant.Items = order.Items

	reqBody, err := json.Marshal(orderForRestaurant)
	if err != nil {
		log.Fatal(err.Error())

	}
	resp, err := http.Post(restaurantAddress.(string)+"v2/order", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Sending Online Order Request to Restaurant Failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Dinning Hall Response reading Failed: %s", err.Error())
		return
	}

	var responseAboutOrder ResponseFromRestaurant
	_ = json.Unmarshal(body, &responseAboutOrder)

	orderId = responseAboutOrder.OrderId
	estimatedTime = responseAboutOrder.EstimatedWaitingTime
	registeredTime = responseAboutOrder.RegisteredTime

	log.Printf("Order %d was sent to restaurant %d \n", orderId, order.RestaurantId)

	return
}
