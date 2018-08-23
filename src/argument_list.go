package main

import (
	"errors"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Username               string `default:"" help:"The Microsoft SQL Server connection user name"`
	Password               string `default:"" help:"The Microsoft SQL Server connection password"`
	Instance               string `default:"" help:"The Microsoft SQL Server instance to connect to"`
	Hostname               string `default:"127.0.0.1" help:"The Microsoft SQL Server connection host name"`
	Port                   string `default:"1443" help:"The Microsoft SQL Server port to connect to. Only needed when instance not specified"`
	EnableSSL              bool   `default:"false" help:"If true will use SSL encryption, false will not use encryption"`
	TrustServerCertificate bool   `default:"false" help:"If true server certificate is not verified for SSL. If false certificate will be verified against supplied certificate"`
	CertificateLocation    string `default:"" help:"Certificate file to verify SSL encryption against"`
	Timeout                string `default:"30" help:"Timeout in seconds for a single SQL Query. Set 0 for no timeout"`
}

// Validate validates SQL specific arguments
func (al argumentList) Validate() error {
	if al.Username == "" {
		return errors.New("invalid configuration: must specify a username")
	}

	if al.Hostname == "" {
		return errors.New("invalid configuration: must specify a hostname")
	}

	if (al.Port != "" && al.Instance != "") || (al.Port == "" && al.Instance == "") {
		return errors.New("invalid configuration: must specify one of port or instance")
	}

	if al.EnableSSL && (!al.TrustServerCertificate && al.CertificateLocation == "") {
		return errors.New("invalid configuration: must specify a certificate file when using SSL and not trusting server certificate")
	}

	return nil
}
