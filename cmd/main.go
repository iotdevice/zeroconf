//go:generate echo "generate"
package main

import (
	"context"
	"fmt"
	"github.com/iotdevice/zeroconf"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	var service = "_http._tcp"
	var domain = "local"
	var waitTime = 5
	myApp := cli.NewApp()
	myApp.Name = "mdns"
	myApp.Version = buildVersion(version, commit, date, builtBy)
	myApp.Commands = []*cli.Command{
		{
			Name:    "types",
			Aliases: []string{"t"},
			Usage:   "find all mdns types",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "service",
					Aliases:     []string{"s"},
					Value:       service,
					Usage:       "service name",
					EnvVars:     []string{"SERVICE"},
					Destination: &service,
				},
				&cli.StringFlag{
					Name:        "domain",
					Aliases:     []string{"d"},
					Value:       "_services._dns-sd._udp",
					Usage:       "domain",
					EnvVars:     []string{"DOMAIN"},
					Destination: &domain,
				},
				&cli.IntFlag{
					Name:        "waitTime",
					Aliases:     []string{"w"},
					Value:       waitTime,
					Usage:       "wait time(second)",
					EnvVars:     []string{"WAITTIME"},
					Destination: &waitTime,
				},
			},
			Action: func(c *cli.Context) error {
				// Discover all services on the network (e.g. _workstation._tcp)
				resolver, err := zeroconf.NewResolver(nil)
				if err != nil {
					log.Println("Failed to initialize resolver:", err.Error())
					return err
				}

				entries := make(chan *zeroconf.ServiceEntry)
				go func(results <-chan *zeroconf.ServiceEntry) {
					for entry := range results {
						log.Println(entry)
					}
					log.Println("No more entries.")
				}(entries)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
				defer cancel()
				err = resolver.Browse(ctx, service, domain, entries)
				if err != nil {
					log.Println("Failed to browse:", err.Error())
					return err
				}

				<-ctx.Done()
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test this command",
			Action: func(c *cli.Context) error {
				fmt.Println("ok")
				return nil
			},
		},
	}
	myApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "service",
			Aliases:     []string{"s"},
			Value:       service,
			Usage:       "service name",
			EnvVars:     []string{"SERVICE"},
			Destination: &service,
		},
		&cli.StringFlag{
			Name:        "domain",
			Aliases:     []string{"d"},
			Value:       domain,
			Usage:       "domain",
			EnvVars:     []string{"DOMAIN"},
			Destination: &domain,
		},
		&cli.IntFlag{
			Name:        "waitTime",
			Aliases:     []string{"w"},
			Value:       waitTime,
			Usage:       "wait time(second)",
			EnvVars:     []string{"WAITTIME"},
			Destination: &waitTime,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		// Discover all services on the network (e.g. _workstation._tcp)
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			log.Println("Failed to initialize resolver:", err.Error())
			return err
		}

		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				log.Println(entry)
			}
			log.Println("No more entries.")
		}(entries)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
		defer cancel()
		err = resolver.Browse(ctx, service, domain, entries)
		if err != nil {
			log.Println("Failed to browse:", err.Error())
			return err
		}

		<-ctx.Done()
		return nil
	}
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
