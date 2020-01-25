package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/mxpv/podsync/pkg/config"
	"github.com/mxpv/podsync/pkg/ytdl"

	"github.com/gin-gonic/gin"
)

type Opts struct {
	ConfigPath string `long:"config" short:"c" default:"config.toml" env:"PODSYNC_CONFIG_PATH"`
	Debug      bool   `long:"debug"`
	NoBanner   bool   `long:"no-banner"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	// Parse args
	opts := Opts{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.WithError(err).Fatal("failed to parse command line arguments")
	}

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}

	downloader, err := ytdl.New(ctx)
	if err != nil {
		log.WithError(err).Fatal("youtube-dl error")
	}

	// Load TOML file
	log.Debugf("loading configuration %q", opts.ConfigPath)
	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		log.WithError(err).Fatal("failed to load configuration file")
	}

	r := gin.Default()

	r.GET("/rss", func(c *gin.Context) {
		// Run updater thread
		log.Debug("creating updater")
		updater, err := NewUpdater(cfg, downloader)

		if err != nil {
			log.WithError(err).Fatal("failed to create updater")
		}

		feed := config.Feed{
			ID:      "ID1",
			URL:     "https://www.youtube.com/playlist?list=PLklI4fp4DoMj4Vz4W8Q7d8gycP_a9hPaE",
			Quality: "high",
			Format:  "audio",
		}

		if err := updater.Update(ctx, feed); err != nil {
			log.WithError(err).Errorf("failed to update feed: %s", feed.URL)
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
