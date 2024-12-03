package provider

//import (
//	//"fmt"
//	//"os"
//	//"testing"
//
//	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
//	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//
//    "github.com/hashicorp/terraform-plugin-framework/providerserver"
//    "github.com/hashicorp/terraform-plugin-go/tfprotov6"
//)
//
//const (
//	PrismacloudcomputeJsonConfigFileEnvVar = "PRISMACLOUDCOMPUTE_CONFIG_FILE"
//)
//
//var (
//	//testAccProviders                   map[string]*schema.Provider
//	//TestAccProviders                   map[string]*schema.Provider
//	//testAccProvider                    *schema.Provider
//	//TestAccProvider                    *schema.Provider
//	//sessionTimeoutOrig, sessionTimeout int
//
//    //testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
//    TestAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
//    //TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
//    //    "prismacloudcompute": providerserver.NewProtocol6WithError(New("test")()),
//    //}
//)
//
//func init() {
//    TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
//        "prismacloudcompute": providerserver.NewProtocol6WithError(New("test")()), 
//    }
//}

//var ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error) = protoV6ProviderFactoriesInit(context.Background(), "prismacloudcompute")
//
//func protoV6ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
//	factories := make(map[string]func() (tfprotov6.ProviderServer, error), len(providerNames))
//
//	for _, name := range providerNames {
//
//		factories[name] = func() (tfprotov6.ProviderServer, error) {
//			providerServerFactory, err := provider.ProviderServerFactoryV6(ctx, getProviderTestingVersion())
//
//			if err != nil {
//				return nil, err
//			}
//
//			return providerServerFactory(), nil
//		}
//	}
//
//	return factories
//}

//func init() {
//	var err error
//
//	//testAccProvider = Provider()
//	TestAccProvider = Provider()
//	//testAccProviders = map[string]*schema.Provider{
//	TestAccProviders = map[string]*schema.Provider{
//		"prismacloudcompute": testAccProvider,
//	}
//
//	client := &api.Client{}
//	if err = client.Initialize(os.Getenv(PrismacloudcomputeJsonConfigFileEnvVar)); err == nil {
//		if err != nil {
//			fmt.Printf("Error initializing client")
//		}
//	}
//}
//
//func TestProvider(t *testing.T) {
//	if err := Provider().InternalValidate(); err != nil {
//		t.Fatalf("err: %s", err)
//	}
//}
//
//func TestProvider_impl(t *testing.T) {
//	var _ *schema.Provider = Provider()
//}
//
//func testAccPreCheck(t *testing.T) {
//	if os.Getenv(PrismacloudcomputeJsonConfigFileEnvVar) == "" {
//		t.Fatalf("%s must be set for acceptance tests", PrismacloudcomputeJsonConfigFileEnvVar)
//	}
//}
