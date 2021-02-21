# Tube2R
GO web application that converts YouTube channel or public playlist to mp3 RSS feed.

Based on project [PodSync](https://github.com/mxpv/podsync)

## Run as binary:
```
$ ./podsync --config config.toml
```

### Query string options
- **`playlist:`** To get playlist as RSS feed use this query string value. Like YouTube
  - https://www.youtube.com/playlist?list=[PlayListValue]
- **`channelId:`** To get channel as RSS feed use this query string value.
  - https://studio.youtube.com/channel/[ChannelId]

[MIT License](https://opensource.org/licenses/MIT)
