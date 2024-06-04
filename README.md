# go-docker-run
This Repository contains the Application which simulates the docker run command.

The Application runs the command passed during execution in a separate PID, HOSTNAME and Mount namespace. Also runs the process in the a separate process Cgroup.

Command to the run the application:
 ./main run <linux comamnd>
 Example: ./main run /bin/bash
 
Download the rootfs from the below link in the same directory as the application to run the code.
 wget https://github.com/ericchiang/containers-from-scratch/releases/download/v0.1.0/rootfs.tar.gz
