package radiko

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Client is a radiko API client.
type Client struct {
	station  string
	username string
	password string

	debug *log.Logger

	// AExp corresponds to the cookie value named 'a_exp'.
	AExp string

	// AreaName corresponds to the value of geo-restriction area name.
	AreaName string

	// RadikoSession corresponds to the cookie value named 'radiko_session'.
	RadikoSession string

	// FullKey is a seed string for authentication.
	FullKey string

	// RawPartialKey is a raw partial key.
	RawPartialKey string

	// Base64EncodedPartialKey is a partial key which is base64 encoded.
	Base64EncodedPartialKey string

	// AuthToken corresponds to the response header value named 'X-Radiko-AuthToken'.
	AuthToken string

	// KeyOffset is a value need when generating partial key.
	KeyOffset int

	// KeyLength is a value need when generating a partial key.
	KeyLength int

	// AllStations holds a result of GetAllStations method.
	AllStations []Station

	// PlaylistM3U8s holds a result of Playlist method.
	PlaylistM3U8s []PlaylistM3U8
}

// New returns a client.
func New(station, username, password string) *Client {
	now := time.Now()
	sum := md5.Sum([]byte(fmt.Sprint(now.Unix())))

	return &Client{
		station:  station,
		username: username,
		password: password,
		debug:    log.New(io.Discard, "", 0),
		AExp:     hex.EncodeToString(sum[:]),
	}
}

