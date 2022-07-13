#!/usr/bin/bash
# install go if required
if which go; then
    echo go installed, continuing
else
    wget https://go.dev/dl/go1.18.4.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
    # extra export just for the duration of the script
    export PATH=$PATH:/usr/local/go/bin
fi
