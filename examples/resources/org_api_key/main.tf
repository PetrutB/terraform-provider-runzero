resource "runzero_org_api_key" "current" {
  # Change this value to rotate the organization API key.
  rotate = "2026-01-01"
}

output "org_api_key_token" {
  value     = runzero_org_api_key.current.token
  sensitive = true
}