// GetAllStations fetches a list of all radio stations.
//
// The result is stored in AllStations field.
//
// You can call this method without any authentication.
func (c *Client) GetAllStations(ctx context.Context) error {
	const u = "https://radiko.jp/v3/station/region/full.xml"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("stations: failed to create request: %w", err)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("stations: failed to fetch full.xml: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("stations: status code:", res.Status)

	body := &bytes.Buffer{}

	data, err := io.ReadAll(io.TeeReader(res.Body, body))

	if err != nil {
		return fmt.Errorf("stations: failed to copy response body: %w", err)
	}

	c.debug.Printf("stations: response body: %q\n", data)

	allStations, err := ParseFullStationXML(body)

	if err != nil {
		return fmt.Errorf("stations: failed to parse response body: %w", err)
	}

	c.AllStations = allStations

	return nil
}

// GetAreaName fetches an area name based off your IP.
//
// The result is stored AreaName field.
//
// You can call this method without any authentication.
func (c *Client) GetAreaName(ctx context.Context) error {
	const u = "https://radiko.jp/area"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("area: failed to create request: %w", err)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("area: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("area: status code:", res.Status)

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("area: failed to read response body: %w", err)
	}

	tokens := strings.Split(string(body), "\"")

	if len(tokens) < 3 {
		return fmt.Errorf("area: unexpected response body: %q", string(body))
	}

	c.AreaName = tokens[1]

	if c.AreaName == "" {
		return fmt.Errorf("area: empty area")
	}

	c.debug.Printf("area: your area name: %q\n", c.AreaName)

	if !strings.HasPrefix(c.AreaName, "JP") {
		return fmt.Errorf("area: your IP is geo-restricted")
	}

	return nil
}

// GetSeed fetches a seed string which is required for authentication.
//
// You can call this method without any authentication.
func (c *Client) GetSeed(ctx context.Context) error {
	const u = "https://radiko.jp/apps/js/playerCommon.js"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("seed: failed to create request: %w", err)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("seed: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("seed: status code:", res.Status)

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("seed: failed to read response body: %w", err)
	}

	tokens := strings.Split(string(body), " ")

	for i, token := range tokens {
		if strings.HasPrefix(token, "RadikoJSPlayer(") {
			fullKey := tokens[i+2]

			fullKey = strings.TrimPrefix(fullKey, "'")
			fullKey = strings.TrimSuffix(fullKey, "',")

			c.FullKey = fullKey

			break
		}
	}
	if c.FullKey == "" {
		return fmt.Errorf("seed: empty full key")
	}

	return nil
}

// Login performs login step.
func (c *Client) Login(ctx context.Context) error {
	if c.username == "" || c.password == "" {
		c.debug.Println("login: continue as normal member")

		return nil
	}

	c.debug.Println("login: continue as premium member")

	values := &url.Values{}

	values.Add("mail", c.username)
	values.Add("pass", c.password)

	const u = "https://radiko.jp/v4/api/member/login"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBufferString(values.Encode()))

	if err != nil {
		return fmt.Errorf("login: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("login: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("login: status code:", res.Status)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("login: failed to login")
	}

	var response struct {
		RadikoSession string `json:"radiko_session"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("login: failed to parse response body: %w", err)
	}

	c.debug.Printf("login: radiko_session: %q\n", response.RadikoSession)

	c.RadikoSession = response.RadikoSession

	// Dummy wait
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Check performs a step validating your radiko member status.
func (c *Client) Check(ctx context.Context) error {
	if c.RadikoSession != "" {
		c.debug.Println("check: skip this step")

		return nil
	}

	const u = "https://radiko.jp/ap/member/webapi/v2/member/login/check"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("check: failed to create request: %w", err)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("check: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("check: status code:", res.Status)
	c.debug.Println("check: hint: this API returns 400 Bad Request when continue as normal member")

	for _, cookie := range res.Cookies() {
		if cookie.Name == "radiko_session" {
			c.RadikoSession = cookie.Value

			c.debug.Printf("check: success: radiko_session=%q\n", c.RadikoSession)

			return nil
		}
	}

	return fmt.Errorf("check: radiko_session not found")
}

// Auth1 performs a authentication step required at first.
func (c *Client) Auth1(ctx context.Context) error {
	const u = "https://radiko.jp/v2/api/auth1"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("auth1: failed to create request: %w", err)
	}

	cookie := fmt.Sprintf(
		"a_exp=%s; default_area_id=%s; radiko_session=%s; tracking_area_id=%s",
		c.AExp,
		c.AreaName,
		c.RadikoSession,
		c.AreaName,
	)

	c.debug.Printf("auth1: cookie=%q\n", cookie)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-radiko-app", "pc_html5")
	req.Header.Set("x-radiko-user", "dummy_user")
	req.Header.Set("x-radiko-device", "pc")
	req.Header.Set("x-radiko-app-version", "0.0.1")

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("auth1: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("auth1: status code:", res.Status)

	c.AuthToken = res.Header.Get("X-Radiko-Authtoken")

	if c.AuthToken == "" {
		return fmt.Errorf("auth1: empty auth token")
	}

	keyLength, err := strconv.ParseInt(res.Header.Get(`X-Radiko-KeyLength`), 10, 64)

	if err != nil {
		return fmt.Errorf("auth1: failed to parse key length: %w", err)
	}
	if keyLength == 0 {
		return fmt.Errorf("auth1: key length is 0")
	}

	keyOffset, err := strconv.ParseInt(res.Header.Get(`X-Radiko-KeyOffset`), 10, 64)

	if err != nil {
		return fmt.Errorf("auth1: failed to parse key offset: %w", err)
	}
	if int(keyOffset+keyLength) > len(c.FullKey) {
		return fmt.Errorf("auth1: invalid key length and offset: length=%v, offset=%v", keyLength, keyOffset)
	}

	c.KeyLength = int(keyLength)
	c.KeyOffset = int(keyOffset)

	c.RawPartialKey = c.FullKey[keyOffset : keyOffset+keyLength]
	c.Base64EncodedPartialKey = base64.StdEncoding.EncodeToString([]byte(c.RawPartialKey))

	c.debug.Println("auth1: raw partial key:", c.RawPartialKey)
	c.debug.Println("auth1: base64 encoded partial key:", c.Base64EncodedPartialKey)

	return nil
}

// Auth2 performs a authentication step required at second.
func (c *Client) Auth2(ctx context.Context) error {
	const u = "https://radiko.jp/v2/api/auth2"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("auth2: failed to create request: %w", err)
	}

	cookie := fmt.Sprintf(
		"a_exp=%s; default_area_id=%s; radiko_session=%s; tracking_area_id=%s",
		c.AExp,
		c.AreaName,
		c.RadikoSession,
		c.AreaName,
	)

	c.debug.Printf("auth2: cookie=%q\n", cookie)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-radiko-authtoken", c.AuthToken)
	req.Header.Set("x-radiko-partialkey", c.Base64EncodedPartialKey)
	req.Header.Set("x-radiko-user", "dummy_user")
	req.Header.Set("x-radiko-device", "pc")

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("auth2: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("auth2: status code:", res.Status)

	return nil
}

// Playlist fetches a list of URL used for generating an HLS source.
func (c *Client) Playlist(ctx context.Context) error {
	u := fmt.Sprintf("https://radiko.jp/v3/station/stream/pc_html5/%s.xml", c.station)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return fmt.Errorf("playlist: failed to create request: %w", err)
	}

	cookie := fmt.Sprintf(
		"a_exp=%s; default_area_id=%s; radiko_session=%s; tracking_area_id=%s",
		c.AExp,
		c.AreaName,
		c.RadikoSession,
		c.AreaName,
	)

	c.debug.Printf("playlist: cookie=%q\n", cookie)

	req.Header.Set("Cookie", cookie)

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("playlist: error response: %w", err)
	}

	defer res.Body.Close()

	c.debug.Println("playlist: status code:", res.Status)

	body := &bytes.Buffer{}

	data, err := io.ReadAll(io.TeeReader(res.Body, body))

	if err != nil {
		return fmt.Errorf("playlist: failed to copy response body: %w", err)
	}

	c.debug.Printf("playlist: response body: %q\n", data)

	playlistM3U8s, err := ParsePlaylistCreateXML(body)

	if err != nil {
		return fmt.Errorf("playlist: failed to parse response body: %w", err)
	}

	c.PlaylistM3U8s = playlistM3U8s

	return nil
}

// Play launches ffmpeg command which plays live streaming.
func (c *Client) Play(ctx context.Context, playbackVolume int) error {
	if err := c.GetAreaName(ctx); err != nil {
		return fmt.Errorf("radiko: failed to get area name: %w", err)
	}
	if err := c.GetSeed(ctx); err != nil {
		return fmt.Errorf("radiko: failed to get seed: %w", err)
	}
	if err := c.Login(ctx); err != nil {
		return fmt.Errorf("radiko: failed to login: %w", err)
	}
	if err := c.Check(ctx); err != nil {
		return fmt.Errorf("radiko: failed to check member status: %w", err)
	}
	if err := c.Auth1(ctx); err != nil {
		return fmt.Errorf("radiko: failed to complete first authentication: %w", err)
	}
	if err := c.Auth2(ctx); err != nil {
		return fmt.Errorf("radiko: failed to complete second authentication: %w", err)
	}
	if err := c.Playlist(ctx); err != nil {
		return fmt.Errorf("radiko: failed to get playlist: %w", err)
	}

	headers := fmt.Sprintf("X-Radiko-AuthToken: %s", c.AuthToken)
	input := fmt.Sprintf(
		"https://rd-wowza-radiko.radiko-cf.com/so/playlist.m3u8?station_id=%s&l=15&lsid=%s&type=c",
		c.station,
		c.AExp,
	)

	ffmpeg := exec.CommandContext(
		ctx, "ffmpeg",
		"-headers", headers,
		"-i", input,
		"-f", "matroska", "-",
	)
	ffplay := exec.CommandContext(
		ctx, "ffplay",
		"-volume", fmt.Sprint(playbackVolume),
		"-i", "-",
	)

	pr, pw := io.Pipe()

	ffmpeg.Stdout = pw
	ffplay.Stdin = pr

	defer pw.Close()
	defer pr.Close()

	if err := ffmpeg.Start(); err != nil {
		return fmt.Errorf("radiko: failed to start ffmpeg command: %w", err)
	}
	if err := ffplay.Start(); err != nil {
		return fmt.Errorf("radiko: failed to start ffplay command: %w", err)
	}
	if err := ffmpeg.Wait(); err != nil {
		return fmt.Errorf("radiko: ffmpeg: unexpected error: %w", err)
	}
	if err := ffplay.Wait(); err != nil {
		return fmt.Errorf("radiko: ffplay: unexpected error: %w", err)
	}

	return nil
}

// Rec launches ffmpeg command and records a specified radio program.
func (c *Client) Rec(ctx context.Context, date time.Time, length time.Duration, outputFile string) error {
	if err := c.GetAreaName(ctx); err != nil {
		return fmt.Errorf("radiko: failed to get area name: %w", err)
	}
	if err := c.GetSeed(ctx); err != nil {
		return fmt.Errorf("radiko: failed to get seed: %w", err)
	}
	if err := c.Login(ctx); err != nil {
		return fmt.Errorf("radiko: failed to login: %w", err)
	}
	if err := c.Check(ctx); err != nil {
		return fmt.Errorf("radiko: failed to check member status: %w", err)
	}
	if err := c.Auth1(ctx); err != nil {
		return fmt.Errorf("radiko: failed to complete first authentication: %w", err)
	}
	if err := c.Auth2(ctx); err != nil {
		return fmt.Errorf("radiko: failed to complete second authentication: %w", err)
	}

	header := fmt.Sprintf("X-Radiko-AuthToken: %s", c.AuthToken)
	input := fmt.Sprintf(
		"https://rd-wowza-radiko.radiko-cf.com/tf/playlist.m3u8?station_id=%s&start_at=%s&ft=%s&end_at=%s&to=%s&l=15&lsid=%s&type=c",
		c.station,
		date.Format("20060102150405"),
		date.Format("20060102150405"),
		date.Add(length).Format("20060102150405"),
		date.Add(length).Format("20060102150405"),
		c.AExp,
	)
	ffmpeg := exec.CommandContext(
		ctx, "ffmpeg",
		"-headers", header,
		"-i", input,
		"-acodec", `copy`,
		"-vn",
		"-bsf:a", "aac_adtstoasc",
		"-y", outputFile,
	)

	if err := ffmpeg.Run(); err != nil {
		return fmt.Errorf("radiko: failed to complete ffmpeg command: %w", err)
	}

	return nil
}

// SetLogger sets a logger for printing debug messages.
func (c *Client) SetLogger(logger *log.Logger) {
	if logger == nil {
		return
	}

	c.debug = logger
}
