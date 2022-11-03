package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/codegen/exp"
	"github.com/deepmap/oapi-codegen/pkg/util"
)

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var (
	// flagConfigFile specifies a path to an oapi-codegen configuration file
	flagConfigFile string
	// flagLogFile specifies the path to a file which will contain debug logs.
	// An empty path means log messages will be discarded.
	flagLogFile string
)

func main() {
	flag.StringVar(&flagConfigFile, "config", "", "A YAML config file that controls oapi-codegen behavior.")
	flag.StringVar(&flagLogFile, "log-file", "", "A path where log messages will be written")

	flag.Parse()

	if flag.NArg() < 1 {
		errExit("Please specify a path to a OpenAPI 3.0 spec file\n")
	} else if flag.NArg() > 1 {
		errExit("Only one OpenAPI 3.0 spec file is accepted and it must be the last CLI argument\n")
	}

	// Read the config file
	var config exp.Configuration
	if flagConfigFile != "" {
		buf, err := os.ReadFile(flagConfigFile)
		if err != nil {
			errExit("error reading config file '%s': %v\n", flagConfigFile, err)
		}
		config, err = exp.LoadConfiguration(buf)
		if err != nil {
			errExit("error loading configuration file: %v", err)
		}
	}

	spec, err := util.LoadSwagger(flag.Arg(0))
	if err != nil {
		errExit("error loading OpenAPI spec in %s\n: %s", flag.Arg(0), err)
	}

	var logWriter io.Writer
	if flagLogFile != "" {
		logFile, err := os.OpenFile(flagLogFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			errExit("error opening log file '%s': %v", logWriter, err)
		}
		logWriter = logFile
		defer logFile.Close()
	} else {
		logWriter = io.Discard
	}

	internalLogger := log.New(logWriter, "", log.Lshortfile)

	var outputBuffer bytes.Buffer
	err = exp.Generate(spec, config, &outputBuffer, internalLogger)
	if err != nil {
		errExit("error generating code: %v\n", err)
	}
	fmt.Println(outputBuffer.String())
}
