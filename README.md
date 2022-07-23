# MyOps
DevOps helper for single host deployment of multiple containers. Docker and git are required. Useful for developers who want to run multiple webapps for free and only own one domain. It uses Caddy as a proxy to enable easy HTTPS, as well as matching subdomains so you can run multiple projects off one domain and one VM.

No installation necessary, can be run as a container. Developed for GCP's container-optimized OS running on a free-tier `e2-micro` instance.

## Overview
To make it easier to deploy multiple side projects on a free cloud provider VM, MyOps is a tool for setting up an https Caddy proxy and configuring your containers for you. MyOps aims to strike a balance between configurability and ease of setup, with a priority on making it easy to deploy new projects and update existing ones.

The config file looks like this:
```
# /var/myops/myops_config.json
{
    "ping": {
        "domainMatcher": "ping.jessmuir.com",
        "repoUrl": "https://github.com/jessm/ping.git",
        "branch": "ping",
        "dockerfile": "go/Dockerfile",
        "envVars": {
            "LISTENPORT": "8000"
        },
        "port": "8000",
        "volumePath": "/app",
        "hostPort": "8001"
    }
}
```
MyOps will automatically configure containers, images, and volumes when run. It will also automatically rebuild the image and rerun the container based on the latest git commit in the repository specified, or if the config changes.

## MyOps setup
To get started, append the contents of `aliases.txt` to your `~/.bash_profile` (or other profile file), run `myops_build`, and `myops_run`.

`myops_build` will build the MyOps image from this GitHub repository and tag it appropriately.

`myops_run` will run the MyOps image as a container with the required name and required bind mounts. The container will create a sample config file for use at `/var/myops/myops_config.json`, and set up a Caddy instance.

## Usage

To deploy your project, put it in a git repository and dockerize it. Add a new config to `/var/myops/myops_config.json`, and fill in a new project config similar to the sample. The `domainMatcher` field is passed straight to Caddy, so check the [Caddy documentation](https://caddyserver.com/docs/caddyfile/matchers) for details. Finally, run `myops_run`.

To update your project, push a new commit to the repository, then run `myops_run`.

To use a private GitHub repository, generate a personal access token and include it in the `repoUrl` like so:
```
"repoUrl": "https://<personal_access_token>@github.com/username/repo.git"
```

Using ssh to pull git repositories isn't supported at the moment.