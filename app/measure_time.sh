#!/bin/bash

# URL to test
url="localhost:9090/echo?iter=1000"

# Use curl to measure the time of the request
start=$(date +%s%3N)  # Start time in milliseconds
curl -s $url > /dev/null  # Perform the request
end=$(date +%s%3N)    # End time in milliseconds

# Calculate duration
duration=$((end - start))

# Output the duration
echo "Request time: ${duration} ms"
