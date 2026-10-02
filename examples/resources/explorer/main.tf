# Adopt an existing RunZero explorer into Terraform management
resource "runzero_explorer" "primary" {
  id   = "00000000-0000-0000-0000-000000000000"
  name = "Primary Explorer"
}
