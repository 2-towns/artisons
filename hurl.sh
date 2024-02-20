#!/bin/bash 

redis-cli -h localhost -p 6379 < web/redis/integration.redis
hurl --variable time=$(date +%s) --variable host=http://localhost:8080 --test */**/cartprocess.hurl
