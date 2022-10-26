package managerElem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func RegisterRestaurant(registerInfo RestaurantRegister) {
	var restaurInfo RestaurantInfo
	RestaurList.RestaurantsNum++

	restaurInfo.Id = registerInfo.Id
	restaurInfo.Name = registerInfo.Name
	restaurInfo.MenuItems = registerInfo.MenuItems
	restaurInfo.Menu = registerInfo.Menu
	restaurInfo.Rating = registerInfo.Rating

	RestaurantsAddress.Store(registerInfo.Id, registerInfo.Address)
	RestaurList.RestaurantsData = append(RestaurList.RestaurantsData, restaurInfo)

}

// here i have to send to dinning_hall each order, after i send  to dinning_hall, i receive the orderId,  which is sent together with estimatedTime, as a response to Client
func ConstructResponseOrder(clientOrder ClientOrder) (clientOrderResponse ClientOrderResponse, errors error) {
	//var clientOrderResponse ClientOrderResponse
	//log.Fatal(len(clientOrderResponse.Orders))
	//creating already null struct values
	var ordersListResponse = make([]OrderResponse, len(clientOrder.Orders))

	atomic.AddInt32(&orderIdInc, 1)
	clientOrderResponse.OrderId = orderIdInc

	for i, order := range clientOrder.Orders {

		//first of all perfom order post request to dinning hall in order to receive important information
		orderId, estimatedTime, registeredTime, err := SendOrderToRestaurant(order)
		if err != nil {
			errors = err
			ordersListResponse[i] = OrderResponse{}
			continue
		}

		ordersListResponse[i].RestaurantId = order.RestaurantId

		address, _ := RestaurantsAddress.Load(order.RestaurantId)
		ordersListResponse[i].RestaurantAddress = address.(string)

		//the order is get from the received info from dinning hall
		ordersListResponse[i].OrderId = orderId
		ordersListResponse[i].EstimatedWaitingTime = estimatedTime
		ordersListResponse[i].CreatedTime = order.CreatedTime
		//NO, IT IS FROM THE RESPONSE OF THE DINNING HALL
		ordersListResponse[i].RegisteredTime = registeredTime
	}

	clientOrderResponse.Orders = ordersListResponse

	return

}

func SendOrderToRestaurant(order ReceivedFromClientOrder) (orderId int, estimatedTime int, registeredTime time.Time, err error) {
	var orderForRestaurant OrderForRestaurant

	restaurantAddress, ok := RestaurantsAddress.Load(order.RestaurantId)
	//here error
	if ok == false {
		err = fmt.Errorf("Could not find address of restaurant")
		return
		//log.Errorf("Could not find address of restaurant ", order.RestaurantId)
	}

	orderForRestaurant.Priority = order.Priority
	orderForRestaurant.MaxWait = order.MaxWait
	orderForRestaurant.CreatedTime = order.CreatedTime
	orderForRestaurant.Items = order.Items

	reqBody, err := json.Marshal(orderForRestaurant)
	//here error
	if err != nil {
		err = fmt.Errorf("%v", err.Error())
		return
	}
	resp, err := http.Post(restaurantAddress.(string)+"v2/order", "application/json", bytes.NewBuffer(reqBody))

	// TODO IF ERROR UNANOUNCE CLIENT AND CANCEL ORDER
	if err != nil {
		err = fmt.Errorf("Sending Online Order Request to Restaurant Failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("Dinning Hall Response reading Failed: %s", err.Error())
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
