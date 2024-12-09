package api

import (
	"log"
	"fmt"
	"smslive2/api/modules"
//	"smslive2/api/validators" 
)

type QueryData struct {
	UUID string                  `json:"uuid"`
	Platform string              `json:"platform"`
	Query    QueryDetails        `json:"query"`
}

type QueryDetails struct {
	Proc  string                 `json:"proc"`
	Data  map[string]interface{} `json:"data"`
}

type TaskProcessor interface {
	GetNumber(platform, service, country, price, realPrice, maxPrice string, optional map[string]interface{}) (map[string]interface{}, error)
	GetStatus(ids []string) (map[string]string, error)
	SetStatus(id, status string) (string, error)
	GetBalance() (interface{}, error)
}

// Platform mapping
var platformProcessors = map[string]func() TaskProcessor{
	"alladin":  modules.NewAlladin,
	// remaining services
	// ...
}

func GetTaskProcessor(platform string) TaskProcessor {
	processorCreator, exists := platformProcessors[platform]
	if !exists {
		log.Printf("No processor found for platform: %s", platform)
		return nil
	}

	return processorCreator()
}

func sendResponse(responseChannel chan QueryData, data map[string]interface{}) {
	var response QueryData
	response.UUID = task.UUID
	response.Platform = task.Platform
	response.Query.Proc = task.Query.Proc
	response.Query.Data = data
	
	responseChannel <- response
}

func StartTaskProcessor(requestChannel, responseChannel chan QueryData) {
	for {
		task := <-requestChannel 
		log.Printf("Received task for platform: %s", task.Platform)

		processor := GetTaskProcessor(task.Platform)
		if processor == nil {
			sendResponse(responseChannel, map[string]interface{}{"error": "Unknown platform or invalid processor"})
			continue
		}

		var Data map[string]interface{}

		switch task.Query.Proc {
			case "getNumber":
				if err := validators.ValidateGetNumberFields(task.Query.Data); err != nil {
					log.Printf("Validation failed: %v", err)
					Data = map[string]interface{}{"error": "Unknown platform or invalid processor"}
				} else {
					result, err := processor.GetNumber(
						task.Query.Data["platform"].(string),
						task.Query.Data["service"].(string),
						task.Query.Data["country"].(string),
						task.Query.Data["price"].(string),
						task.Query.Data["realPrice"].(string),
						task.Query.Data["maxPrice"].(string),
						map[string]interface{}{}, 
					)
					if err != nil {
						log.Printf("Error calling GetNumber: %v", err)
						Data = map[string]interface{}{"error": err.Error()}
					} else {
						log.Printf("GetNumber result: %v", result)
						response.Query.Data = result
					}
				}

			case "getStatus":
				if err := validators.ValidateGetStatusFields(task.Query.Data); err != nil {
					log.Printf("Validation failed: %v", err)
					Data = map[string]interface{}{"error": err.Error()}
				} else {
					ids := task.Query.Data["ids"].([]string)
					result, err := processor.GetStatus(ids)
					if err != nil {
						log.Printf("Error calling GetStatus: %v", err)
						Data = map[string]interface{}{"error": err.Error()}
					} else {
						log.Printf("GetStatus result: %v", result)
						Data = result
					}
				}

			case "setStatus":
				if err := validators.ValidateSetStatusFields(task.Query.Data); err != nil {
					log.Printf("Validation failed: %v", err)
					Data = map[string]interface{}{"error": err.Error()}
				} else {
					id := task.Query.Data["id"].(string)
					status := task.Query.Data["status"].(string)
					result, err := processor.SetStatus(id, status)
					if err != nil {
						log.Printf("Error calling SetStatus: %v", err)
						Data = map[string]interface{}{"error": err.Error()}
					} else {
						log.Printf("SetStatus result: %s", result)
						Data = map[string]interface{}{"result": result}
					}
				}

			case "getBalance":
				result, err := processor.GetBalance()
				if err != nil {
					log.Printf("Error calling GetBalance: %v", err)
					Data = map[string]interface{}{"error": err.Error()}
				} else {
					log.Printf("GetBalance result: %v", result)
					Data = result
				}

			default:
				log.Printf("Unknown method: %s", task.Query.Proc)
				Data = map[string]interface{}{"error": "Unknown method"}
		}

		sendResponse(responseChannel, Data)
	}
}
