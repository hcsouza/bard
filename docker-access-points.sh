#!/bin/bash

echo ""
echo "|-----------------------|"
echo "  Service Access Points  "
echo "|-----------------------|"
echo ""
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bard_web_1 |
awk '{print "RestApi Sample: http://"  $1 ":8088/musics/city?name=campinas"  }'
echo " "
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bard_web_1 |
awk '{print "Metrics Stream: http://"  $1 ":81/hystrix.stream"  }'
echo ""
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bard_dashboard_1 |
awk '{print "Metrics Dashboard: http://"  $1 ":8080"  }'


