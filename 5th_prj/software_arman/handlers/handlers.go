package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"software_arman/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeHandler struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

func RecipeCollection(collection *mongo.Collection, redisclient *redis.Client) *RecipeHandler {
	return &RecipeHandler{
		collection:  collection,
		redisClient: redisclient,
	}
}

func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	val, er := handler.redisClient.Get("recipes").Result()

	if er == redis.Nil {
		log.Println("Request to MongoDB")
		curr, err := handler.collection.Find(c, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
			return
		}
		defer curr.Close(c)
		var list_of_recipes []models.Recipe

		for curr.Next(c) {
			var recipe models.Recipe
			curr.Decode(&recipe)
			list_of_recipes = append(list_of_recipes, recipe)
		}
		data, _ := json.Marshal(list_of_recipes)
		handler.redisClient.Set("recipes", string(data), time.Hour)
		c.JSON(http.StatusOK, list_of_recipes)

	} else if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": er.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		var list_of_recipes []models.Recipe
		json.Unmarshal([]byte(val), &list_of_recipes)
		c.JSON(http.StatusOK, list_of_recipes)
	}
}
