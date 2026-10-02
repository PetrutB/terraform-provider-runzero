# Adopt an existing RunZero agent into Terraform management
resource "runzero_agent" "scanner" {
  id   = "00000000-0000-0000-0000-000000000000"
  name = "Primary Scanner"
}
