terraform {
  required_providers {
    runzero = {
      source  = "petrutb/runzero"
      version = "~> 0.1.0"
    }
  }
}

provider "runzero" {
  # api_url    = "https://console.runzero.com/api/v1.0"  # or set RUNZERO_API_URL
  # org_token  = "..."  # or set RUNZERO_ORG_TOKEN
  # account_token = "..."  # or set RUNZERO_ACCOUNT_TOKEN

  # OAuth2 (alternative to tokens):
  # client_id     = "..."  # or set RUNZERO_CLIENT_ID
  # client_secret = "..."  # or set RUNZERO_CLIENT_SECRET
}
