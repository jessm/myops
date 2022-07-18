# myops
DevOps helper for single host deployment

To build example:
```
docker build -t ping -f Dockerfile.ping .
```

To run example:
```
docker run -p 8123:8123 -v pingvol:/app ping
```
