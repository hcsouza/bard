#!/bin/bash

cluster=$1
echo ""
echo "Cluster Name: $cluster"
echo ""

echo "|-----------------------|"
echo " Service Access Points"
echo "|-----------------------|"
echo ""
minikube -p $cluster service bard-api-service --url | cut -d ' ' -f 2 | awk '{print "RestApi: "  $1 "/musics/city"  }' &&
echo ""
minikube -p $cluster service bard-metrics-service --url | cut -d ' ' -f 2 | awk '{print "Metrics Stream: "  $1 "/hystrix.stream"  }'
echo ""
minikube -p $cluster service bard-metrics-dash --url | cut -d ' ' -f 2 | awk '{print "Metrics Dashboard: "  $1   }'
