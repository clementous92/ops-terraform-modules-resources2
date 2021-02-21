output "output" {
  value = element(
    concat(
      google_compute_ssl_certificate.cert,
      google_compute_region_ssl_certificate.cert,
    ),
    0,
  )
}
