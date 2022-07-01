package handlers

import (
	"2nd_prj/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const jwt_key = "123456"

type AuthHandler struct {
	collection *mongo.Collection
}

type Claims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
	}
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error:": err.Error()})
		return
	}
	if len(user.Username) == 0 || len(user.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error:": "The Username or Password has not been specified."})
		return
	}
	_, err := handler.collection.InsertOne(c, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message:": "The registration has been completed."})
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	//https://pkg.go.dev/crypto/sha256#example-New

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error:": err.Error()})
		return
	}
	cur := handler.collection.FindOne(c, bson.M{"username": user.Username, "password": user.Password})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Username or Password"})
		return
	}
	sessionToken := xid.New().String()
	session := sessions.Default(c)

	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"Message:": "User signed in"})
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"Message:": "The User signed out..."})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")

		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{"Message:": "No logged"})
			c.Abort()
		}
		c.Next()
	}
}
