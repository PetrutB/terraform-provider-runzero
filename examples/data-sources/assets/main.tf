data "runzero_assets" "all" {}

data "runzero_assets" "windows" {
  search = "os:windows"
}

output "total_assets" {
  value = length(data.runzero_assets.all.assets)
}

output "windows_assets" {
  value = length(data.runzero_assets.windows.assets)
}
