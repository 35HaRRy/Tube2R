package feed

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mxpv/podsync/pkg/config"
	"github.com/mxpv/podsync/pkg/link"
)

var (
	testCtx = context.Background()
	ytKey   = os.Getenv("YOUTUBE_TEST_API_KEY")
)

func TestYT_QueryChannel(t *testing.T) {
	if ytKey == "" {
		t.Skip("YouTube API key is not provided")
	}

	builder, err := NewYouTubeBuilder(ytKey)
	require.NoError(t, err)

	channel, err := builder.listChannels(testCtx, link.TypeChannel, "UC2yTVSttx7lxAOAzx1opjoA", "id")
	require.NoError(t, err)
	require.Equal(t, "UC2yTVSttx7lxAOAzx1opjoA", channel.Id)

	channel, err = builder.listChannels(testCtx, link.TypeUser, "fxigr1", "id")
	require.NoError(t, err)
	require.Equal(t, "UCr_fwF-n-2_olTYd-m3n32g", channel.Id)
}

func TestYT_BuildFeed(t *testing.T) {
	if ytKey == "" {
		t.Skip("YouTube API key is not provided")
	}

	builder, err := NewYouTubeBuilder(ytKey)
	require.NoError(t, err)

	urls := []string{
		"https://youtube.com/user/fxigr1",
		"https://www.youtube.com/channel/UCupvZG-5ko_eiXAupbDfxWw",
		"https://www.youtube.com/playlist?list=PLfVk3KMh3VX1yJShGRsJmsqAjvMIviJYQ",
		"https://www.youtube.com/channel/UCK9lZ2lHRBgx2LOcqPifukA",
		"https://youtube.com/user/WylsaLive",
		"https://www.youtube.com/playlist?list=PLUVl5pafUrBydT_gsCjRGeCy0hFHloec8",
	}

	for _, addr := range urls {
		t.Run(addr, func(t *testing.T) {
			feed, err := builder.Build(testCtx, &config.Feed{URL: addr})
			require.NoError(t, err)

			assert.NotEmpty(t, feed.Title)
			assert.NotEmpty(t, feed.Description)
			assert.NotEmpty(t, feed.Author)
			assert.NotEmpty(t, feed.ItemURL)

			assert.NotZero(t, len(feed.Episodes))

			for _, item := range feed.Episodes {
				assert.NotEmpty(t, item.Title)
				assert.NotEmpty(t, item.VideoURL)
				assert.NotZero(t, item.Duration)

				assert.NotEmpty(t, item.Title)
				assert.NotEmpty(t, item.Thumbnail)
			}
		})
	}
}

func TestYT_GetVideoCount(t *testing.T) {
	if ytKey == "" {
		t.Skip("YouTube API key is not provided")
	}

	builder, err := NewYouTubeBuilder(ytKey)
	require.NoError(t, err)

	feeds := []*link.Info{
		{Provider: link.ProviderYoutube, LinkType: link.TypeUser, ItemID: "fxigr1"},
		{Provider: link.ProviderYoutube, LinkType: link.TypeChannel, ItemID: "UCupvZG-5ko_eiXAupbDfxWw"},
		{Provider: link.ProviderYoutube, LinkType: link.TypePlaylist, ItemID: "PLfVk3KMh3VX1yJShGRsJmsqAjvMIviJYQ"},
		{Provider: link.ProviderYoutube, LinkType: link.TypeChannel, ItemID: "UCK9lZ2lHRBgx2LOcqPifukA"},
		{Provider: link.ProviderYoutube, LinkType: link.TypeUser, ItemID: "WylsaLive"},
		{Provider: link.ProviderYoutube, LinkType: link.TypePlaylist, ItemID: "PLUVl5pafUrBydT_gsCjRGeCy0hFHloec8"},
	}

	for _, f := range feeds {
		feed := f
		t.Run(f.ItemID, func(t *testing.T) {
			count, err := builder.GetVideoCount(testCtx, feed)
			assert.NoError(t, err)
			assert.NotZero(t, count)
		})
	}
}
