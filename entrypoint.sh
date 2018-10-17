#!/usr/bin/env bash

set -a
. $1
horizon /config.env "${@:2}"
