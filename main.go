package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"log"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	// "go.mongodb.org/mongo-go-driver/mongo"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lasttname,omitempty" bson:"firstname,omitempty"`
}

var client *mongo.Client

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("Cluster0").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	// json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("Cluster0").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// result, _ := collection.InsertOne(ctx, person)
	cursor, err := collection.Find(ctx, bson.M{})
	// json.NewEncoder(response).Encode(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message":"` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(people)
}

func main() {
	fmt.Println("starting the app...")
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")
	// client, _ = mongo.Connect(ctx, "mongodb+srv://TEST00:Abhishek@007@test00.24gdy.mongodb.net/TEST00?retryWrites=true&w=majority")
	// router := mux.NewRouter()
	// http.ListenAndServe(":12345",router)

	// clientOptions := options.Client().
	// 	ApplyURI("mongodb+srv://TEST00:Abhishek@007@test00.24gdy.mongodb.net/TEST00?retryWrites=true&w=majority")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// client, _ = mongo.Connect(ctx, clientOptions)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// clientOptions := options.Client().
	// 	ApplyURI("mongodb+srv://abhishek:abhishek@cluster0.24gdy.mongodb.net/Cluster0?retryWrites=true&w=majority")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, clientOptions)
	// _ = client
	// if err != nil {
	// 	log.Fatal(err)
	// }

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://abhishek:abhishek@cluster0.24gdy.mongodb.net/Cluster0?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// defer client.Disconnect(ctx)
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)

}
