# -*- mode: ruby -*-
# vi: set ft=ruby :

$vm_memory ||= 2048
$vm_cpus ||= 2
$vm_box ||= "ubuntu/focal64"

$provision_dependency ||= <<-SHELL
  apt update
  apt install -y \
    bison build-essential cmake flex git python wget \
    libedit-dev libllvm7 llvm-7-dev libclang-7-dev \
    zlib1g-dev libelf-dev libfl-dev \
    luajit libluajit-5.1-dev \
    netperf arping iperf3
SHELL
$provision_go ||= <<-SHELL
  wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
  echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc
  mkdir go
  rm go1.16.3.linux-amd64.tar.gz #Clean up directory
SHELL
$provision_bcc ||= <<-SHELL
  git clone https://github.com/iovisor/bcc.git
  mkdir bcc/build; cd bcc/build
  cmake ..
  make
  sudo make install
SHELL
$provision_bcc_python3_binding ||= <<-SHELL
  cd bcc/build
  cmake -DPYTHON_CMD=python3 ..
  pushd src/python/
  make
  sudo make install
  popd
SHELL

Vagrant.require_version ">= 2.0.0"
Vagrant.configure("2") do |config|
  config.vm.provider "virtualbox" do |vb|
		vb.memory = $vm_memory
		vb.cpus = $vm_cpus
	end
  config.vm.box = $vm_box
  config.vm.synced_folder "./dev", "/home/vagrant/dev", create: "true"

  # provision
  # Comment out as appropriate according to your plan.

  # provision_dependency: Install dependent packages.
  # If you do not use bcc python3 binding, comment out the following package.
  # luajit libluajit-5.1-dev netperf arping iperf3
  config.vm.provision :shell, inline: $provision_dependency

  # provision_go: Install golang
  # golang ver1.16.3
  config.vm.provision :shell, privileged: false, inline: $provision_go

  # provision_bcc: Build bcc
  config.vm.provision :shell, privileged: false, inline: $provision_bcc
  # If you do not use bcc python3 binding, comment out this provision.
  config.vm.provision :shell, privileged: false, inline: $provision_bcc_python3_binding
end
