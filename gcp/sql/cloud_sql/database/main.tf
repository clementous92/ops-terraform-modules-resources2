resource "google_sql_database" "database" {
  name     = var.name
  instance = var.instance
  project  = var.project_id
  charset  = var.charset
}
