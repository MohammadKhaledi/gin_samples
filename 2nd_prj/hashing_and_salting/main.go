package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Users struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	var users = make([]Users, 0)
	file, _ := ioutil.ReadFile("user.json")
	_ = json.Unmarshal([]byte(file), &users)
	//func Background() Context. Background returns a non-nil, empty Context.
	//It is never canceled, has no values, and has no deadline.
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongo")
	collection := client.Database("admin").Collection("users")

	var list_of_users []interface{}
	for _, user := range users {
		list_of_users = append(list_of_users, user)
	}
	insert_many, err := collection.InsertMany(ctx, list_of_users)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted Users: ", len(insert_many.InsertedIDs))
}
