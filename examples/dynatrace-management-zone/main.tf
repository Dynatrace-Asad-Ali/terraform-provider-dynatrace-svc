# Configure the dynatrace provider
provider "dynatrace" {
    dt_env_url    = var.environment 
    dt_api_token  = var.token
}

resource "dynatrace_management_zones" "sockshop_prod" {

  name = "sockshop_prod"

  rule{
    type = "SERVICE"
    enabled = true
    propagation_types = [
      "SERVICE_TO_HOST_LIKE", 
      "SERVICE_TO_PROCESS_GROUP_LIKE"
    ]
    condition {
      key {
        attribute = "SERVICE_TAGS"
      }
      comparison_info {
        type = "TAG"
        operator = "EQUALS"
        value = {
          "context" = "CONTEXTLESS"
          "key" = "project"
          "value" = "sockshop"   
        }
        negate = false
      }
    }
    condition {
      key {
        attribute = "SERVICE_TAGS"
      }
      comparison_info {
        type = "TAG"
        operator = "EQUALS"
        value = {
          "context" = "CONTEXTLESS"
          "key" = "app"
          "value" = "carts"
        }
        negate = false
      }
    }
  }

  rule{
    type = "SERVICE"
    enabled = true
    propagation_types = [
      "SERVICE_TO_HOST_LIKE", 
      "SERVICE_TO_PROCESS_GROUP_LIKE"
    ]
    condition {
      key {
        attribute = "SERVICE_TAGS"
      }
      comparison_info {
        type = "TAG"
        operator = "EQUALS"
        value = {
          "context" = "CONTEXTLESS"
          "key" = "env"
          "value" = "prod"
        }
        negate = false
      }
    }
  }

}