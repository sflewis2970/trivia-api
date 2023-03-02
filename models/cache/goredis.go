package cache

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
	REDIS_DB_NAME_MSG      string = "GO_REDIS: "
	REDIS_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	// REDIS_GET_CONFIG_ERROR      string = "Getting config error...: "
	// REDIS_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	// REDIS_OPEN_ERROR            string = "Open method not implemented..."
	REDIS_MARSHAL_ERROR        string = "Marshaling error...: "
	REDIS_UNMARSHAL_ERROR      string = "Unmarshalling error...: "
	REDIS_INSERT_ERROR         string = "Insert error...: "
	REDIS_ITEM_NOT_FOUND_ERROR string = "Item not found...: "
	REDIS_GET_ERROR            string = "Get error...: "
	// REDIS_UPDATE_ERROR          string = "Update error..."
	REDIS_DELETE_ERROR string = "Delete error...: "
	// REDIS_RESULTS_ERROR         string = "Results error...: "
	// REDIS_ROWS_AFFECTED_ERROR   string = "Rows affected error...: "
	REDIS_PING_ERROR string = "Error pinging in-memory cache server...: "
	// REDIS_CONVERSION_ERROR      string = "Conversion error...: "
)

var cacheModel *DataModel

type DataModel struct {
	cfgData  *config.CfgData
	memCache *redis.Client
}

// Ping database server, since this is local to the server make sure the object for storing data is created
func (dm *DataModel) Ping() error {
	ctx := context.Background()

	statusCmd := dm.memCache.Ping(ctx)
	pingErr := statusCmd.Err()
	if pingErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_PING_ERROR, pingErr)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dm *DataModel) Insert(trivia messages.Trivia) error {
	ctx := context.Background()

	var tTable messages.TriviaTable
	tTable.Question = trivia.Question
	tTable.Category = trivia.Category
	tTable.Answer = trivia.Answer

	byteStream, marshalErr := json.Marshal(tTable)
	if marshalErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_MARSHAL_ERROR, marshalErr)
		return marshalErr
	}

	log.Print("Adding a new record to map, ID: ", trivia.QuestionID)
	setErr := dm.memCache.Set(ctx, trivia.QuestionID, byteStream, time.Duration(0)).Err()
	if setErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_INSERT_ERROR, setErr)
		return setErr
	}

	return nil
}

// Get a single record from table
func (dm *DataModel) Get(questionID string) (messages.TriviaTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	var tTable messages.TriviaTable
	ctx := context.Background()
	getResult, getErr := dm.memCache.Get(ctx, questionID).Result()
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
func (dm *DataModel) Update(updatedRec messages.Trivia) {
	log.Println("Updating record in the map")

	ctx := context.Background()

	var tTable messages.TriviaTable
	tTable.Question = updatedRec.Question
	tTable.Category = updatedRec.Category
	tTable.Answer = updatedRec.Answer

	// Send update message to cache
	dm.memCache.Set(ctx, updatedRec.QuestionID, tTable, 0)
}

// Delete a single record from table
func (dm *DataModel) Delete(questionID string) error {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	ctx := context.Background()
	delErr := dm.memCache.Del(ctx, questionID).Err()
	if delErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_DELETE_ERROR, delErr)
		return delErr
	}

	return nil
}

func NewCacheModel(cfgData *config.CfgData) *DataModel {
	// Initialize go-cache in-memory cache model
	log.Print("Creating goRedis dbModel object...")
	cacheModel = new(DataModel)

	// Assign config data
	cacheModel.cfgData = cfgData

	// Define go-redis cache settings
	log.Print(REDIS_DB_NAME_MSG + REDIS_CREATE_CACHE_MSG)

	// Define connection variables
	var redisOptions *redis.Options

	// The config package handles reading the environment variables and parsing the url.
	// Once the external packages access the values, the environment has already been taken
	// care of.
	addr := cacheModel.cfgData.Redis.URL + ":" + cacheModel.cfgData.Redis.Port
	log.Print("redis URL...: ", cacheModel.cfgData.Redis.URL)
	log.Print("redis Port...: ", cacheModel.cfgData.Redis.Port)
	log.Print("The redis address is...: ", addr)

	redisOptions = &redis.Options{
		Addr:     addr,                              // redis Server Address,
		Password: cacheModel.cfgData.Redis.Password, // set password
		DB:       0,                                 // use default DB
	}

	// Create go-redis in-memory cache
	cacheModel.memCache = redis.NewClient(redisOptions)

	return cacheModel
}
