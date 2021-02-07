/*
   Copyright 2021 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nokute78/go-gpt/pkg/gpt"
	"io"
	"os"
)

const version string = "0.0.2"

// Exit status
const (
	ExitOK int = iota
	ExitArgError
	ExitCmdError
)

// CLI has In/Out/Err streams.
type CLI struct {
	OutStream io.Writer
	InStream  io.Reader
	ErrStream io.Writer
	quiet     bool // for testing to suppress output
}

// Run executes real main function.
func (cli *CLI) Run(args []string) (ret int) {
	cnf, err := Configure(args[1:], cli.quiet)
	if err != nil {
		if err == flag.ErrHelp {
			return ExitOK
		}
		fmt.Fprintf(cli.ErrStream, "%s\n", err)
		return ExitArgError
	}

	if cnf.showVersion {
		fmt.Fprintf(cli.OutStream, "Ver: %s\n", version)
		return ExitOK
	}
	if len(cnf.devices) == 0 {
		fmt.Fprintf(cli.ErrStream, "no input\n")
		return ExitArgError
	}

	gpts := []gpt.RGpt{}

	for _, v := range cnf.devices {
		f, err := os.Open(v)
		if err != nil {
			fmt.Fprintf(cli.ErrStream, "os.Open err:%s\n", err)
			continue
		}
		defer f.Close()
		g, err := gpt.ReadGpt(f)
		if err != nil {
			fmt.Fprintf(cli.ErrStream, "ReadGpt err:%s\n", err)
			continue
		}
		jg := gpt.NewRGpt(*g)
		gpts = append(gpts, *jg)
	}

	enc := json.NewEncoder(cli.OutStream)
	err = enc.Encode(gpts)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "Encode err:%s\n", err)
		return ExitCmdError
	}

	return ExitOK
}

func main() {
	cli := &CLI{OutStream: os.Stdout, InStream: os.Stdin, ErrStream: os.Stderr}

	os.Exit(cli.Run(os.Args))
}
