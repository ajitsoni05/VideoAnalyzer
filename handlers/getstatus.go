package handlers

import (
	"fmt"
	"net/http"

	"log"

	"github/poornachandra7707/myboilerplate/models"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStatusFromDb(key string) (string, error) {

	col := mongoClient.Database("boilerplate").Collection("videos")
	filter := bson.D{primitive.E{Key: "key", Value: key}}
	var result models.AnalysisResult
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("analysis not found")
		}
		return "", fmt.Errorf("error fetching the status: %v", err)
	}
	return result.Status, nil
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	fmt.Println(key)

	/*
		//Pre-check in db
		precheck, _ := GetStatusFromDb(key)
		if precheck == "completed" {
			w.Write([]byte(fmt.Sprintln("Status of video: Already done previously")))

			return
		}
	*/
	// check the status in cache
	status, err := cacheClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			// status not found in cache, look for it in MongoDB
			status, err = GetStatusFromDb(key)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		} else {
			log.Println("Error fetching the status from cache:", err)
			http.Error(w, "Error fetching the status", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte(fmt.Sprintf("Status of video:%s", status)))
}
