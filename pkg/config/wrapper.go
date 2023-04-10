package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

type Wrapper struct {
	credentials         *ini.File
	CredentialsLocation string
}

func (w *Wrapper) LoadCredentials() error {
	credentialsConfig, err := ini.Load(w.CredentialsLocation)
	if err != nil {
		return fmt.Errorf("could not load credentials config file: %v ", err)
	}
	w.credentials = credentialsConfig
	return nil
}

func (w *Wrapper) SetToken(token string) error {
	w.credentials.Section(CredentialsAuthSection).Key(CredentialsTokenKey).SetValue(token)
	if err := w.credentials.SaveTo(w.CredentialsLocation); err != nil {
		return fmt.Errorf("could not save token to %v: %v", w.CredentialsLocation, err)
	}

	return nil
}

func (w *Wrapper) GetToken() (string, error) {
	key, err := w.credentials.Section(CredentialsAuthSection).GetKey(CredentialsTokenKey)
	if err != nil {
		return "", fmt.Errorf("could not get token from auth section: %v", err)
	}

	return key.Value(), nil
}

func (w *Wrapper) GetSSOStartURL() string {
	return viper.GetStringMap("sso")[SsoStartUrlKey].(string)
}

func (w *Wrapper) GetSSOLoginURL() string {
	return fmt.Sprintf("%v/sign_in", w.GetSSOStartURL())
}

func (w *Wrapper) GetHost() string {
	return viper.GetStringMap("general")[HostKey].(string)
}

func (w *Wrapper) GetBaseUrl() string {
	host := viper.GetStringMap("general")[HostKey].(string)
	rootUrlPath := viper.GetStringMap("general")[RootUrlPathKey].(string)
	return fmt.Sprintf("%v%v", host, rootUrlPath)
}

func (w *Wrapper) GetCALocation() string {
	return viper.GetStringMap("tls")["ca_cert_location"].(string)
}
