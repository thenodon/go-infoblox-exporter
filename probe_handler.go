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
// Copyright 2023-2025 Anders Håål

package main

import (
	"context"
	"fmt"
	"go-infoblox-exporter/probes"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var probeSuccessGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "probe_success",
	Help: "Probe call success (1=Up,0=Down)",
})
var probeDurationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "probe_duration_seconds",
	Help: "How many seconds the probe call took to complete",
})

func ProbeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	module := r.URL.Query().Get("module")

	if target == "" || module == "" {
		http.Error(w, "target and module parameters are required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(30)*time.Second)
	defer cancel()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)

	start := time.Now()
	pc := &probes.ProbeCollector{}
	registry.MustRegister(pc)

	success, err := pc.Probe(ctx, target, module)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Probe request rejected")
		http.Error(w, fmt.Sprintf("probe: %v", err), http.StatusBadRequest)
		return
	}
	duration := time.Since(start).Seconds()
	probeDurationGauge.Set(duration)

	if success {
		probeSuccessGauge.Set(1)
	} else {
		probeSuccessGauge.Set(0)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
