#!/bin/bash
cat prefix.txt > docker-compose-dev.yaml
for i in $(seq $1)
do
	sed 's/REPLACE/'$i'/' repeat.txt >>  docker-compose-dev.yaml
done
cat sufix.txt >> docker-compose-dev.yaml
