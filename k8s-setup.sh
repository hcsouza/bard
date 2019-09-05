#!/bin/bash

kubectl apply -f k8s-app.yml &&
kubectl apply -f k8s-cache.yml &&
kubectl apply -f k8s-dashboard.yml
