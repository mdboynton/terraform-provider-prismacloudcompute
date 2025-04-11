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
    "math"
)

type PrismaCloudComputeAPIClientConfig struct {
    ConsoleURL      *string `tfsdk:"console_url" json:"console_url"`
    Username        *string `tfsdk:"username" json:"username"`
    Password        *string `tfsdk:"password" json:"password"`
    Insecure        *bool   `tfsdk:"insecure" json:"insecure"`
    RequestTimeout  *int    `tfsdk:"request_timeout" json:"request_timeout"`
    ConfigFile      *string `tfsdk:"config_file" json:"config_file"`
}

type PrismaCloudComputeAPIClient struct {
	Config     PrismaCloudComputeAPIClientConfig
	HTTPClient *http.Client
	JWT        string
}

type ErrResponse struct {
	Err string
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func Client(config PrismaCloudComputeAPIClientConfig) (*PrismaCloudComputeAPIClient, error) {
	apiClient := &PrismaCloudComputeAPIClient{
		Config: config,
	}

    // Parse request timeout value
    if config.RequestTimeout == nil {
        defaultTimeout := 60
        config.RequestTimeout = &defaultTimeout
    } else if *config.RequestTimeout > math.MaxInt {
        return nil, fmt.Errorf("Error occured while creating API client: Invalid value supplied for request_timeout. Value must be an integer between 1 and %d.", math.MaxInt)
    }

    requestTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", *config.RequestTimeout))
    if err != nil {
        return nil, fmt.Errorf("Error occured while creating API client: Failed to parse request timeout value\n%s", err.Error())
    }

    //            //    fmt.Sprintf("Error configuring provider: Invalid value specified for \"request_timeout\" in configuration file. Value must be an integer between 1 and %d", math.MaxInt),

    // Instantiate HTTP client
    httpClient := &http.Client{
        Timeout: requestTimeout,
    }

    // If the insecure flag is set to true, add TLS configuration with InsecureSkipVerify enabled 
    if (config.Insecure != nil && *config.Insecure) {
        transport := http.Transport{
		    TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
        (*httpClient).Transport = &transport
    }

    apiClient.HTTPClient = httpClient

    // Authenticate to API
	if err := apiClient.Authenticate(); err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c *PrismaCloudComputeAPIClient) Authenticate() (err error) {
	res := AuthResponse{}

    if c == nil {
        return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: client uninitialized")
    }

    if c.Config.ConsoleURL== nil {
        return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: nil console URL")
    }

    if c.Config.Username == nil {
        return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: nil username")
    }

    if c.Config.Password == nil {
        return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: nil password")
    }
    
	if err := c.Request(http.MethodPost, "api/v1/authenticate", nil, AuthRequest{*c.Config.Username, *c.Config.Password}, &res); err != nil {
		return fmt.Errorf("Error occured while authenticating to Prisma Cloud Compute API: %v", err)
	}
	c.JWT = res.Token

	return nil
}

func (c *PrismaCloudComputeAPIClient) Request(method, endpoint string, query, data, response interface{}) (err error) {
    // Parse console URL from config
    consoleUrl, err := url.Parse(*c.Config.ConsoleURL)
	if err != nil {
		return err
	}

    // Append endpoint to URL
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

    // Add query parameters to endpoint URL, if any are provided
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
