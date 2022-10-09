package managerElem

import "time"

type ClientOrder struct {
	ClientId int
	Orders   []Order
}

type ClientOrderResponse struct {
	OrderId int32           `json:"order_id"`
	Orders  []OrderResponse `json:"orders"`
}

type Order struct {
	RestaurantId int       `json:"restaurant_id"`
	Items        []int     `json:"items"`
	Priority     int       `json:"priority"`
	MaxWait      float64   `json:"max_wait"`
	CreatedTime  time.Time `json:"created_time"`
}

type OrderResponse struct {
	RestaurantId         int       `json:"restaurant_id"`
	RestaurantAddress    string    `json:"restaurant_address"`
	OrderId              int       `json:"order_id"`
	EstimatedWaitingTime int       `json:"estimated_waiting_time"`
	CreatedTime          time.Time `json:"created_time"`
	RegisteredTime       time.Time `json:"registered_time"`
}
