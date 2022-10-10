package managerElem

import "time"

type OrderForRestaurant struct {
	Items       []int     `json:"items"`
	Priority    int       `json:"priority"`
	MaxWait     float32   `json:"max_wait"`
	CreatedTime time.Time `json:"created_time"`
}

type ResponseFromRestaurant struct {
	RestaurantId         int       `json:"restaurant_id"`
	OrderId              int       `json:"order_id"`
	EstimatedWaitingTime int       `json:"estimated_waiting_time"`
	CreatedTime          time.Time `json:"created_time"`
	RegisteredTime       time.Time `json:"registered_time"`
}
