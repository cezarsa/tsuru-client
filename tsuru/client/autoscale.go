// Copyright 2020 tsuru-client authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tsuru/gnuflag"
	"github.com/tsuru/go-tsuruclient/pkg/client"
	"github.com/tsuru/go-tsuruclient/pkg/tsuru"
	"github.com/tsuru/tsuru/cmd"
)

type int32Value int32

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}
func (i *int32Value) Get() interface{} { return int32(*i) }
func (i *int32Value) String() string   { return fmt.Sprintf("%v", *i) }

type AutoScaleSet struct {
	cmd.GuessingCommand
	fs        *gnuflag.FlagSet
	autoscale tsuru.AutoScaleSpec
}

func (c *AutoScaleSet) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "unit-autoscale-set",
		Usage:   "unit-autoscale-set [-a/--app appname] [-p/--process processname] [--cpu targetCPU] [--min minUnits] [--max maxUnits]",
		Desc:    `Sets a unit auto scale configuration.`,
		MinArgs: 0,
		MaxArgs: 0,
	}
}

func (c *AutoScaleSet) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = c.GuessingCommand.Flags()

		c.fs.StringVar(&c.autoscale.Process, "process", "", "Process name")
		c.fs.StringVar(&c.autoscale.Process, "p", "", "Process name")

		c.fs.StringVar(&c.autoscale.AverageCPU, "cpu", "", "Target CPU value in percent of a single core. Example: 50%")

		c.autoscale.MinUnits = 1
		c.fs.Var((*int32Value)(&c.autoscale.MinUnits), "min", "Minimum Units")

		c.fs.Var((*int32Value)(&c.autoscale.MaxUnits), "max", "Maximum Units")
	}
	return c.fs
}

func (c *AutoScaleSet) Run(ctx *cmd.Context, cli *cmd.Client) error {
	apiClient, err := client.ClientFromEnvironment(&tsuru.Configuration{
		HTTPClient: cli.HTTPClient,
	})
	if err != nil {
		return err
	}
	appName, err := c.Guess()
	if err != nil {
		return err
	}

	if strings.HasSuffix(c.autoscale.AverageCPU, "%") {
		rawCPU := strings.TrimSuffix(c.autoscale.AverageCPU, "%")
		var percentCPU int
		percentCPU, err = strconv.Atoi(rawCPU)
		if err == nil {
			c.autoscale.AverageCPU = fmt.Sprintf("%dm", percentCPU*10)
		}
	}

	_, err = apiClient.AppApi.AutoScaleAdd(context.TODO(), appName, c.autoscale)
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Stdout, "Unit auto scale successfully set.")
	return nil
}

type AutoScaleUnset struct {
	cmd.GuessingCommand
	fs      *gnuflag.FlagSet
	process string
}

func (c *AutoScaleUnset) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "unit-autoscale-unset",
		Usage:   "unit-autoscale-unset [-a/--app appname] [-p/--process processname]",
		Desc:    `Unsets a unit auto scale configuration.`,
		MinArgs: 0,
		MaxArgs: 0,
	}
}

func (c *AutoScaleUnset) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = c.GuessingCommand.Flags()
		c.fs.StringVar(&c.process, "process", "", "Process name")
		c.fs.StringVar(&c.process, "p", "", "Process name")
	}
	return c.fs
}

func (c *AutoScaleUnset) Run(ctx *cmd.Context, cli *cmd.Client) error {
	apiClient, err := client.ClientFromEnvironment(&tsuru.Configuration{
		HTTPClient: cli.HTTPClient,
	})
	if err != nil {
		return err
	}
	appName, err := c.Guess()
	if err != nil {
		return err
	}
	_, err = apiClient.AppApi.AutoScaleRemove(context.TODO(), appName, c.process)
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Stdout, "Unit auto scale successfully unset.")
	return nil
}
