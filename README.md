go-radiko
=========

The command line tool for playing and recording the radiko.jp live stream. Also comes with the Go package.

## Warning

Illegal use of this tool, such as capturing the live stream and uploading it on the public network are not considered personal use.

## Installation

First, download and install the latest `ggmpeg` and `go`.

- [Download FFmpeg](https://ffmpeg.org/download.html)
- [Downloads - The Go Programming Language](https://golang.org/dl/)

Open Terminal app, and then run the floowing command:

```console
go install github.com/moutend/go-radiko/cmd/radiko@latest
```

## Usage

### List all radio stations

Use `station` sub command to print all radio stations.

```console
radiko station
```

### Play Live Stream

#### As a normal member

To listen the Tokyo FM (its identifier is `FMT`), run the following command.

```console
radiko play FMT
```

If you are in Tokyo or near location, you can listen the live stream. If not, the playback will never start.

#### As a premium member

You need set the environment variable or create configuration file.

##### Environment variable

```console
export RADIKO_USERNAME="you@example.com"
export RADIKO_PASSWORD="xxxxxxxx"
```

Then run the command. You are able to listen the live stream regardless the location.

```console
radiko play FMT
```

##### Configuration file

Create the text file with the following content and then save as `radiko.toml`.

```toml
radiko_username = "you@example.com"
radiko_password = "xxxxxxxx"
```

Run the command with `-c` flag.

```console
radiko -c radiko.toml play FMT
```

## Author

Yoshiyuki Koyanagi <moutend@gmail.com>

## LICENSE

MIT
