package gcloud

import (
	"time"

	"github.com/golang/protobuf/jsonpb"
	"go.pedge.io/proto/time"
	"go.pedge.io/protolog"
	"google.golang.org/api/logging/v1beta3"
)

const customServiceName = "compute.googleapis.com"

var (
	marshaler = &jsonpb.Marshaler{}

	// https://cloud.google.com/logging/docs/api/ref/rest/v1beta3/projects.logs.entries/write#LogSeverity
	severityName = map[protolog.Level]string{
		protolog.Level_LEVEL_NONE:  "DEFAULT",
		protolog.Level_LEVEL_DEBUG: "DEBUG",
		protolog.Level_LEVEL_INFO:  "INFO",
		protolog.Level_LEVEL_WARN:  "WARNING",
		protolog.Level_LEVEL_ERROR: "ERROR",
		protolog.Level_LEVEL_FATAL: "ERROR",
		protolog.Level_LEVEL_PANIC: "ALERT",
	}
)

type pusher struct {
	service   *logging.ProjectsLogsEntriesService
	projectID string
	logName   string
}

func newPusher(
	service *logging.ProjectsLogsEntriesService,
	projectID string,
	logName string,
) *pusher {
	return &pusher{
		service,
		projectID,
		logName,
	}
}

func (p *pusher) Push(entry *protolog.Entry) error {
	logEntry, err := newLogEntry(entry)
	if err != nil {
		return err
	}
	_, err = p.service.Write(
		p.projectID,
		p.logName,
		&logging.WriteLogEntriesRequest{
			Entries: []*logging.LogEntry{logEntry},
		},
	).Do()
	return err
}

func (p *pusher) Flush() error {
	return nil
}

func newLogEntry(entry *protolog.Entry) (*logging.LogEntry, error) {
	payload, err := marshaler.MarshalToString(entry)
	if err != nil {
		return nil, err
	}
	metadata := &logging.LogEntryMetadata{
		ServiceName: customServiceName,
		Severity:    severityName[entry.Level],
	}
	if entry.Timestamp != nil {
		metadata.Timestamp = prototime.TimestampToTime(entry.Timestamp).Format(time.RFC3339)
	}
	return &logging.LogEntry{
		InsertId:    entry.Id,
		TextPayload: payload,
		Metadata:    metadata,
	}, nil
}
