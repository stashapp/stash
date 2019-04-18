FROM ubuntu:18.04 as prep
LABEL MAINTAINER="leopere [at] nixc [dot] us"

RUN apt-get update && \
    apt-get -y install curl && \
    apt-get autoclean -y && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -L -o /stash $(curl -s https://api.github.com/repos/stashapp/stash/releases | grep -F 'stash-linux' | grep download | head -n 1 | cut -d'"' -f4)

FROM ubuntu:18.04 as app
RUN adduser stash --gecos GECOS --shell /bin/bash --disabled-password --home /home/stash
COPY --from=prep /stash /usr/bin
# COPY ../../scripts/start.sh /usr/local/bin/
EXPOSE 9998
CMD ["stash"]
