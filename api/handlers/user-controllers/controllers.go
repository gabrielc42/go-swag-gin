package usercontrollers

import (
	"encoding/json"
	"net/http"
	"time"

	configuration "github.com/gabrielc42/api/configuration"
	CONSTANTS "github.com/gabrielc42/api/constant"
	helper "github.com/gabrielc42/api/helpers"
	models "github.com/gabrielc42/api/models/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var userData models.User

	params := struct {
		ProfilePhoto string `json:"ProfilePhoto,default=pawan"`
	}{}

	c.Bind(&params)

	email := c.PostForm("email")
	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	password := c.PostForm("password")
	mobile := c.PostForm("mobile")

	userData.Email = email
	userData.FirstName = firstName
	userData.LastName = lastName
	userData.Password = password
	userData.Mobile = mobile

	validateInputErr := helper.ValidateSignupInput(userData)
	if validateInputErr != "" {
		helper.RespondWithError(c, http.StatusBadRequest, validateInputErr)
		return
	}
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if hashError != nil {
		helper.RespondWithError(c, http.StatusBadRequest, hashError)
		return
	}
	userData.Password = string(hashedPassword)
	userData.Date = time.Now().Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	mongoSession := configuration.ConnectDb(CONSTANTS.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(CONSTANTS.Database).C("user")

	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	indexErr := getCollection.EnsureIndex(index)

	if indexErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, indexErr)
		return
	}

	userData.ID = bson.NewObjectId()

	token, expire, _ := helper.CreateToken(userData.ID.Hex())
	userData.Token = token
	userData.TokenExpiresAt = expire

	err := getCollection.Insert(userData)

	if err != nil {
		if mgo.IsDup(err) == true {
			helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.AccountAlreadyExists)
			return
		}
		helper.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.CreatedSuccessfully, userData)
}

func Login(c *gin.Context) {
	var Login models.UserLogin
	var userData models.User

	userErr := json.NewDecoder(c.Request.Body).Decode(&Login)
	LoginErr := helper.ValidateLoginInput(Login)
	if LoginErr != "" {
		helper.RespondWithError(c, http.StatusBadRequest, LoginErr)
		return
	}
	if userErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, userErr)
		return
	}

	mongoSession := configuration.ConnectDb(CONSTANTS.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(CONSTANTS.Database).C("user")

	err := getCollection.Find(bson.M{"email": Login.Email}).One(&userData)
	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.AccountNotExists)
		return
	}

	PassErr := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(Login.Password))
	if PassErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.IncorrectPassword)
		return
	}
	token, expire, _ := helper.CreateToken(userData.ID.Hex())
	userData.TokenExpiresAt = expire
	userData.Token = token
	err = getCollection.UpdateId(bson.ObjectIdHex(userData.ID.Hex()), userData)
	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.LoggedInSuccess, userData)
}
