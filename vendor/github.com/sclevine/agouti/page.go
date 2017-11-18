package agouti

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sclevine/agouti/api"
)

// A Page represents an open browser session. Pages may be created using the
// *WebDriver.Page() method or by calling the NewPage or SauceLabs functions.
type Page struct {
	selectable
	logs map[string][]Log
}

// A Log represents a single log message
type Log struct {
	// Message is the text of the log message.
	Message string

	// Location is the code location of the log message, if present
	Location string

	// Level is the log level ("DEBUG", "INFO", "WARNING", or "SEVERE").
	Level string

	// Time is the time the message was logged.
	Time time.Time
}

// NewPage opens a Page using the provided WebDriver URL. This method takes
// the same Options as *WebDriver.NewPage. Unlike *WebDriver.NewPage, this
// method will respect the HTTPClient Option if provided.
func NewPage(url string, options ...Option) (*Page, error) {
	pageOptions := config{}.Merge(options)
	session, err := api.OpenWithClient(url, pageOptions.Capabilities(), pageOptions.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %s", err)
	}
	return newPage(session), nil
}

// JoinPage creates a Page using existing session URL.
func JoinPage(url string) *Page {
	session := api.New(url)
	return newPage(session)
}

func newPage(session *api.Session) *Page {
	return &Page{selectable{session, nil}, nil}
}

// String returns a string representation of the Page. Currently: "page"
func (p *Page) String() string {
	return "page"
}

// Session returns a *api.Session that can be used to send direct commands
// to the WebDriver. See: https://code.google.com/p/selenium/wiki/JsonWireProtocol
func (p *Page) Session() *api.Session {
	return p.session.(*api.Session)
}

