package main

import (
	"errors"
	"fmt"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/lycoris0731/evans/env"
	"github.com/lycoris0731/evans/lib/parser"
	"github.com/lycoris0731/evans/repl"

	"io"
	"os"
)

type Meta struct {
	Title, Version string
}

type UI struct {
	Reader            io.Reader
	Writer, ErrWriter io.Writer
}

func NewUI() *UI {
	return &UI{
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
	}
}

type Options struct {
	Proto []string `arg:"positional,help:.proto files"`

	Port        int    `arg:"-p,help:gRPC port"`
	Interactive bool   `arg:"-i,help:use interactive mode"`
	Package     string `arg:"help:default package"`
	Service     string `arg:"help:default service. evans parse package from this if --package is nothing."`
}

type CLI struct {
	meta    *Meta
	ui      *UI
	options *Options
}

func NewCLI(title, version string) *CLI {
	return &CLI{
		meta: &Meta{
			Title:   title,
			Version: version,
		},
		ui: NewUI(),
		options: &Options{
			Port: 50051,
		},
	}
}

func (c *CLI) Error(err error) {
	fmt.Fprintln(c.ui.ErrWriter, err)
}

func (c *CLI) Run(args []string) int {
	arg.MustParse(c.options)

	desc, err := parser.ParseFile(c.options.Proto, []string{})
	if err != nil {
		c.Error(err)
		return 1
	}

	config := &repl.Config{
		Port: c.options.Port,
	}
	env := env.NewEnv(desc)
	if c.options.Package != "" {
		if err := env.UsePackage(c.options.Package); err != nil {
			c.Error(err)
			return 1
		}

		if c.options.Service != "" {
			if err := env.UseService(c.options.Service); err != nil {
				c.Error(err)
				return 1
			}
		}
	} else if c.options.Service != "" {
		s := strings.SplitN(c.options.Service, ".", 2)
		if len(s) != 2 {
			c.Error(errors.New("please set package (package_name.service_name or set --package flag)"))
			return 1
		}
		if err := env.UsePackage(s[0]); err != nil {
			c.Error(err)
			return 1
		}
		if err := env.UseService(s[1]); err != nil {
			c.Error(err)
			return 1
		}
	}

	if err := repl.NewREPL(config, env, repl.NewUI()).Start(); err != nil {
		c.Error(err)
		return 1
	}

	return 0
}