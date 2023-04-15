package config

import (
	"errors"
	"fmt"
	"github.com/borealisdb/cli/pkg/templates"
	"github.com/sirupsen/logrus"
	"os"
	"text/template"
)

type ClITemplate struct {
	ProjectFolder string
	Host          string
	RootUrlPath   string
}

type Credentials struct {
}

type CliConfig struct {
	cliConfigPath string
	Log           *logrus.Entry
	Params        Params
}

func New(params Params, log *logrus.Entry) *CliConfig {
	return &CliConfig{
		Log:    log,
		Params: params,
	}
}

func (k *CliConfig) GenerateDefaultConfigFiles() error {
	if err := os.Mkdir(CliConfigDefaultPath, 0700); err != nil {
		if os.IsExist(err) {
			k.Log.Warningf("folder %v already exists, skipping creation", CliConfigDefaultPath)
			return nil
		}
	}
	if err := k.generateDefaultCliConfig(); err != nil {
		return err
	}
	if err := k.generateDefaultCliCredentials(); err != nil {
		return err
	}

	return nil
}

func (k *CliConfig) generateDefaultCliConfig() error {
	templateConfig := ClITemplate{
		ProjectFolder: k.Params.ProjectFolder,
		RootUrlPath:   k.Params.RootPathUrl,
		Host:          k.Params.Host,
	}
	templateFileSource := fmt.Sprintf("%v/%v", TemplatesPath, CliTemplateFilename)
	fileDestination := fmt.Sprintf("%v/%v", CliConfigDefaultPath, CliConfigFilename)
	if err := LoadParseAndSaveTemplate(templateFileSource, fileDestination, templateConfig); err != nil {
		return err
	}

	return nil
}

func (k *CliConfig) generateDefaultCliCredentials() error {
	templateFileSource := fmt.Sprintf("%v/%v", TemplatesPath, CliCredentialsTemplateFilename)
	fileDestination := fmt.Sprintf("%v/%v", CliConfigDefaultPath, CliCredentialsFilename)
	// Additional check, we don't want to risk to overwrite credentials
	if _, err := os.Stat(fileDestination); errors.Is(err, os.ErrNotExist) {
		templateConfig := Credentials{}
		if err := LoadParseAndSaveTemplate(templateFileSource, fileDestination, templateConfig); err != nil {
			return err
		}
	} else {
		k.Log.Warningf("Credential file %v, was not written as it already exists", CliCredentialsFilename)
	}

	return nil
}

func LoadParseAndSaveTemplate(templateFileSource, templateFileDest string, config any) error {
	data, err := templates.Asset(templateFileSource)
	if err != nil {
		return fmt.Errorf("could not load asset template: %v", err)
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	f, err := os.Create(templateFileDest)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}

	err = t.Execute(f, config)
	if err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}

	return nil
}
