# Docker install on Ubuntu 18.04
Installing StashApp can likely work on others if your OS either has it's own package manager or comes shipped with Docker and docker-compose.

## Dependencies
The goal is to avoid as many dependencies as possible so for now the only pre-requisites you are required to have are `curl`, `docker`, and `docker-compose` for the most part your understanding of the technologies can be superficial so long as you can follow commands and are open to reading a bit you should be fine.

### Docker

Docker is effectively the cross-platform software package repository it allows you to ship an entire environment in what's referred to as a container. Containers are intended to hold everything that is needed to ship what's required to run an application from one place to another with a degree of a standard that makes it easy for everyone along the way to reproduce the environment for their step in the chain.

The other side of docker is it brings everything that we would typically have to teach you about the individual components of your soon to be installed StashApp and ffmpeg, docker-compose wraps it up nicely in a handful of easy to follow steps that should result in the same environment on everyone's host.

The installation method we recommend is via the `docker.com` website however if your specific operating system's repository versions are at the latest along with docker you should be good to launch with you using whatever instructions you wish.  The version of Docker we used in our deployment for testing this process was `Docker version 17.05.0-ce, build 89658be` however any versions later than this will be sufficient.  At the writing of this tutorial, this was not the latest version of Docker.

#### Just the link to installation instructions, please
Instructions for installing on Ubuntu are at the link that follows:
https://docs.docker.com/install/linux/docker-ce/ubuntu/

If you plan on using other versions of OS you should at least aim to be a Linux base with an x86_64 CPU and the appropriate minimum version of the dependencies.

### Docker-compose
Docker Compose's role in this deployment is to get you a fully working instance of StashApp exactly as you would need it to have a reasonable instance for testing / developing on, you could technically deploy a live instance with this, but without a reverse proxy, is not recommended.  You are encouraged to learn how to use the Docker-Compose format, but it's not a required prerequisite for getting this running you need to have it installed successfully.

Install Docker Compose via this guide below, and it is essential if you're using an older version of Linux to use the official documentation from Docker.com because you require the more recent version of docker-compose at least version 3.4 aka 1.22.0 or newer.

#### Just the link to installation instructions, please
https://docs.docker.com/compose/install/

### Install curl
This one's easy, copy paste.

```
apt update -y && \
apt install -f curl
```

### Get the docker-compose.yml file

Now you can either navigate to the [docker-compose.yml](https://raw.githubusercontent.com/stashapp/stash/master/docker/production/docker-compose.yml) in the repository, OR you can make your Linux console do it for you with this.

```
curl -o ~/docker-compose.yml https://raw.githubusercontent.com/stashapp/stash/master/docker/production/docker-compose.yml
```

Once you have that file where you want it, you can either modify the settings as you please OR you can run the following to get it up and running instantly.

```
cd ~ && docker-compose up -d
```

Installing StashApp this way will by default bind stash to port 9999 or in web browser terms.  http://YOURIP:9999 or if you're doing this on your machine locally which is the only recommended production version of this container as is with no security configurations set at all is http://localhost:9999

Good luck and have fun!
