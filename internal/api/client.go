package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

func (c *PrismaCloudComputeAPIClient) Request(method, endpoint string, query, data, response interface{}) (err error) {
    // Parse console URL from config
    consoleUrl, err := url.Parse(*c.Config.ConsoleURL)
	if err != nil {
		return err
	}

    // Set URL scheme to HTTPS if undefined
	if consoleUrl.Scheme == "" {
		consoleUrl.Scheme = "https"
	}

    // Append endpoint to URL"
	consoleUrl.Path = path.Join(consoleUrl.Path, endpoint)

    // Marshal request payload into buffer, if not nil
	var buf bytes.Buffer
	if data != nil {
		data_json, err := json.Marshal(data)
		if err != nil {
			return err
		}

		buf = *bytes.NewBuffer(data_json)
	}

    // Create new HTTP request object
	req, err := http.NewRequest(method, consoleUrl.String(), &buf)
	if err != nil {
		return err
	}

    // Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	req.Header.Set("Content-Type", "application/json")

    if query != nil {
		queryParams := req.URL.Query()
		if queryMap, ok := query.(map[string]string); ok {
			for key, val := range queryMap {
				queryParams.Add(key, val)
			}
		}
		req.URL.RawQuery = queryParams.Encode()
	}

    // Execute request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If API responds with HTTP 429 (Too Many Requests), sleep 3 seconds and try again
	if res.StatusCode == 429 {
		time.Sleep(3 * time.Second)
		return c.Request(method, endpoint, query, data, &response)
	}

    // If API responds with a non-OK status, return error
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body from non-OK response: %s", err)
		}

		var response ErrResponse
		if err = json.Unmarshal(body, &response); err != nil {
			return err
		}

		return fmt.Errorf("Non-OK status: %d (%s)", res.StatusCode, response.Err)
	}

    // Parse response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

    // If response body is non-empty, unmarshal into response object 
	if len(body) > 0 && response != nil {
		if err = json.Unmarshal(body, response); err != nil {
			return err
		}
	}

	return nil
}

func (c *PrismaCloudComputeAPIClient) Authenticate() (err error) {
	res := AuthResponse{}
	if err := c.Request(http.MethodPost, "api/v1/authenticate", nil, AuthRequest{*c.Config.Username, *c.Config.Password}, &res); err != nil {
		return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: %v", err)
	}
	c.JWT = res.Token

	return nil
}

// Create Client and authenticate.
func Client(config PrismaCloudComputeAPIClientConfig) (*PrismaCloudComputeAPIClient, error) {
	apiClient := &PrismaCloudComputeAPIClient{
		Config: config,
	}

	if config.Insecure {
		apiClient.HTTPClient = &http.Client{
            Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		apiClient.HTTPClient = &http.Client{
            Timeout: 60 * time.Second,
        }
	}

	if err := apiClient.Authenticate(); err != nil {
		return nil, err
	}

	return apiClient, nil
}
