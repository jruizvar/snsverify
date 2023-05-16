package main

import (
	"apichi/sns"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type snsWriter struct {
	f io.Writer
}

func (w *snsWriter) Write(b []byte) (n int, err error) {
	n, err = w.f.Write(b)
	return
}

func (f *snsWriter) writeMessage(w http.ResponseWriter, r *http.Request) {
	// write header to file
	f.Write([]byte("Header\n\n"))
	for key, field := range r.Header {
		f.Write([]byte(key + ":\t"))
		for _, v := range field {
			f.Write([]byte(v + ", "))
		}
		f.Write([]byte("\n"))
	}
	f.Write([]byte("\n"))
	// decode payload
	payload := sns.GetData(r.Body)
	defer r.Body.Close()
	// confirm subscription
	if payload.Type == "SubscriptionConfirmation" {
		res, _ := payload.Confirm()
		fmt.Println("Subscription status code:", res)
	}
	// verify sns signature
	if err := payload.Verify(); err != nil {
		fmt.Printf("Error from verification: %s\n", err)
	} else {
		fmt.Println("Verification succeded")
	}
	// write payload to file
	u, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatal("unable to encode payload")
	}
	fmt.Fprintf(f, "Payload\n\n%s\n\n", u)
}

func main() {
	f, err := os.OpenFile(
		"/home/ec2-user/apichi/data/messages.txt",
		os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	// Write header and payload in messages.txt
	w := snsWriter{f}
	r.Post("/sns", w.writeMessage)
	http.ListenAndServe(":80", r)
}
