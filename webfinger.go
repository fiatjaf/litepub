package litepub

import (
	"fmt"
	"net/http"
	"strings"
)

type WebfingerResponse struct {
	Subject string          `json:"subject"`
	Links   []WebfingerLink `json:"links"`
}

type WebfingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func HandleWebfingerRequest(r *http.Request) (string, error) {
	rsc := r.URL.Query().Get("resource")
	parts := strings.Split(rsc, "acct:")
	if len(parts) != 2 {
		return "", fmt.Errorf("resource querystring param is wrong: %s", rsc)
	}

	account := parts[1]
	parts = strings.Split(account, "@")
	if len(parts) != 2 || parts[0] == "" {
		return "", fmt.Errorf("account not formatted correctly: %s", account)
	}

	name := parts[0]
	return name, nil
}
