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
	"strconv"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var infobloxApi InfoBloxApi

func SetInfobloxApi(api InfoBloxApi) {
	infobloxApi = api
}

type InfoBloxConfiguration struct {
	Master              string
	Version             string
	Port                int64
	Username            string
	Password            string
	SSLVerify           bool
	HTTPRequestTimeout  int
	HTTPPoolConnections int
}

func NewInfoBloxConfiguration() InfoBloxConfiguration {
	return InfoBloxConfiguration{
		Master:              viper.GetString("infoblox.master"),
		Version:             viper.GetString("infoblox.wapi_version"),
		Port:                viper.GetInt64("infoblox.master_port"),
		Username:            viper.GetString("infoblox.username"),
		Password:            viper.GetString("infoblox.password"),
		SSLVerify:           viper.GetBool("infoblox.ssl_verify"),
		HTTPRequestTimeout:  viper.GetInt("infoblox.http_request_timeout"),
		HTTPPoolConnections: viper.GetInt("infoblox.http_pool_connections"),
	}
}

type Member struct {
	ibclient.IBBase
	Ref                      string                   `json:"_ref,omitempty"`
	HostName                 string                   `json:"host_name,omitempty"`
	ConfigAddrType           string                   `json:"config_addr_type,omitempty"`
	PLATFORM                 string                   `json:"platform,omitempty"`
	ServiceTypeConfiguration string                   `json:"service_type_configuration,omitempty"`
	Nodeinfo                 []ibclient.Nodeinfo      `json:"node_info,omitempty"`
	TimeZone                 string                   `json:"time_zone,omitempty"`
	ServiceStatus            []ibclient.Servicestatus `json:"service_status,omitempty"`
}

func (m *Member) ObjectType() string {
	return "member"
}

func NewMember(nodeName string) *Member {
	return &Member{
		HostName: nodeName,
	}
}

type Range struct {
	ibclient.IBBase
	Ref         string      `json:"_ref,omitempty"`
	Cidr        string      `json:"network,omitempty"`
	Ea          ibclient.EA `json:"extattrs"`
	Comment     string      `json:"comment"`
	Utilization int64       `json:"dhcp_utilization"`
}

func (r *Range) ObjectType() string {
	return "range"
}

func NewRange(cidr string, comment string, ea ibclient.EA) *Range {
	return &Range{
		Cidr:    cidr,
		Ea:      ea,
		Comment: comment,
	}
}

type InfoBloxApi struct {
	Conn *ibclient.Connector
}

func NewInfobloxApi() InfoBloxApi {
	config := NewInfoBloxConfiguration()

	hostConfig := ibclient.HostConfig{
		Host:    config.Master,
		Version: config.Version,
	}

	authConfig := ibclient.AuthConfig{
		Username:   config.Username,
		Password:   config.Password,
		ClientCert: nil,
		ClientKey:  nil,
	}
	transportConfig := ibclient.NewTransportConfig(strconv.FormatBool(config.SSLVerify), config.HTTPRequestTimeout,
		config.HTTPPoolConnections)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
	conn, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		log.Error("Failed to connect", err)
	}

	return InfoBloxApi{Conn: conn}
}

func (i InfoBloxApi) GetDhcpUtilization(network string) (Range, error) {
	var res []Range
	net := NewRange(network, "", nil)

	queryAttribute := map[string]string{
		"network":        network,
		"_return_fields": "extattrs,network,dhcp_utilization,comment",
	}
	qp := ibclient.NewQueryParams(false, queryAttribute)
	err := i.Conn.GetObject(net, "", qp, &res)

	if err != nil {
		log.Error("Failed to get network", err)
		return *net, err
	}

	return res[0], nil
}

func (i InfoBloxApi) GetMember(nodeName string) (Member, error) {
	var res []Member
	net := NewMember(nodeName)

	queryAttribute := map[string]string{
		"host_name":      nodeName,
		"_return_fields": "extattrs,host_name,node_info,service_status",
	}
	qp := ibclient.NewQueryParams(false, queryAttribute)
	err := i.Conn.GetObject(net, "", qp, &res)

	if err != nil {
		log.Error("Failed to get node", err)
		return *net, err
	}
	return res[0], nil
}

func (i InfoBloxApi) Logout() {
	i.Conn.Logout()
}
