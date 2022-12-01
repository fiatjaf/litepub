package litepub

import (
	"encoding/json"
	"time"
)

type DummyBaseContext struct{}

func (c *DummyBaseContext) UnmarshalJSON([]byte) error {
	return nil
}

func (c DummyBaseContext) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		"https://pleroma.site/schemas/litepub-0.1.jsonld",
	})
}

type Base struct {
	Context DummyBaseContext `json:"@context"`
	Id      string           `json:"id"`
	Type    string           `json:"type"`
}

type Actor struct {
	Base

	Name                      string     `json:"name"`
	PreferredUsername         string     `json:"preferredUsername"`
	ManuallyApprovesFollowers bool       `json:"manuallyApprovesFollowers"`
	Image                     ActorImage `json:"image,omitempty"`
	Icon                      ActorImage `json:"icon,omitempty"`
	Summary                   string     `json:"summary,omitempty"`
	URL                       string     `json:"url"`
	Inbox                     string     `json:"inbox"`
	Outbox                    string     `json:"outbox"`
	Followers                 string     `json:"followers"`
	Following                 string     `json:"following"`
	Published                 time.Time  `json:"published"`

	PublicKey PublicKey `json:"publicKey"`
}

type ActorImage struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PublicKey struct {
	Id           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPEM string `json:"publicKeyPem"`
}

type Accept struct {
	Base

	Object interface{} `json:"object"`
}

type OrderedCollection struct {
	Base

	TotalItems int             `json:"totalItems"`
	First      json.RawMessage `json:"first"`
}

type OrderedCollectionPage[I any] struct {
	Base

	TotalItems   int    `json:"totalItems"`
	PartOf       string `json:"partOf"`
	OrderedItems []I    `json:"orderedItems"`
	Next         string `json:"next"`
}

type Follow struct {
	Base

	Actor  string `json:"actor"`
	Object string `json:"object"`
}

type Create[O any] struct {
	Base

	Actor  string `json:"actor"`
	Object O      `json:"object"`
}

type Note struct {
	Base

	Published    time.Time `json:"published"`
	AttributedTo string    `json:"attributedTo"`
	InReplyTo    string    `json:"InReplyToAtomUri"`
	Content      string    `json:"content"`
	To           []string  `json:"to"`
	CC           []string  `json:"cc"`
}
