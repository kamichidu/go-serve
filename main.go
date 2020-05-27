package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/comail/colog"
	"github.com/spf13/pflag"
)

const (
	progName = "go-serve"
)

func initLogger(level, outName string) error {
	if lvl, err := colog.ParseLevel(level); err != nil {
		return err
	} else {
		colog.SetMinLevel(lvl)
	}
	// do not init more
	if outName == "-" {
		return nil
	}
	var w io.Writer
	if f, err := os.OpenFile(outName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		w = f
	}
	colog.SetOutput(w)
	return nil
}

func run(stdout, stderr io.Writer, args []string) int {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetOutput(stderr)

	var (
		configPath      string
		listenAddress   string
		autoOpenBrowser bool
		loglevel        string
		logfile         string
	)
	flgs := pflag.NewFlagSet(progName, pflag.ContinueOnError)
	flgs.StringVarP(&configPath, "config", "", DefaultConfigPath, "The config file `path`")
	flgs.StringVarP(&listenAddress, "listen", "l", "", "Listen `address` for tcp (default "+DefaultListenAddress+" or from config file)")
	flgs.BoolVarP(&autoOpenBrowser, "open-browser", "o", false, "Open browser with serving directory")
	flgs.StringVarP(&loglevel, "log-level", "", "info", "The minimul log `level` to output (choices: trace, debug, info, warn, error, alert)")
	flgs.StringVarP(&logfile, "log-file", "", "-", "The log output file `path`")
	if err := flgs.Parse(args[1:]); errors.Is(err, pflag.ErrHelp) {
		return 0
	} else if err != nil {
		fmt.Fprintln(stderr, err)
		return 128
	}
	if flgs.NArg() > 1 {
		flgs.Usage()
		return 128
	}

	if err := initLogger(loglevel, logfile); err != nil {
		fmt.Fprintf(stderr, "invalid --log-level %q or --log-file %q: %v", loglevel, logfile, err)
		return 128
	}

	var dir string
	if flgs.NArg() > 0 {
		dir = flgs.Arg(0)
	} else {
		dir = "."
	}

	cfg, err := LoadConfigFromFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load config from %q: %v\n", configPath, err)
		return 1
	}
	if listenAddress != "" {
		cfg.ListenAddress = listenAddress
	}

	if err := doServe(cfg, dir, autoOpenBrowser); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args))
}
