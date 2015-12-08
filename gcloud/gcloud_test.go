package gcloud_test

import (
	"log"
	"os"

	"go.pedge.io/protolog"
	"go.pedge.io/protolog/gcloud"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/logging/v1beta3"
)

func ExampleExamples() {
	projectID, _ := os.LookupEnv("GCLOUD_PROJECT_ID")
	logName := "protolog"
	client, err := google.DefaultClient(
		context.Background(),
		logging.LoggingWriteScope,
		logging.CloudPlatformScope,
	)
	if err != nil {
		log.Fatal(err)
	}
	service, err := logging.New(client)
	if err != nil {
		log.Fatal(err)
	}
	logger := protolog.NewStandardLogger(
		gcloud.NewPusher(
			service.Projects.Logs.Entries,
			projectID,
			logName,
		),
	)
	logger.Infoln("Hello from protolog!")
}
