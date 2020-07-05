package radiko

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetAuthKey returns auth key.
func GetAuthKey() (string, error) {
	res, err := http.Get(fmt.Sprintf("http://radiko.jp/apps/js/playerCommon.js?_=%d", time.Now().Unix()))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", fmt.Errorf("radiko: body is empty")
	}

	tokens := strings.Split(string(body), " ")
	authKey := ""

	for i, token := range tokens {
		if strings.HasPrefix(token, "RadikoJSPlayer(") {
			authKey = tokens[i+2]
			authKey = strings.TrimPrefix(authKey, "'")
			authKey = strings.TrimSuffix(authKey, "',")
		}
	}

	return authKey, nil
}

// Station represents radio station information.
type Station struct {
	Identifier string
	Name       string
}

// GetStations returns all available radio stations.
func GetStations() ([]Station, error) {
	res, err := http.Get(`http://radiko.jp/index/`)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("<a href=\"/index/.+/\">.+</a>")

	matches := re.FindAllString(string(body), -1)
	stations := make([]Station, len(matches))

	for n, match := range matches {
		{
			m := strings.TrimPrefix(match, "<a href=\"/index/")
			i := strings.Index(m, "/")

			stations[n].Identifier = m[:i]
		}
		{
			i := strings.Index(match, ">")
			j := strings.Index(match[i:], "<")

			stations[n].Name = match[i+1 : i+j]
		}
	}

	return stations, nil
}

// Session holds parameters they required during authentication.
type Session struct {
	RadikoSession string
	AuthToken     string
	AuthKey       string
	PartialKey    string
	KeyLength     int
	KeyOffset     int
	debug         *log.Logger
	username      string
	password      string
}

// NewSession creates a session.
func NewSession(username, password string) *Session {
	return &Session{
		debug:    log.New(ioutil.Discard, "", 0),
		username: username,
		password: password,
	}
}

// Login performs login action.
func (s *Session) Login() error {
	if s.username == "" || s.password == "" {
		s.debug.Println("continue as normal member")

		goto GET_AUTH_KEY
	}

	s.debug.Println("continue as premium member")

	if err := s.loginWithWebForm(); err != nil {
		s.debug.Println("try another login method because the first login method failed")

		if err := s.loginWithJSONAPI(); err != nil {
			return err
		}
	}

GET_AUTH_KEY:

	authKey, err := GetAuthKey()

	s.debug.Println("auth key:", authKey)

	if err != nil {
		return err
	}

	s.AuthKey = authKey

	return nil
}

func (s *Session) loginWithJSONAPI() error {
	s.debug.Println("try login with JSON API")

	values := &url.Values{}

	values.Add("mail", s.username)
	values.Add("pass", s.password)

	req, err := http.NewRequest(
		http.MethodPost,
		`https://radiko.jp/v4/api/member/login`,
		bytes.NewBufferString(values.Encode()),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	s.debug.Println("login status:", res.Status)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("radiko: failed to login")
	}

	var result struct {
		RadikoSession string `json:"radiko_session"`
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	s.debug.Println("radiko session:", result.RadikoSession)

	s.RadikoSession = result.RadikoSession

	// Dummy wait
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (s *Session) loginWithWebForm() error {
	s.debug.Println("try login with web form")

	values := &url.Values{}

	values.Add("mail", s.username)
	values.Add("pass", s.password)

	req, err := http.NewRequest(
		http.MethodPost,
		`https://radiko.jp/ap/member/login/login`,
		bytes.NewBufferString(values.Encode()),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", `application/x-www-form-urlencoded`)

	res, err := (&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	s.debug.Println("login status:", res.Status)

	if res.StatusCode != http.StatusFound {
		return fmt.Errorf("radiko: failed to login")
	}

	radikoSession := ""

	for _, cookie := range res.Cookies() {
		if cookie.Name == "radiko_session" {
			radikoSession = cookie.Value
		}
	}

	s.debug.Println("radiko session:", radikoSession)

	s.RadikoSession = radikoSession

	return nil
}

// Auth1 performs first authentication.
func (s *Session) Auth1() error {
	req, err := http.NewRequest(http.MethodGet, `https://radiko.jp/v2/api/auth1`, nil)

	if err != nil {
		return err
	}

	req.Header.Set(`X-Radiko-App`, "pc_html5")
	req.Header.Set(`X-Radiko-App-Version`, "0.0.1")
	req.Header.Set(`X-Radiko-User`, "dummy_user")
	req.Header.Set(`X-Radiko-Device`, "pc")

	jar, err := cookiejar.New(nil)

	if err != nil {
		return err
	}

	res, err := (&http.Client{Jar: jar}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	s.debug.Println("first authentication:", res.Status)

	authToken := res.Header.Get(`X-Radiko-Authtoken`)

	s.debug.Println("auth token:", authToken)

	if authToken == "" {
		return fmt.Errorf("radiko: failed to obtain auth token")
	}

	keyLength, err := strconv.ParseInt(res.Header.Get(`X-Radiko-KeyLength`), 10, 64)

	s.debug.Println("key length:", keyLength)

	if err != nil || keyLength == 0 {
		return fmt.Errorf("radiko: failed to obtain key length")
	}

	keyOffset, err := strconv.ParseInt(res.Header.Get(`X-Radiko-KeyOffset`), 10, 64)

	s.debug.Println("key offset:", keyOffset)

	if err != nil {
		return fmt.Errorf("radiko: failed to obtain key offset")
	}
	if int(keyOffset+keyLength) > len(s.AuthKey) {
		return fmt.Errorf("radiko: invalid partial key: (offset=%d, length=%d)", keyOffset, keyLength)
	}

	s.AuthToken = authToken
	s.KeyLength = int(keyLength)
	s.KeyOffset = int(keyOffset)
	s.PartialKey = base64.StdEncoding.EncodeToString([]byte(s.AuthKey[keyOffset : keyOffset+keyLength]))

	s.debug.Println("auth partial key:", s.PartialKey)

	return nil
}

// Auth2 performs second authentication.
func (s *Session) Auth2() error {
	req, err := http.NewRequest(http.MethodGet, `https://radiko.jp/v2/api/auth2`, nil)

	if err != nil {
		return err
	}

	req.Header.Set("x-radiko-user", "dummy_user")
	req.Header.Set("x-radiko-device", "pc")
	req.Header.Set("x-radiko-authtoken", s.AuthToken)
	req.Header.Set("x-radiko-partialkey", s.PartialKey)

	req.AddCookie(&http.Cookie{Name: "radiko_session", Value: s.RadikoSession})

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	s.debug.Println("auth2:", res.Status)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("radiko: failed to complete 2nd authentication step")
	}

	return nil
}

// SetLogger sets default logger for debugging.
func (s *Session) SetLogger(logger *log.Logger) {
	if logger == nil {
		return
	}

	s.debug = logger
}
