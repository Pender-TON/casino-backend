package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type VerifyRequest struct {
    Hash string `json:"hash"`
    Data string `json:"data"`
}

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
    router.HandleFunc("/verifySignature", verifySignatureHandler).Methods("POST")
	router.HandleFunc("/dbInit", dbInit).Methods("POST")
    router.HandleFunc("/getPosition", getPositionHandler).Methods("POST")

    http.ListenAndServe(":8080", router)
}

func updateFieldHandler(w http.ResponseWriter, r *http.Request) {
	type UpdateRequest struct {
        UserID   int    `json:"userId"`
        UserName string `json:"userName"`
		Count   int    `json:"count"`
    }

    var req UpdateRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    collection := client.Database("pender-clicks").Collection("clicks-01")
    filter := bson.M{"userId": req.UserID, "userName": req.UserName}

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

    if req.Count <= currentDoc.Count {
        fmt.Fprint(w, currentDoc.Count)
        return
    }

    update := bson.M{"$set": bson.M{
        "userId":   req.UserID,
        "userName": req.UserName,
        "count":    req.Count,
    }}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, req.Count)
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

    err = createDocument(req.UserID, req.UserName, req.Count)
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
    UserID   int    `json:"userId"`
    UserName string `json:"userName"`
    Count    int    `json:"count"`
}

func getPositionHandler(w http.ResponseWriter, r *http.Request) {
    // Define a struct to hold the incoming JSON data
    type RequestData struct {
        UserId   int    `json:"userId"`
    }

    // Decode the JSON request body into the struct
    var requestData RequestData
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Use the userId from the request data
    userId := requestData.UserId

    // Query the MongoDB collection and sort the documents in descending order by the 'count' field
    collection := client.Database("pender-clicks").Collection("clicks-01")
    cursor, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.D{{"count", -1}}))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var results []bson.M
    if err = cursor.All(context.TODO(), &results); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Iterate over the sorted documents to find the position of the specified userId
    position := -1
    for i, result := range results {
        dbUserId, ok := result["userId"].(int64)
        if !ok {
            dbUserId32, ok := result["userId"].(int32)
            if ok {
                dbUserId = int64(dbUserId32)
            } else {
                fmt.Printf("Unexpected type for userId: %T\n", result["userId"])
                continue
            }
        }
    
        if dbUserId == int64(userId) {
            position = i + 1
            break
        }
    }

    // Return the position of the specified userId
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"position": position})
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
        "userId":   req.UserID,
        "userName": req.UserName,
    }

    var result struct {
        Count int `bson:"count"`
    }
    err = collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            err = createDocument(req.UserID, req.UserName, req.Count)
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

func verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
    // Read the body of the request
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusInternalServerError)
        return
    }

    // Parse the body as a URL-encoded string
    values, err := url.ParseQuery(string(body))
    if err != nil {
        http.Error(w, "Error parsing request body", http.StatusInternalServerError)
        return
    }

    // Get the hash from the parsed values
    hash := values.Get("hash")

    // Generate the data-check-string
    var dataCheckString strings.Builder
    keys := make([]string, 0, len(values))
    for key := range values {
        if key != "hash" {
            keys = append(keys, key)
        }
    }
    sort.Strings(keys)
    for _, key := range keys {
        // Ensure that the 'user' field is treated as a single string
        value := values.Get(key)
        if key == "user" {
            value, _ = url.QueryUnescape(value)
        }
        dataCheckString.WriteString(fmt.Sprintf("%s=%s\n", key, value))
    }

    // Print the data-check-string for debugging
    fmt.Println("Data-check-string:", dataCheckString.String())

    // Replace with your bot's token
    botToken := os.Getenv("BOT_TOKEN")

    // Generate the secret key
    secretKeyHMAC := hmac.New(sha256.New, []byte(botToken))
    secretKeyHMAC.Write([]byte("WebAppData"))
    secretKey := secretKeyHMAC.Sum(nil)

    // Generate the HMAC-SHA-256 signature of the data-check-string
    hmacData := hmac.New(sha256.New, secretKey)
    hmacData.Write([]byte(dataCheckString.String()))
    dataSignature := hmacData.Sum(nil)

    // Convert the received hash and computed signature to lowercase hex strings for comparison
    receivedHash, err := hex.DecodeString(hash)
    if err != nil {
        http.Error(w, "Error decoding hash", http.StatusInternalServerError)
        return
    }

    // Compare the received hash with the data signature
    if hmac.Equal(receivedHash, dataSignature) {
        w.Write([]byte("true"))
    } else {
        w.Write([]byte("false"))
    }
}