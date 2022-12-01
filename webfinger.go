package litepub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func FetchActivityPubURL(identifier string) (string, error) {
	spl := strings.Split(identifier, "@")
	if len(spl) != 2 {
		return "", fmt.Errorf("'%s' is not a valid identifier", identifier)
	}
	url := "https://" + spl[1] + "/.well-known/webfinger?resource=acct:" + identifier

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var wf WebfingerResponse
	if err := json.Unmarshal(b, &wf); err != nil {
		str := string(b)
		if len(str) > 100 {
			str = str[:100] + "..."
		}
		return "", fmt.Errorf("got invalid webfinger response (%s): %w", str, err)
	}

	for _, link := range wf.Links {
		if link.Type == "application/activity+json" {
			return link.Href, nil
		}
	}

	return "", fmt.Errorf("couldn't find any activitypub matching records")
}
