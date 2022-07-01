package handlers

import (
	"2nd_test/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		UserName: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwt_key))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message:": "The registration has been completed.",
		"Token:": tokenString})
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
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := Claims{
		UserName: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwt_key))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}
	jwtout := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtout)
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenValue, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwt_key), nil
			})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusInternalServerError, gin.H{"Error1:": err.Error()})
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusInternalServerError, gin.H{"Error2:": err.Error()})
		}

		c.Next()
	}
}

func (handler *AuthHandler) RefreshAuthHandler(c *gin.Context) {
	token_value := c.GetHeader("Authorization")
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token_value, claims, func(token *jwt.Token) (interface{}, error) {
		return ([]byte(jwt_key)), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
	}

	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"Error:": "Invalid Token"})
	}
	// time.Unix(dateTime obj, 0) determines how many seconds passed from 0 to the time speicfied by dataTime object
	// https://www.geeksforgeeks.org/time-time-sub-function-in-golang-with-examples/
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"Error:": "Token is not expired yet!",
			"Seconds to Expiration Time:": time.Unix(claims.ExpiresAt, 0).Sub(time.Now())})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwt_key))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}
	jwtout := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtout)
}
