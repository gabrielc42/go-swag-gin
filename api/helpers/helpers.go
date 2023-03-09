package helper

import (
	"time"

	modelUser "github.com/gabrielc42/api/models/user"

	modelTodo "github.com/gabrielc42/api/models/todo"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func createToken(_id string) (string, int64, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = _id
	atClaims["exp"] = time.Now().Add(time.Minute*15).Unix() * 1000
	expires := atClaims["exp"].(int64)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", expires, err
	}
	return token, expires, nil
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error":  message,
		"status": false,
	})
}

func respondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"data":    data,
		"message": message,
		"status":  true,
	})
}

func validateSignupInput(user modelUser.User) string {
	var errMsg string

	if user.FirstName == "" {
		errMsg = "Please enter a valid first name."
		return errMsg
	}

	if user.LastName == "" {
		errMsg = "Please enter a valid last name."
		return errMsg

	}

	if user.Email == "" {
		errMsg = "Please enter a valid email id."
		return errMsg
	}

	if user.Password == "" {
		errMsg = "Please enter a valid password"
		return errMsg
	}

	if user.Mobile == "" {
		errMsg = "Please enter a valid mobile"
		return errMsg
	}
	return ""
}

func validateLoginInput(user modelUser.UserLogin) string {
	var errMsg string

	if user.Email == "" {
		errMsg = "Please provide email"
		return errMsg
	}

	if user.Password == "" {
		errMsg = "Please provide password"
		return errMsg
	}

	return ""
}

func validateTodoInput(todo modelTodo.Todo) string {
	var errMsg string

	if todo.Name == "" {
		errMsg = "please provide valid name!"
		return errMsg
	}

	return ""
}
