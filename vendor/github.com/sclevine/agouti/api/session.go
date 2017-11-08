package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/sclevine/agouti/api/internal/bus"
)

type Session struct {
	Bus
}

type Bus interface {
	Send(method, endpoint string, body, result interface{}) error
}

func New(sessionURL string) *Session {
	busClient := &bus.Client{sessionURL, http.DefaultClient}
	return &Session{busClient}
}

func Open(url string, capabilities map[string]interface{}) (*Session, error) {
	return OpenWithClient(url, capabilities, nil)
}

func OpenWithClient(url string, capabilities map[string]interface{}, client *http.Client) (*Session, error) {
	busClient, err := bus.Connect(url, capabilities, client)
	if err != nil {
		return nil, err
	}
	return &Session{busClient}, nil
}

func (s *Session) Delete() error {
	return s.Send("DELETE", "", nil, nil)
}

func (s *Session) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := s.Send("POST", "element", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := s.Send("POST", "elements", selector, &results); err != nil {
		return nil, err
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, s})
	}

	return elements, nil
}

func (s *Session) GetActiveElement() (*Element, error) {
	var result struct{ Element string }

	if err := s.Send("POST", "element/active", nil, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetWindow() (*Window, error) {
	var windowID string
	if err := s.Send("GET", "window_handle", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{windowID, s}, nil
}

func (s *Session) GetWindows() ([]*Window, error) {
	var windowsID []string
	if err := s.Send("GET", "window_handles", nil, &windowsID); err != nil {
		return nil, err
	}

	var windows []*Window
	for _, windowID := range windowsID {
		windows = append(windows, &Window{windowID, s})
	}
	return windows, nil
}

func (s *Session) SetWindow(window *Window) error {
	if window == nil {
		return errors.New("nil window is invalid")
	}

	request := struct {
		Name string `json:"name"`
	}{window.ID}

	return s.Send("POST", "window", request, nil)
}

func (s *Session) SetWindowByName(name string) error {
	request := struct {
		Name string `json:"name"`
	}{name}

	return s.Send("POST", "window", request, nil)
}

func (s *Session) DeleteWindow() error {
	if err := s.Send("DELETE", "window", nil, nil); err != nil {
		return err
	}
	return nil
}

func (s *Session) GetCookies() ([]*Cookie, error) {
	var cookies []*Cookie
	if err := s.Send("GET", "cookie", nil, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

func (s *Session) SetCookie(cookie *Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}
	request := struct {
		Cookie *Cookie `json:"cookie"`
	}{cookie}

	return s.Send("POST", "cookie", request, nil)
}

func (s *Session) DeleteCookie(cookieName string) error {
	return s.Send("DELETE", "cookie/"+cookieName, nil, nil)
}

func (s *Session) DeleteCookies() error {
	return s.Send("DELETE", "cookie", nil, nil)
}

func (s *Session) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := s.Send("GET", "screenshot", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (s *Session) GetURL() (string, error) {
	var url string
	if err := s.Send("GET", "url", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (s *Session) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return s.Send("POST", "url", request, nil)
}

func (s *Session) GetTitle() (string, error) {
	var title string
	if err := s.Send("GET", "title", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (s *Session) GetSource() (string, error) {
	var source string
	if err := s.Send("GET", "source", nil, &source); err != nil {
		return "", err
	}

	return source, nil
}

func (s *Session) MoveTo(region *Element, offset Offset) error {
	request := map[string]interface{}{}

	if region != nil {
		request["element"] = region.ID
	}

	if offset != nil {
		if xoffset, present := offset.x(); present {
			request["xoffset"] = xoffset
		}

		if yoffset, present := offset.y(); present {
			request["yoffset"] = yoffset
		}
	}

	return s.Send("POST", "moveto", request, nil)
}

func (s *Session) Frame(frame *Element) error {
	var elementID interface{}

	if frame != nil {
		elementID = struct {
			Element string `json:"ELEMENT"`
		}{frame.ID}
	}

	request := struct {
		ID interface{} `json:"id"`
	}{elementID}

	return s.Send("POST", "frame", request, nil)
}

func (s *Session) FrameParent() error {
	return s.Send("POST", "frame/parent", nil, nil)
}

func (s *Session) Execute(body string, arguments []interface{}, result interface{}) error {
	if arguments == nil {
		arguments = []interface{}{}
	}

	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := s.Send("POST", "execute", request, result); err != nil {
		return err
	}

	return nil
}

func (s *Session) Forward() error {
	return s.Send("POST", "forward", nil, nil)
}

func (s *Session) Back() error {
	return s.Send("POST", "back", nil, nil)
}

func (s *Session) Refresh() error {
	return s.Send("POST", "refresh", nil, nil)
}

func (s *Session) GetAlertText() (string, error) {
	var text string
	if err := s.Send("GET", "alert_text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Session) SetAlertText(text string) error {
	request := struct {
		Text string `json:"text"`
	}{text}
	return s.Send("POST", "alert_text", request, nil)
}

func (s *Session) AcceptAlert() error {
	return s.Send("POST", "accept_alert", nil, nil)
}

func (s *Session) DismissAlert() error {
	return s.Send("POST", "dismiss_alert", nil, nil)
}

func (s *Session) NewLogs(logType string) ([]Log, error) {
	request := struct {
		Type string `json:"type"`
	}{logType}

	var logs []Log
	if err := s.Send("POST", "log", request, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Session) GetLogTypes() ([]string, error) {
	var types []string
	if err := s.Send("GET", "log/types", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}

func (s *Session) DoubleClick() error {
	return s.Send("POST", "doubleclick", nil, nil)
}

func (s *Session) Click(button Button) error {
	request := struct {
		Button Button `json:"button"`
	}{button}
	return s.Send("POST", "click", request, nil)
}

func (s *Session) ButtonDown(button Button) error {
	request := struct {
		Button Button `json:"button"`
	}{button}
	return s.Send("POST", "buttondown", request, nil)
}

func (s *Session) ButtonUp(button Button) error {
	request := struct {
		Button Button `json:"button"`
	}{button}
	return s.Send("POST", "buttonup", request, nil)
}

func (s *Session) TouchDown(x, y int) error {
	request := struct {
		X int `json:"x"`
		Y int `json:"y"`
	}{x, y}
	return s.Send("POST", "touch/down", request, nil)
}

func (s *Session) TouchUp(x, y int) error {
	request := struct {
		X int `json:"x"`
		Y int `json:"y"`
	}{x, y}
	return s.Send("POST", "touch/up", request, nil)
}

func (s *Session) TouchMove(x, y int) error {
	request := struct {
		X int `json:"x"`
		Y int `json:"y"`
	}{x, y}
	return s.Send("POST", "touch/move", request, nil)
}

func (s *Session) TouchClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}

	request := struct {
		Element string `json:"element"`
	}{element.ID}
	return s.Send("POST", "touch/click", request, nil)
}

func (s *Session) TouchDoubleClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}

	request := struct {
		Element string `json:"element"`
	}{element.ID}
	return s.Send("POST", "touch/doubleclick", request, nil)
}

func (s *Session) TouchLongClick(element *Element) error {
	if element == nil {
		return errors.New("nil element is invalid")
	}

	request := struct {
		Element string `json:"element"`
	}{element.ID}
	return s.Send("POST", "touch/longclick", request, nil)
}

func (s *Session) TouchFlick(element *Element, offset Offset, speed Speed) error {
	if speed == nil {
		return errors.New("nil speed is invalid")
	}

	if (element == nil) != (offset == nil) {
		return errors.New("element must be provided if offset is provided and vice versa")
	}

	var request interface{}
	if element == nil {
		xSpeed, ySpeed := speed.vector()
		request = struct {
			XSpeed int `json:"xspeed"`
			YSpeed int `json:"yspeed"`
		}{xSpeed, ySpeed}
	} else {
		xOffset, yOffset := offset.position()
		request = struct {
			Element string `json:"element"`
			XOffset int    `json:"xoffset"`
			YOffset int    `json:"yoffset"`
			Speed   uint   `json:"speed"`
		}{element.ID, xOffset, yOffset, speed.scalar()}
	}

	return s.Send("POST", "touch/flick", request, nil)
}

func (s *Session) TouchScroll(element *Element, offset Offset) error {
	if element == nil {
		element = &Element{}
	}

	if offset == nil {
		return errors.New("nil offset is invalid")
	}

	xOffset, yOffset := offset.position()
	request := struct {
		Element string `json:"element,omitempty"`
		XOffset int    `json:"xoffset"`
		YOffset int    `json:"yoffset"`
	}{element.ID, xOffset, yOffset}
	return s.Send("POST", "touch/scroll", request, nil)
}

func (s *Session) Keys(text string) error {
	splitText := strings.Split(text, "")
	request := struct {
		Value []string `json:"value"`
	}{splitText}
	return s.Send("POST", "keys", request, nil)
}

func (s *Session) DeleteLocalStorage() error {
	return s.Send("DELETE", "local_storage", nil, nil)
}

func (s *Session) DeleteSessionStorage() error {
	return s.Send("DELETE", "session_storage", nil, nil)
}

func (s *Session) SetImplicitWait(timeout int) error {
	request := struct {
		MS int `json:"ms"`
	}{timeout}
	return s.Send("POST", "timeouts/implicit_wait", request, nil)
}

func (s *Session) SetPageLoad(timeout int) error {
	request := struct {
		MS   int    `json:"ms"`
		Type string `json:"type"`
	}{timeout, "page load"}
	return s.Send("POST", "timeouts", request, nil)
}

func (s *Session) SetScriptTimeout(timeout int) error {
	request := struct {
		MS int `json:"ms"`
	}{timeout}
	return s.Send("POST", "timeouts/async_script", request, nil)
}
