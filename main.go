package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/figassis/mysql-backup/app"
	"github.com/figassis/mysql-backup/app/config"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const (
	ENV_CONFIG_FILE     = "CONFIG_FILE"
	DEFAULT_CONFIG_FILE = "config.yaml"

	ENV_LOG_LEVEL = "LOG_LEVEL"
)

func main() {

	// configure logger
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// get log level
	logLevel := strings.Trim(os.Getenv(ENV_LOG_LEVEL), "\n\r ")
	if len(logLevel) > 0 {
		tLogLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Fatalf("failed to parse log level, level given: %v, error: %+v", logLevel, err)
		}

		logrus.SetLevel(tLogLevel)
	}

	// get config file
	configFile := strings.Trim(os.Getenv(ENV_CONFIG_FILE), "\n\r ")
	if len(configFile) <= 0 {
		configFile = DEFAULT_CONFIG_FILE
	}

	configFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		logrus.Fatalf("failed to convert config file to absolute path: %v", configFile)
	}

	// load config
	cfg := config.NewWithDefaults()
	logrus.Infof("loading config from: %v", configFile)
	cfg.LoadYaml(configFileAbs)

	// print cfg
	//logrus.Infof("configuration: \n%v", cfg.ToString())

	// validate cfg
	validationErrors := cfg.Validate()
	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			logrus.Errorf("validation failed: %+v", err)
		}

		logrus.Fatalf("validation failed, exiting")
	}

	instance := app.NewApp(cfg)
	signals := make(chan bool, 1)

	if cfg.Schedule == "" {
		logrus.Print("Starting one off backup")
		instance.Run()
	} else {
		logrus.Printf("Starting scheduled backups: %s", cfg.Schedule)
		var wg sync.WaitGroup
		wg.Add(1)
		c := cron.New()

		//This is important to avoid locks by different PIDs
		go runSchedule(&wg, signals, instance)

		c.AddFunc(cfg.Schedule, func() { signals <- true })
		c.Start()
		wg.Wait()
	}

}

func runSchedule(Wg *sync.WaitGroup, signals <-chan bool, instance *app.App) {
	defer Wg.Done()

	logrus.Printf("Running initial backup")
	instance.Run()
	for {
		select {
		case signal := <-signals:
			switch signal {
			case true:
				instance.Run()
			}
		}
	}
}
