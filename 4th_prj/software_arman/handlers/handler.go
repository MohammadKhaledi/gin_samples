package handlers

import (
	"2nd_prj/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

func RecipeCollection(collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		redisClient: redisClient,
	}
}

func (handler *RecipesHandler) PostNewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := handler.collection.InsertOne(c, recipe)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	log.Println("Erase Cache")
	handler.redisClient.Del("recipes")
	c.JSON(http.StatusOK, gin.H{"message: ": "New recipe posted"})
}

func (handler *RecipesHandler) ListRcipesHandler(c *gin.Context) {

	val, er := handler.redisClient.Get("recipes").Result()
	if er == redis.Nil {
		log.Printf("Request to MongoDb")
		curr, err := handler.collection.Find(c, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
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
		//https://pkg.go.dev/github.com/go-redis/redis/v8#Client.Set
		//Zero expiration means the key has no expiration time.
		handler.redisClient.Set("recipes", string(data), time.Hour)
		c.JSON(http.StatusOK, list_of_recipes)
	} else if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": er.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		var list_of_recipes []models.Recipe
		json.Unmarshal([]byte(val), &list_of_recipes)
		c.JSON(http.StatusOK, list_of_recipes)
	}

}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	set_elements := bson.D{{Key: "_id", Value: objectId}, {Key: "name", Value: recipe.Name},
		{Key: "tags", Value: recipe.Tags}, {Key: "ingredients", Value: recipe.Instructions},
		{Key: "instructions", Value: recipe.Ingredients}, {Key: "publishedAt", Value: time.Now()}}

	update_elements := bson.D{{Key: "$set", Value: set_elements}}
	_, err := handler.collection.UpdateOne(c, bson.D{{Key: "_id", Value: objectId}}, update_elements)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "Recipe has been updated"})
}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.DeleteOne(c, bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "Recipe has been deleted"})
}

func (handler *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")

	curr, err := handler.collection.Find(c, bson.D{{Key: "tags", Value: tag}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	defer curr.Close(c)
	var list_of_recipes []models.Recipe
	for curr.Next(c) {
		var recipe models.Recipe
		curr.Decode(&recipe)
		list_of_recipes = append(list_of_recipes, recipe)
	}
	c.JSON(http.StatusOK, list_of_recipes)
}
