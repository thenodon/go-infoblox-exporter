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
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var prefixDhcpUtilization = fmt.Sprintf("%s_%s", prefix, "dhcp")

var (
	dhcpUtilization = prometheus.NewDesc(
		fmt.Sprintf("%s_%s", prefixDhcpUtilization, "utilization_ratio"),
		"Dhcp utilization",
		nil, nil,
	)
)

func probeDhcpUtilization(target string) ([]prometheus.Metric, bool) {

	var m []prometheus.Metric

	api := NewInfobloxApi()
	defer api.Logout()
	utilization, err := api.GetDhcpUtilization(target)
	if err != nil {
		return m, false
	}

	m = metricsDevice(target, utilization, m)

	return m, true
}

func metricsDevice(target string, utilization Range, m []prometheus.Metric) []prometheus.Metric {

	m = append(m, prometheus.MustNewConstMetric(dhcpUtilization, prometheus.GaugeValue, float64(utilization.Utilization)/1000.0))

	return m
}
