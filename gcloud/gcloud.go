package gcloud

import (
	"go.pedge.io/protolog"
	"google.golang.org/api/logging/v1beta3"
)

func NewPusher(
	service *logging.ProjectsLogsEntriesService,
	projectID string,
	logName string,
) protolog.Pusher {
	return newPusher(
		service,
		projectID,
		logName,
	)
}
