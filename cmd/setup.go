package cmd

import (
	"fmt"
	"os"

	input "github.com/natsukagami/go-input"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup of aws-cli-oidc",
	Long:  `Interactive setup of aws-cli-oidc. Will prompt you for OIDC provider URL and other settings.`,
	Run:   setup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func setup(cmd *cobra.Command, args []string) {
	runSetup()
}

func runSetup() {
	providerName, _ := ui.Ask("OIDC provider name:", &input.Options{
		Required: true,
		Loop:     true,
	})
	server, _ := ui.Ask("OIDC provider metadata URL (https://your-oidc-provider/.well-known/openid-configuration):", &input.Options{
		Required: true,
		Loop:     true,
	})
	additionalQuery, _ := ui.Ask("Additional query for OIDC authentication request (Default: none):", &input.Options{
		Default:  "",
		Required: false,
	})
	successfulRedirectURL, _ := ui.Ask("Successful redirect URL (Default: none):", &input.Options{
		Default:  "",
		Required: false,
	})
	failureRedirectURL, _ := ui.Ask("Failure redirect URL (Default: none):", &input.Options{
		Default:  "",
		Required: false,
	})
	clientID, _ := ui.Ask("Client ID which is registered in the OIDC provider:", &input.Options{
		Required: true,
		Loop:     true,
	})
	clientSecret, _ := ui.Ask("Client secret which is registered in the OIDC provider (Default: none):", &input.Options{
		Default:  "",
		Required: false,
	})
	answerFedType, _ := ui.Ask(fmt.Sprintf("Choose type of AWS federation [%s/%s]:", AWS_FEDERATION_TYPE_OIDC, AWS_FEDERATION_TYPE_SAML2), &input.Options{
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			if s != AWS_FEDERATION_TYPE_SAML2 && s != AWS_FEDERATION_TYPE_OIDC {
				return errors.New(fmt.Sprintf("Input must be '%s' or '%s'", AWS_FEDERATION_TYPE_OIDC, AWS_FEDERATION_TYPE_SAML2))
			}
			return nil
		},
	})

	config := map[string]string{}

	config[OIDC_PROVIDER_METADATA_URL] = server
	config[OIDC_AUTHENTICATION_REQUEST_ADDITIONAL_QUERY] = additionalQuery
	config[SUCCESSFUL_REDIRECT_URL] = successfulRedirectURL
	config[FAILURE_REDIRECT_URL] = failureRedirectURL
	config[CLIENT_ID] = clientID
	config[CLIENT_SECRET] = clientSecret
	config[AWS_FEDERATION_TYPE] = answerFedType

	if answerFedType == AWS_FEDERATION_TYPE_OIDC {
		oidcSetup(config)
	} else if answerFedType == AWS_FEDERATION_TYPE_SAML2 {
		saml2Setup(config)
	}

	viper.Set(providerName, config)

	os.MkdirAll(ConfigPath(), 0700)
	configPath := ConfigPath() + "/config.yaml"
	viper.SetConfigFile(configPath)
	err := viper.WriteConfig()

	if err != nil {
		Writeln("Failed to write %s", configPath)
		Exit(err)
	}

	Writeln("Saved %s", configPath)
}

func oidcSetup(config map[string]string) {
	awsRole, _ := ui.Ask("AWS federation role (arn:aws:iam::<Account ID>:role/<Role Name>):", &input.Options{
		Required: true,
		Loop:     true,
	})
	awsRoleSessionName, _ := ui.Ask("AWS federation roleSessionName:", &input.Options{
		Required: true,
		Loop:     true,
	})
	config[AWS_FEDERATION_ROLE] = awsRole
	config[AWS_FEDERATION_ROLE_SESSION_NAME] = awsRoleSessionName
}

func saml2Setup(config map[string]string) {
	answer, _ := ui.Ask(`Select the subject token type to exchange for SAML2 assertion:
	1. Access Token (urn:ietf:params:oauth:token-type:access_token)
	2. ID Token (urn:ietf:params:oauth:token-type:id_token)
  `, &input.Options{
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			if s != "1" && s != "2" {
				return errors.New("Input must be number")
			}
			return nil
		},
	})
	var subjectTokenType string
	if answer == "1" {
		subjectTokenType = TOKEN_TYPE_ACCESS_TOKEN
	} else if answer == "2" {
		subjectTokenType = TOKEN_TYPE_ID_TOKEN
	}
	config[OIDC_PROVIDER_TOKEN_EXCHANGE_SUBJECT_TOKEN_TYPE] = subjectTokenType

	audience, _ := ui.Ask("Audience for token exchange:", &input.Options{
		Required: true,
		Loop:     true,
	})
	config[OIDC_PROVIDER_TOKEN_EXCHANGE_AUDIENCE] = audience
}
