# Block specific paths from being accessed
resource "sonatyperepo_routing_rule" "block_example" {
  name        = "block-example-paths"
  description = "Block requests matching example.com paths"
  mode        = "BLOCK"
  matchers = [
    "^/com/example/.*",
    "^/org/example/.*"
  ]
}

# Allow only specific approved paths
resource "sonatyperepo_routing_rule" "allow_approved" {
  name        = "allow-approved-only"
  description = "Allow only approved organization paths"
  mode        = "ALLOW"
  matchers = [
    "^/com/approved/.*",
    "^/org/approved/.*",
    "^/io/approved/.*"
  ]
}