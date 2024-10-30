package api

import (
	"net/http"
)

type PrismaCloudComputeAPIClientConfig struct {
    ConsoleURL  *string `tfsdk:"console_url" json:"console_url"`
    Username    *string `tfsdk:"username" json:"username"`
    Password    *string `tfsdk:"password" json:"password"`
    Insecure    bool   `tfsdk:"insecure" json:"insecure"`
    ConfigFile  *string `tfsdk:"config_file" json:"config_file"`
	//Project   string `tfsdk:"project"`
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
