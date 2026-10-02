resource "runzero_site" "example" {
  name        = "Production Network"
  description = "Production environment scan site"
  scope       = "10.0.0.0/8"
  excludes    = "10.0.0.1"
}

output "site_id" {
  value = runzero_site.example.id
}
