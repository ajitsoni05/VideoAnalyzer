package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github/poornachandra7707/myboilerplate/models"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cacheClient *redis.Client
var mongoClient *mongo.Client
var ctx context.Context

func init() {
	cacheClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // replace with your redis server's address
		Password: "",               // replace with your redis server's password if any
		DB:       0,
	})
	_, err := cacheClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://poornachandra:apple123@cluster0.kbuqyo3.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer mongoClient.Disconnect(ctx)
}

func AnalyzeVideo(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into an AnalysisResult struct
	var analysisResult models.AnalysisResult
	err = json.Unmarshal(body, &analysisResult)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	//create the key using platformName_assetID
	key := analysisResult.PlatformName + "_" + analysisResult.AssetID
	analysisResult.Key = key

	//return if analysis for the asset is already done
	//Pre-check in db
	precheck, _ := GetStatusFromDb(key)
	if precheck == "completed" {
		w.Write([]byte(fmt.Sprintln("Status of video:Already done previously")))
		return
	}

	//check in cache as well
	status, err := cacheClient.Get(key).Result()
	if err == nil {
		w.Write([]byte(fmt.Sprintf("Status of video:%v\n", status)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Please note down the asset key for further reference:%v\n", analysisResult.Key)))

	go analyzeVideoUtil(analysisResult)
}

func analyzeVideoUtil(analysisResult models.AnalysisResult) {

	col := mongoClient.Database("boilerplate").Collection("videos")

	// Set the status to "in progress"
	analysisResult.Status = "in progress"
	err := cacheClient.Set(analysisResult.Key, analysisResult.Status, 0).Err()
	if err != nil {
		log.Println("Error setting the status in cache:", err)
	}

	// Print a message indicating that the analysis has started
	println("Analysis started for asset with key:", analysisResult.Key)

	// Perform the analysis
	output, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", analysisResult.FilePath).Output()
	if err != nil {
		// If there is an error, set the status to "failed" and print an error message
		analysisResult.Status = "failed"
		println("Analysis failed for asset with key:", analysisResult.Key)
		//http.Error(w, "Error running ffprobe", http.StatusInternalServerError)

		analysisResult.Status = "failed"
		err = cacheClient.Set(analysisResult.Key, analysisResult.Status, 0).Err()
		if err != nil {
			log.Println("Error setting the status in cache:", err)
		}
		println("Analysis failed for asset with key:", analysisResult.Key)
		return
	}

	// If the analysis was successful, set the status to "success" in cache
	analysisResult.Status = "completed"
	err = cacheClient.Set(analysisResult.Key, analysisResult.Status, 0).Err()
	if err != nil {
		log.Println("Error setting the status in cache:", err)
	}

	// Update the analysisResult with the results and set the status to "completed"
	analysisResult.Data = string(output)
	analysisResult.Timestamp = time.Now().Unix()
	analysisResult.Status = "completed"
	//analysisResults[analysisResult.AssetID] = &analysisResult

	// Print a message indicating that the analysis has completed
	println("Analysis completed for asset ID:", analysisResult.AssetID)
	// Print the complete metadata
	println(analysisResult.Data)

	// Write the analysis result to the response
	response, err := json.Marshal(analysisResult)
	if err != nil {
		fmt.Println("Error marshalling analysis result")
		//http.Error(w, "Error marshalling analysis result", http.StatusInternalServerError)
		return
	}
	_ = response
	//w.Write([]byte(analysisResult.Data))

	doc := []interface{}{analysisResult}
	result, err := col.InsertMany(context.TODO(), doc)
	_ = result

	if err != nil {
		panic(err)
	}

	contactIDs := result.InsertedIDs
	fmt.Printf("Asset with key %v inserted into the MongoDB with ObjectID:%v \n", analysisResult.Key, contactIDs)
	//w.Write([]byte(fmt.Sprintf("Asset with key %s inserted into the MongoDB with ObjectID : %s \n", analysisResult.Key, contactIDs)))

}
