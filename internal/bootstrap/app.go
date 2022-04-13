package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/aliykh/reddit-feed/internal/config"
	"github.com/aliykh/reddit-feed/internal/http/server"

	"github.com/aliykh/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	errr "errors"
)

type App struct {

	//	main config
	config *config.Config

	//	http server
	http *http.Server

	//	log
	log *log.Factory

	// mongodb client
	mongoClient *mongo.Client

	//	tearDowns -> for graceful shutdown
	tearDowns []func()
}

func New(cfg *config.Config, ctx context.Context) *App {
	app := &App{}
	app.config = cfg

	app.log = log.NewFactory(log.ZapLogger, cfg.LogLevel)

	// todo mongodb init
	if err := app.initMongoDb(ctx); err != nil {
		app.compensate()
		app.log.Default().Fatal("mongodb init", zap.String("err", err.Error()))
	}

	if err := app.initHTTPServer(); err != nil {
		app.compensate()
		app.log.Default().Fatal("http server init", zap.String("err", err.Error()))
	}

	return app
}

func (a *App) initMongoDb(ctx context.Context) error {

	//startedCommands := make(map[int64]bson.Raw)
	//cmdMonitor := &event.CommandMonitor{
	//	Started: func(_ context.Context, evt *event.CommandStartedEvent) {
	//		startedCommands[evt.RequestID] = evt.Command
	//	},
	//	Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
	//		a.log.Default().Debug(fmt.Sprintf("Command: %v Reply: %v\n",
	//			startedCommands[evt.RequestID].String(),
	//			evt.Reply.String()),
	//		)
	//	},
	//	Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
	//		a.log.Default().Debug(fmt.Sprintf("Command: %v Failure: %v\n",
	//			startedCommands[evt.RequestID].String(),
	//			evt.Failure,
	//		))
	//	},
	//}

	heartBeatFailEvent := &event.ServerMonitor{
		ServerHeartbeatFailed: func(failedEvent *event.ServerHeartbeatFailedEvent) {
			a.log.Default().Error(fmt.Sprintf("mongo db server is not responding, cause: %s", failedEvent.Failure.Error()))
		},
		ServerHeartbeatStarted: func(startedEvent *event.ServerHeartbeatStartedEvent) {
			a.log.Default().Debug(fmt.Sprintf("checking mongo db server %v", startedEvent.ConnectionID))
		},
		ServerHeartbeatSucceeded: func(succeededEvent *event.ServerHeartbeatSucceededEvent) {
			a.log.Default().Debug(fmt.Sprintf("mongo db server is online %v", succeededEvent.ConnectionID))
		},
	}

	clientOpts := options.Client().SetHeartbeatInterval(time.Second * 30)
	clientOpts.ApplyURI(a.config.MongoAddr).SetServerMonitor(heartBeatFailEvent) //.SetMonitor(cmdMonitor)
	clientOpts.SetMaxPoolSize(20)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}

	tr := func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = client.Disconnect(ctx); err != nil {
			a.log.Default().Debug("mongo client disconnect", zap.String("err", err.Error()))
			return
		}

		a.log.Default().Info("mongo-db client shutdown")

	}

	a.tearDowns = append(a.tearDowns, tr)

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		tr()
		return err
	}

	a.mongoClient = client

	// todo add health checker for db connection availability and also for the case when the mongo db shuts down unexpectedly

	return err
}

// initHTTPServer initializes http server.
func (app *App) initHTTPServer() error {

	hs, err := server.New(app.config, app.log, app.mongoClient)

	if err != nil {
		return err
	}

	address := fmt.Sprintf(":%v", app.config.ServerPort)
	app.http = &http.Server{
		Addr:         address,
		Handler:      hs,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	return nil
}

func (app *App) compensate() {
	for _, v := range app.tearDowns {
		v()
	}
}

func (app *App) Run(ctx context.Context) {

	// run
	go func() {
		app.log.Default().Info(fmt.Sprintf("REST Server started at port: %v", app.config.ServerPort))

		if err := app.http.ListenAndServe(); err != nil && !errr.Is(err, http.ErrServerClosed) {
			app.log.Default().Error(fmt.Sprintf("Failed To Run REST Server: %s\n", err.Error()))
		}

		app.log.Default().Debug("http server has been shut down")
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := app.http.Shutdown(ctxShutDown); err != nil {
		app.log.Default().Error(fmt.Sprintf("http server shutdown failed:%s\n", err.Error()))
	}

	for _, v := range app.tearDowns {
		v()
	}

	app.log.Default().Info("server shutdown")
}
