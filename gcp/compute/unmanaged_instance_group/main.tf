# Creates a group of dissimilar Compute Engine virtual machine instances.
resource "google_compute_instance_group" "instance_group" {
  project   = var.project_id
  name      = var.group_name
  zone      = var.zone
  instances = var.instances

  named_port {
    name = var.port_name
    port = var.port
  }

  # Recreating an instance group that's in use by another resource will give a
  # resourceInUseByAnotherResource error.
  # lifecycle.create_before_destroy will avoid this type of error
  lifecycle {
    create_before_destroy = true
  }
}
