# Docker Installation (for most 64-bit GNU/Linux systems)
StashApp is supported on most systems that support Docker. Your OS likely ships with or makes available the necessary packages.

## Dependencies
Only `docker` is required. For the most part your understanding of the technologies can be superficial. So long as you can follow commands and are open to reading a bit, you should be fine.

Installation instructions are available below, and if your distributions's repository ships a current version of docker, you may use that.
https://docs.docker.com/engine/install/

On some distributions, `docker compose` is shipped seperately, usually as `docker-cli-compose`. docker-compose is not recommended.

### Get the docker-compose.yml file

Now you can either navigate to the [docker-compose.yml](https://raw.githubusercontent.com/stashapp/stash/develop/docker/production/docker-compose.yml) in the repository, or if you have curl, you can make your Linux console do it for you:

```
mkdir stashapp && cd stashapp
curl -o docker-compose.yml https://raw.githubusercontent.com/stashapp/stash/develop/docker/production/docker-compose.yml
```

Once you have that file where you want it, modify the settings as you please, and then run:

```
docker compose up -d
```

Installing StashApp this way will by default bind stash to port 9999. This is available in your web browser locally at http://localhost:9999 or on your network as http://YOUR-LOCAL-IP:9999

Good luck and have fun!

### Docker
Docker is effectively a cross-platform software package repository. It allows you to ship an entire environment in what's referred to as a container. Containers are intended to hold everything that is needed to run an application from one place to another, making it easy for everyone along the way to reproduce the environment.

The StashApp docker container ships with everything you need to automatically run stash, including ffmpeg.

### docker compose
Docker Compose lets you specify how and where to run your containers, and to manage their environment. The docker-compose.yml file in this folder gets you a fully working instance of StashApp exactly as you would need it to have a reasonable instance for testing / developing on. If you are deploying a live instance for production, a [reverse proxy](https://docs.stashapp.cc/guides/reverse-proxy/) (such as NGINX or Traefik) is recommended, but not required.

The latest version is always recommended.
