package adaptee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alcounit/adaptee/selenoid"
	"github.com/koding/websocketproxy"
	"github.com/sirupsen/logrus"
)

var (
	defaults = struct {
		serviceType, testName, browserName, browserVersion, screenResolution, enableVNC, timeZone, session string
	}{
		serviceType:      "type",
		testName:         "testName",
		browserName:      "browserName",
		browserVersion:   "browserVersion",
		screenResolution: "SCREEN_RESOLUTION",
		enableVNC:        "ENABLE_VNC",
		timeZone:         "TZ",
		session:          "session",
	}
)

//Service ...
type Service struct {
	SessionID string            `json:"id"`
	Labels    map[string]string `json:"labels"`
}

//Selenosis ...
type Selenosis struct {
	Total    int                 `json:"total"`
	Used     int                 `json:"used"`
	Browsers map[string][]string `json:"config,omitempty"`
	Sessions []*Service          `json:"sessions,omitempty"`
}

//HandleStatus ...
func (app *App) HandleStatus(w http.ResponseWriter, r *http.Request) {
	logger := app.logger.WithFields(logrus.Fields{
		"request": fmt.Sprintf("%s %s", r.Method, r.URL.Path),
	})

	state := selenoid.State{
		Browsers: make(selenoid.Browsers),
	}

	w.Header().Set("Content-Type", "application/json")

	resp, err := http.Get(app.selenosisURL + "/status")
	if err != nil {
		json.NewEncoder(w).Encode(state)
		logger.Errorf("failed to get stats from selenosis: %v", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		json.NewEncoder(w).Encode(state)
		logger.Errorf("failed to get read response from selenosis: %v", err)
		return
	}

	defer r.Body.Close()

	type status struct {
		Status    int       `json:"status"`
		Error     error     `json:"err"`
		Selenosis Selenosis `json:"selenosis"`
	}

	stats := status{}
	err = json.Unmarshal(body, &stats)
	if err != nil {
		json.NewEncoder(w).Encode(state)
		logger.Errorf("failed to unmarshal response from selenosis: %v", err)
		return
	}

	state.Total = stats.Selenosis.Total

	for n, b := range stats.Selenosis.Browsers {
		state.Browsers[n] = make(selenoid.Version)
		for _, v := range b {
			state.Browsers[n][v] = make(selenoid.Quota)
		}
	}

	for _, entry := range stats.Selenosis.Sessions {
		state.Used++
		browserName := entry.Labels[defaults.browserName]
		browserVersion := entry.Labels[defaults.browserVersion]
		_, ok := state.Browsers[browserName]
		if !ok {
			state.Browsers[browserName] = make(selenoid.Version)
		}
		_, ok = state.Browsers[browserName][browserVersion]
		if !ok {
			state.Browsers[browserName][browserVersion] = make(selenoid.Quota)
		}

		v, ok := state.Browsers[browserName][browserVersion]["unknown"]
		if !ok {
			v = &selenoid.Sessions{0, []selenoid.Session{}}
			state.Browsers[browserName][browserVersion]["unknown"] = v
		}
		var vnc bool
		vnc, err = strconv.ParseBool(entry.Labels[defaults.enableVNC])
		if err != nil {
			vnc = false
		}
		v.Count++
		sess := selenoid.Session{
			ID:        entry.SessionID,
			Container: entry.SessionID,
			Caps: selenoid.Caps{
				BrowserName:      browserName,
				BrowserVersion:   browserVersion,
				ScreenResolution: entry.Labels[defaults.screenResolution],
				VNC:              vnc,
				TestName:         entry.Labels[defaults.testName],
				TimeZone:         entry.Labels[defaults.timeZone],
			},
		}
		sess.Container = entry.SessionID
		v.Sessions = append(v.Sessions, sess)
	}

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(state)
	logger.Info("status response OK")

}

//HandleWs ...
func (app *App) HandleWs(w http.ResponseWriter, r *http.Request) {
	logger := app.logger.WithFields(logrus.Fields{
		"request": fmt.Sprintf("%s %s", r.Method, r.URL.Path),
	})
	u, _ := url.Parse(app.selenosisURL)
	ws := &url.URL{Scheme: "ws", Host: u.Host, Path: r.URL.Path}
	websocketproxy.DefaultUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wsProxy := websocketproxy.NewProxy(ws)
	logger.Infof("ws proxy: %v", ws.Path)
	wsProxy.ServeHTTP(w, r)
}
