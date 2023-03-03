package controllers

import (
	"github.com/gorilla/mux"
	"github.com/sflewis2970/trivia-api/handlers"
	"log"
)

// Controller structure defines teh layout of the Controller
type Controller struct {
	Router        *mux.Router
	triviaHandler *handlers.TriviaHandler
}

// Package controllers object
var controller *Controller

func (c *Controller) setupRoutes() {
	// Display log message
	log.Print("Setting up trivia api service routes")

	// Trivia routes
	c.Router.HandleFunc("/api/v1/api/getquestion", c.triviaHandler.GetQuestion).Methods("GET")
	c.Router.HandleFunc("/api/v1/api/answerquestion", c.triviaHandler.AnswerQuestion).Methods("POST")
}

// NewController function create a new Controller and initializes new Controller object
func NewController() *Controller {
	// Create controllers component
	log.Print("Creating controllers object...")
	controller = new(Controller)

	// Trivia handler
	controller.triviaHandler = handlers.NewTriviaHandler()

	// Set controllers routes
	controller.Router = mux.NewRouter()
	controller.setupRoutes()

	return controller
}
