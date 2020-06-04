[X] Add status http endpoint and handler  
[X] Add k8s pod/service build integration  
[X] add functions to read the user input  
[X] make sure TraceId is returned as part of the first trigger  
[X] Add support for multiple docker files in case the user has that requirement  


DEMO: 
1. Create the AWS instance [Prep]
2. Make sure that the keys are plugged into the script [Prep]
3. Run the `rad_management_cluster.sh` script 
4. Run the `rad_workload_cluster.sh` script
5. Run `setup_go.sh` script (this installs go, its dependencies, and also installs and runs `redis`)
6. Run `trigger/main.go`  
7. /trigger `collect the traceId`
8. /status