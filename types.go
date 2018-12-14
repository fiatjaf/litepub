package litepub

type Base struct {
	Context []string `json:"@context,omitempty"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`
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

	TotalItems int                   `json:"totalItems"`
	First      OrderedCollectionPage `json:"first"`
}

type OrderedCollectionPage struct {
	Base

	TotalItems   int         `json:"totalItems"`
	PartOf       string      `json:"partOf"`
	OrderedItems interface{} `json:"orderedItems"`
}

type Follow struct {
	Base

	Actor  string `json:"actor"`
	Object string `json:"object"`
}

type Create struct {
	Base

	Actor  string      `json:"actor"`
	Object interface{} `json:"object"`
}

type Note struct {
	Base

	Published    string `json:"published"`
	AttributedTo string `json:"attributedTo"`
	Content      string `json:"content"`
	To           string `json:"to"`
}
