package trivia

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sflewis2970/trivia-api/common"
	"github.com/sflewis2970/trivia-api/messages"
	"io"
	"io/ioutil"
	"log"
	"time"
)

//goland:noinspection SpellCheckingInspection,SpellCheckingInspection
const (
	RapidAPIHostKey string = "X-RapidAPI-Host"
	RapidAPIKey     string = "X-RapidAPI-Key"
	RapidAPIValue   string = "1f8720c0c7msh43fe783209a6813p1833b2jsnc2300c30b9a9"

	TriviaURL          string = "https://trivia-by-api-ninjas.p.rapidapi.com/v1/trivia"
	TriviaAPIHostValue string = "trivia-by-api-ninjas.p.rapidapi.com"

	TriviaCategoryCount  int = 14
	EmptyRecordCount     int = 0
	TriviaMaxRecordCount int = 5
)

var CategoryList = [TriviaCategoryCount]string{"artliterature", "language", "sciencenature", "general", "fooddrink", "peopleplaces",
	"geography", "historyholidays", "entertainment", "toysgames", "music", "mathematics", "religionmythology", "sportsleisure"}

type TriviaResponse struct {
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type API struct {
}

// GetTrivia exported type method
func (a *API) GetTrivia(category string) (messages.Trivia, error) {
	// Initialize data store when needed
	categoryLen := len(category)
	limit := 0
	requestComplete := false
	var apiResponseErr error
	var apiResponses []TriviaResponse
	apiResponsesSize := 0
	timestamp := ""

	// validate category
	if categoryLen > 0 && !isItemInCategoryList(category) {
		errMsg := fmt.Sprintf("%s is invalid", category)
		log.Print(errMsg)
		return messages.Trivia{}, errors.New(errMsg)
	}

	// Check for duplicates for marking the request as complete
	for !requestComplete {
		// Send request to API
		apiResponses, timestamp, apiResponseErr = a.triviaRequest(category, limit)

		// Get API Response size
		apiResponsesSize = len(apiResponses)

		if apiResponsesSize > 0 {
			// When results are returned, make sure there are no duplicate answers
			if !a.containsDuplicates(apiResponses) {
				log.Print("No duplicates found...")
				requestComplete = true
			} else {
				log.Print("Found duplicates...")
			}
		} else {
			// An error occurred or no results found
			requestComplete = true
		}
	}

	// Question (Request) Response message
	var trivia messages.Trivia

	// Build API Response
	trivia.Timestamp = timestamp

	if apiResponseErr != nil {
		// If an error occurs let the client know
		return messages.Trivia{}, apiResponseErr
	} else {
		// Since the client is no longer allowed to supply a limit
		// there should be five items returned from the API
		// After getting a valid response from the API, generate a question ID
		trivia.QuestionID = uuid.New().String()
		trivia.QuestionID = common.BuildUUID(trivia.QuestionID, messages.DASH, messages.ONE_SET)
		trivia.Category = apiResponses[0].Category
		trivia.Question = apiResponses[0].Question
		trivia.Answer = apiResponses[0].Answer

		// Build choices string
		var choiceList []string
		for idx := 0; idx < apiResponsesSize; idx++ {
			choiceList = append(choiceList, apiResponses[idx].Answer)
		}

		// Shuttle list
		choiceList = common.ShuffleList(choiceList)

		// Add a message filler to the beginning of the list
		trivia.Choices = append(trivia.Choices, messages.MAKE_SELECTION_MSG)
		trivia.Choices = append(trivia.Choices, choiceList...)
	}

	return trivia, nil
}

// unexported type method
// triviaRequest is a function that sends a request to the API to retrieve the trivia
func (a *API) triviaRequest(category string, limit int) ([]TriviaResponse, string, error) {
	// Build URL string
	url := TriviaURL

	// Add optional parameters string
	// Get category string
	categoryLength := len(category)
	if categoryLength > 0 {
		url = url + "?category=" + category
	}

	// Set limit default value
	if limit == 0 {
		limit = TriviaMaxRecordCount
	}

	// Add limit string to the end of the url
	if categoryLength > 0 {
		url = url + "&limit=" + fmt.Sprint(limit)
	} else {
		url = url + "?limit=" + fmt.Sprint(limit)
	}

	headers := []common.HTTPHeader{
		{Key: RapidAPIHostKey, Value: TriviaAPIHostValue},
		{Key: RapidAPIKey, Value: RapidAPIValue},
	}

	// Create a http request
	method := "GET"
	request, requestErr := common.CreateRequest(method, url, headers, nil)
	if requestErr != nil {
		log.Print("Error creating request...")
		return nil, "", requestErr
	}

	// Execute request
	response, responseErr := common.ExecuteRequest(request)
	if responseErr != nil {
		log.Print("Error executing request...")
		return nil, "", responseErr
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)

	// Get timestamp right after receiving a valid request
	timestamp := common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")

	// Parse request body
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Print("Error reading response...", readErr)
		return nil, "", readErr
	}

	// Parse response into JSON format
	responses := make([]TriviaResponse, 0)
	unmarshalErr := json.Unmarshal(body, &responses)
	if unmarshalErr != nil {
		log.Print("Error unmarshalling response...")
		return nil, "", unmarshalErr
	}

	// Return a valid response (in JSON format) as well as a timestamp
	return responses, timestamp, nil
}

// containsDuplicates checks the slice for any duplicate items
func (a *API) containsDuplicates(items []TriviaResponse) bool {
	// Initialize the map for usage
	itemsMap := make(map[string]int)

	// Since maps uses unique keys, use the string value of answer to be the key
	for idx, item := range items {
		itemsMap[item.Answer] = idx + 1
	}

	// If the size of the map is the same size of the slice, then there are no duplicates
	if len(itemsMap) != len(items) {
		return true
	}

	// Otherwise return false
	return false
}

func New() *API {
	log.Print("Creating API object...")
	api := new(API)

	return api
}

// unexported functions
func isItemInCategoryList(item string) bool {
	for _, category := range CategoryList {
		if item == category {
			return true
		}
	}

	return false
}
