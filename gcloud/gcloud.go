package gcloud

import (
	"net/http"

	"go.pedge.io/protolog"
)

func NewPusher(
	client *http.Client,
	projectID string,
	logName string,
) protolog.Pusher {
	return newPusher(
		client,
		projectID,
		logName,
	)
}
