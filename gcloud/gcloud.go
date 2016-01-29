// Package protolog_gcloud integrates protolog with the Google Cloud Logging service.
//
// See https://cloud.google.com/logging/docs/ for more details.
package protolog_gcloud

import (
	"github.com/sr/protolog"
	"google.golang.org/api/logging/v1beta3"
)

// NewPusher creates a new protolog.Pusher that logs using the Google Cloud Logging API.
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
