
# must have roles/resourcemanager.folderCreator.

resource "google_folder" "org_folder" {
  display_name = var.folder_name
  parent       = var.parent
}
