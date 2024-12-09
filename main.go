package main

import (
	"log"
	"smslive2/api"           
	// "smslive/kafka"         
	// "smslive/memory"        
	"smslive2/api/modules"       
)

func main() {
	inputChannel := make(chan string)
	outputChannel := make(chan string)

	type Query struct {
		Proc string                 `json:"proc"`
		Data map[string]interface{} `json:"data"`
	}
	
	type Task struct {
		UUID     string `json:"uuid"`
		Platform string `json:"platform"`
		Query    Query  `json:"query"`
	}

	// test channel data
	task := Task{
		UUID:     "12345",
		Platform: "alladin",
		Query: Query{
			Proc: "getNumber",
			Data: map[string]interface{}{
				"platform":  "alladin",
				"service":   "37",
				"country":   "US",
				"price":     "10",
				"realPrice": "9",
				"maxPrice":  "100",
				"optional":  nil, 
			},
		},
	}

	inputChannel <- task   // test task

	// go memory.MonitorMemoryUsage()

	go api.StartAPIServer(inputChannel, outputChannel)

	// go kafka.StartKafkaConsumer(inputChannel)

	// producer, err := kafka.StartKafkaProducer(outputChannel)
	// if err != nil {
	// 	log.Fatal("Failed to start Kafka producer: ", err)
	// }
	// defer producer.Close()

	modules.RegisterModules(resultChannel)

	select {}
}
