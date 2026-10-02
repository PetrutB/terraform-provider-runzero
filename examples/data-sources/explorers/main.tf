data "runzero_explorers" "all" {}

output "connected_explorers" {
  value = [for e in data.runzero_explorers.all.explorers : e.name if e.connected]
}
