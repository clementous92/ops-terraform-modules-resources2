variable "project_id" {
  type        = string
  description = "GCP project ID"
}

variable "network_name" {
  description = "The name of the network where subnets will be created"
  type        = string
}

variable "subnets" {
  type        = list(map(string))
  description = "The list of subnets being created"
}

# example for subnets lists
#  subnets = [
#         {
#             subnet_name               = "subnet-02"
#             subnet_ip                 = "10.10.20.0/24"
#             subnet_region             = "us-west1"
#             subnet_private_access     = "true"
#             subnet_flow_logs          = "true"
#             subnet_flow_logs_interval = "INTERVAL_10_MIN"
#             subnet_flow_logs_sampling = 0.7
#             subnet_flow_logs_metadata = "INCLUDE_ALL_METADATA"
#             description               = "This subnet has a description"
#         }
#     ]
