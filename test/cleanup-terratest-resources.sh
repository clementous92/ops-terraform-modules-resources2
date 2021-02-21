#!/bin/bash

echo "Deleting terratest projects"
gcloud projects list --filter=parent.id=976583563296 --format="value(PROJECT_ID)"  | grep terratest | xargs -I {}  gcloud projects delete {} --quiet

echo "Deleting terratest Cloud SQL"
gcloud sql instances list --project ami-doit-playground | grep terra | awk '{print $1}' | xargs -I {} gcloud sql instances delete --project ami-doit-playground {}
