package models

import (
	"github.com/sflewis2970/trivia-api/common"
	"github.com/sflewis2970/trivia-api/config"
	"github.com/sflewis2970/trivia-api/messages"
	"github.com/sflewis2970/trivia-api/models/cache"
	"log"
	"time"
)

type TriviaModel struct {
	cfgData    *config.CfgData
	redisModel *cache.DataModel
}

func (m *TriviaModel) AddTriviaQuestion(qRequest messages.Trivia) error {
	insertErr := m.redisModel.Insert(qRequest)
	if insertErr != nil {
		errMsg := "Error inserting record...: "
		log.Print(errMsg, insertErr)
	}

	return insertErr
}

func (m *TriviaModel) GetTriviaAnswer(aRequest messages.AnswerRequest) (messages.AnswerResponse, error) {
	// AnswerResponse
	var aResponse messages.AnswerResponse

	// Get timestamp right after receiving a valid request
	timestamp := common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")

	// Send request to get question from Redis cache
	triviaTable, getErr := m.redisModel.Get(aRequest.QuestionID)
	if getErr != nil {
		errMsg := "Get record error...: "
		log.Print(errMsg, getErr)
		aResponse.Error = errMsg
		return aResponse, getErr
	} else {
		// Build AnswerResponse message
		if len(triviaTable.Question) > 0 {
			aResponse.Question = triviaTable.Question
			aResponse.Category = triviaTable.Category
			aResponse.Answer = triviaTable.Answer
			aResponse.Response = aRequest.Response
			aResponse.Timestamp = timestamp

			if aRequest.Response == triviaTable.Answer {
				aResponse.Correct = true
				aResponse.Message = m.cfgData.Messages.CongratsMsg
			} else {
				aResponse.Correct = false
				aResponse.Message = m.cfgData.Messages.TryAgainMsg
			}
		}
	}

	return aResponse, nil
}

func (m *TriviaModel) DeleteTriviaQuestion(questionID string) error {
	// Send request to delete question from Redis cache
	deleteErr := m.redisModel.Delete(questionID)
	if deleteErr != nil {
		errMsg := "Delete record error...: "
		log.Print(errMsg, deleteErr)
	}

	return deleteErr
}

func (m *TriviaModel) GetCfgData() *config.CfgData {
	return m.cfgData
}

func NewTriviaModel() *TriviaModel {
	log.Print("Creating model object...")
	model := new(TriviaModel)

	// Get config data
	var cfgDataErr error
	model.cfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print("Error getting config data: ", cfgDataErr)
		return nil
	}

	// New model (cacheModel)
	model.redisModel = cache.NewCacheModel(model.cfgData)

	return model
}
