#!/usr/bin/env zsh
 curl -O https://raw.githubusercontent.com/keploy/keploy/main/keploy.sh && source keploy.sh
 keploy
 docker volume create --driver local --opt type=debugfs --opt device=debugfs debugfs
 set -a && source .env.docker && set +a
 docker build -t go-template-app:1.0 .