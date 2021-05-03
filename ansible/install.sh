#!/bin/bash

set -e
set -o pipefail

sudo apt update
sudo apt install -y ansible git

if [ -e bhive ]; then
  git -C bhive/ pull
else
  git clone https://github.com/queensaver/bhive.git
fi

ansible-playbook bhive/ansible/bhive.yml
ansible-playbook bhive/ansible/read_only.yml
sudo reboot
