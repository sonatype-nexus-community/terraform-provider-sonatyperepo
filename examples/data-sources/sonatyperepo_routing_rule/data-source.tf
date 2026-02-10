data "sonatyperepo_routing_rule" "example" {
  name = "my-routing-rule"
}

output "routing_rule" {
  value = data.sonatyperepo_routing_rule.example
}
