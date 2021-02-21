variable "secondary_ranges" {
  type        = map(list(object({ range_name = string, ip_cidr_range = string })))
  description = "Secondary ranges that will be used in some of the subnets"
  default     = {}
}

# Example for secondary ranges input
# secondary_ranges = {
#     subnet-01 = [
#         {
#             range_name    = "subnet-01-secondary-01"
#             ip_cidr_range = "192.168.64.0/24"
#         },
#     ]

#     subnet-02 = []
# }


