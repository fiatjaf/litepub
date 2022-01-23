package litepub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func FetchActor(theirId string) (*Actor, error) {
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
