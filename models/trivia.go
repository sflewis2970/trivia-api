package models

import (
	"github.com/sflewis2970/trivia-api/common"
	"github.com/sflewis2970/trivia-api/config"
	"github.com/sflewis2970/trivia-api/messages"
	"log"
	"time"
)

type TriviaModel struct {
	cfgData    *config.CfgData
	redisModel *RedisModel
}

var triviaModel *TriviaModel

func (tm *TriviaModel) AddQuestion(qRequest messages.Trivia) error {
	insertErr := tm.redisModel.Insert(qRequest)
	if insertErr != nil {
		errMsg := "Error inserting record...: "
		log.Print(errMsg, insertErr)
	}

	return insertErr
}

func (tm *TriviaModel) GetAnswer(aRequest messages.AnswerRequest) (messages.AnswerResponse, error) {
	// AnswerResponse
	var aResponse messages.AnswerResponse

	// Get timestamp right after receiving a valid request
	common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")

	// Send request to get question from Redis cache
	_, getErr := tm.redisModel.Get(aRequest.QuestionID)
	if getErr != nil {
		errMsg := "Get record error...: "
		log.Print(errMsg, getErr)
		aResponse.Error = errMsg
		return aResponse, getErr
	} else {
		// Build AnswerResponse message
	}

	return aResponse, nil
}

func (tm *TriviaModel) DeleteQuestion(questionID string) error {
	// Send request to delete question from Redis cache
	deleteErr := tm.redisModel.Delete(questionID)
	if deleteErr != nil {
		errMsg := "Delete record error...: "
		log.Print(errMsg, deleteErr)
	}

	return deleteErr
}

func (tm *TriviaModel) CfgData() *config.CfgData {
	return tm.cfgData
}

func NewTriviaModel() *TriviaModel {
	log.Print("Creating model object...")
	triviaModel := new(TriviaModel)

	// Get config data
	triviaModel.cfgData = config.NewConfig().LoadCfgData()

	// New model (cacheModel)
	triviaModel.redisModel = NewRedisModel()

	return triviaModel
}
