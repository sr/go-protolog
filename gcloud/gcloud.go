package gcloud

import (
	"net/http"

	"go.pedge.io/protolog"
)

func NewPusher(
	client *http.Client,
	projectId string,
	logName string,
) protolog.Pusher {
	return newPusher(client, projectId, logName)
}
