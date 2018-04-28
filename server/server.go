package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Init the configuration needed to start the server
func (s *Server) Init() bool {
	s.router = mux.NewRouter()
	s.router.PathPrefix("/hrs/")
	//Reading configuration file
	dat, err := ioutil.ReadFile("jsons/mongoconfig.json")
	check(err)
	var result MongoConf
	check(json.Unmarshal(dat, &result))

	s.initialized = true
	return true
}

// Start the server
func (s *Server) Start(config map[string]string) chan bool {

	// Recovering config server
	addr, ok := config["addr"]
	if !ok {
		addr = ":8080"
	}
	s.Addr = addr

	if s.initialized != true {
		err := s.Init()
		if err {
			return nil
		}
	}

	s.logger.Infof("Starting server....")
	s.worker = &Worker{}
	s.worker.Init(s.Ctx, s.GetLogger())

	s.addRoutes()

	exitChan := make(chan bool)

	// Go routines and channel to orchestate
	go func() {
		<-exitChan
		s.logger.Infoln("Stopping server")
		// Server shutdown
		err := s.Server.Shutdown(s.Ctx)

		if err != nil {
			s.logger.Errorln("Error shutdowning server - error: ", err.Error())
		}
	}()
	go func() {
		log.Printf("Listening on.... %s", s.Addr)
		log.Fatal(http.ListenAndServe(addr, s.router))
	}()

	return exitChan
}

// addRoutes - Define API routes
func (s *Server) addRoutes() {

	/** RECIPES ENDPOINTS**/
	s.router.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var recipe Recipe
		err := decoder.Decode(&recipe)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		hrsResp := s.worker.CreateRecipe(recipe)

		if hrsResp.GetError() != nil {
			s.logger.Errorln("[POST] - ERROR: %s", hrsResp.GetError())
		}

		// TODO:return response object
		fmt.Fprintf(w, "Recipe created\n")
	}).Methods("POST")

	s.router.HandleFunc("/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.GetRecipeByID(id)
		s.logger.Infoln("Recipe returned: %s", hrsResp.RespObj)

		// TODO:return response object

		fmt.Fprintf(w, "You've requested the recipe: %s\n", id)
	}).Methods("GET")

	/** OTHER ENDPOINTS **/
	s.router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "WTF\n")
	}).Methods("GET")
}

//TODO: Change error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}
