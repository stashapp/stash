#!/usr/bin/env bash

USER_ID=${LOCAL_UID:-1000}

USER_NAME="$(getent passwd | awk -F: '$3 == '${USER_ID}' { print $1 }')"

if [ "$USER_NAME" == "" ]; then
  USER_NAME=stash

  if [ ! -d "/home/{USER_NAME}" ]; then
    ARGS='-m'
  fi

  useradd --shell /bin/bash -u "$USER_ID" $ARGS "$USER_NAME"
  export HOME="/home/${USER_NAME}"

  chown "$USER_NAME" $HOME
fi

gosu "$USER_NAME" "$@"
