#!/bin/bash
APP_DIR="$(dirname "$0")"
cd "$APP_DIR" || exit 1

export LD_LIBRARY_PATH=$APP_DIR/resources/lib

export CFW=KNULLI

./grout
