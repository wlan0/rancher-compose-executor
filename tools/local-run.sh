#!/bin/bash
godep go clean; godep go build

# Asssumes you have a default service account type setup
CATTLE_AGENT_LOCALHOST_REPLACE="10.0.3.2" RANCHER_ACCESS_KEY="service" RANCHER_SECRET_KEY="servicepass" RANCHER_URL=http://localhost:8080/v1 ./rancher-compose-executor
