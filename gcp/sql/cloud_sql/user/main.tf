
resource "google_sql_user" "sql_user" {
  name     = var.username
  instance = var.instance_name
  host     = var.host
  password = var.password
  project  = var.project_id
}
