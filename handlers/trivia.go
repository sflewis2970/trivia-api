package handlers

import (
	"encoding/json"
	"github.com/sflewis2970/trivia-api/external/OpenTriviaAPI"
	"github.com/sflewis2970/trivia-api/messages"
	"github.com/sflewis2970/trivia-api/models"
	"log"
	"net/http"
)

type TriviaHandler struct {
	openTrivia  *OpenTriviaAPI.OpenTrivia
	triviaModel *models.TriviaModel
}

var triviaHandler *TriviaHandler

// GetTriviaQuestion is a http handler that receives a client "GET" request.
// Clients will send a request when they want to receive a api question from the api API.
// The format used is: 'http://<server-name>:8080/api/add?category=name'. category is optional
// When 'category' is supplied the api API returns a question related to the requested category
// When 'category' is omitted, the api API determines whether not the selected question is related
// to a category.
// The request returns a QuestionResponse object.
// The format for QuestionResponse is:
//       {"questionid": "<random_id>",
//        "question": "<question from api API>",
//        "category": "<category is not required and could be blank>",
//        "choices": "<choices are generated from API. One answer is correct, the others are incorrect>",
//        "timestamp": "<formatted string of when the API returned the question>",
//        "warning": "<optional warning message>",
//        "error": "<optional error message>"}
func (th *TriviaHandler) GetQuestion(rw http.ResponseWriter, r *http.Request) {
	// Display a log message
	log.Print("data received from client...")

	// Get category from query parameter
	category := r.URL.Query().Get("category")

	var qResponse messages.QuestionResponse

	// Process API Get Request
	triviaData, triviaErr := th.openTrivia.GetTrivia(category)
	if triviaErr != nil {
		log.Print("Error encoding json...:", triviaErr)

		// Update QuestionResponse struct
		qResponse.Error = triviaErr.Error()

		// Update HTTP header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		encodeResponse(rw, qResponse)
		return
	}

	// Send request to model to insert api question
	insertErr := th.triviaModel.AddQuestion(triviaData)

	// Add question to data store
	if insertErr != nil {
		log.Print("Error encoding json...:", insertErr)

		// Update QuestionResponse struct
		qResponse.QuestionID = ""
		qResponse.Category = ""
		qResponse.Question = ""
		qResponse.Choices = []string{}
		qResponse.Error = insertErr.Error()

		// Update HTTP header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		encodeResponse(rw, qResponse)
		return
	}

	// Update HTTP Header
	rw.WriteHeader(http.StatusCreated)

	// Write JSON to stream
	encodeResponse(rw, qResponse)

	// Display a log message
	log.Print("data sent back to client...")
}

// SubmitTriviaAnswer is a http handler that receives a response message from the client.
// The client is responding to question received from the api API.
// The request uses the form of: 'http://<server-name>:8080//api/v1/api/questions' including a
// json object:
//        "questionid": "<id received in the question response>",
//        "response": "<answer question from list of choices>"
// The client will receive a response in the form of the following:
//       "question": "<the question the client provided the answer for>",
//       "timestamp": "<formatted string of when the API returned the question>",
//       "category": "<if the question is linked to a category that information will be provided here>",
//       "response": "<the response the client provided>",
//       "answer": "<the answer to the question>",
//       "message": "<message to client whether question was answered correctly>",
//       "warning": "<optional warning message>",
//       "error": "<optional error message>"
func (th *TriviaHandler) AnswerQuestion(rw http.ResponseWriter, r *http.Request) {
	var aRequest messages.AnswerRequest
	var aResponse messages.AnswerResponse

	// Read JSON from stream
	decodeErr := json.NewDecoder(r.Body).Decode(&aRequest)
	if decodeErr != nil {
		log.Print("Error decoding json...: ", decodeErr)

		// Update AnswerResponse
		aResponse.Error = decodeErr.Error()

		// Update HTTP Header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		encodeResponse(rw, aResponse)
		return
	}

	// Send a request to the model for the answer
	var getErr error
	aResponse, getErr = th.triviaModel.GetAnswer(aRequest)

	if getErr != nil {
		log.Print("Error getting api answer...: ", getErr)

		// Update AnswerResponse
		aResponse.Error = getErr.Error()

		// Update HTTP Header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		encodeResponse(rw, aResponse)
		return
	}

	// Send a request to the model to delete the question
	deleteErr := th.triviaModel.DeleteQuestion(aRequest.QuestionID)

	if deleteErr != nil {
		log.Print("Error deleting api question...: ", deleteErr)

		// Update AnswerResponse
		aResponse.Error = deleteErr.Error()

		// Update HTTP Header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		encodeResponse(rw, aResponse)
		return
	}

	// Send OK status
	rw.WriteHeader(http.StatusOK)

	// Encode response
	encodeResponse(rw, aResponse)

	// Display a log message
	log.Print("data sent back to client...")
}

type MessageSet interface {
	messages.QuestionResponse | messages.AnswerResponse
}

func encodeResponse[T MessageSet](rw http.ResponseWriter, response T) {
	// Write JSON to stream
	encodeErr := json.NewEncoder(rw).Encode(response)
	if encodeErr != nil {
		log.Print("Error encoding json...:", encodeErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func NewTriviaHandler() *TriviaHandler {
	triviaHandler := new(TriviaHandler)

	// Create api api
	triviaHandler.openTrivia = OpenTriviaAPI.NewOpenTrivia()

	// Create api model
	triviaHandler.triviaModel = models.NewTriviaModel()

	return triviaHandler
}
