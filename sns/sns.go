package sns

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Payload struct {
	Message          string `json:"Message"`
	MessageId        string `json:"MessageId"`
	Signature        string `json:"Signature"`
	SignatureVersion string `json:"SignatureVersion"`
	SigningCertURL   string `json:"SigningCertURL"`
	SubscribeURL     string `json:"SubscribeURL"`
	Subject          string `json:"Subject"`
	Timestamp        string `json:"Timestamp"`
	Token            string `json:"Token"`
	TopicArn         string `json:"TopicArn"`
	Type             string `json:"Type"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

func GetData(r io.ReadCloser) Payload {
	var payload Payload
	err := json.NewDecoder(r).Decode(&payload)
	if err != nil {
		log.Fatal("unable to decode payload")
	}
	return payload
}

func (payload Payload) Confirm() (int, error) {
	res, err := http.Get(payload.SubscribeURL)
	if err != nil {
		log.Fatal("unable to confirm the subscription")
		return http.StatusNoContent, err
	}
	return res.StatusCode, nil
}

func (payload Payload) Verify() error {
	pub := payload.downloadCertificate()
	hashed := payload.applyHash()
	sig := payload.decodeSignature()
	if payload.SignatureVersion == "1" {
		return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed, sig)
	}
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, sig)
}
