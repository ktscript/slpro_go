package api

import "fmt"

func ValidateGetNumberFields(data map[string]interface{}) error {
	requiredFields := []string{
		"platform",
		"service",
		"country",
		"price",
		"realPrice",
		"maxPrice",
	}

	for _, field := range requiredFields {
		if value, ok := data[field].(string); !ok || value == "" {
			return fmt.Errorf("missing or invalid field: %s", field)
		}
	}

	return nil
}

func ValidateGetStatusFields(data map[string]interface{}) error {
	if ids, ok := data["ids"].([]string); !ok || len(ids) == 0 {
		return fmt.Errorf("missing or invalid field: ids")
	}
	return nil
}

func ValidateSetStatusFields(data map[string]interface{}) error {
	if id, ok := data["id"].(string); !ok || id == "" {
		return fmt.Errorf("missing or invalid field: id")
	}
	if status, ok := data["status"].(string); !ok || status == "" {
		return fmt.Errorf("missing or invalid field: status")
	}
	return nil
}

func ValidateGetBalanceFields(data map[string]interface{}) error {
	return nil
}
