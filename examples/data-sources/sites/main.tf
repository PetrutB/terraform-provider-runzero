data "runzero_sites" "all" {}

output "site_names" {
  value = [for s in data.runzero_sites.all.sites : s.name]
}
