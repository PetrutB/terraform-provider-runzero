resource "runzero_credential" "example" {
  name   = "scanner-cred"
  type   = "snmpv2_community"
  global = true
  secret = {
    community = "public"
  }
}
