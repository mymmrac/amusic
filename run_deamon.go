package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"github.com/urfave/cli/v3"
)

func runDeamon(ctx context.Context, cmd *cli.Command) error {
	app := fiber.New(fiber.Config{
		AppName: "amusic",
	})

	app.Use(func(fCtx fiber.Ctx) error {
		requestCtx, cancel := context.WithCancel(fCtx.Context())
		defer cancel()

		conn := fCtx.RequestCtx().Conn()
		go func() {
			_, _ = conn.Read(nil)
			cancel()
		}()

		fCtx.SetContext(requestCtx)
		return fCtx.Next()
	})

	dm := newDeamon(cmd.String("token"))

	app.Get("/", func(fCtx fiber.Ctx) error {
		return fCtx.SendString("OK")
	})

	app.Post("/api/v1/play", dm.play)

	addr := cmd.String("host") + ":" + cmd.String("port")
	log.Infof("listening on %s", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	done := make(chan struct{})
	if err = app.Listener(listener{Listener: ln}, fiber.ListenConfig{
		GracefulContext: ctx,
		OnShutdownError: func(err error) {
			log.Errorf("shutdown deamon: %s", err)
			close(done)
		},
		OnShutdownSuccess: func() {
			log.Info("shutting down deamon")
			close(done)
		},
		ShutdownTimeout:       time.Second * 10,
		DisableStartupMessage: true,
	}); err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	<-done

	return nil
}

type deamon struct {
	player Player
	token  string
}

func newDeamon(token string) *deamon {
	return &deamon{
		player: NewAutoPlayer(NewPipeWire(), NewPulseAudio()),
		token:  token,
	}
}

type playRequest struct {
	File   string   `json:"file"`
	Volume *float64 `json:"volume,omitempty"`
	Wait   bool     `json:"wait,omitempty"`
}

func (d *deamon) play(fCtx fiber.Ctx) error {
	if d.token != "" && string(fCtx.Request().Header.Peek(fiber.HeaderAuthorization)) != d.token {
		return fCtx.SendStatus(fiber.StatusUnauthorized)
	}

	request := &playRequest{}
	if err := fCtx.Bind().JSON(request); err != nil {
		log.Warnf("invalid request: %s", err)
		return fCtx.Status(fiber.StatusBadRequest).SendString("invalid request")
	}

	if request.File == "" {
		return fCtx.Status(fiber.StatusBadRequest).SendString("invalid file")
	}

	volume := 0.5
	if request.Volume != nil {
		volume = *request.Volume
	}

	if volume < 0.0 || volume > 1.0 {
		return fCtx.Status(fiber.StatusBadRequest).SendString("invalid volume")
	}

	if !d.player.Supported() {
		log.Errorf("no supported player found")
		return fCtx.Status(fiber.StatusInternalServerError).SendString("no supported player found")
	}

	if request.Wait {
		log.Debugf("playing file: %s", request.File)
		err := d.player.Play(fCtx.Context(), request.File, volume)
		if err != nil {
			log.Errorf("play: %s", err)
			return fCtx.Status(fiber.StatusInternalServerError).SendString("failed to play")
		}
		log.Debug("done playing")
	} else {
		log.Debugf("playing file in background: %s", request.File)
		err := d.player.PlayInBackground(request.File, volume)
		if err != nil {
			log.Errorf("play in background: %s", err)
			return fCtx.Status(fiber.StatusInternalServerError).SendString("failed to play in background")
		}
	}

	return fCtx.SendStatus(fiber.StatusOK)
}
