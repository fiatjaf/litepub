package litepub

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LitePub struct {
	PrivateKey  *rsa.PrivateKey
	PublicKeyId string
}

var CONTEXT = []string{
	"https://www.w3.org/ns/activitystreams",
	"https://w3id.org/security/v1",
	"https://pleroma.site/schemas/litepub-0.1.jsonld",
}

func (l LitePub) SendSigned(url string, data interface{}) (*http.Response, error) {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	r, _ := http.NewRequest("POST", url, body)

	date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	b, _ := ioutil.ReadAll(r.Body)
	digest := sha256.Sum256(b)
	signed := fmt.Sprintf(
		"(request-target): post %s\nhost: %s\ndate: %s",
		r.URL.Path, r.Host, date)
	hashed := sha256.Sum256([]byte(signed))
	signature, err := rsa.SignPKCS1v15(rand.Reader, l.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	sigheader := fmt.Sprintf(
		`keyId="%s",headers="(request-target) host date",signature="%s",algorithm="rsa-sha256"`,
		l.PublicKeyId, base64.StdEncoding.EncodeToString(signature))

	r.Header.Set("Content-Type", "application/activity+json")
	r.Header.Set("Digest", fmt.Sprintf("SHA2-256=%x", digest))
	r.Header.Set("Signature", sigheader)
	r.Header.Set("Date", date)
	r.Header.Set("Host", r.Host)

	return http.DefaultClient.Do(r)
}

func (l LitePub) WrapCreate(note Note, createId string) (create Create) {
	return Create{
		Base: Base{
			Type: "Create",
			Id:   createId,
		},
		Actor:  note.AttributedTo,
		Object: note,
	}
}
