# yourcast-streamer

It's youtube restreamer written in Go.

## Prerequisites

1. [Go](https://golang.org/)
2. [Make](https://www.gnu.org/software/make/)
3. [Glide](https://github.com/Masterminds/glide)

## Installation

```
mkdir -p $GOPATH/src/github.com/maddevsio/
cd $GOPATH/src/github.com/maddevsio
git clone https://github.com/maddevsio/yourcast-streamer
cd yourcast-streamer
make depends
make
```

Or golang way

```
mkdir -p $GOPATH/src/github.com/maddevsio/
cd $GOPATH/src/github.com/maddevsio
git clone https://github.com/maddevsio/yourcast-streamer
cd yourcast-streamer
go get -v
go build -v
go install
```
## Configure 

```
GLOBAL OPTIONS:
   --download_limit value   (default: 100) [$RESTREAMER_DOWNLOAD_LIMIT]
   --rtmp_server_url value  (default: "rtmp://localhost/hls") [$RESTREAMER_RTMP_ROOT_SERVER_URL]
   --log_level value        (default: "debug") [$RESTREAMER_LOG_LEVEL]
   --web_ui_url value       (default: "http://localhost:8000") [$RESTREAMER_WEB_UI_URL]
   --ffmpeg_path value      (default: "/usr/local/bin/ffmpeg") [$RESTREAMER_FFMPEG_PATH]
   --http_bind_addr value   (default: ":8080") [$RESTREAMER_HTTP_BIND_ADDR]
   --youtube_api_key value  (default: "Aiza...") [$RESTREAMER_YOUTUBE_API_KEY]
   --root_path value        (default: "./storage") [$RESTREAMER_FILE_ROOT_PATH]
   --disable_streaming       [$RESTREAMER_DISABLE_STREAMING]
   --help, -h               show help
   --version, -v            print the version   
```
