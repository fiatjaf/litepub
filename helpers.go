package litepub

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

func FetchInbox(theirId string) (string, error) {
	r, _ := http.NewRequest("GET", theirId, nil)
	r.Header.Set("Accept", "application/activity+json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	inbox := gjson.GetBytes(b, "inbox").String()
	if inbox == "" {
		return "", errors.New("didn't find .inbox property on " + string(b)[:100] + "...")
	}

	return inbox, nil
}
