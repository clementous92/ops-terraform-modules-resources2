output "name_servers" {
  value = element(
    concat(
      google_dns_managed_zone.private.*.name_servers,
      google_dns_managed_zone.public.*.name_servers,
    ),
    0,
  )
}

output "name" {
  value = var.name
}

output "dns_name" {
  value = var.dns_name
}

output "output" {
  value = element(
    concat(
      google_dns_managed_zone.private,
      google_dns_managed_zone.public,
    ),
    0,
  )
}
