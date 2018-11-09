package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ninh0gauch0/homerecipes/hrscli"
	"github.com/ninh0gauch0/homerecipes/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	version = "0.9.0-beta"
)

var (
	baseContext   = context.Background()
	contextLogger *log.Entry
	exitChan      chan bool
	app           *cli.App
)

// init - always called at the begining
func init() {
	logger := log.StandardLogger()
	logger.Formatter = &log.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}
	logger.Out = os.Stdout
	logger.SetLevel(log.InfoLevel)
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	contextLogger = logger.WithFields(log.Fields{
		"service": "Home Recipes Service",
	})
}

func main() {
	// Create a cli app
	contextLogger.Infof("Starting app...")

	app = cli.NewApp()
	app.Version = version
	app.Description = "Home Recipes App is an application made for my girlfriend; she's an awesome cook :D"
	app.Authors = []cli.Author{
		{Name: "Julianinho", Email: "julianinhogaucho10@gmail.com"},
	}

	// Commands definition
	app.Commands = hrscli.GetCommands()

	// Flags definition
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug, d", Usage: "If set, the log level is set to DEBUG"},
		cli.BoolFlag{Name: "error, e", Usage: "If set, the log level is set to ERROR"},
		cli.BoolFlag{Name: "quiet, q", Usage: "If set, the log level is set to FATAL"},
	}

	// Define the server command. This command is in charge of server start
	serverCommnad := cli.Command{}
	serverCommnad.Name = "start"
	serverCommnad.Usage = "Starts the HR Server"

	// Additional command flags
	serverCommnad.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "8089",
			Usage: "Server port",
		},
	}

	// Starts the server with a given configuration
	serverCommnad.Action = func(c *cli.Context) {
		// Object server declaration & context initialization
		serverContext, cancelFunc := context.WithCancel(baseContext)
		s := server.Server{
			Ctx: serverContext,
		}

		if c.GlobalBool("debug") {
			contextLogger.Logger.SetLevel(log.DebugLevel)
		}

		if c.GlobalBool("error") {
			contextLogger.Logger.SetLevel(log.ErrorLevel)
		}

		if c.GlobalBool("quiet") {
			contextLogger.Logger.SetLevel(log.FatalLevel)
		}

		// Sets the logger to the server
		s.SetLogger(contextLogger)
		defer cancelFunc()

		// Config definition
		config := map[string]string{
			"addr": fmt.Sprintf(":%s", c.String("port")),
		}
		// Init the server
		if s.Init() {
			// Starting the server
			exitChan = s.Start(config)

			if exitChan != nil {
				// Waiting for terminal signals
				sigs := make(chan os.Signal)
				signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
				sig := <-sigs
				switch sig {
				case syscall.SIGINT:
					fallthrough
				case os.Interrupt:
					fallthrough
				case os.Kill:
					fallthrough
				case syscall.SIGTERM:
					exitChan <- true
				}
			}
		} else {
			contextLogger.Fatalf("Fail on initialization...")
		}
	}

	app.Commands = append(app.Commands, serverCommnad)
	app.EnableBashCompletion = true
	app.Run(os.Args)
}
