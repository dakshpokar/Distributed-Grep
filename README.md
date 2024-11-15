# MP1-Distributed-Systems

This is a distributed grep system that allows for efficient parallel search across multiple nodes. The system is implemented using Go, leveraging concurrency and distributed processing to search files across multiple nodes in a network.

## Project Setup for all machines

0. Ensure that you have Go Lang Installed on your machine.
1. Clone this repo on each machine.
2. cd into mp1-distributed-systems folder on all machines.
```bash
cd mp1-distributed-systems
```
3. Start the server on all machines by executing following command - 
```go
go run server.go
```
4. Ensure that your log files are present in the same directory on all machines. Example - `/root/`

## How to run Distributed Grep
Now that all machines are running the server, you can execute client.go from any of the machines but before that you have to follow these steps -

1. Create a `.env` file in mp1-distributed-system folder and add the following -
```env
REMOTE_IP=XXX.XXX.XXX.XXX
```
Enter the common path of the IP / Hostname. For example - 
```env
REMOTE_IP=fa24-cs425-92
```
This has to be performed once on the machine where the distributed grep is to be run. If you run distributed grep on any other machine then you have to add this environment file again.
2. Run client.go on the machine
```bash
go run client.go <<PATTERN>> <<COMMON_PATH_TO_FILE>>
```
Example - 
```bash
go run client.go PUT /root/vm
```
In the above example all the machines have vmX.log file in root directory.

## Group G92

Rishi Mundada (rishirm3@illinois.edu)
Daksh Pokar (dakshp2@illinois.edu)