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
// Copyright 2023 Anders Håål

package probes

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "infoblox"

type ProbeCollector struct {
	metrics []prometheus.Metric
}

type TargetMetadata struct {
	VersionMajor int
	VersionMinor int
}

type probeFunc func(target string) ([]prometheus.Metric, bool)

type probeDetailedFunc struct {
	name     string
	function probeFunc
}

func (p *ProbeCollector) Probe(ctx context.Context, target string, modules string) (bool, error) {

	success := true
	var aProbe probeDetailedFunc

	switch modules {
	case "member_services":
		aProbe = probeDetailedFunc{"member_services", probeMember}
	case "dhcp_utilization":
		aProbe = probeDetailedFunc{"dhcp_utilization", probeDhcpUtilization}
	default:
		return false, nil
	}

	m, ok := aProbe.function(target)
	if !ok {
		success = false
	}
	p.metrics = append(p.metrics, m...)

	return success, nil
}

func (p *ProbeCollector) Collect(c chan<- prometheus.Metric) {
	// Collect result of new probe functions
	for _, m := range p.metrics {
		c <- m
	}
}

func (p *ProbeCollector) Describe(c chan<- *prometheus.Desc) {
}
