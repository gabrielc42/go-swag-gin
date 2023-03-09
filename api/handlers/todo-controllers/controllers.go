package todocontrollers

import (
	configuration "/api/configuration"
	constants "/api/constant"
	helper "/api/helpers"
	models "/api/models/todo"
	"net/http"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func createTodo() {
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

	if Todo.name == "" {
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
		helper.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, constants.TodoCreatedSuccess, Todo)
}
