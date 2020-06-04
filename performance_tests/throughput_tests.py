import requests
import time
import os


times_taken = []

for i in range(20):
    os.system("kubectl delete --all pods")
    os.system("kubectl delete --all deployments")
    os.system("kubectl delete service nats-svc")
    os.system("kubectl delete service pub-sub-svc")

    url = "http://localhost:8080/trigger"

    payload = "{\n\t\"gitrepo\": \"https://github.com/AartiJivrajani/Auto-Pub-Sub.git\"\n}"
    headers = {
        'Content-Type': 'application/json'
    }

    uid = requests.request("POST", url, headers=headers, data=payload)

    print(uid.text)

    # Sleep for 60 seconds
    time.sleep(30)

    responseurl = "http://localhost:8080/status?TRACE_ID=" + str(uid.text)

    payload = {}
    headers = {}

    while True:
        response = requests.request("GET", responseurl, headers=headers, data=payload)
        final_time = response.text.split("\n")
        if len(final_time) > 2:
            time_taken = final_time[-2]
            components = time_taken.split()
            if len(components) == 4 and components[0] == "time" and components[1] == "taken:":
                times_taken.append(float(components[2]))
                break
            else:
                print("Not done yet:", final_time[-2])
        time.sleep(10)


print(times_taken)
avg = sum(times_taken) / len(times_taken)
print("AVG:", avg)
