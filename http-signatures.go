package litepub

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	drand "math/rand"
	"net/http"
	"strings"
	"time"
)

func SendSigned(
	privateKey *rsa.PrivateKey,
	publicKeyId string,
	target string,
	data interface{},
) (*http.Response, error) {
	if spl := strings.Split(target, "@"); len(spl) == 2 {
		// it's an identifier like name@domain.com
		// turn it into an url/id
		if url, err := FetchActivityPubURL(target); err != nil {
			return nil, fmt.Errorf("webfinger fetch failed for '%s': %w", target, err)
		} else {
			// then from that grab the url of their AP inbox
			if actor, err := FetchActor(url); err != nil {
				return nil, fmt.Errorf("actor fetch failed for '%s': %w", url, err)
			} else {
				target = actor.Inbox
			}
		}
	} else {
		// otherwise it's already an url, which we assume it's the url of an AP inbox
	}

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(data)
	r, _ := http.NewRequest("POST", target, body)

	date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	b, _ := ioutil.ReadAll(r.Body)
	digest := sha256.Sum256(b)
	payload := fmt.Sprintf(
		"(request-target): post %s\nhost: %s\ndate: %s",
		r.URL.Path, r.Host, date)
	hashed := sha256.Sum256([]byte(payload))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	sigheader := fmt.Sprintf(
		`keyId="%s",headers="(request-target) host date",signature="%s",algorithm="rsa-sha256"`,
		publicKeyId, base64.StdEncoding.EncodeToString(signature))

	r.Header.Set("Content-Type", "application/activity+json")
	r.Header.Set("Digest", fmt.Sprintf("SHA2-256=%x", digest))
	r.Header.Set("Signature", sigheader)
	r.Header.Set("Date", date)
	r.Header.Set("Host", r.Host)

	return http.DefaultClient.Do(r)
}

func CheckSignature(r *http.Request) error {
	sigheader := r.Header.Get("Signature")
	var signature []byte
	var key *rsa.PublicKey
	var payload string
	for _, entry := range strings.Split(sigheader, ",") {
		kv := strings.Split(strings.TrimSpace(entry), "=")
		if len(kv) == 2 {
			k := kv[0]
			v := strings.Trim(kv[1], `"`)

			switch k {
			case "keyId":
				actor, err := FetchActor(v)
				if err != nil {
					return fmt.Errorf("failed to fetch actor '%s': %w", v, err)
				}
				pk, err := ParsePublicKeyFromPEM(actor.PublicKey.PublicKeyPEM)
				if err != nil {
					return fmt.Errorf("failed to parse key '%s': %w", actor.PublicKey.PublicKeyPEM, err)
				}

				key = pk
			case "headers":
				payload = ""
				for _, h := range strings.Split(v, " ") {
					h = strings.TrimSpace(h)
					if h == "(request-target)" {
						payload += "(request-target): " + r.URL.Path + "\n"
					} else {
						payload += h + ": " + r.Header.Get(h) + "\n"
					}
				}
				payload = strings.TrimSuffix(payload, "\n")
			case "signature":
				sig, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					return fmt.Errorf("signature '%s' is invalid base64: %w", v, err)
				}

				signature = sig
			case "algorithm":
				if v != "" {
					return fmt.Errorf("we only support rsa-sha256, not %s", v)
				}
			}
		}
	}

	hashed := sha256.Sum256([]byte(payload))
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
}

func GeneratePrivateKey(seed [4]byte) (*rsa.PrivateKey, error) {
	sourceSeed := binary.BigEndian.Uint32(seed[:])
	source := drand.NewSource(int64(sourceSeed))
	generator := drand.New(source)
	return rsa.GenerateKey(generator, 2048)
}

func ParsePrivateKeyFromPEM(pemString string) (*rsa.PrivateKey, error) {
	decoded, _ := pem.Decode([]byte(pemString))
	if decoded == nil {
		return nil, fmt.Errorf("failed to decode PEM private key from string")
	}
	sk, err := x509.ParsePKCS1PrivateKey(decoded.Bytes)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func ParsePublicKeyFromPEM(pemString string) (*rsa.PublicKey, error) {
	decoded, _ := pem.Decode([]byte(pemString))
	pk, err := x509.ParsePKCS1PublicKey(decoded.Bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func PublicKeyToPEM(pk *rsa.PublicKey) (string, error) {
	key, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: key,
	})), nil
}
