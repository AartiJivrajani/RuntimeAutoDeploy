#!/bin/sh

start=$SECONDS
sh ../setup/setup_rad_wordload_cluster.sh
echo $(( SECONDS - start )) seconds


kubectl delete -f ../../../../capi-quickstart.yaml 
