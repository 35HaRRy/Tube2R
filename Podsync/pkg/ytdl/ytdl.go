package ytdl

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxpv/podsync/pkg/config"
	"github.com/mxpv/podsync/pkg/model"
)

const DownloadTimeout = 10 * time.Minute

type YoutubeDl struct{}

type DownloadResult struct {
	Result string
	Error  error
}

func New(ctx context.Context) (*YoutubeDl, error) {
	ytdl := &YoutubeDl{}

	// Make sure youtube-dl exists
	version, err := ytdl.exec(ctx, "--version")
	if err != nil {
		return nil, errors.Wrap(err, "could not find youtube-dl")
	}

	log.Infof("using youtube-dl %s", version)

	// Make sure ffmpeg exists
	output, err := exec.CommandContext(ctx, "ffmpeg", "-version").CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "could not find ffmpeg")
	}

	log.Infof("using ffmpeg %s", output)

	return ytdl, nil
}

func (dl YoutubeDl) Download(ctx context.Context, feedConfig *config.Feed, episode *model.Episode, feedPath string) (string, error) /*<-chan (DownloadResult)*/ {
	// result := make(chan DownloadResult)

	// go func() {
	// defer close(result)

	var (
		outputTemplate = makeOutputTemplate(feedPath, episode)
		url            = episode.VideoURL
	)

	if feedConfig.Format == model.FormatAudio {
		// Audio
		if feedConfig.Quality == model.QualityHigh {
			// High quality audio (encoded to mp3)
			/*result, err :=*/
			return dl.exec(ctx,
				"--extract-audio",
				"--audio-format",
				"mp3",
				"--format",
				"bestaudio",
				"--output",
				outputTemplate,
				url,
			)

			// if err != nil {
			// 	// YouTube might block host with HTTP Error 429: Too Many Requests
			// 	// We still need to generate XML, so just stop sending download requests and retry next time
			// 	if strings.Contains(result, "HTTP Error 429") {
			// 		logger.WithError(err).Warnf("got too many requests error, will retry download next time")
			// 		break
			// 	}

			// 	logger.WithError(err).Errorf("youtube-dl error: %s", result)
			// }
		} else {
			//nolint
			// Low quality audio (encoded to mp3)
			/*result, err :=*/
			return dl.exec(ctx,
				"--extract-audio",
				"--audio-format",
				"mp3",
				"--format",
				"worstaudio",
				"--output",
				outputTemplate,
				url,
			)
		}
	} else {
		/*
			Video
		*/
		if feedConfig.Quality == model.QualityHigh {
			// High quality
			/*result, err :=*/
			return dl.exec(ctx,
				"--format",
				"bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
				"--output",
				outputTemplate,
				url,
			)
		} else { //nolint
			// Low quality
			/*result, err :=*/
			return dl.exec(ctx,
				"--format",
				"worstvideo[ext=mp4]+worstaudio[ext=m4a]/worst[ext=mp4]/worst",
				"--output",
				outputTemplate,
				url,
			)
		}
	}
	// }()

	// return result
}

func (YoutubeDl) execAsync(ctx context.Context, args ...string) DownloadResult {
	ctx, cancel := context.WithTimeout(ctx, DownloadTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "youtube-dl", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return DownloadResult{
			Result: string(output),
			Error:  errors.Wrap(err, "failed to execute youtube-dl"),
		}
	}

	return DownloadResult{
		Result: string(output),
		Error:  nil,
	}
}
func (YoutubeDl) exec(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, DownloadTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "youtube-dl", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), errors.Wrap(err, "failed to execute youtube-dl")
	}

	return string(output), nil
}

func makeOutputTemplate(feedPath string, episode *model.Episode) string {
	filename := fmt.Sprintf("%s.%s", episode.ID, "%(ext)s")
	return filepath.Join(feedPath, filename)
}
