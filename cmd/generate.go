package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MUlt1mate/geny/app"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Run generation commands from yaml or go:generate",
	RunE:  generate,
}

func generate(cmd *cobra.Command, _ []string) (err error) {
	var batch *app.CommandBatch
	batch, err = loadBatch(runFlags.File)
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true
	if err = geny.Exec(batch); err != nil {
		return err
	}
	return nil
}

func loadBatch(path string) (batch *app.CommandBatch, err error) {
	if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
		return loadBatchFromYaml(path)
	}
	if strings.HasSuffix(path, ".go") {
		return loadBatchFromGoFile(path)
	}
	return nil, nil
}

func loadBatchFromYaml(path string) (batch *app.CommandBatch, err error) {
	var yamlContent []byte
	if yamlContent, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	if batch, err = geny.ParseYAML(yamlContent); err != nil {
		return nil, err
	}
	return batch, nil
}

func loadBatchFromGoFile(path string) (batch *app.CommandBatch, err error) {
	var goContent []byte
	if goContent, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	if batch, err = geny.ParseGoFile(string(goContent)); err != nil {
		return nil, err
	}
	return batch, nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&runFlags.File, "file", "f", "geny.yaml", "yaml or go file with commands")
}
