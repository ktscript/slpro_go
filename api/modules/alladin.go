package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"strings"
	"log"
	"dario.cat/mergo"
	"sync"
	"os"
)

// Alladin struct
type Alladin struct {
	URL        string
	AppKey     string
	SecretKey  string
	HttpClient *http.Client
}

// New instance of Alladin
func NewAlladin() *Alladin {

	url := os.Getenv("MODULE_ALLADIN_URL")
	appKey := os.Getenv("MODULE_ALLADIN_APIKEY")
	secretKey := os.Getenv("MODULE_ALLADIN_SECRETKEY")

	return &Alladin{
		URL:       url,
		AppKey:    appKey,
		SecretKey: secretKey,
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetBalance
func (a *Alladin) GetBalance() (interface{}, error) {
	return nil, nil
}

func (a *Alladin) GetNumber(platform, service, country, price, realPrice, maxPrice string, optional map[string]interface{}) (map[string]interface{}, error) {
    params := map[string]interface{}{
        "infos": []map[string]interface{}{
            {
                "productId": service,
                "abbr":      strings.ToUpper(country), 
            },
        },
    }

	response, err := a.Request("phone", params)
    if err != nil {
        return nil, err
    }

    phonesIface, phonesOk := response["phones"].([]interface{})
    if !phonesOk || len(phonesIface) == 0 {
		return map[string]interface{}{ "status": "NO_NUMBERS" }, nil
    }

    phones := phonesIface 

    phoneNodesIface, phoneNodesOk := phones[0].(map[string]interface{})["phoneNodes"].([]interface{})
    if !phoneNodesOk || len(phoneNodesIface) == 0 {
        return nil, errors.New("invalid or missing phoneNodes data")
    }

    phoneNodes := phoneNodesIface 

    phoneNode := phoneNodes[0].(map[string]interface{})
    taskID, taskIDOk := phoneNode["taskId"].(string)
    phone, phoneOk := phoneNode["phone"].(string)

    if taskIDOk && taskID != "" && phoneOk && phone != "" {
		phone = strings.TrimLeft(phone, "+")

        return map[string]interface{}{
            "platform":   platform,
            "update_id":  taskID,
            "number":     phone, 
            "price":      price,
            "real_price": realPrice,
        }, nil
    }

    return nil, errors.New("missing taskId or phone")
}


// GetStatus 
func (a *Alladin) GetStatus(ids []string) (map[string]string, error) {
	responses, err := a.Pool("code", ids)
	if err != nil {
		return nil, err
	}

	responseJSON1, err := json.MarshalIndent(responses, "", "  ") // for print debug responce info
	log.Println("Response:", string(responseJSON1))

	results := make(map[string]string)
	for _, response := range responses {
		if response != nil {
			codes, ok := response["codes"].([]interface{})
			if ok && len(codes) > 0 {
				code, ok := codes[0].(map[string]interface{})["code"].(string)
				if ok && code != "" {
					results["task_id"] = code
				}
			}
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no valid results found")
	}
	return results, nil
}

// SetStatus 
func (a *Alladin) SetStatus(id, status string) (string, error) {
	statusMap := map[string]int{
		"6": 1, // ACCESS_ACTIVATION
		"8": 2, // ACCESS_CANCEL
	}

	if code, exists := statusMap[status]; exists {
		params := map[string]interface{}{
			"type":   2,
			"taskId": id,
			"status": code,
			"msg":    "",
		}

		_, err := a.Request("feedback", params)
		if err != nil {
			return "", err
		}

		statusMapRev := map[int]string{
			1: "ACCESS_ACTIVATION",
			2: "ACCESS_CANCEL",
		}

		return statusMapRev[code], nil
	}

	return "BAD_STATUS", nil
}

func (a *Alladin) prepareRequestBody(params map[string]interface{}) ([]byte, error) {
	bodyData := map[string]interface{}{
		"appKey":    a.AppKey,
		"secretKey": a.SecretKey,
	}

	if err := mergo.Merge(&bodyData, params); err != nil {
		return nil, err
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	return body, nil
}


// Request 
func (a *Alladin) Request(method string, params map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", a.URL, method)

	body, err := a.prepareRequestBody(params)
	if err != nil {
		return nil, err
	}

	resp, err := a.HttpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, err
	}

	return responseMap, nil
}

// Pool 
func (a *Alladin) Pool(method string, ids []string) ([]map[string]interface{}, error) {
	var responses []map[string]interface{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			params := map[string]interface{}{
				"taskIds": []map[string]interface{}{
					{"taskId": id},
				},
			}

			response, err := a.Request(method, params) 
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}

			mu.Lock()
			responses = append(responses, response)
			mu.Unlock()
		}(id) 
	}

	wg.Wait()

	return responses, nil
}