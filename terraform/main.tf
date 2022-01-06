terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      /*version = "~> 3.17.0"*/
    }
    dotenv = {
      source  = "jrhouston/dotenv"
      version = "~> 1.0"
    }
  }
}

module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 3.0"

  name = "single-instance"

  ami                    = var.ami
  instance_type          = "t2.micro"
  key_name               = "kp_lab"
  monitoring             = true
  vpc_security_group_ids = ["sg-0bbd5c6b35d23a7d9"]
  subnet_id              = "subnet-0b0fa07c52ffc4553"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

#via export TF_VAR
 output "ami_id" {
   value = var.ami
 }