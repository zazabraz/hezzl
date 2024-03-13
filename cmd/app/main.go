package main

import (
	"context"
	customLog "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"hezzl/internal/application"
	"os"
	"time"
)

func main() {
	time.Local = time.UTC
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	logger := logrus.New()
	logger.SetFormatter(&customLog.Formatter{
		FieldsOrder: []string{"component"},
	})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	logger.Infoln("Creating app")
	app, err := application.New(ctx, *logger.WithContext(ctx))
	if err != nil {
		cancel()
		logger.Fatalf("creating app: %s", err)
	}
	logger.Infoln("Run app!")
	err = app.Run(ctx)
	if err != nil {
		logger.Fatalf("running app: %s\n", err)
	}
}