// Destroy closes any open browsers by ending the session.
func (p *Page) Destroy() error {
	if err := p.session.Delete(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

// Reset deletes all cookies set for the current domain and navigates to a blank page.
// Unlike Destroy, Reset will permit the page to be re-used after it is called.
// Reset is faster than Destroy, but any cookies from domains outside the current
// domain will remain after a page is reset.
func (p *Page) Reset() error {
	p.ConfirmPopup()

	url, err := p.URL()
	if err != nil {
		return err
	}
	if url == "about:blank" {
		return nil
	}

	if err := p.ClearCookies(); err != nil {
		return err
	}

	if err := p.session.DeleteLocalStorage(); err != nil {
		if err := p.RunScript("localStorage.clear();", nil, nil); err != nil {
			return err
		}
	}

	if err := p.session.DeleteSessionStorage(); err != nil {
		if err := p.RunScript("sessionStorage.clear();", nil, nil); err != nil {
			return err
		}
	}

	return p.Navigate("about:blank")
}

// Navigate navigates to the provided URL.
func (p *Page) Navigate(url string) error {
	if err := p.session.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

// GetCookies returns all cookies on the page.
func (p *Page) GetCookies() ([]*http.Cookie, error) {
	apiCookies, err := p.session.GetCookies()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %s", err)
	}
	cookies := []*http.Cookie{}
	for _, apiCookie := range apiCookies {
		expSeconds := int64(apiCookie.Expiry)
		expNano := int64(apiCookie.Expiry-float64(expSeconds)) * 1000000000
		cookie := &http.Cookie{
			Name:     apiCookie.Name,
			Value:    apiCookie.Value,
			Path:     apiCookie.Path,
			Domain:   apiCookie.Domain,
			Secure:   apiCookie.Secure,
			HttpOnly: apiCookie.HTTPOnly,
			Expires:  time.Unix(expSeconds, expNano),
		}
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}

// SetCookie sets a cookie on the page.
func (p *Page) SetCookie(cookie *http.Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}

	var expiry int64
	if !cookie.Expires.IsZero() {
		expiry = cookie.Expires.Unix()
	}

	apiCookie := &api.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		Expiry:   float64(expiry),
	}

	if err := p.session.SetCookie(apiCookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

// DeleteCookie deletes a cookie on the page by name.
func (p *Page) DeleteCookie(name string) error {
	if err := p.session.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

// ClearCookies deletes all cookies on the page.
func (p *Page) ClearCookies() error {
	if err := p.session.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

// URL returns the current page URL.
func (p *Page) URL() (string, error) {
	url, err := p.session.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

// Size sets the current page size in pixels.
func (p *Page) Size(width, height int) error {
	window, err := p.session.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

// Screenshot takes a screenshot and saves it to the provided filename.
// The provided filename may be an absolute or relative path.
func (p *Page) Screenshot(filename string) error {
	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to find absolute path for filename: %s", err)
	}

	screenshot, err := p.session.GetScreenshot()
	if err != nil {
		return fmt.Errorf("failed to retrieve screenshot: %s", err)
	}

	if err := ioutil.WriteFile(absFilePath, screenshot, 0666); err != nil {
		return fmt.Errorf("failed to save screenshot: %s", err)
	}

	return nil
}

// Title returns the page title.
func (p *Page) Title() (string, error) {
	title, err := p.session.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

// HTML returns the current contents of the DOM for the entire page.
func (p *Page) HTML() (string, error) {
	html, err := p.session.GetSource()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page HTML: %s", err)
	}
	return html, nil
}

// RunScript runs the JavaScript provided in the body. Any keys present in
// the arguments map will be available as variables in the body.
// Values provided in arguments are converted into javascript objects.
// If the body returns a value, it will be unmarshalled into the result argument.
// Simple example:
//    var number int
//    page.RunScript("return test;", map[string]interface{}{"test": 100}, &number)
//    fmt.Println(number)
// -> 100
func (p *Page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
	var (
		keys   []string
		values []interface{}
	)

	for key, value := range arguments {
		keys = append(keys, key)
		values = append(values, value)
	}

	argumentList := strings.Join(keys, ", ")
	cleanBody := fmt.Sprintf("return (function(%s) { %s; }).apply(this, arguments);", argumentList, body)

	if err := p.session.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

// PopupText returns the current alert, confirm, or prompt popup text.
func (p *Page) PopupText() (string, error) {
	text, err := p.session.GetAlertText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %s", err)
	}
	return text, nil
}

// EnterPopupText enters text into an open prompt popup.
func (p *Page) EnterPopupText(text string) error {
	if err := p.session.SetAlertText(text); err != nil {
		return fmt.Errorf("failed to enter popup text: %s", err)
	}
	return nil
}

// ConfirmPopup confirms an alert, confirm, or prompt popup.
func (p *Page) ConfirmPopup() error {
	if err := p.session.AcceptAlert(); err != nil {
		return fmt.Errorf("failed to confirm popup: %s", err)
	}
	return nil
}

// CancelPopup cancels an alert, confirm, or prompt popup.
func (p *Page) CancelPopup() error {
	if err := p.session.DismissAlert(); err != nil {
		return fmt.Errorf("failed to cancel popup: %s", err)
	}
	return nil
}

// Forward navigates forward in history.
func (p *Page) Forward() error {
	if err := p.session.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

// Back navigates backwards in history.
func (p *Page) Back() error {
	if err := p.session.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

// Refresh refreshes the page.
func (p *Page) Refresh() error {
	if err := p.session.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}

// SwitchToParentFrame focuses on the immediate parent frame of a frame selected
// by Selection.Frame. After switching, all new and existing selections will refer
// to the parent frame. All further Page methods will apply to this frame as well.
//
// This method is not supported by PhantomJS. Please use SwitchToRootFrame instead.
func (p *Page) SwitchToParentFrame() error {
	if err := p.session.FrameParent(); err != nil {
		return fmt.Errorf("failed to switch to parent frame: %s", err)
	}
	return nil
}

// SwitchToRootFrame focuses on the original, default page frame before any calls
// to Selection.Frame were made. After switching, all new and existing selections
// will refer to the root frame. All further Page methods will apply to this frame
// as well.
func (p *Page) SwitchToRootFrame() error {
	if err := p.session.Frame(nil); err != nil {
		return fmt.Errorf("failed to switch to original page frame: %s", err)
	}
	return nil
}

// SwitchToWindow switches to the first available window with the provided name
// (JavaScript `window.name` attribute).
func (p *Page) SwitchToWindow(name string) error {
	if err := p.session.SetWindowByName(name); err != nil {
		return fmt.Errorf("failed to switch to named window: %s", err)
	}
	return nil
}

// NextWindow switches to the next available window.
func (p *Page) NextWindow() error {
	windows, err := p.session.GetWindows()
	if err != nil {
		return fmt.Errorf("failed to find available windows: %s", err)
	}

	var windowIDs []string
	for _, window := range windows {
		windowIDs = append(windowIDs, window.ID)
	}

	// order not defined according to W3 spec
	sort.Strings(windowIDs)

	activeWindow, err := p.session.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to find active window: %s", err)
	}

	for position, windowID := range windowIDs {
		if windowID == activeWindow.ID {
			activeWindow.ID = windowIDs[(position+1)%len(windowIDs)]
			break
		}
	}

	if err := p.session.SetWindow(activeWindow); err != nil {
		return fmt.Errorf("failed to change active window: %s", err)
	}

	return nil
}

// CloseWindow closes the active window.
func (p *Page) CloseWindow() error {
	if err := p.session.DeleteWindow(); err != nil {
		return fmt.Errorf("failed to close active window: %s", err)
	}
	return nil
}

// WindowCount returns the number of available windows.
func (p *Page) WindowCount() (int, error) {
	windows, err := p.session.GetWindows()
	if err != nil {
		return 0, fmt.Errorf("failed to find available windows: %s", err)
	}
	return len(windows), nil
}

// LogTypes returns all of the valid log types that may be used with a LogReader.
func (p *Page) LogTypes() ([]string, error) {
	types, err := p.session.GetLogTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve log types: %s", err)
	}
	return types, nil
}

// ReadNewLogs returns new log messages of the provided log type. For example,
// page.ReadNewLogs("browser") returns browser console logs, such as JavaScript
// logs and errors. Only logs since the last call to ReadNewLogs are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadNewLogs(logType string) ([]Log, error) {
	if p.logs == nil {
		p.logs = map[string][]Log{}
	}

	clientLogs, err := p.session.NewLogs(logType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %s", err)
	}

	messageMatcher := regexp.MustCompile(`^(?s:(.+))\s\(([^)]*:\w*)\)$`)

	var logs []Log
	for _, clientLog := range clientLogs {
		matches := messageMatcher.FindStringSubmatch(clientLog.Message)
		message, location := clientLog.Message, ""
		if len(matches) > 2 {
			message, location = matches[1], matches[2]
		}

		log := Log{message, location, clientLog.Level, msToTime(clientLog.Timestamp)}
		logs = append(logs, log)
	}
	p.logs[logType] = append(p.logs[logType], logs...)

	return logs, nil
}

// ReadAllLogs returns all log messages of the provided log type. For example,
// page.ReadAllLogs("browser") returns browser console logs, such as JavaScript logs
// and errors. All logs since the session was created are returned.
// Valid log types may be obtained using the LogTypes method.
func (p *Page) ReadAllLogs(logType string) ([]Log, error) {
	if _, err := p.ReadNewLogs(logType); err != nil {
		return nil, err
	}

	return append([]Log(nil), p.logs[logType]...), nil
}

func msToTime(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}

// MoveMouseBy moves the mouse by the provided offset.
func (p *Page) MoveMouseBy(xOffset, yOffset int) error {
	if err := p.session.MoveTo(nil, api.XYOffset{X: xOffset, Y: yOffset}); err != nil {
		return fmt.Errorf("failed to move mouse: %s", err)
	}

	return nil
}

// DoubleClick double clicks the left mouse button at the current mouse
// position.
func (p *Page) DoubleClick() error {
	if err := p.session.DoubleClick(); err != nil {
		return fmt.Errorf("failed to double click: %s", err)
	}

	return nil
}

// Click performs the provided Click event using the provided Button at the
// current mouse position.
func (p *Page) Click(event Click, button Button) error {
	var err error
	switch event {
	case SingleClick:
		err = p.session.Click(api.Button(button))
	case HoldClick:
		err = p.session.ButtonDown(api.Button(button))
	case ReleaseClick:
		err = p.session.ButtonUp(api.Button(button))
	default:
		err = errors.New("invalid touch event")
	}
	if err != nil {
		return fmt.Errorf("failed to %s %s: %s", event, button, err)
	}

	return nil
}

// SetImplicitWait sets the implicit wait timeout (in ms)
func (p *Page) SetImplicitWait(timeout int) error {
	return p.session.SetImplicitWait(timeout)
}

// SetPageLoad sets the page load timeout (in ms)
func (p *Page) SetPageLoad(timeout int) error {
	return p.session.SetPageLoad(timeout)
}

// SetScriptTimeout sets the script timeout (in ms)
func (p *Page) SetScriptTimeout(timeout int) error {
	return p.session.SetScriptTimeout(timeout)
}
