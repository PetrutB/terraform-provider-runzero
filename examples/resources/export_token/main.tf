resource "runzero_export_token" "example" {
  org_id = runzero_organization.example.id
}

output "export_token_key" {
  value     = runzero_export_token.example.key
  sensitive = true
}
