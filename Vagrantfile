Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"

  # Set up a synced folder (host to guest)
  config.vm.synced_folder ".", "/home/vagrant/app"

  # Optional: Set resources
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
    vb.cpus = 2
  end

  # Optional: Provision Go 1.18
  config.vm.provision "shell", inline: <<-SHELL
    # Update system
    sudo apt-get update
    sudo apt-get install -y wget tar

    # Download and install Go 1.18
    wget https://golang.org/dl/go1.18.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz

    # Clean up downloaded tar file
    rm go1.18.linux-amd64.tar.gz

    # Set Go environment variables
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.bashrc

    # Apply the changes to the current shell session
    source ~/.bashrc
  SHELL
end
