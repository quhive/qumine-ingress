package api

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qumine/qumine-ingress/internal/k8s"
	"github.com/qumine/qumine-ingress/internal/server"
	"github.com/sirupsen/logrus"
)

var (
	k *k8s.K8S
	s *server.Server
)

// API represents the api server
type API struct {
	httpServer *http.Server
}

// NewAPI creates a new api instance with the given host and port
func NewAPI() *API {
	r := http.NewServeMux()
	r.HandleFunc("/healthz", getHealthz)
	r.Handle("/metrics", promhttp.Handler())
	return &API{
		httpServer: &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: r,
		},
	}
}

// Start the Api
func (api *API) Start(context context.Context, k8s *k8s.K8S, server *server.Server) {
	defer api.httpServer.Close()
	logrus.WithField("addr", api.httpServer.Addr).Info("starting api...")

	k = k8s
	s = server

	go logrus.WithError(api.httpServer.ListenAndServe()).Fatal("api failed to start")

	for {
		select {
		case <-context.Done():
			return
		}
	}
}

func getHealthz(writer http.ResponseWriter, request *http.Request) {
	details := make(map[string]string)
	details["k8s"] = k.Status
	details["server"] = s.Status

	if k.Status == "up" && s.Status == "up" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte{})
	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte{})
	}
}

type healthz struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}
