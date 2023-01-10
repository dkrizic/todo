package main

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	todo "github.com/dkrizic/todo/api/todo"
	"github.com/gorilla/mux"
	muxlogrus "github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	listenAddress = "0.0.0.0:8000"
)

func main() {
	log.Info("Starting app")
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/notification", NotificationHandler).Methods("POST", "OPTIONS")
	http.Handle("/", r)
	r.Use(muxlogrus.NewLogger().Middleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("listenAddress", listenAddress).Info("Starting listener")
	log.Fatal(srv.ListenAndServe())
	log.Info("Stopping listener")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	event, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		log.Print("failed to parse CloudEvent from request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	log.WithField("event", event).WithField("data", event.Data()).Info("Received event")
	var change todo.Change
	json.Unmarshal(event.Data(), &change)
	log.WithFields(log.Fields{
		"beforeId":          change.Before.Id,
		"beforeTitle":       change.Before.Title,
		"beforeDescription": change.Before.Description,
		"afterId":           change.After.Id,
		"afterTitle":        change.After.Title,
		"afterDescription":  change.After.Description,
	}).Info("Received event")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
