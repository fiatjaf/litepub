package litepub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func request(url string, result interface{}) error {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Accept", "application/activity+json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, result); err != nil {
		str := string(b)
		if len(str) > 100 {
			str = str[:100] + "..."
		}

		objName := reflect.TypeOf(result).Name()
		return fmt.Errorf("error unmarshaling %s (\"%s\"): %w", objName, str, err)
	}

	return nil
}

func FetchActorFromIdentifier(theirId string) (*Actor, error) {
	if spl := strings.Split(theirId, "@"); len(spl) == 2 {
		// it's an identifier like name@domain.com
		// turn it into an url/id
		if url, err := FetchActivityPubURL(theirId); err != nil {
			return nil, fmt.Errorf("webfinger fetch failed for '%s': %w", theirId, err)
		} else {
			return FetchActor(url)
		}
	} else {
		return nil, fmt.Errorf("'%s' is not an identifier like name@domain.com", theirId)
	}
}

func FetchActor(url string) (*Actor, error) {
	var actor Actor
	err := request(url, &actor)
	return &actor, err
}

func FetchNote(url string) (*Note, error) {
	var note Note
	err := request(url, &note)
	return &note, err
}

func FetchFollowing(url string) ([]string, error) {
	var following []string

	var collection OrderedCollection
	err := request(url, &collection)
	if err != nil {
		return nil, err
	}

	// "first" may be an OrderedCollectionPage object or a URL
	var page OrderedCollectionPage[string]
	json.Unmarshal(collection.First, &page)
	if page.Id == "" {
		var pageUrl string
		json.Unmarshal(collection.First, &pageUrl)
		err := request(pageUrl, &page)
		if err != nil {
			return nil, err
		}
	}

	// here we know page is the correct object
	following = append(following, page.OrderedItems...)

	// do we need to fetch more pages?
	if page.TotalItems > len(following) {
		for len(following) < 400 /* hard limit at 400 */ {
			if page.Next == "" || len(page.OrderedItems) == 0 {
				break
			}

			var nextPage OrderedCollectionPage[string]
			err := request(page.Next, &page)
			if err != nil {
				// ignore the error and return what we got
				return following, nil
			}

			following = append(following, nextPage.OrderedItems...)
			page = nextPage
		}
	}

	return following, nil
}

func FetchNotes(url string) ([]Note, error) {
	var notes []Note

	var collection OrderedCollection
	err := request(url, &collection)
	if err != nil {
		return nil, err
	}

	// "first" may be an OrderedCollectionPage object or a URL
	var page OrderedCollectionPage[Create[Note]]
	json.Unmarshal(collection.First, &page)
	if page.Id == "" {
		var pageUrl string
		json.Unmarshal(collection.First, &pageUrl)
		err := request(pageUrl, &page)
		if err != nil {
			return nil, err
		}
	}

	// here we know page is the correct object
	notes = append(notes, mapSlice(page.OrderedItems, unwrapCreate)...)

	// do we need to fetch more pages?
	if page.TotalItems > len(notes) {
		for len(notes) < 100 /* hard limit at 100 */ {
			if page.Next == "" || len(page.OrderedItems) == 0 {
				break
			}

			var nextPage OrderedCollectionPage[Create[Note]]
			err := request(page.Next, &page)
			if err != nil {
				// ignore the error and return what we got
				return notes, nil
			}

			notes = append(notes, mapSlice(nextPage.OrderedItems, unwrapCreate)...)
			page = nextPage
		}
	}

	return notes, nil
}
