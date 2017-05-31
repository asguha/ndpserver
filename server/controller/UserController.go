package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/asguha/ndpserver/server/models"
	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ndpdbname = "ndpulse"

//UserController struct
type UserController struct {
	//mongo session
	session *mgo.Session
}

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.USER{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.ID = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB(ndpdbname).C("users").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// GetUser retrieves an individual user resource
func (uc UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub user
	u := []models.USER{}

	// Fetch user
	log.Info("fetching all users")
	if err := uc.session.DB(ndpdbname).C("users").Find(bson.M{}).All(&u); err != nil {
		log.Info(err)
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")
	log.Info(id)

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.USER{}

	if err := uc.session.DB(ndpdbname).C("users").FindId(oid).One(&u); err != nil {
		log.Info(err)
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
