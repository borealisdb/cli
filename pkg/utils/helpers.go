package utils

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/templates"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

func RenderTemplate(obj any, assetPath string, out string) error {
	data, err := templates.Asset(assetPath)
	if err != nil {
		return fmt.Errorf("could not load asset template: %v", err)
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}

	err = t.Execute(f, obj)
	if err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}

	return nil
}

func WriteAssetToFile(assetPath, filename string) error {
	asset, err := templates.Asset(assetPath)
	if err != nil {
		return fmt.Errorf("could find template asset secrets.yaml: %v", err)
	}
	if err := os.WriteFile(filename, asset, 0700); err != nil {
		cobra.CheckErr(fmt.Sprintf("could write file secrets.yaml into destination folder: %v", err))
	}

	return nil
}