package captcha

import (
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Println(err)
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(string(body))
	turnstileRes := turnstileResponse{}
	err = json.Unmarshal(body, &turnstileRes)
	if err != nil {
		return false
	}
	return turnstileRes.Success
}
