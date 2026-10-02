resource "runzero_api_key" "example" {
  name   = "ci-pipeline-key"
  org_id = runzero_organization.example.id
}

output "api_key" {
  value     = runzero_api_key.example.key
  sensitive = true
}
