# Adopt a discovered asset and manage its metadata
resource "runzero_asset" "webserver" {
  id          = "00000000-0000-0000-0000-000000000000"
  comments    = "Production web server"
  criticality = "high"
}
