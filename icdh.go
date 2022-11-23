// Package icdh is a Go client library for NGINX Ingress Controller Deep Health API.
package icdh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type option func(*client) error

// WithHTTPClient is an option for setting
// up a custom http client.
func WithHTTPClient(h *http.Client) option {
	return func(c *client) error {
		if h == nil {
			return errors.New("nil http client")
		}
		c.httpClient = h
		return nil
	}
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new health service client.
func NewClient(path string, opts ...option) (*client, error) {
	c := client{
		baseURL:    path,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, fmt.Errorf("creating a client: %w", err)
		}
	}
	return &c, nil
}

// Stats holds statistics for services
// associated with the hostname.
type Stats struct {
	Total     int
	Up        int
	Unhealthy int
}

type healthServiceResp struct {
	Total     int
	Up        int
	Unhealthy int
}

// GetStats takes hostane and returns health service stats
// for all upstreams associated with the given hostname.
func (c client) GetStats(ctx context.Context, hostname string) (Stats, error) {
	u := fmt.Sprintf("%s/probe/%s", c.baseURL, hostname)
	var h healthServiceResp
	if err := c.get(ctx, u, &h); err != nil {
		return Stats{}, fmt.Errorf("retriving stats for host %s, %w", hostname, err)
	}
	s := Stats(h)
	return s, nil
}

func (c client) get(ctx context.Context, url string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response code: %v", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("decoding response body: %w", err)
	}
	return nil
}
