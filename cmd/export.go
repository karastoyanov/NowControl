package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/karastoyanov/nowcontrol/internal/export"
	"github.com/spf13/cobra"
)

// addExportFlags registers the --format and --output flags shared by
// commands that export tabular data, returning pointers to their values.
func addExportFlags(cmd *cobra.Command) (format, output *string) {
	format = cmd.Flags().String("format", "json", "output format: json, csv, xml, xlsx")
	output = cmd.Flags().String("output", "", "write output to this file instead of stdout (required for xlsx)")
	return format, output
}

// exportRecords validates the requested format/output, then writes records
// to the resolved destination, confirming on stderr when writing to a file.
func exportRecords(formatStr, outputPath string, records []map[string]any) error {
	format, err := export.ParseFormat(formatStr)
	if err != nil {
		return err
	}
	if format == export.FormatXLSX && outputPath == "" {
		return fmt.Errorf("--format xlsx requires --output <file> (binary format can't be written to stdout)")
	}

	w, closeFn, err := openOutput(outputPath)
	if err != nil {
		return err
	}
	defer closeFn()

	if err := export.WriteRecords(w, format, records); err != nil {
		return err
	}

	if outputPath != "" {
		fmt.Fprintf(os.Stderr, "Wrote %d record(s) to %s\n", len(records), outputPath)
	}
	return nil
}

func openOutput(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, func() {}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("creating %s: %w", path, err)
	}
	return f, func() { f.Close() }, nil
}
