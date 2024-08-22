provider "aws" {
  region     = var.region
  access_key = "test"
  secret_key = "test"


  endpoints {
    ec2 = "http://localhost:4566"
    ecs = "http://localhost:4566"
  }

  default_tags {
    tags = {
      Environment = "Local"
      Service     = "LocalStack"
      Company     = var.company_name
    }
  }
}