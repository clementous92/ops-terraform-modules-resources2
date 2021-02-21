# Collet rules for allow and rule for deny from input and map them to custom_allow and custom_deny
locals {
  rules-allow = {
    for name, attrs in var.custom_rules : name => attrs if attrs.action == "allow"
  }
  rules-deny = {
    for name, attrs in var.custom_rules : name => attrs if attrs.action == "deny"
  }
}

resource "google_compute_firewall" "custom_allow" {
  provider                = google-beta
  for_each                = local.rules-allow
  name                    = each.key
  description             = each.value.description
  direction               = each.value.direction
  network                 = var.network
  project                 = var.project_id
  source_ranges           = each.value.direction == "INGRESS" ? each.value.ranges : null
  destination_ranges      = each.value.direction == "EGRESS" ? each.value.ranges : null
  source_tags             = each.value.use_service_accounts || each.value.direction == "EGRESS" ? null : each.value.sources
  source_service_accounts = each.value.use_service_accounts && each.value.direction == "INGRESS" ? each.value.sources : null
  target_tags             = each.value.use_service_accounts ? null : each.value.targets
  target_service_accounts = each.value.use_service_accounts ? each.value.targets : null
  disabled                = lookup(each.value.extra_attributes, "disabled", false)
  priority                = lookup(each.value.extra_attributes, "priority", 1000)
  # enable_logging          = lookup(each.value.extra_attributes, "enable_logging", false)
  dynamic "allow" {
    for_each = each.value.rules
    iterator = rule
    content {
      protocol = rule.value.protocol
      ports    = rule.value.ports
    }
  }
}

resource "google_compute_firewall" "custom_deny" {
  provider                = google-beta
  for_each                = local.rules-deny
  name                    = each.key
  description             = each.value.description
  direction               = each.value.direction
  network                 = var.network
  project                 = var.project_id
  source_ranges           = each.value.direction == "INGRESS" ? each.value.ranges : null
  destination_ranges      = each.value.direction == "EGRESS" ? each.value.ranges : null
  source_tags             = each.value.use_service_accounts || each.value.direction == "EGRESS" ? null : each.value.sources
  source_service_accounts = each.value.use_service_accounts && each.value.direction == "INGRESS" ? each.value.sources : null
  target_tags             = each.value.use_service_accounts ? null : each.value.targets
  target_service_accounts = each.value.use_service_accounts ? each.value.targets : null
  disabled                = lookup(each.value.extra_attributes, "disabled", false)
  priority                = lookup(each.value.extra_attributes, "priority", 1000)
  # enable_logging          = lookup(each.value.extra_attributes, "enable_logging", false)

  dynamic "deny" {
    for_each = each.value.rules
    iterator = rule
    content {
      protocol = rule.value.protocol
      ports    = rule.value.ports
    }
  }
}
