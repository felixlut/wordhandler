# wordhandler
Small Golang application developed to learn kubernetes. It contains 2 parts, a server (receiver) and an emitter. The emitter sends a random word to the receiver, which catches it and stores some data about it (times seen, last seen, ...). 


## Run instructions
The program can be ran in multiple environments

### Default
Start the reciever and emitter by running `go build .` and invoking the compiled binary in their respective sub-directories. 

### Docker Compose
Run the application by running `docker-compose up`

### Kubernetes


## Cli
The project also contains a cli for interacting with the receiver. Currently only the `get word` command is available, which returns the stored statistics for the word (first seen, ...). 