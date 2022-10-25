package managerElem

type RaitingsResponses struct {
	ClientId int           `json:"client_id"`
	OrderId  int           `json:"order_id"`
	Orders   []orderRating `json:"orders"`
}

type orderRating struct {
	RestaurantId  int `json:"restaurant_id,omitempty"`
	OrderId       int `json:"order_id"`
	Rating        int `json:"rating"`
	EstimatedTime int `json:"estimated_waiting_time"`
	WaitedTime    int `json:"waiting_time"`
}

type restaurantRating struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float32 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
