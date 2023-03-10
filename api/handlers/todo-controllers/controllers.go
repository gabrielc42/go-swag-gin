package todocontrollers

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	configuration "github.com/gabrielc42/api/configuration"
	constants "github.com/gabrielc42/api/constant"
	helper "github.com/gabrielc42/api/helpers"
	models "github.com/gabrielc42/api/models/todo"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateTodo(c *gin.Context) {
	var Todo models.Todo
	val := reflect.ValueOf(c.Keys["user_id"])

	Todo.ID = bson.NewObjectId()
	Todo.User = bson.ObjectIdHex(val.String())
	Todo.Date = time.Now()
	Todo.Status = false

	jsonErr := c.BindJSON(&Todo)

	if jsonErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, jsonErr)
		return
	}

	if Todo.Name == "" {
		helper.RespondWithError(c, http.StatusBadRequest, "Please provide valid name.")
		return
	}

	mongoSession := configuration.ConnectDb(constants.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(constants.Database).C("todo")
	err := getCollection.Insert(Todo)

	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, constants.TodoCreatedSuccess, Todo)
}

func GetTodos(c *gin.Context) {
	val := reflect.ValueOf(c.Keys["user_id"])

	resp := []bson.M{}
	mongoSession := configuration.ConnectDb(constants.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(constants.Database).C("todo")

	err := getCollection.Find(bson.M{"user": bson.ObjectIdHex(val.String())}).All(&resp)

	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, constants.ListFetchedSuccess, resp)
}

func UpdateTodo(c *gin.Context) {
	todoID := c.Param("todoId")
	var Todo models.Todo

	jsonErr := c.BindJSON(&Todo)
	if jsonErr != nil {
		fmt.Print(jsonErr)
		helper.RespondWithError(c, http.StatusBadRequest, jsonErr)
		return
	}

	_id := bson.ObjectIdHex(todoID)
	if Todo.Name == "" {
		helper.RespondWithError(c, http.StatusBadRequest, "Please provide valid name!")
		return
	}

	resp := bson.M{}
	mongoSession := configuration.ConnectDb(constants.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(constants.Database).C("todo")

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"name": Todo.Name}},
		ReturnNew: true,
	}

	info, err := getCollection.Find(bson.M{"_id": _id}).Apply(change, &resp)

	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Print(info)
	helper.RespondWithSuccess(c, http.StatusOK, constants.ListFetchedSuccess, resp)
}

func DeleteTodo(c *gin.Context) {
	todoID := c.Param("todoId")
	_id := bson.ObjectIdHex(todoID)
	resp := bson.M{}
	mongoSession := configuration.ConnectDb(constants.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(constants.Database).C("todo")

	err := getCollection.RemoveId(_id)

	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, constants.DeletedSuccessfully, resp)
}
