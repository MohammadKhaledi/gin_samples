package main

import (
	handlers "2nd_test/handlers"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client

var recipeHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {

	ctx := context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	rediClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := rediClient.Ping()
	fmt.Println(status)
	collection1 := client.Database("demo").Collection("recipes")
	collection2 := client.Database("admin").Collection("users")
	recipeHandler = handlers.RecipeCollection(collection1, rediClient)
	authHandler = handlers.NewAuthHandler(collection2)

	log.Println("Connected to MongoDB")

}

func main() {

	router := gin.Default()

	router.POST("/recipes", recipeHandler.PostNewRecipeHandler)

	router.GET("/recipes", authHandler.AuthMiddleware(), recipeHandler.ListRcipesHandler)

	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)

	router.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)

	router.GET("/recipes/search", recipeHandler.SearchRecipeHandler)

	router.POST("/signin", authHandler.AuthMiddleware(), authHandler.SignInHandler)

	router.POST("/refresh", authHandler.RefreshAuthHandler)

	router.POST("/signup", authHandler.SignUpHandler)

	router.Run()
}
