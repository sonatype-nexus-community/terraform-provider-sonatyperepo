data "sonatyperepo_routing_rules" "all" {
}

output "routing_rules" {
  value = data.sonatyperepo_routing_rules.all.routing_rules
}
