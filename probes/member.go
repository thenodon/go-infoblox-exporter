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

package probes

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var prefixMember = fmt.Sprintf("%s_%s", prefix, "member")
var memberLabels = []string{"service"}
var memberNodeLabels = []string{"service", "node_ip"}
var memberNodeInfoLabels = []string{"ha_status", "hwid", "hwtype", "node_ip", "platform"}

var (
	nodeInfo = prometheus.NewDesc(
		fmt.Sprintf("%s_%s", prefixMember, "node_info"),
		"Node info",
		memberNodeInfoLabels, nil,
	)
	nodeService = prometheus.NewDesc(
		fmt.Sprintf("%s_%s", prefixMember, "node_service"),
		"Node service (0=Failed, 1=Working, 2=Unknown)",
		memberNodeLabels, nil,
	)
	service = prometheus.NewDesc(
		fmt.Sprintf("%s_%s", prefixMember, "service"),
		"Service (0=Failed, 1=Working, 2=Unknown)",
		memberLabels, nil,
	)
)

func probeMember(target string) ([]prometheus.Metric, bool) {

	var m []prometheus.Metric

	member, err := infobloxApi.GetMember(target)
	if err != nil {
		return m, false
	}

	m = metricsMember(member, m)

	return m, true
}

func metricsMember(member Member, m []prometheus.Metric) []prometheus.Metric {

	for _, mem := range member.ServiceStatus {
		if mem.Status != "INACTIVE" {
			m = append(m, prometheus.MustNewConstMetric(service, prometheus.GaugeValue, getStatus(mem.Status), mem.Service))
		}
	}

	for _, mem := range member.Nodeinfo {

		dup := make(map[string]string)
		ip := mem.LanHaPortSetting.MgmtLan
		m = append(m, prometheus.MustNewConstMetric(nodeInfo, prometheus.GaugeValue, 1.0,
			mem.HaStatus, mem.Hwid, mem.Hwtype, ip, mem.Hwplatform))
		for _, node := range mem.ServiceStatus {
			if node.Status != "INACTIVE" {
				_, ok := dup[node.Service]
				if ok {
					continue
				} else {
					dup[node.Service] = node.Service
				}
				m = append(m, prometheus.MustNewConstMetric(nodeService, prometheus.GaugeValue, getStatus(node.Status), node.Service, ip))
			}
		}
	}

	return m
}

func getStatus(status string) float64 {
	if status == "WORKING" {
		return 1.0
	} else if status == "UNKNOWN" {
		return 2.0
	}
	return 0.0
}
