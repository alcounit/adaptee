package selenoid

/* -------------- *
 * SELENOID TYPES *
 * ------v------- */

// Caps - user capabilities
type Caps struct {
	BrowserName      string `json:"browserName"`
	BrowserVersion   string `json:"version"`
	ScreenResolution string `json:"screenResolution"`
	VNC              bool   `json:"enableVNC"`
	TestName         string `json:"name"`
	TimeZone         string `json:"timeZone"`
}

// Session - session id and vnc flag
type Session struct {
	ID        string `json:"id"`
	Container string `json:"container"`
	Caps      Caps   `json:"caps"`
}

// Sessions - used count and individual sessions for quota user
type Sessions struct {
	Count    int       `json:"count"`
	Sessions []Session `json:"sessions"`
}

// Quota - list of sessions for quota user
type Quota map[string]*Sessions

// Version - browser version for quota
type Version map[string]Quota

// Browsers - browser names for versions
type Browsers map[string]Version

// State - current state
type State struct {
	Total    int      `json:"total"`
	Used     int      `json:"used"`
	Pending  int      `json:"pending"`
	Browsers Browsers `json:"browsers"`
}

// sessionInfo - extended session information
type sessionInfo struct {
	ID        string `json:"id"`
	Container string `json:"container"`
	Caps      Caps   `json:"caps"`
	Quota     string `json:"quota"`
}

// result - processed selenoid state
type result struct {
	State    State                  `json:"state"`
	Origin   string                 `json:"origin"`
	Browsers map[string]int         `json:"browsers"`
	Sessions map[string]sessionInfo `json:"sessions"`
	Version  string                 `json:"version"`
	Errors   []interface{}          `json:"errors"`
}
