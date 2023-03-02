package common

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type HTTPHeader struct {
	Key   string
	Value string
}

// GetFormattedTime Build formatted time string
func GetFormattedTime(timeNow time.Time, timeFormat string) string {
	return timeNow.Format(timeFormat)
}

// GetWorkingDir Get working directory
func GetWorkingDir() (string, error) {
	workingDir, getErr := os.Getwd()
	if getErr != nil {
		log.Print("Error getting working directory...")
		return "", getErr
	}

	return workingDir, nil
}

// BuildUUID Build UUID string
func BuildUUID(uuid string, delimiter string, nbrOfGroups int) string {
	newUUID := ""

	uuidList := strings.Split(uuid, delimiter)
	for key, value := range uuidList {
		if key < nbrOfGroups {
			newUUID = newUUID + value
		}
	}

	return newUUID
}

// BuildDelimitedStr Utility to build strings seperated by a delimiter

// ShuffleList Utility to move string item to a different position within the list
func ShuffleList(strList []string) []string {
	rand.Shuffle(len(strList), func(idx1, idx2 int) {
		strList[idx1], strList[idx2] = strList[idx2], strList[idx1]
	})

	return strList
}

func CreateRequest(method string, url string, headers []HTTPHeader, httpBody io.Reader) (*http.Request, error) {
	// Create new http request
	request, requestErr := http.NewRequest(method, url, httpBody)
	if requestErr != nil {
		log.Print("A request error has occurred...")
		return nil, requestErr
	}

	// Setup request headers
	for _, header := range headers {
		request.Header.Add(header.Key, header.Value)
	}

	return request, nil
}

func ExecuteRequest(request *http.Request) (*http.Response, error) {
	// Get response from http request
	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		log.Print("A response error has occurred...")
		return nil, responseErr
	}

	return response, nil
}
