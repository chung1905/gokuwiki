package captcha

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type turnstileResponse struct {
	Success bool `json:"success"`
}

func Validate(token string, secretKey string) bool {
	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	method := "POST"

	payload := strings.NewReader("secret=" + secretKey + "&response=" + token)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(string(body))
	turnstileRes := turnstileResponse{}
	err = json.Unmarshal(body, &turnstileRes)
	if err != nil {
		return false
	}
	return turnstileRes.Success
}
