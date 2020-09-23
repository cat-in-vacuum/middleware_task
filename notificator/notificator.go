package notificator

import (
	"context"
	"encoding/json"
	"errors"
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
	ctxCancelErr = "shutdown_server_or_canceling_request"
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
	URL   string      `json:"url,omitempty"`
	Body  interface{} `json:"body,omitempty"`
	Code  string      `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (n *Notificator) Send(ctx context.Context, notification Notification) Response {
	var (
		out = Response{
			URL: notification.URL,
		}
	)

	err := validURL(out.URL)
	if err != nil {
		out.Error = err.Error()
		log.Error(err)
		return out
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, notification.URL, nil)
	if err != nil {
		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		log.Error(err)
		return out
	}

	// time.Sleep(time.Second*5)

	httpResp, err := n.client.Do(req)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			out.Code = fmt.Sprintf("%d %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
			out.Error = externalApiError
			log.Error(err)
			return out
		}

		if ctx.Err() != nil {
			out.Error = ctxCancelErr
			log.Error(err)
			return out
		}

		out.Code = fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		out.Error = serverErr
		log.Error(err)
		return out
	}

	out.Code = httpResp.Status
	if httpResp.StatusCode != http.StatusOK {
		out.Error = "non_200_resp_status"
		log.Error(err)
		return out
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&out.Body); err != nil {
		out.Error = jsonDecodeErr
		log.Error(err)
		return out
	}

	return out
}

func validURL(u string) error {
	rawUrl, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}
	scheme := rawUrl.Scheme

	if scheme != "https" && scheme != "http" {
		log.Debug(scheme)
		return errors.New(fmt.Sprintf("unsupported protocol scheme: %s", scheme))
	}
	return nil
}
