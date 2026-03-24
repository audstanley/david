package cli

import (
	"errors"
	"fmt"
	syslog "log"
	"net/http"

	"github.com/audstanley/david/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/webdav"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the WebDAV server",
	Long:  `Start the WebDAV server with optional configuration overrides.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunServer(cmd, args)
	},
}

func init() {

	serverCmd.Flags().StringP("config", "c", "", "Path to configuration file")
	serverCmd.Flags().StringP("host", "H", "", "Override host address")
	serverCmd.Flags().StringP("port", "p", "", "Override port")
	serverCmd.Flags().BoolP("debug", "d", false, "Enable debug logging")
	serverCmd.Flags().BoolP("production", "", false, "Enable production (JSON) logging")
	serverCmd.Flags().String("hash-algorithm", "", "Hash algorithm: bcrypt, argon2, scrypt (default: from config)")

	viper.BindPFlags(serverCmd.Flags())
}

func RunServer(cmd *cobra.Command, args []string) {
	configPath, _ := cmd.Flags().GetString("config")
	hostOverride, _ := cmd.Flags().GetString("host")
	portOverride, _ := cmd.Flags().GetString("port")
	debugEnabled, _ := cmd.Flags().GetBool("debug")
	productionEnabled, _ := cmd.Flags().GetBool("production")
	hashAlgorithm, _ := cmd.Flags().GetString("hash-algorithm")

	ProductionFormatter := &logrus.JSONFormatter{}
	NonProductionFormatter := &logrus.TextFormatter{}
	logrus.SetFormatter(ProductionFormatter)
	logrus.SetLevel(logrus.DebugLevel)

	config := app.ParseConfig(configPath)

	if hostOverride != "" {
		config.Address = hostOverride
	}
	if portOverride != "" {
		config.Port = portOverride
	}

	if productionEnabled {
		config.Log.Production = true
	}
	if debugEnabled {
		config.Log.Debug = true
	}

	if hashAlgorithm != "" {
		config.Hash.Algorithm = hashAlgorithm
	}

	logger := logrus.New()
	if config.Log.Production {
		logger.Formatter = ProductionFormatter
		logrus.WithField("production", config.Log.Production).Debug("Production mode enabled")
	} else {
		logger.Formatter = NonProductionFormatter
		logrus.WithField("production", config.Log.Production).Debug("Production mode disabled")
		logrus.SetFormatter(NonProductionFormatter)
	}
	if config.Log.Debug {
		logrus.WithField("debug", config.Log.Debug).Debug("Debug mode enabled")
	} else {
		logrus.WithField("debug", config.Log.Debug).Debug("Debug mode has now been disabled from config")
		logrus.SetLevel(logrus.InfoLevel)
	}
	writer := logger.Writer()
	defer writer.Close()
	syslog.SetOutput(writer)

	wdHandler := webdav.Handler{
		Prefix: config.Prefix,
		FileSystem: &app.Dir{
			Config: config,
		},
		LockSystem: webdav.NewMemLS(),
		Logger: func(request *http.Request, err error) {
			if config.Log.Error && err != nil {
				logrus.Error(err)
			}
		},
	}

	a := &app.App{
		Config:  config,
		Handler: &wdHandler,
	}

	http.Handle("/", wrapRecovery(app.NewBasicAuthWebdavHandler(a), config))
	connAddr := fmt.Sprintf("%s:%s", config.Address, config.Port)

	if config.TLS != nil {
		logrus.WithFields(logrus.Fields{
			"address":  config.Address,
			"port":     config.Port,
			"security": "TLS",
		}).Info("Server is starting and listening")
		logrus.Fatal(http.ListenAndServeTLS(connAddr, config.TLS.CertFile, config.TLS.KeyFile, nil))
	} else {
		logrus.WithFields(logrus.Fields{
			"address":  config.Address,
			"port":     config.Port,
			"security": "none",
		}).Info("Server is starting and listening")
		logrus.Fatal(http.ListenAndServe(connAddr, nil))
	}
}

func wrapRecovery(handler http.Handler, config *app.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch t := err.(type) {
				case string:
					logrus.Printf("panic type: %T, value: %v", err, err)
					logrus.WithFields(logrus.Fields{"error": err, "writer": w}).Warn("An error occurred handling a webdav request")
					logrus.WithError(errors.New(t)).Error("An error occurred handling a webdav request")
				case error:
					logrus.Printf("panic type: %T, value: %v", err, err)
					logrus.WithFields(logrus.Fields{"error": err, "writer": w}).Warn("An error occurred handling a webdav request")
					logrus.WithError(t).Error("An error occurred handling a webdav request")
				}
			}
		}()

		if len(config.Cors.Origin) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", config.Cors.Origin)
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			if config.Cors.Credentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		handler.ServeHTTP(w, r)
	})
}
