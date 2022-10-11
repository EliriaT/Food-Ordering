package managerElem

import "sync"

var (
	RestaurList        RestaurantsList
	orderIdInc         int32
	RestaurantsAddress sync.Map
	Port               = ":8084"
)
