# install go if required
if which go; then
else
    wget https://go.dev/dl/go1.18.4.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz
fi
