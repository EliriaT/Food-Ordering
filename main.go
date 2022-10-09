package Food_Ordering

import (
	"encoding/json"
	managerElem "github.com/EliriaT/Food-Ordering/orders-manager-elem"
	"sync"

	"github.com/gorilla/mux"
	"net/http"
)

func getMenus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonMenus, err := json.Marshal(managerElem.RestaurList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//by default sends 200
	w.Write(jsonMenus)
}

func registerRestaurant(w http.ResponseWriter, r *http.Request) {
	var registeredRestaurant managerElem.RestaurantRegister
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&registeredRestaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managerElem.RestaurList.RestaurantsNum++
	managerElem.RestaurList.RestaurantsData = append(managerElem.RestaurList.RestaurantsData, registeredRestaurant)
	managerElem.RestaurantsAddress.Store(registeredRestaurant.Id, registeredRestaurant.Address)

	defer r.Body.Close()
	jsonCookedOrder, _ := json.Marshal(registeredRestaurant)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
}

func receiveOrder(w http.ResponseWriter, r *http.Request) {
	var clientOrder managerElem.ClientOrder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&clientOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	clientOrderResponse := managerElem.ConstructResponseOrder(clientOrder)

	jsonCookedOrder, _ := json.Marshal(clientOrderResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
}

func main() {
	//for init
	managerElem.RestaurantsAddress = sync.Map{}

	r := mux.NewRouter()
	r.HandleFunc("/menu", getMenus).Methods("GET")
	r.HandleFunc("/register", registerRestaurant).Methods("POST")
	r.HandleFunc("/order", receiveOrder).Methods("POST")
}
