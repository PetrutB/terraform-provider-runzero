data "runzero_org" "current" {}

output "org_name" {
  value = data.runzero_org.current.name
}
