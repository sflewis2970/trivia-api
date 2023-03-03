package models

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sflewis2970/trivia-api/config"
	"github.com/sflewis2970/trivia-api/messages"
	"log"
	"time"
)

const (
	// REDIS_TLS_URL Redis Constants
	REDIS_PASSWORD         string = "REDIS_PASSWORD"
	REDIS_DB_NAME_MSG      string = "GO_REDIS: "
	REDIS_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	REDIS_MARSHAL_ERROR        string = "Marshaling error...: "
	REDIS_UNMARSHAL_ERROR      string = "Unmarshalling error...: "
	REDIS_INSERT_ERROR         string = "Insert error...: "
	REDIS_ITEM_NOT_FOUND_ERROR string = "Item not found...: "
	REDIS_GET_ERROR            string = "Get error...: "
	REDIS_DELETE_ERROR         string = "Delete error...: "
	REDIS_PING_ERROR           string = "Error pinging in-memory cache server...: "
)

type Redis struct {
	TLS_URL  string `json:"tls_url"`
	URL      string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type RedisModel struct {
	cfgData  *config.CfgData
	memCache *redis.Client
}

var redisModel *RedisModel

// Ping database server, since this is local to the server make sure the object for storing data is created
func (rm *RedisModel) Ping() error {
	ctx := context.Background()

	statusCmd := rm.memCache.Ping(ctx)
	pingErr := statusCmd.Err()
	if pingErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_PING_ERROR, pingErr)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (rm *RedisModel) Insert(trivia messages.Trivia) error {
	ctx := context.Background()

	// var tTable messages.TriviaTable
	// tTable.Question = trivia.Question
	// tTable.Category = trivia.Category
	// tTable.Answer = trivia.Answer

	byteStream, marshalErr := json.Marshal(messages.Trivia{})
	if marshalErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_MARSHAL_ERROR, marshalErr)
		return marshalErr
	}

	// log.Print("Adding a new record to map, ID: ", trivia.QuestionID)
	setErr := rm.memCache.Set(ctx, "", byteStream, time.Duration(0)).Err()
	if setErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_INSERT_ERROR, setErr)
		return setErr
	}

	return nil
}

// Get a single record from table
func (rm *RedisModel) Get(questionID string) (messages.TriviaTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	var tTable messages.TriviaTable
	ctx := context.Background()
	getResult, getErr := rm.memCache.Get(ctx, questionID).Result()
	if getErr == redis.Nil {
		log.Print(REDIS_DB_NAME_MSG + REDIS_ITEM_NOT_FOUND_ERROR)
		return messages.TriviaTable{}, nil
	} else if getErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_GET_ERROR, getErr)
		return messages.TriviaTable{}, getErr
	} else {
		unmarshalErr := json.Unmarshal([]byte(getResult), &tTable)
		if unmarshalErr != nil {
			log.Print(REDIS_DB_NAME_MSG+REDIS_UNMARSHAL_ERROR, unmarshalErr)
			return messages.TriviaTable{}, unmarshalErr
		}
	}

	return tTable, nil
}

// Update a single record in table
func (rm *RedisModel) Update(updatedRec messages.Trivia) {
	log.Println("Updating record in the map")

	ctx := context.Background()

	var tTable messages.TriviaTable
	// tTable.Question = updatedRec.Question
	// tTable.Category = updatedRec.Category
	// tTable.Answer = updatedRec.Answer

	// Send update message to cache
	rm.memCache.Set(ctx, "", tTable, 0)
}

// Delete a single record from table
func (rm *RedisModel) Delete(questionID string) error {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	ctx := context.Background()
	delErr := rm.memCache.Del(ctx, questionID).Err()
	if delErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_DELETE_ERROR, delErr)
		return delErr
	}

	return nil
}

func NewRedisModel() *RedisModel {
	// Initialize go-cache in-memory cache model
	log.Print("Creating goRedis dbModel object...")
	redisModel = new(RedisModel)

	// Get config data
	redisModel.cfgData = config.NewConfig().LoadCfgData()

	// Define go-redis cache settings
	log.Print(REDIS_DB_NAME_MSG + REDIS_CREATE_CACHE_MSG)

	// Define connection variables
	var redisOptions *redis.Options

	// The config package handles reading the environment variables and parsing the url.
	// Once the external packages access the values, the environment has already been taken
	// care of.
	redisAddr := redisModel.cfgData.RedisURL + ":" + redisModel.cfgData.RedisPort
	log.Print("The redis address is...: ", redisAddr)

	redisOptions = &redis.Options{
		Addr:     redisAddr, // redis Server Address,
		Password: "",        // cacheModel.cfgData.Redis.Password, // set password
		DB:       0,         // use default DB
	}

	// Create go-redis in-memory cache
	redisModel.memCache = redis.NewClient(redisOptions)

	return redisModel
}
