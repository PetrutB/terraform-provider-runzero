resource "runzero_user" "example" {
  email      = "engineer@example.com"
  first_name = "Jane"
  last_name  = "Doe"
  org_id     = runzero_organization.example.id
}
