package main

import (
	handlers "2nd_prj/handlers"
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)

	collection1 := client.Database("demo").Collection("recipes")
	collection2 := client.Database("admin").Collection("users")
	recipeHandler = handlers.RecipeCollection(collection1, redisClient)
	authHandler = handlers.NewAuthHandler(collection2)

	log.Println("Connected to MongoDB")

}

// https://github.com/gorilla/sessions
// https://github.com/boj/redistore
// https://golangexample.com/gorilla-sessions-provides-cookie-and-filesystem-sessions-and-infrastructure-for-custom-session-backends/

func main() {

	router := gin.Default()

	store := cookie.NewStore([]byte("secret_key"))
	router.Use((sessions.Sessions("recipes_api", store)))

	authorized := router.Group("/users")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipeHandler.PostNewRecipeHandler)

		authorized.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)

		authorized.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)

		authorized.GET("/recipes/search", recipeHandler.SearchRecipeHandler)

		// authorized.GET("/recipes", recipeHandler.ListRcipesHandler)

		authorized.POST("/signout", authHandler.SignOutHandler)
	}
	router.GET("/recipes", recipeHandler.ListRcipesHandler)

	router.POST("/signin", authHandler.SignInHandler)

	router.POST("/signup", authHandler.SignUpHandler)

	router.RunTLS(":443", "./certs/localhost.crt", "./certs/localhost.key")
}
