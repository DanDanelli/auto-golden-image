# auto-golden-image
Golden image created from packer with ansible installing apache and terraform running an ec2 from this image (used Go as wrapper)<br/>
1 - Don't forget to export the programmatic variables first<br/>
2 - Make sure the subnet you're informing is linked to a VPC with Internet Gateway, Route Table and a Security Group with internet access<br/>
3 - To access the instance created: {instance_public_ip}:8081<br/>
