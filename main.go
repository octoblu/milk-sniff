package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/milk-sniff/sniffer"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("milk-sniff:main")

func main() {
	app := cli.NewApp()
	app.Name = "milk-sniff"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "iterations, i",
			EnvVar: "MILK_SNIFF_ITERATIONS",
			Usage:  "Redis key pattern to inspect",
			Value:  10,
		},
		cli.StringFlag{
			Name:   "redis-uri, r",
			EnvVar: "MILK_SNIFF_REDIS_URI",
			Usage:  "Redis URI to sniff test",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	iterations, redisURI := getOpts(context)
	nose := sniffer.New(redisURI)

	for i := 1; i <= iterations; i++ {
		result, err := nose.Sniff()
		if err != nil {
			log.Fatalln("nose.Sniff errored:", err.Error())
		}
		fmt.Println(result)
	}
}

func getOpts(context *cli.Context) (int, string) {
	iterations := context.Int("iterations")
	redisURI := context.String("redis-uri")

	if redisURI == "" {
		cli.ShowAppHelp(context)

		if redisURI == "" {
			color.Red("  Missing required flag --redis-uri or MILK_SNIFF_REDIS_URI")
		}
		os.Exit(1)
	}

	return iterations, redisURI
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
