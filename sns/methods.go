package sns

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"log"
	"net/http"
)

func (p Payload) downloadCertificate() (pub *rsa.PublicKey) {
	res, err := http.Get(p.SigningCertURL)
	if err != nil {
		log.Fatal("unable to download certificate")
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(body)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("failed to decode PEM block containing certificate")
	}
	cert, _ := x509.ParseCertificate(block.Bytes)
	pub = cert.PublicKey.(*rsa.PublicKey)
	return pub
}

func (p Payload) applyHash() (hashed []byte) {
	var concatenate string
	switch p.Type {
	case "SubscriptionConfirmation", "UnsubscribeConfirmation":
		concatenate = "Message" + "\n" + p.Message + "\n" +
			"MessageId" + "\n" + p.MessageId + "\n" +
			"SubscribeURL" + "\n" + p.SubscribeURL + "\n" +
			"Timestamp" + "\n" + p.Timestamp + "\n" +
			"Token" + "\n" + p.Token + "\n" +
			"TopicArn" + "\n" + p.TopicArn + "\n" +
			"Type" + "\n" + p.Type + "\n"
	case "Notification":
		concatenate = "Message" + "\n" + p.Message + "\n" +
			"MessageId" + "\n" + p.MessageId + "\n" +
			"Subject" + "\n" + p.Subject + "\n" +
			"Timestamp" + "\n" + p.Timestamp + "\n" +
			"TopicArn" + "\n" + p.TopicArn + "\n" +
			"Type" + "\n" + p.Type + "\n"
	}
	if p.SignatureVersion == "1" {
		hash := sha1.Sum([]byte(concatenate))
		hashed = hash[:]
		return
	}
	hash := sha256.Sum256([]byte(concatenate))
	hashed = hash[:]
	return hashed
}

func (p Payload) decodeSignature() []byte {
	signature, err := base64.StdEncoding.DecodeString(p.Signature)
	if err != nil {
		log.Fatal("unable to decode signature")
	}
	return []byte(signature)
}
