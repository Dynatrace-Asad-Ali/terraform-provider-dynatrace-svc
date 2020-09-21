output "ap" {
  value = data.dynatrace_alerting_profiles.keptn
}

output "ap_sockshop_errors" {
  value = dynatrace_alerting_profiles.sockshop_errors
}
