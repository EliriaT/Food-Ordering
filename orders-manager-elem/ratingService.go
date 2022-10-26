package managerElem

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SendRatingsToRestaurants(ratings RaitingsResponses) {
	for _, rating := range ratings.Orders {
		localAdress, _ := RestaurantsAddress.Load(rating.RestaurantId)

		rating.RestaurantId = 0

		reqBody, err := json.Marshal(rating)
		if err != nil {
			log.Fatal(err.Error())

		}
		resp, err := http.Post(localAdress.(string)+"v2/rating", "application/json", bytes.NewBuffer(reqBody))

		if err != nil {
			log.Fatal("Sending Rating Request to Restaurant Failed: %s", err.Error())
		}
		defer resp.Body.Close()

		//as a response i get restaurant average
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Dinning Hall Rating Response reading Failed: %s", err.Error(), body)
			return
		}

		var individualRating restaurantRating
		_ = json.Unmarshal(body, &individualRating)
		CalculateAvgSimRating(individualRating)

	}
}

func CalculateAvgSimRating(individualRating restaurantRating) {
	var avg float32 = 0
	for i, restaurant := range RestaurList.RestaurantsData {
		if restaurant.Id == individualRating.RestaurantId {
			RestaurList.RestaurantsData[i].Rating = individualRating.RestaurantAvgRating
		}
		avg = avg + RestaurList.RestaurantsData[i].Rating
	}
	avg = avg / float32(RestaurList.RestaurantsNum)
	log.Println("----------------------SIMULATION AVERAGE IS: ", avg, "----------------------")
}