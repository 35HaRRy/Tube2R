package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
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

	route := gin.Default()

	// Create routes
	group.Go(func() error {
		// route.Static("/files", cfg.Server.DataDir)
		route.GET("/files", func(c *gin.Context) {
			// Download file
			log.Debug("creating updater")
			updater, err := NewUpdater(cfg, downloader)

			pathPrefix := ""

			id := c.Request.URL.Query().Get("channelId")
			if id != "" {
				pathPrefix = "channel/"
			} else {
				id = c.Request.URL.Query().Get("playlist")
				pathPrefix = "playlist?list="
			}

			url := []string{"https://www.youtube.com/", pathPrefix, id}

			feed := config.Feed{
				ID: id,
				// URL:     "https://www.youtube.com/playlist?list=PLklI4fp4DoMj4Vz4W8Q7d8gycP_a9hPaE",
				URL:     strings.Join(url, ""),
				Quality: "high",
			}

			logger := log.WithFields(log.Fields{
				"episode_id": episode.ID,
			})

			feedPath := filepath.Join(updater.config.Server.DataDir, feed.ID)
			episodePath := filepath.Join(feedPath, updater.episodeName(&feed, episode))

			_, err := os.Stat(episodePath)
			if err != nil && !os.IsNotExist(err) {
				return errors.Wrap(err, "failed to check whether episode exists"), nil
			}

			if os.IsNotExist(err) {
				// There is no file on disk, download episode
				logger.Infof("! downloading episode %s", episode.VideoURL)
				updater.downloader.Download(ctx, &feed, episode, feedPath)
			} else {
				// Episode already downloaded
				logger.Debug("skipping download of episode")
			}

			// Record file size
			// if size, err := u.fileSize(episodePath); err != nil {
			// 	// Don't return on error, use estimated file size provided by builders
			// 	logger.WithError(err).Error("failed to get episode file size")
			// } else { //nolint
			// 	logger.Debugf("file size %d", size)
			// 	sizes[episode.ID] = size
			// }
		})
		route.GET("/rss", func(c *gin.Context) {
			// Run updater thread
			log.Debug("creating updater")
			updater, err := NewUpdater(cfg, downloader)

			if err != nil {
				log.WithError(err).Fatal("failed to create updater")
			}

			pathPrefix := ""

			id := c.Request.URL.Query().Get("channelId")
			if id != "" {
				pathPrefix = "channel/"
			} else {
				id = c.Request.URL.Query().Get("playlist")
				pathPrefix = "playlist?list="
			}

			url := []string{"https://www.youtube.com/", pathPrefix, id}

			feed := config.Feed{
				ID: id,
				// URL:     "https://www.youtube.com/playlist?list=PLklI4fp4DoMj4Vz4W8Q7d8gycP_a9hPaE",
				URL:     strings.Join(url, ""),
				Quality: "high",
				Format:  "audio",
			}

			err, podcast := updater.Update(ctx, &feed)
			if err != nil {
				log.WithError(err).Errorf("failed to update feed: %s", feed.URL)
			}

			c.String(200, podcast.String())
			// c.XML(200, podcast)
			// c.JSON(200, gin.H{
			// 	"message": "pong",
			// })
		})

		return route.Run()
	})

	// Terminate app
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-stop:
				cancel()
				return nil
			}
		}
	})

	if err := group.Wait(); err != nil && err != context.Canceled {
		log.WithError(err).Error("wait error")
	}

	log.Info("gracefully stopped")
}
