package notificator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	serverErr        = "server_error"
	externalApiError = "external_service_unavailable"

	jsonDecodeErr = "invalid_json"
)

type Notificator struct {
	client *http.Client
}

func New(client *http.Client) *Notificator {
	return &Notificator{client: client}
}

type Notification struct {
	URL string
}

type Response struct {
	URL   string
	Body  map[string]interface{}
	Code  string
	Error string
}

func (n *Notificator) Send(ctx context.Context, notification Notification) *Response {
	var (
		out = &Response{
			URL: notification.URL,
		}
	)
	req, err := http.NewRequest(http.MethodGet, notification.URL, nil)
	if err != nil {
		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		return out
	}

	httpResp, err := n.client.Do(req)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			out.Code = fmt.Sprintf("%d %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
			out.Error = externalApiError
			return out
		}

		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		return out
	}

	out.Code = httpResp.Status

	d := json.NewDecoder(httpResp.Body)
	for {
		if err := d.Decode(&out.Body); err != io.EOF {
			out.Error = jsonDecodeErr
		}
		break
	}

	return out
}
