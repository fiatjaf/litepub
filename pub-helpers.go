package litepub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var CONTEXT = []string{
	"https://www.w3.org/ns/activitystreams",
	"https://w3id.org/security/v1",
	"https://pleroma.site/schemas/litepub-0.1.jsonld",
}

func WrapCreate(note Note, createId string) (create Create) {
	return Create{
		Base: Base{
			Type: "Create",
			Id:   createId,
		},
		Actor:  note.AttributedTo,
		Object: note,
	}
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

func FetchActor(theirId string) (*Actor, error) {
	if spl := strings.Split(theirId, "@"); len(spl) == 2 {
		// it's an identifier like name@domain.com
		// turn it into an url/id
		if url, err := FetchActivityPubURL(theirId); err != nil {
			return nil, fmt.Errorf("webfinger fetch failed for '%s': %w", theirId, err)
		} else {
			theirId = url
		}
	} else {
		// otherwise it's already an url/id like https://domain.com/pub/actor/name
	}

	r, err := http.NewRequest("GET", theirId, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/activity+json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var actor Actor
	if err := json.Unmarshal(b, &actor); err != nil {
		str := string(b)
		if len(str) > 100 {
			str = str[:100] + "..."
		}
		return nil, fmt.Errorf("error unmarshaling actor (\"%s\"): %w", str, err)
	}

	return &actor, nil
}
