terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-east-1"
}

resource "aws_dynamodb_table" "users_dynamo_table" {
  name = "users"

  hash_key = "userID"

  write_capacity = 5
  read_capacity  = 5

  attribute {
    name = "userID"
    type = "S"
  }

  attribute {
    name = "login"
    type = "S"
  }

  attribute {
    name = "email"
    type = "S"
  }

  global_secondary_index {
    hash_key        = "login"
    name            = "LoginIndex"
    projection_type = "ALL"
    write_capacity = 5
    read_capacity = 5
  }

  global_secondary_index {
    hash_key        = "email"
    name            = "EmailIndex"
    projection_type = "ALL"
    write_capacity = 5
    read_capacity = 5
  }
}