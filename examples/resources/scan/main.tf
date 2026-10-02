# Trigger a scan on a site
resource "runzero_scan" "network" {
  site_id = runzero_site.production.id
  name    = "Weekly Network Scan"
  targets = ["10.0.0.0/24", "192.168.1.0/24"]
}
