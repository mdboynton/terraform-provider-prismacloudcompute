package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type PrismaCloudComputeAPIClientConfig struct {
	ConsoleURL           string `tfsdk:"console_url"`
	//Project              string `tfsdk:"project"`
	Username             string `tfsdk:"username"`
	Password             string `tfsdk:"password"`
	Insecure             bool   `tfsdk:"insecure"`
	//ConfigFile           string `tfsdk:"config_file"`
}

type PrismaCloudComputeAPIClient struct {
	Config     PrismaCloudComputeAPIClientConfig
	HTTPClient *http.Client
	JWT        string
}

type ErrResponse struct {
	Err string
}

func (c *PrismaCloudComputeAPIClient) Initialize(filename string) error {
    // TODO: add logic to re-use API token if its still valid, instead of authing every time
	c2 := PrismaCloudComputeAPIClient{}

	if filename != "" {
		var (
			b   []byte
			err error
		)

		b, err = ioutil.ReadFile(filename)

		if err != nil {
			return err
		}

		if err = json.Unmarshal(b, &c2); err != nil {
			return err
		}
	}

	if c.Config.ConsoleURL == "" && c2.Config.ConsoleURL != "" {
		c.Config.ConsoleURL = c2.Config.ConsoleURL
	}

	//if c.Config.Project == "" && c2.Config.Project != "" {
	//	c.Config.Project = c2.Config.Project
	//}

	if c.Config.Username == "" && c2.Config.Username != "" {
		c.Config.Username = c2.Config.Username
	}

	if c.Config.Password == "" && c2.Config.Password != "" {
		c.Config.Password = c2.Config.Password
	}

	c.HTTPClient = &http.Client{}

	return c.Authenticate()
}

func (c *PrismaCloudComputeAPIClient) Request(method, endpoint string, query, data, response interface{}) (err error) {
    parsedURL, err := url.Parse(c.Config.ConsoleURL)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}
	parsedURL.Path = path.Join(parsedURL.Path, endpoint)
    //fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")
    //fmt.Println("sending " + method + " request to")
    //fmt.Println(parsedURL)
    //fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")

	var buf bytes.Buffer

	if data != nil {
		data_json, err := json.Marshal(data)
		if err != nil {
			return err
		}
		buf = *bytes.NewBuffer(data_json)
	}

	req, err := http.NewRequest(method, parsedURL.String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.JWT)
	req.Header.Set("Content-Type", "application/json")

	// TODO: simplify logic
    //if c.Config.Project != "" {
	//	queryParams := req.URL.Query()
	//	//queryParams.Set("project", c.Config.Project)
	//	if query != nil {
	//		if queryMap, ok := query.(map[string]string); ok {
	//			for key, val := range queryMap {
	//				queryParams.Add(key, val)
	//			}
	//		}
	//	}
	//	req.URL.RawQuery = queryParams.Encode()
	//} else if query != nil {
    if query != nil {
		queryParams := req.URL.Query()
		if queryMap, ok := query.(map[string]string); ok {
			for key, val := range queryMap {
				queryParams.Add(key, val)
			}
		}
		req.URL.RawQuery = queryParams.Encode()
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Retry in case backend responds with HTTP 429
	// sleep for 3 seconds before retry
	if res.StatusCode == 429 {
		time.Sleep(3 * time.Second)
		return c.Request(method, endpoint, query, data, &response)
	}
    fmt.Println(res.StatusCode)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

    fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")
    fmt.Println("unmarshalling response body")
    fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")
	if len(body) > 0 && response != nil {
		if err = json.Unmarshal(body, response); err != nil {
			return err
		}
        
        if endpoint != "/api/v1/authenticate" && endpoint != "api/v1/static/vulnerabilities" {
            fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")
            fmt.Println("recieved response from endpoint: ")
            fmt.Printf("%+v\n", response)
            fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&")
        }
	}
	return nil
}

func (c *PrismaCloudComputeAPIClient) Authenticate() (err error) {
	type AuthRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type AuthResponse struct {
		Token string `json:"token"`
	}

	res := AuthResponse{}
	if err := c.Request(http.MethodPost, "api/v1/authenticate", nil, AuthRequest{c.Config.Username, c.Config.Password}, &res); err != nil {
		return fmt.Errorf("error POSTing to authenticate endpoint: %v", err)
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
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		apiClient.HTTPClient = &http.Client{}
	}

	if err := apiClient.Authenticate(); err != nil {
		return nil, err
	}

	return apiClient, nil
}
