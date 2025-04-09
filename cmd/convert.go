package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/MUlt1mate/geny/app"
)

const outputStdout = "stdout"

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert commands from/to yaml format",
	RunE:  convert,
}

func convert(_ *cobra.Command, _ []string) (err error) {
	var batch *app.CommandBatch
	batch, err = loadBatch(runFlags.Input)
	if err != nil {
		return err
	}
	if runFlags.Output == outputStdout {
		var output []byte
		if output, err = yaml.Marshal(batch); err != nil {
			return err
		}
		fmt.Println(string(output))
	}
	return saveBatch(batch, runFlags.Output)
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&runFlags.Input, "input", "i", ".", "yaml or go file with commands or path for go generate command")
	convertCmd.Flags().StringVarP(&runFlags.Output, "output", "o", outputStdout, "yaml or go file to save commands")
	convertCmd.Flags().StringVarP(&runFlags.Package, "package", "p", "main", "package name for generated go file")
	convertCmd.Flags().BoolVarP(&runFlags.HideGeneratedComment, "hideComment", "g", false, "hide Code generated comment in go file")
}

func saveBatch(batch *app.CommandBatch, path string) (err error) {
	var output []byte
	if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
		if output, err = yaml.Marshal(batch); err != nil {
			return err
		}
	}
	if strings.HasSuffix(path, ".go") {
		output = []byte(geny.FormatGoFile(batch, runFlags.Package, !runFlags.HideGeneratedComment))
	}
	return os.WriteFile(path, output, 0644)
}
