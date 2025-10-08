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
	"strings"

	"github.com/spf13/viper"
)

const (
	// Name of the program
	ExporterName = "infoblox-exporter"

	// MetricsPrefix the prefix for all internal metrics
	MetricsPrefix = "infoblox_exporter_"
)

// ExporterNameAsEnv return the ExportName as an env prefix
func ExporterNameAsEnv() string {
	return strings.ToUpper(strings.ReplaceAll(ExporterName, "-", "_"))
}

// SetDefaultValues define all default values
func SetDefaultValues() {

	// If set as env vars use the ExporterName as prefix like ACI_STREAMER_PORT for the port var
	viper.SetEnvPrefix(ExporterNameAsEnv())

	// All fields with . will be replaced with _ for ENV vars
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// infoblox-exporter
	viper.SetDefault("exporter.port", 9597)
	viper.BindEnv("exporter.port")
	viper.SetDefault("exporter.logfile", "")
	viper.BindEnv("exporter.logfile")
	viper.SetDefault("exporter.logformat", "json")
	viper.BindEnv("exporter.logformat")
	viper.SetDefault("exporter.config", "config")
	viper.BindEnv("exporter.config")

	// Basic auth exporter
	//viper.SetDefault("exporter.basic_auth.username)
	viper.BindEnv("exporter.basic_auth.username")
	//viper.SetDefault("exporter.basic_auth.password", "")
	viper.BindEnv("exporter.basic_auth.password")

	// Infoblox master
	viper.SetDefault("infoblox.master", "")
	viper.BindEnv("infoblox.master")
	viper.SetDefault("infoblox.master_port", "")
	viper.BindEnv("infoblox.master_port")
	viper.SetDefault("infoblox.wapi_version", "")
	viper.BindEnv("infoblox.wapi_version")
	viper.SetDefault("infoblox.username", "")
	viper.BindEnv("infoblox.username")
	viper.SetDefault("infoblox.password", "")
	viper.BindEnv("infoblox.password")
	viper.SetDefault("infoblox.ssl_verify", false)
	viper.BindEnv("infoblox.ssl_verify")
	viper.SetDefault("infoblox.http_request_timeout", 20)
	viper.BindEnv("infoblox.http_request_timeout")
	viper.SetDefault("infoblox.http_pool_connections", 10)
	viper.BindEnv("infoblox.http_pool_connections")

}
