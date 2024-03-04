#!/usr/bin/env zsh
 curl -O https://raw.githubusercontent.com/keploy/keploy/main/keploy.sh && source keploy.sh
 keploy
 docker volume create --driver local --opt type=debugfs --opt device=debugfs debugfs