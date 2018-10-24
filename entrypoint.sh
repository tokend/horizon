#!/usr/bin/env bash

set -a
. /config.env
horizon "${@:1}"
