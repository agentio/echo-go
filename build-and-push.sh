#!/bin/sh

docker buildx build --sbom=true --push --provenance=true -t docker.io/agentio/echo:latest .
