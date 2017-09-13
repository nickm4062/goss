package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aelsabbahy/goss"
	"github.com/aelsabbahy/goss/outputs"
	"github.com/urfave/cli"
	//"time"
)

var version string

func main() {
	startTime := time.Now()
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = version
	app.Name = "goss"
	app.Usage = "Quick and Easy server validation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gossfile, g",
			Value:  "./goss.yaml",
			Usage:  "Goss file to read from / write to",
			EnvVar: "GOSS_FILE",
		},
		cli.StringFlag{
			Name:   "vars",
			Usage:  "json/yaml file containing variables for template",
			EnvVar: "GOSS_VARS",
		},
		cli.StringFlag{
			Name:  "package",
			Usage: "Package type to use [rpm, deb, apk, pacman]",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "Validate system",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "rspecish",
					Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
					EnvVar: "GOSS_FMT",
				},
				cli.BoolFlag{
					Name:   "color",
					Usage:  "Force color on",
					EnvVar: "GOSS_COLOR",
				},
				cli.BoolFlag{
					Name:   "no-color",
					Usage:  "Force color off",
					EnvVar: "GOSS_NOCOLOR",
				},
				cli.DurationFlag{
					Name:   "sleep,s",
					Usage:  "Time to sleep between retries, only active when -r is set",
					Value:  1 * time.Second,
					EnvVar: "GOSS_SLEEP",
				},
				cli.DurationFlag{
					Name:   "retry-timeout,r",
					Usage:  "Retry on failure so long as elapsed + sleep time is less than this",
					Value:  0,
					EnvVar: "GOSS_RETRY_TIMEOUT",
				},
				cli.IntFlag{
					Name:   "max-concurrent",
					Usage:  "Max number of tests to run concurrently",
					Value:  50,
					EnvVar: "GOSS_MAX_CONCURRENT",
				},
			},
			Action: func(c *cli.Context) error {
				goss.Validate(c, startTime)
				return nil
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "Serve a health endpoint",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "rspecish",
					Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
					EnvVar: "GOSS_FMT",
				},
				cli.DurationFlag{
					Name:   "cache,c",
					Usage:  "Time to cache the results",
					Value:  5 * time.Second,
					EnvVar: "GOSS_CACHE",
				},
				cli.StringFlag{
					Name:   "listen-addr,l",
					Value:  ":8080",
					Usage:  "Address to listen on [ip]:port",
					EnvVar: "GOSS_LISTEN",
				},
				cli.StringFlag{
					Name:   "endpoint,e",
					Value:  "/healthz",
					Usage:  "Endpoint to expose",
					EnvVar: "GOSS_ENDPOINT",
				},
				cli.IntFlag{
					Name:   "max-concurrent",
					Usage:  "Max number of tests to run concurrently",
					Value:  50,
					EnvVar: "GOSS_MAX_CONCURRENT",
				},
			},
			Action: func(c *cli.Context) error {
				goss.Serve(c)
				return nil
			},
		},
		{
			Name:    "render",
			Aliases: []string{"r"},
			Usage:   "render gossfile after imports",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: fmt.Sprintf("Print debugging info when rendering"),
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Print(goss.RenderJSON(c))
				return nil
			},
		},
		{
			Name:    "autoadd",
			Aliases: []string{"aa"},
			Usage:   "automatically add all matching resource to the test suite",
			Action: func(c *cli.Context) error {
				return goss.AutoAddResources(c.GlobalString("gossfile"), c.Args(), c)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a resource to the test suite",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "exclude-attr",
					Usage: "Exclude the following attributes when adding a new resource",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "package",
					Usage: "add new package",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Packages", c.Args(), c)
					},
				},
				{
					Name:  "file",
					Usage: "add new file",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Files", c.Args(), c)
					},
				},
				{
					Name:  "addr",
					Usage: "add new remote address:port - ex: google.com:80",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 500 * time.Millisecond,
						},
					},
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Addrs", c.Args(), c)
					},
				},
				{
					Name:  "port",
					Usage: "add new listening [protocol]:port - ex: 80 or udp:123",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Ports", c.Args(), c)
					},
				},
				{
					Name:  "service",
					Usage: "add new service",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Services", c.Args(), c)
					},
				},
				{
					Name:  "user",
					Usage: "add new user",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Users", c.Args(), c)
					},
				},
				{
					Name:  "group",
					Usage: "add new group",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Groups", c.Args(), c)
					},
				},
				{
					Name:  "command",
					Usage: "add new command",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 10 * time.Second,
						},
					},
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Commands", c.Args(), c)
					},
				},
				{
					Name:  "dns",
					Usage: "add new dns",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 500 * time.Millisecond,
						},
						cli.StringFlag{
							Name:  "server",
							Usage: "The IP address of a DNS server to query",
						},
					},
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "DNS", c.Args(), c)
					},
				},
				{
					Name:  "process",
					Usage: "add new process name",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Processes", c.Args(), c)
					},
				},
				{
					Name:  "http",
					Usage: "add new http",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name: "insecure, k",
						},
						cli.BoolFlag{
							Name: "no-follow-redirects, r",
						},
						cli.DurationFlag{
							Name:  "timeout",
							Value: 5 * time.Second,
						},
					},
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "HTTPs", c.Args(), c)
					},
				},
				{
					Name:  "goss",
					Usage: "add new goss file, it will be imported from this one",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Gossfiles", c.Args(), c)
					},
				},
				{
					Name:  "kernel-param",
					Usage: "add new goss kernel param",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "KernelParams", c.Args(), c)
					},
				},
				{
					Name:  "mount",
					Usage: "add new mount",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Mounts", c.Args(), c)
					},
				},
				{
					Name:  "interface",
					Usage: "add new interface",
					Action: func(c *cli.Context) error {
						return goss.AddResources(c.GlobalString("gossfile"), "Interfaces", c.Args(), c)
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
