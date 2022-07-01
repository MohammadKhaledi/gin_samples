package main

import (
	"context"
	handlers "first_test/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//https://goswagger.io/generate/spec/operation.html
//https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-gin-gonic-version-269m

// For swagger.json indentation, DO NOT USE TAB!

// var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client

var recipeHandler *handlers.RecipesHandler

// type Recipe struct {
// 	ID           primitive.ObjectID `json:"id" bson:"_id"`
// 	Name         string             `json:"name" bson:"name"`
// 	Tags         []string           `json:"tags" bson:"tags"`
// 	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
// 	Instructions []string           `json:"instructions" bson:"instructions"`
// 	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
// }

func init() {

	ctx := context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://localhost:27017/admin"))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database("demo").Collection("recipes")
	recipeHandler = handlers.RecipeCollection(collection)
	log.Println("Connected to MongoDB")

}

//////////////////////////////////////////////////////////
// func NewRecipeHandler(ctx *gin.Context) {
// 	var recipe Recipe

// 	if err := ctx.ShouldBindJSON(&recipe); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
// 		return
// 	}
// 	recipe.ID = primitive.NewObjectID()
// 	recipe.PublishedAt = time.Now()
// 	collection := client.Database("demo").Collection("recipes")

// 	_, err = collection.InsertOne(ctx, recipe)
// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"Error: ": "Error while inserting new recipe!"})
// 		return
// 	}

// 	//	ctx.JSON(http.StatusOK, recipe)
// }

/////////////////////////////////////////////////
// func ListRcipesHandler(ctx *gin.Context) {
// 	collection := client.Database("demo").Collection("recipes")

// 	cur, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
// 		return
// 	}
// 	//https://www.mongodb.com/docs/v4.0/reference/method/cursor.next/
// 	//https://www.w3resource.com/mongodb/shell-methods/cursor/cursor-next.php
// 	defer cur.Close(ctx)
// 	recipes := make([]Recipe, 0)
// 	for cur.Next(ctx) {
// 		var recipe Recipe
// 		cur.Decode(&recipe)
// 		recipes = append(recipes, recipe)
// 	}
// 	// cur.Close(ctx)
// 	ctx.JSON(http.StatusOK, recipes)
// }
//////////////////////////////////////////////////////////////////

//https://www.mongodb.com/docs/drivers/go/v1.8/quick-reference/
// func UpdateRecipeHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	var recipe Recipe

// 	if err := ctx.ShouldBindJSON(&recipe); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
// 		return
// 	}

// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	collection := client.Database("demo").Collection("recipes")
// 	// Explanation for the warning: https://stackoverflow.com/questions/54548441/composite-literal-uses-unkeyed-fields
// 	// To solve warnings: bson.D{primitive.E{Key:"_id", Value:"objectId"},....}; It is not necessary
// 	set_elements := bson.D{{"_id", objectId}, {"name", recipe.Name},
// 		{"tags", recipe.Tags}, {"ingredients", recipe.Instructions},
// 		{"instructions", recipe.Ingredients}, {"publishedAt", time.Now()}}

// 	update_elements := bson.D{{"$set", set_elements}}
// 	// I do not know why, but the UpdateById method did not work for me; So I used UpdateOne method instead
// 	_, err := collection.UpdateOne(ctx, bson.D{{"_id", objectId}}, update_elements)

// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message: ": "Recipe has been updated"})

// }
////////////////////////////////////////////////////////////

// func DeleteRecipeHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	collection := client.Database("demo").Collection("recipes")

// 	_, err := collection.DeleteOne(ctx, bson.D{{"_id", objectId}})

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message: ": "Recipe has been deleted"})
// }

/////////////////////////////////////////////////////////////

// func SearchRecipeHandler(ctx *gin.Context) {
// 	tag := ctx.Query("tag")

// 	collection := client.Database("demo").Collection("recipes")

// 	cur, err := collection.Find(ctx, bson.D{{"tags", tag}})

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
// 		return
// 	}
// 	//https://www.mongodb.com/docs/v4.0/reference/method/cursor.next/
// 	//https://www.w3resource.com/mongodb/shell-methods/cursor/cursor-next.php
// 	defer cur.Close(ctx)
// 	var list_of_recipes []Recipe
// 	for cur.Next(ctx) {
// 		var recipe Recipe
// 		cur.Decode(&recipe)
// 		list_of_recipes = append(list_of_recipes, recipe)
// 	}

// 	ctx.JSON(http.StatusOK, list_of_recipes)
// }

// Json Unmarshal work: https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
// Json Marshal work: https://golang.cafe/blog/golang-json-marshal-example.html
func main() {

	router := gin.Default()

	router.POST("/recipes", recipeHandler.PostNewRecipeHandler)

	router.GET("/recipes", recipeHandler.ListRcipesHandler)

	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)

	router.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)

	router.GET("/recipes/search", recipeHandler.SearchRecipeHandler)

	router.Run()
}
