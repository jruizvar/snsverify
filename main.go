package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"snsverify/sns"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type snsWriter struct {
	output io.Writer
}

func (w *snsWriter) Write(b []byte) (n int, err error) {
	n, err = w.output.Write(b)
	return
}

func (w *snsWriter) writeMessage(_ http.ResponseWriter, r *http.Request) {
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
	// encode payload and write to file
	u, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatal("unable to encode payload")
	}
	fmt.Fprintf(w, "Payload\n\n%s\n\n", u)
}

func main() {
	output, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	// Write payload to output file
	w := snsWriter{output}
	r.Post("/sns", w.writeMessage)
	http.ListenAndServe(":80", r)
}
