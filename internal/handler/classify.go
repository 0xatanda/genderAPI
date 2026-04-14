package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

func ClassifyHandler(w http.ResponseWriter, r *http.Request) {

	// 1. CORS (always first)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 2. READ QUERY PARAM
	name := r.URL.Query().Get("name")

	// 3. VALIDATION (👉 PUT YOUR CODE HERE)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Name query parameter is required",
		})
		return
	}

	// 4. AFTER VALIDATION → CALL EXTERNAL API
	apiURL := fmt.Sprintf("https://api.genderize.io?name=%s", url.QueryEscape(name))
	resp, err := http.Get(apiURL)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Upstream API error",
		})
		return
	}

	defer resp.Body.Close()

	// 5. PROCESS RESPONSE
	var data GenderizeResponse
	json.NewDecoder(resp.Body).Decode(&data)

	// 6. EDGE CASE CHECK
	if data.Gender == nil || data.Count == 0 {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "No prediction available for the provided name",
		})
		return
	}

	// 7. BUSINESS LOGIC
	isConfident := data.Probability >= 0.7 && data.Count >= 100

	// 8. FINAL RESPONSE
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"name":         name,
			"gender":       *data.Gender,
			"probability":  data.Probability,
			"sample_size":  data.Count,
			"is_confident": isConfident,
			"processed_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}
