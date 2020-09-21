# Configure the dynatrace provider
provider "dynatrace" {
    dt_env_url    = var.environment 
    dt_api_token  = var.token
}

data "dynatrace_alerting_profiles" "keptn" {
  id = var.alerting_profile_id
}

resource "dynatrace_alerting_profiles" "sockshop_errors" {

  display_name = "sockshop_errors"
  mz_id = dynatrace_management_zones.sockshop_prod.id

  rule{
    severity_level = "AVAILABILITY"
    tag_filters {
      include_mode = "INCLUDE_ALL"
      tag_filter {
        context = "ENVIRONMENT"
        key = "product"
        value = "sockshop"
      }
    }
    delay_in_minutes = 2
  }

  rule{
    severity_level = "ERROR"
    tag_filters {
      include_mode = "INCLUDE_ALL"
      tag_filter {
        context = "CONTEXTLESS"
        key = "app"
        value = "carts"
      }
      tag_filter {
        context = "CONTEXTLESS"
        key = "env"
        value = "prod"
      }
    }
    delay_in_minutes = 1
  }

  event_type_filter{
    predefined_event_filter{
      negate = true
      event_type = "EC2_HIGH_CPU"
    }
  }

  event_type_filter{
    custom_event_filter{
      custom_title_filter{
        enabled = true
        value = "sockshop"
        operator = "CONTAINS"
        negate = true
        case_insensitive = false
      }
    }
  }

  event_type_filter{
    custom_event_filter{
      custom_description_filter{
        enabled = true
        value = "carts"
        operator = "CONTAINS"
        negate = false
        case_insensitive = true
      }
    }
  }
    
}