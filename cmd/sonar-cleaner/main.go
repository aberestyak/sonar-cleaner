package main

import (
	"os"
	"time"

	sonarapi "github.com/aberestyak/sonar-cleaner/internal/sonarapi"
	sonarProject "github.com/aberestyak/sonar-cleaner/internal/sonarproject"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var Version = "local"

func main() {
	app := &cli.App{
		Name:     "sonar-cleaner",
		Version:  Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name: "Aleskandr Berestyak",
			},
		},
		HelpName:               "sonar-cleaner",
		Usage:                  "Tool to delete branches and analysis in sonarqube",
		UsageText:              "sonar-cleaner [global options]",
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "log-level",
				Value:    "info",
				Usage:    "Log level",
				EnvVars:  []string{"SONAR_CLEANER_LOG_LEVEL"},
				Required: false,
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Show projects to be cleaned up",
				EnvVars: []string{"SONAR_CLEANER_DRY_RUN"},
			},
			&cli.StringFlag{
				Name:     "address",
				Required: true,
				Usage:    "Sonarqube address",
				Aliases:  []string{"a"},
				EnvVars:  []string{"SONARQUBE_ADDRESS"},
			},
			&cli.StringFlag{
				Name:     "token",
				Required: true,
				Usage:    "Sonarqube token",
				Aliases:  []string{"t"},
				EnvVars:  []string{"SONARQUBE_TOKEN"},
			},
			&cli.IntFlag{
				Name:     "days",
				Required: true,
				Usage:    "Limit of days since last analysis to keep project's branches and analysis",
				Aliases:  []string{"d"},
				EnvVars:  []string{"KEEP_DAYS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			logLevel, err := log.ParseLevel(ctx.String("log-level"))
			if err != nil {
				log.Fatalf("Incorrect log-level")
			}
			log.SetLevel(logLevel)

			config := sonarapi.QueryConfig{
				SonarqubeAddress: ctx.String("address"),
				// curl -u THIS_IS_MY_TOKEN: https://sonarqube.com/api/user_tokens/search
				SonarqubeToken: ctx.String("token") + ":",
				KeepDays:       ctx.Int("days"),
			}
			if ctx.Bool("dry-run") {
				return sonarProject.ShowProjects(config)
			}
			return sonarProject.CleanProjects(config)
		},
	}

	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		ShowFullLevel:   true,
		TimestampFormat: ">",
	})

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Couldn't execute %v:\n%+v", os.Args[0], err)
	}
}
