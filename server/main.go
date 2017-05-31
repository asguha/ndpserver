package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/asguha/ndpserver/server/controller"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	Id     int
	Name   string
	School string
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(201)
	fmt.Fprint(w, "Welcome!\n")
}

func IndexPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := User{Id: 1, Name: ps.ByName("name")}
	js, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//fmt.Fprintf(w, "%s", js)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/", IndexPost)
	router.GET("/hello/:name", Hello)

	uc := controller.NewUserController(getSession())
	router.POST("/user", uc.CreateUser)
	router.GET("/user/:id", uc.GetUser)
	router.GET("/users", uc.GetAllUsers)

	log.Info("Started server on localhost:8080")
	http.ListenAndServe(":8080", &APIServer{router})
}

type APIServer struct {
	r *httprouter.Router
}

func (s *APIServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//Catching unexpected error.
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)

		}
	}()

	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Atmosphere-Remote-User,X-Atmosphere-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	if req.Method == "GET" && (req.URL.String() == "/version" || strings.HasPrefix(req.URL.String(), "/swaggerui/") || strings.HasPrefix(req.URL.String(), "/dist/")) {
		//no log here, health check is very frequent
	} else {
		log.Info(req.Method, req.URL)

	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
