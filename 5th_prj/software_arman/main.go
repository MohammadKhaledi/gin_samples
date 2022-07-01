package main

import (
	"context"
	"fmt"
	"log"
	handlers "software_arman/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var err error
var client *mongo.Client

var recipeHandler *handlers.RecipeHandler

func init() {
	ctx := context.Background()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	fmt.Println(status)

	collection1 := client.Database("demo").Collection("recipes")
	recipeHandler = handlers.RecipeCollection(collection1, redisClient)

	log.Println("Connected to MongoDB")

}
func IndexHandler(c *gin.Context) {
	c.File("./html_files/index.html")
}

func main() {

	router := gin.Default()

	router.GET("/", IndexHandler)

	router.GET("/recipes", recipeHandler.ListRecipesHandler)

	router.Run()

}
