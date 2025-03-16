package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(false)
	log.SetTimeFormat("[02.01.2006 15:04:05]")
	styles := log.DefaultStyles()
	styles.Levels = map[log.Level]lipgloss.Style{
		log.DebugLevel: styles.Levels[log.DebugLevel].UnsetMaxWidth(),
		log.InfoLevel:  styles.Levels[log.InfoLevel].UnsetMaxWidth(),
		log.WarnLevel:  styles.Levels[log.WarnLevel].UnsetMaxWidth(),
		log.ErrorLevel: styles.Levels[log.ErrorLevel].UnsetMaxWidth(),
		log.FatalLevel: styles.Levels[log.FatalLevel].UnsetMaxWidth(),
	}
	log.SetStyles(styles)
}

func main() {
	cmd := &cli.Command{
		Name:                  "amusic",
		Usage:                 "Simple tool for playing music (sound) files from CLI or API.",
		Version:               version(),
		Authors:               []any{"Artem Yadelskyi (mymmrac)"},
		EnableShellCompletion: true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose logging.",
				Sources: cli.EnvVars("AMUSIC_VERBOSE"),
				Value:   false,
				Action: func(_ context.Context, _ *cli.Command, verbose bool) error {
					if verbose {
						log.SetLevel(log.DebugLevel)
						log.SetReportCaller(true)
						log.SetReportTimestamp(true)

						log.Debug("verbose logging enabled")
					}
					return nil
				},
			},
		},

		Commands: []*cli.Command{
			{
				Name:      "play",
				Aliases:   []string{"p"},
				Usage:     "Play a music (sound) file.",
				ArgsUsage: "[file]",
				Flags: []cli.Flag{
					&cli.FloatFlag{
						Name:  "volume",
						Usage: "Volume level (0.0-1.0).",
						Value: 0.5,
						Validator: func(volume float64) error {
							if volume < 0.0 || volume > 1.0 {
								return fmt.Errorf("volume must be between 0.0 and 1.0")
							}
							return nil
						},
					},
				},
				Action: runPlay,
			},
			{
				Name:    "deamon",
				Aliases: []string{"d"},
				Usage:   "Run a music (sound) deamon.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host to listen on.",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "port",
						Usage: "Port to listen on.",
						Value: "2567",
					},
					&cli.StringFlag{
						Name:  "token",
						Usage: "Token to use for authentication.",
						Value: "",
					},
				},
				Action: runDeamon,
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}
