#!/bin/bash

echo "start screenshot"

sudo docker run --rm -it -d \
  -v /vagrant:/tmp/outputs \
  --name goheadless \
  -v $(pwd):/go/src/github.com/tjtjtjtj/web-zabbix \
  -e ZABBIX_ENV=local \
  -e ZABBIX_USER=admin \
  -e ZABBIX_PASSWORD=zabbix \
  goheadless

echo "end screenshot"
