package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	var err error
	err = godotenv.Load()
	if err != nil {
    	log.Fatal("Error loading .env file")
  	}	
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING"))
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB Atlas!")

	err = godotenv.Load()
	if err != nil {
    	log.Fatal("Error loading .env file")
  	}	

    router := mux.NewRouter()
    router.HandleFunc("/updateField", updateFieldHandler).Methods("POST")
    router.HandleFunc("/createDocument", createDocumentHandler).Methods("POST")
    router.HandleFunc("/getDocument", getDocumentHandler).Methods("GET")
	router.HandleFunc("/dbInit", dbInit).Methods("POST")

    http.ListenAndServe(":8080", router)
}

func updateFieldHandler(w http.ResponseWriter, r *http.Request) {
    type Data struct {
        UserID   int    `json:"userId"`
        UserName string `json:"userName"`
		Count   int    `json:"count"`
    }

	type UpdateRequest struct {
        Data Data `json:"data"`
    }

    var req UpdateRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    collection := client.Database("pender-clicks").Collection("clicks-01")
    filter := bson.M{"userId": req.Data.UserID, "userName": req.Data.UserName}

    var currentDoc struct {
        Count int `bson:"count"`
    }
    err = collection.FindOne(context.TODO(), filter).Decode(&currentDoc)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No document found with the provided userId and userName", http.StatusNotFound)
            return
        }

        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if req.Data.Count <= currentDoc.Count {
        fmt.Fprint(w, currentDoc.Count)
        return
    }

    update := bson.M{"$set": bson.M{
        "userId":   req.Data.UserID,
        "userName": req.Data.UserName,
        "count":    req.Data.Count,
    }}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, req.Data.Count)
}

func createDocument(userId int, userName string, count int) error {
    collection := client.Database("pender-clicks").Collection("clicks-01")
    document := bson.M{
        "userId":   userId,
        "userName": userName,
        "count":    count,
    }

    _, err := collection.InsertOne(context.TODO(), document)
    return err
}

func createDocumentHandler(w http.ResponseWriter, r *http.Request) {
    var req InitRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = createDocument(req.Data.UserID, req.Data.UserName, req.Data.Count)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func getDocumentHandler(w http.ResponseWriter, r *http.Request) {
    collection := client.Database("pender-clicks").Collection("clicks-01")

    opts := options.Find()
    opts.SetSort(bson.D{{"count", -1}})
    opts.SetLimit(5)

    cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var results []bson.M
    if err = cursor.All(context.TODO(), &results); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    formattedResults := make([]string, len(results))
    for i, result := range results {
        formattedResults[i] = fmt.Sprintf("%s : %v", result["userName"], result["count"])
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(formattedResults)
}

type InitRequest struct {
    Data struct {
        UserID   int    `json:"userId"`
        UserName string `json:"userName"`
        Count    int    `json:"count"`
    } `json:"data"`
}

func dbInit(w http.ResponseWriter, r *http.Request) {
    var req InitRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    collection := client.Database("pender-clicks").Collection("clicks-01")
    filter := bson.M{
        "userId":   req.Data.UserID,
        "userName": req.Data.UserName,
    }

    var result struct {
        Count int `bson:"count"`
    }
    err = collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            err = createDocument(req.Data.UserID, req.Data.UserName, req.Data.Count)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusCreated)
            return
        }

        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, result.Count)
}