package controllers

import (
	"github.com/gorilla/mux"
	"github.com/sflewis2970/trivia-api/config"
	"github.com/sflewis2970/trivia-api/handlers/trivia"
	"log"
)

// Controller struct definition
type Controller struct {
	cfgData       *config.CfgData
	Router        *mux.Router
	triviaHandler *trivia.TriviaHandler
}

// Package controller object
var controller *Controller

func (c *Controller) setupRoutes() {
	// Display log message
	log.Print("Setting up trivia service routes")

	// Trivia routes
	c.Router.HandleFunc("/api/v1/trivia/questions", c.triviaHandler.GetTriviaQuestion).Methods("GET")
	c.Router.HandleFunc("/api/v1/trivia/questions", c.triviaHandler.SubmitTriviaAnswer).Methods("POST")
}

// New Export functions
func New() *Controller {
	// Create controller component
	log.Print("Creating controller object...")
	controller = new(Controller)

	// Load config data
	var getCfgDataErr error
	controller.cfgData, getCfgDataErr = config.Get().GetData()
	if getCfgDataErr != nil {
		log.Print("Error getting config data: ", getCfgDataErr)
		return nil
	}

	// Trivia handlers
	controller.triviaHandler = trivia.New()

	// Set controller routes
	controller.Router = mux.NewRouter()
	controller.setupRoutes()

	return controller
}
