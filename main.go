// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Copyright 2023 Anders Håål and Telenor AB

package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

var version = "undefined"

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", ExporterName)
		fmt.Printf("Version %s\n", version)
		flag.PrintDefaults()
	}

	SetDefaultValues()

	flag.Int("p", viper.GetInt("port"), "The port to start on")

	logFile := flag.String("logfile", viper.GetString("logfile"), "Set log file, default stdout")
	logFormat := flag.String("logformat", viper.GetString("logformat"), "Set log format to text or json, default json")

	configFile := flag.String("config", viper.GetString("config"), "Set configuration file, default config.yaml")
	usage := flag.Bool("u", false, "Show usage")
	versionFlag := flag.Bool("v", false, "Show version")
	writeConfig := flag.Bool("default", false, "Write default config")
	//customer := flag.String("customer", viper.GetString("customer"), "The customer to use in the config")
	//output := flag.String("output", viper.GetString("output"), "The output file, default stdout")

	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	if *logFormat == "text" {
		log.SetFormatter(&log.TextFormatter{})
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)

	viper.SetConfigName(*configFile) // name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.infoblox-exporter")
	viper.AddConfigPath("/usr/local/etc/infoblox-exporter")
	viper.AddConfigPath("/etc/infoblox-exporter")

	if *usage {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("infoblox-exporter, version %s\n", version)
		os.Exit(0)
	}

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		log.SetOutput(f)
	}

	if *writeConfig {
		err := viper.WriteConfigAs("./infoblox-exporter_default_config.yaml")
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Can not write default config file")
		}
		os.Exit(0)
	}

	// Find and read the config file
	err = viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Configuration file not valid")
		os.Exit(1)
	}
	// Create a Prometheus histogram for response time of the exporter
	responseTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    MetricsPrefix + "request_duration_seconds",
		Help:    "Histogram of the time (in seconds) each request took to complete.",
		Buckets: []float64{0.001, 0.005, 0.010, 0.020, 0.100, 0.200, 0.500, 1.000, 2.000},
	},
		[]string{"url", "status", "request_type"},
	)
	http.Handle("/alive",
		logCall(promMonitor(http.HandlerFunc(alive), responseTime, "/alive")))

	http.Handle("/probe",
		logCall(promMonitor(basicAuth(http.HandlerFunc(ProbeHandler)), responseTime, "/probe")))

	log.Info(fmt.Sprintf("%s starting on port %d", ExporterName, viper.GetInt("port")))
	s := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(viper.GetInt("port")),
	}
	log.Fatal(s.ListenAndServe())
	/*
		api := NewInfobloxApi()
		defer api.Logout()
		member, err := api.GetMember("vgns0022.vgregion.se")
		fmt.Printf("%s\n", member.HostName)
		utilization, err := api.GetDhcpUtilization("10.199.73.128/26")
		fmt.Printf("%.2f\n", float64(utilization.Utilization)/1000.0)
		utilization, err = api.GetDhcpUtilization("10.31.240.160/28")
		fmt.Printf("%.2f\n", float64(utilization.Utilization)/1000.0)
		utilization, err = api.GetDhcpUtilization("10.93.25.0/24")
		fmt.Printf("%.2f\n", float64(utilization.Utilization)/1000.0)
	*/
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func logCall(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		lrw := loggingResponseWriter{ResponseWriter: w}
		requestId := nextRequestID()

		ctx := context.WithValue(r.Context(), "requestId", requestId)
		next.ServeHTTP(&lrw, r.WithContext(ctx)) // call original

		w.Header().Set("Content-Length", strconv.Itoa(lrw.length))
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
			//"endpoint":  endpoint,
			"status":    lrw.statusCode,
			"length":    lrw.length,
			"requestId": requestId,
			"exec_time": time.Since(start).Microseconds(),
		}).Info("api call")
	})

}

func nextRequestID() ksuid.KSUID {
	return ksuid.New()
}

func promMonitor(next http.Handler, ops *prometheus.HistogramVec, endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestType := "customer"
		if r.URL.Query().Get("device_id") != "" || r.URL.Query().Get("device_name") != "" {
			requestType = "single"
		}
		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)
		response := time.Since(start).Seconds()
		ops.With(prometheus.Labels{"url": endpoint, "status": strconv.Itoa(lrw.Status()), "request_type": requestType}).Observe(response)
	})
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if basic auth should be used
		if !viper.IsSet("exporter.basicAuth") {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected
			// usernames and passwords.
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(viper.GetString("exporter.basicAuth.username")))
			expectedPasswordHash := sha256.Sum256([]byte(viper.GetString("exporter.basicAuth.password")))

			// Use the subtle.ConstantTimeCompare() function to check if
			// the provided username and password hashes equal the
			// expected username and password hashes. ConstantTimeCompare
			// will return 1 if the values are equal, or 0 otherwise.
			// Importantly, we should to do the work to evaluate both the
			// username and password before checking the return values to
			// avoid leaking information.
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func alive(w http.ResponseWriter, r *http.Request) {

	var alive = fmt.Sprintf("infoblox-exporter version %s alive\n", version)
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(alive)))
	lrw := loggingResponseWriter{ResponseWriter: w}
	lrw.WriteHeader(200)

	_, err := w.Write([]byte(alive))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "endpoint": "alive"}).Error("write api response failed")
	}
}
