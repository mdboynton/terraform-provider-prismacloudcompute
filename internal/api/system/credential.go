package system

//import (
//	"fmt"
//	"net/http"
//    "sort"
//    "strconv"
//
//    "github.com/hashicorp/terraform-plugin-framework/diag"
//
//	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
//)

//const CredentialsEndpoint = "api/v1/credentials"
//
//type Credential struct {
//    Id      string  `json:"_id"`
//    AccountID string `json:"accountID"`
//    AccountName string `json:"accountName"`
//    ApiToken Secret `json:"apiToken"`
//    CaCert string `json:"caCert"`
//    Description string `json:"description"`
//    Secret Secret `json:"secret"`
//    SkipVerify bool `json:"skipVerify"`
//    Provider string `json:"type"`
//    UseAwsRole bool `json:"useAWSRole"`
//    UseStsRegionalEndpoint bool `json:"useSTSRegionalEndpoint"`
//}
//
//type Secret struct {
//    Encrypted string `json:"encrypted"`
//    Plain string `json:"plain"`
//}
