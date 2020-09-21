# dynatrace_management_zones Resource

Provides a dynatrace management resource. It allows to create, update, delete management zones in a dynatrace environment. [Management Zones API]

## Example Usage

```hcl
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
          "key" = "project"
          "value" = "carts"
        }
        negate = false
      }
    }
  }

}
```

## Argument Reference

* `name` - (Required) The name of the management zone..
* `rule` - (Optional) The ID of the management zone to which the alerting profile applies.
* `rule` - (Optional) A nested block that contains a list of rules for management zone usage. Each rule is evaluated independently of all other rules. See Nested rule block below for details.

## Attribute Reference

* `id` - The ID of the management zone.

## Nested rule block

* `type` - (Required) The type of Dynatrace entities the management zone can be applied to.
* `enabled` - (Required) The rule is enabled (true) or disabled (false).
* `propagation_types` - (Optional) How to apply the management zone to underlying entities.
* `condition` - (Required) A list of matching rules for the management zone. The management zone applies only if all conditions are fulfilled.
    * `key` - (Required) The key to identify the data we're matching."
        * `attribute` - (Required) The attribute to be used for comparision.
        * `type` - (Optional) Defines the actual set of fields depending on the value.
    * `comparison_info` (Required) Defines how the matching is actually performed: what and how are we comparing.
        * `operator` - (Required) Operator of the comparison. You can reverse it by setting negate to true. Possible values depend on the type of the comparison. Find the list of actual models in the description of the type field and check the description of the model you need.
        * `value` - (Optional) The value to compare to.
        * `negate` - (Required) Reverses the comparison operator. For example it turns the begins with into does not begin with.
        * `type` - (Required) Defines the actual set of fields depending on the value.

## Import

Dynatrace management zones can be imported using their ID, e.g.

```hcl
$ terraform import dynatrace_management_zones.keptn-carts -4638826838889583423
```

[Management Zones API]: (https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/management-zones-api/)
