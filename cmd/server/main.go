package server

import (
	"github.com/rs/cors"
	"github.com/sflewis2970/trivia-api/config"
	controllers "github.com/sflewis2970/trivia-api/controller"
	"log"
	"net/http"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get config data
	cfgData, cfgDataErr := config.Get().GetData(config.UPDATE_CONFIG_DATA)
	if cfgDataErr != nil {
		log.Fatal("Error getting config data: ", cfgDataErr)
	}

	// Create controller
	controller := controllers.New()

	// setup Cors
	log.Print("Setting up CORS...")
	corsOptionsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	corsHandler := corsOptionsHandler.Handler(controller.Router)

	// Server Address info
	addr := cfgData.Host + ":" + cfgData.Port
	log.Print("The address used by the service is: ", addr)

	// Start Server
	log.Print("Web service server is ready...")

	// Listen and Serve
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}
