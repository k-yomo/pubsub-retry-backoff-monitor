package pubsub_retry_backoff_monitor

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"google.golang.org/api/pubsub/v1"
	"log"
	"net/http"
	"os"
	"time"
)

type PubsubRequest struct {
	PublishedAt time.Time
	CreatedAt   time.Time
}

func (p *PubsubRequest) Save() (row map[string]bigquery.Value, insertID string, err error) {
	return map[string]bigquery.Value{
		"published_at": p.PublishedAt,
		"created_at":   p.CreatedAt,
	}, p.CreatedAt.String(), nil
}

func RecordPubsubHandler(w http.ResponseWriter, r *http.Request) {
	client, err := bigquery.NewClient(r.Context(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	pubsubMessage := pubsub.ReceivedMessage{}
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		log.Fatal(err)
	}
	publishTime, err := time.Parse(time.RFC3339, pubsubMessage.Message.PublishTime)
	if err != nil {
		log.Fatal(err)
	}
	pubsubRequest := PubsubRequest{PublishedAt: publishTime, CreatedAt: time.Now()}
	inserter := client.Dataset("playground").Table("pubsub_request").Inserter()
	if err := inserter.Put(r.Context(), &pubsubRequest); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(500)
}
