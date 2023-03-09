package todocontrollers

import (
	"net/http"
	"reflect"
	"time"

	configuration "github.com/gabrielc42/api/configuration"
	constants "github.com/gabrielc42/api/constant"
	helper "github.com/gabrielc42/api/helpers"
	models "github.com/gabrielc42/api/models/todo"
	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2/bson"
)

func createTodo(c *gin.Context) {
	var Todo models.Todo
	val := reflect.ValueOf(c.Keys["user_id"])

	Todo.ID = bson.NewObjectId()
	Todo.User = bson.ObjectIdHex(val.String())
	Todo.Date = time.Now()
	Todo.Status = false

	jsonErr := c.BindJSON(&Todo)

	if jsonErr != nil {
		helper.respondWithError(c, http.StatusBadRequest, jsonErr)
		return
	}

	if Todo.Name == "" {
		helper.respondWithError(c, http.StatusBadRequest, "Please provide valid name.")
		return
	}

	mongoSession := configuration.ConnectDb(constants.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.Db(constants.Database).c("todo")
	err := getCollection.Insert(Todo)

	if err != nil {
		helper.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.respondWithSuccess(c, http.StatusOK, constants.TodoCreatedSuccess, Todo)
}
