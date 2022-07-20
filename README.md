# myops
DevOps helper for single host deployment

To run myops:
```
docker build -t myops -f Dockerfile.ops git@github.com:jessm/myops#myops-container && docker run --rm --name myops -v /var/myops:/var/myops -v /var/run/docker.sock:/var/run/docker.sock myops
```

To build example:
```
docker build -t ping -f Dockerfile.ping .
```

To run example:
```
docker run -p 8123:8123 -v pingvol:/app ping
```
