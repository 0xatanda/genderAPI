package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type APIError struct {
	Code    int
	Message string
}

func GetGenderData(name string) (map[string]interface{}, *APIError) {
	apiURL := fmt.Sprintf("https://api.genderize.io?name=%s", url.QueryEscape(name))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, &APIError{Code: 502, Message: "Upstream API error"}
	}
	defer resp.Body.Close()

	var data struct {
		Gender      *string `json:"gender"`
		Probability float64 `json:"probability"`
		Count       int     `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, &APIError{Code: 500, Message: "Internal server error"}
	}

	if data.Gender == nil || data.Count == 0 {
		return map[string]interface{}{
			"status":  "error",
			"message": "No prediction available for the provided name",
		}, nil
	}

	isConfident := data.Probability >= 0.7 && data.Count >= 100

	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"name":         name,
			"gender":       *data.Gender,
			"probability":  data.Probability,
			"sample_size":  data.Count,
			"is_confident": isConfident,
			"processed_at": time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}
