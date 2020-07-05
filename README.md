go-radiko
=========

Debugger for radiko.jp authentication.

## Warning

Illegal use of this tool, such as capturing the live stream and uploading it on the public network are not considered personal use. For more information, refer to the government web site.

- [AGENCY FOR CULTURAL AFFAIRS](https://www.bunka.go.jp/english/index.html)

## Installation

Download and install the latest `ggmpeg` and `go` at first.

- [Download FFmpeg](https://ffmpeg.org/download.html)
- [Downloads - The Go Programming Language](https://golang.org/dl/)

Then install `radiko`.

```console
go get -u github.com/moutend/go-radiko/cmd/radiko
```

## Usage

### As a normal member

```console
radiko
```

NOTE that the radiko.jp limits the available area for a normal member. If you want to listen TOKYO FM from another area, e.g. Osaka, you need charge and became a premium member.

### As a premium member

```console
export RADIKO_USERNAME=you@example.com
export RADIKO_PASSWORD=your_password

radiko
```

## LICENSE

MIT

## Author

Yoshiyuki Koyanagi <moutend@gmail.com>
