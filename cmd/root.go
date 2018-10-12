package cmd

import (
	"fmt"
	"os"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	input "github.com/tcnksm/go-input"
)

var rootCmd = &cobra.Command{
	Use:   "aws-cli-oidc",
	Short: "CLI tool for retrieving AWS temporary credentials using OIDC provider",
	Long:  `CLI tool for retrieving AWS temporary credentials using OIDC provider`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Writeln(err.Error())
	}
}

var configdir string

const OIDC_PROVIDER_METADATA_URL = "oidc_provider_metadata_url"
const SUCCESSFUL_REDIRECT_URL = "successful_redirect_url"
const FAILURE_REDIRECT_URL = "failure_redirect_url"
const CLIENT_ID = "client_id"
const CLIENT_SECRET = "client_secret"
const AWS_FEDERATION_TYPE = "aws_federation_type"

// OIDC config
const AWS_FEDERATION_ROLE = "aws_federation_role"
const AWS_FEDERATION_ROLE_SESSION_NAME = "aws_federation_role_session_name"

// SAML config
const OIDC_PROVIDER_TOKEN_EXCHANGE_AUDIENCE = "oidc_provider_token_exchange_audience"
const OIDC_PROVIDER_TOKEN_EXCHANGE_SUBJECT_TOKEN_TYPE = "oidc_provider_token_exchange_subject_token_type" // Only support saml2

// OAuth 2.0 Token Exchange
const TOKEN_TYPE_ACCESS_TOKEN = "urn:ietf:params:oauth:token-type:access_token"
const TOKEN_TYPE_ID_TOKEN = "urn:ietf:params:oauth:token-type:id_token"

// Federation Type
const AWS_FEDERATION_TYPE_OIDC = "oidc"
const AWS_FEDERATION_TYPE_SAML2 = "saml2"

func init() {
	cobra.OnInitialize(initConfig)
}

var ui *input.UI
var isTraceEnabled bool

func initConfig() {
	viper.SetConfigFile(ConfigPath() + "/config.yaml")

	if err := viper.ReadInConfig(); err == nil {
		Writeln("Using config file: %s", viper.ConfigFileUsed())
	}

	ui = &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	isTraceEnabled = false // TODO: configuable
}

func ConfigPath() string {
	if configdir != "" {
		return configdir
	}
	path := os.Getenv("AWS_CLI_OIDC_CONFIG")
	if path == "" {
		home, err := homedir.Dir()
		if err != nil {
			Exit(err)
		}
		path = home + "/.aws-cli-oidc"
	}
	return path
}

func Exit(err error) {
	if err != nil {
		Writeln(err.Error())
	}
	os.Exit(1)
}

func CheckInstalled(name string) (*OIDCClient, error) {
	return InitializeClient(name)
}

func Write(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, msg...)
}

func Writeln(format string, msg ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, msg...))
}

func Export(key string, value string) {
	var msg string
	if runtime.GOOS == "windows" {
		msg = fmt.Sprintf("set %s=%s\n", key, value)
	} else {
		msg = fmt.Sprintf("export %s=%s\n", key, value)
	}
	fmt.Fprint(os.Stdout, msg)
}

func Traceln(format string, msg ...interface{}) {
	if isTraceEnabled {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, msg...))
	}
}