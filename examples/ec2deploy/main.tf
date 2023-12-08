terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.30.0"
    }
    provisioner = {
      source  = "terraform-registry.coding-ia.com/coding-ia/provisioner"
      version = "1.0.0-alpha-1"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}

provider "provisioner" {
  sns_topic = "arn:aws:sns:us-east-2:809674927168:PublishDNS"
  region    = "us-east-2"
}

variable "instance_name" {
  description = "Name of the EC2 instance"
  default     = "example-instance"
}

resource "aws_instance" "example" {
  ami           = "ami-0e83be366243f524a"
  instance_type = "t2.micro"

  tags = {
    Name = "example-instance"
  }
}

resource "provisioner_provision" "test" {
  name        = var.instance_name
  instance_id = aws_instance.example.id
  private_ip  = aws_instance.example.private_ip
}