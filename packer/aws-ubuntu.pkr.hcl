packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}
variable "ami-name" {  
  type = string  
  default = "linux-golden_image-from-packer"
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "${var.ami-name}"
  instance_type = "m3.medium"
  region        = "eu-west-1"
  vpc_id        = "vpc-005907dca31d725f6"
  subnet_id     = "subnet-0b0fa07c52ffc4553"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

build {
  name = "learn-packer"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

provisioner "shell" {       
  inline = [        
    "sudo apt-get update"      
  ]  
}

provisioner "ansible" {      
  playbook_file = "./ansible/playbook.yml"    
}
 
 post-processor "manifest" {
  output = "./packer/manifest.json"
  strip_path = true    
  custom_data = {
   my_custom_data = "example"    
  }
 }
}
