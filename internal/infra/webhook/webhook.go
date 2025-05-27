package webhook

import (
	"io"
	"net/http"

	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type DefaultClient struct{}

func (d *DefaultClient) Post(url, contentType string, body io.Reader) (*contracts.HTTPResponse, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}

	return &contracts.HTTPResponse{
		StatusCode: resp.StatusCode,
		Close:      resp.Body.Close,
	}, nil
}
