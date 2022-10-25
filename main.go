package main

import (
	"encoding/json"
	managerElem "github.com/EliriaT/Food-Ordering/orders-manager-elem"
	"log"
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

	managerElem.RegisterRestaurant(registeredRestaurant)

	log.Println("Succesfully registered restaurand with id ", registeredRestaurant.Id)
	log.Printf("Current menu : %v", managerElem.RestaurList)
	defer r.Body.Close()
	//jsonCookedOrder, _ := json.Marshal(registeredRestaurant)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(jsonCookedOrder)
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

	log.Printf("Order from client %d was received : %v \n", clientOrder.ClientId, clientOrder.Orders)

	clientOrderResponse := managerElem.ConstructResponseOrder(clientOrder)

	jsonCookedOrder, _ := json.Marshal(clientOrderResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
}

func receiveRaitings(w http.ResponseWriter, r *http.Request) {
	var clientRatings managerElem.RaitingsResponses
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&clientRatings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	log.Printf("Ratings from client %d was received : %+v \n", clientRatings.ClientId, clientRatings.Orders)

	//sending the ratings to the restaurants
	managerElem.SendRatingsToRestaurants(clientRatings)

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}

func main() {
	//for init
	managerElem.RestaurantsAddress = sync.Map{}

	r := mux.NewRouter()
	r.HandleFunc("/menu", getMenus).Methods("GET")
	r.HandleFunc("/register", registerRestaurant).Methods("POST")
	r.HandleFunc("/order", receiveOrder).Methods("POST")
	r.HandleFunc("/rating", receiveRaitings).Methods("POST")

	log.Println("Food-Ordering server started..")
	log.Fatal(http.ListenAndServe(managerElem.Port, r))
}
