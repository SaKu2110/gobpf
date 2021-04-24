# -*- mode: ruby -*-
# vi: set ft=ruby :

$vm_box ||= "ubuntu/focal64"
$vm_provision ||= <<-SHELL
  apt update
  apt install -y wget
  wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
  tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
  mkdir /home/vagrant/go
  echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc 
SHELL

Vagrant.require_version ">= 2.0.0"
Vagrant.configure("2") do |config|
  config.vm.box = $vm_box
  config.vm.synced_folder "./dev", "/home/vagrant/dev", create: "true"
  config.vm.provision :shell, inline: $vm_provision
end
