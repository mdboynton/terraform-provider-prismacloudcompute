resource "prismacloudcompute_collection" "test" {
    name = "testcollection"
    description = "example collection to test Prisma Cloud Compute Terraform provider"
    account_ids = ["*dev"]
    containers = ["foo", "bar"]
    hosts = ["nginx", "envoy"]
    images = ["*docker*"]
}
