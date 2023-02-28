# go-radiko

The command line tool for playing and recording the radiko.jp live stream. Also comes with the radiko API client for Go.

## Warning

Illegal use of this tool, such as capturing the live stream and uploading it on the public network are not considered personal use.

## Installation

First, download and install the latest `ggmpeg` and `go`.

- [Download FFmpeg](https://ffmpeg.org/download.html)
- [Downloads - The Go Programming Language](https://golang.org/dl/)

And then run the command below:

```console
go install github.com/moutend/go-radiko/cmd/radiko@latest
```

## Usage

### List all radio stations

```console
radiko station
```

### Play Live Stream

As a normal member:

For instance, to listen the Tokyo FM (its identifier is `FMT`), run the following command.

```console
radiko play FMT
```

The radiko.jp is geo-restricted by your IP. If you are in Tokyo or near location, you can listen the live stream. Otherwise, the playback will never start.

As a premium member:

You need set the environment variable or create configuration file.

**Environment variable**

```console
export RADIKO_USERNAME="you@example.com"
export RADIKO_PASSWORD="xxxxxxxx"
```

**Configuration file**

Create the text file with the following content and then save as `radiko.toml`.

```toml
radiko_username = "you@example.com"
radiko_password = "xxxxxxxx"
```

And then run the command with `-c` flag:

```console
radiko -c radiko.toml play FMT
```

## Author

Yoshiyuki Koyanagi <moutend@gmail.com>

## LICENSE

MIT
