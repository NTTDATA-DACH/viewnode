#!/bin/bash
for i in {1..100}
do
   yq e -M '.metadata.name = "test-agent-service-'"$i"'"' service-pod.yaml | kubectl apply -f -
   sleep .$[ ( $RANDOM % 10 ) + 1 ]s
done
