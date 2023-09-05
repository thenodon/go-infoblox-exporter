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
	viper.SetDefault("port", 9597)
	viper.BindEnv("port")
	viper.SetDefault("logfile", "")
	viper.BindEnv("logfile")
	viper.SetDefault("logformat", "json")
	viper.BindEnv("logformat")
	viper.SetDefault("config", "config")
	viper.BindEnv("config")
	viper.SetDefault("output", "")
	viper.BindEnv("output")

	// Infoblox master
	viper.SetDefault("master", "")
	viper.BindEnv("master")
	viper.SetDefault("wapi_version", "")
	viper.BindEnv("wapi_version")
	viper.SetDefault("username", "")
	viper.BindEnv("username")
	viper.SetDefault("password", "")
	viper.BindEnv("password")

	// Basic auth exporter
	viper.SetDefault("basic_auth_username", "")
	viper.BindEnv("basic_auth_username")
	viper.SetDefault("basic_auth_password", "")
	viper.BindEnv("basic_auth_password")

	// HTTPCLient
	viper.SetDefault("HTTPClient.timeout", 3)
	viper.BindEnv("HTTPClient.timeout")

	viper.SetDefault("HTTPClient.keepalive", 10)
	viper.BindEnv("HTTPClient.keepalive")

	viper.SetDefault("HTTPClient.tlshandshaketimeout", 10)
	viper.BindEnv("HTTPClient.tlshandshaketimeout")

	viper.SetDefault("HTTPClient.insecureHTTPS", true)
	viper.BindEnv("HTTPClient.insecureHTTPS")
}
