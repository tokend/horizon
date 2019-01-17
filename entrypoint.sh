#!/usr/bin/env bash

set -a
. ${CONFIG:-/config.env}
horizon "${@:1}"
