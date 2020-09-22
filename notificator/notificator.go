package notificator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cat-in-vacuum/middleware_task/log"
	"net"
	"net/http"
	"net/url"
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

func (n *Notificator) Send(ctx context.Context, notification Notification) Response {
	var (
		out = Response{
			URL: notification.URL,
		}
	)
	_, err := url.ParseRequestURI(out.URL)
	if err != nil {
		out.Error = err.Error()
		log.Error(err)
		return out
	}
	req, err := http.NewRequest(http.MethodGet, notification.URL, nil)
	if err != nil {
		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		log.Error(err)
		return out
	}

	httpResp, err := n.client.Do(req)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			out.Code = fmt.Sprintf("%d %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
			out.Error = externalApiError
			log.Error(err)
			return out
		}

		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		log.Error(err)
		return out
	}

	out.Code = httpResp.Status

	if err := json.NewDecoder(httpResp.Body).Decode(&out.Body); err != nil {
		out.Error = jsonDecodeErr
		log.Error(err)
		return out
	}

	return out
}
