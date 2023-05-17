package sns

import (
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
	var p Payload
	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		log.Fatal("unable to decode payload")
	}
	return p
}

func (p Payload) Confirm() (int, error) {
	res, err := http.Get(p.SubscribeURL)
	if err != nil {
		log.Fatal("unable to confirm the subscription")
		return http.StatusNoContent, err
	}
	return res.StatusCode, nil
}
