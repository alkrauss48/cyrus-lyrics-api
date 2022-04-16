package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleConfig struct {
	Web SubGoogleConfig `json:"web"`
}

type SubGoogleConfig struct {
	ClientId                string    `json:"client_id"`
	ProjectId               string    `json:"project_id"`
	AuthUri                 string    `json:"auth_uri"`
	TokenUri                string    `json:"token_uri"`
	AuthProviderX509CertUrl string    `json:"auth_provider_x509_cert_url"`
	ClientSecret            string    `json:"client_secret"`
	RedirectUris            [1]string `json:"redirect_uris"`
}

func MarshalOAuthCredentials() ([]byte, error) {
	credentials := GoogleConfig{
		Web: SubGoogleConfig{
			AuthProviderX509CertUrl: os.Getenv("AUTH_PROVIDER_CERT_URL"),
			AuthUri:                 os.Getenv("AUTH_URI"),
			ClientId:                os.Getenv("CLIENT_ID"),
			ClientSecret:            os.Getenv("CLIENT_SECRET"),
			ProjectId:               os.Getenv("PROJECT_ID"),
			RedirectUris:            [1]string{os.Getenv("REDIRECT_URI")},
			TokenUri:                os.Getenv("TOKEN_URI"),
		},
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	return jsonCredentials, nil
}

func GetGoogleOAuthConfig() (*oauth2.Config, error) {
	jsonCredentials, err := MarshalOAuthCredentials()
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(jsonCredentials, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	return config, nil
}
