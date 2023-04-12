# !/usr/bin/bash
docker build . -t caos
docker run caos --env-file=.env 