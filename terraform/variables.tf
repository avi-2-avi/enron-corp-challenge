variable "region" {
  description = "AWS region to create resources in"
  default = "us-east-1"
}

variable "availability_zone" {
  description = "AWS availability zone to create resources in"
}

variable "company_name" {
  description = "Company name to identify resources created"
  default = "Enron Corp"
}