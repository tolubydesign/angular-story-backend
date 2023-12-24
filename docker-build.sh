#!/bin/bash
docker build --platform linux/amd64 -t aws-cdk-golang:test .
# docker build --platform linux/amd64 -t docker-image:test .