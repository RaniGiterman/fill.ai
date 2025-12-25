#!/bin/sh

docker rm -f fill_ai

docker build -t fill_ai:v1.0 .

docker run -d -p 8080:8080 --name fill_ai fill_ai:v1.0
