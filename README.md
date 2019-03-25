 # XII-Traits Backend Challenge
 Multi-provider Data Aggregation ([plz read the challlenge discription first](https://github.com/rehehe/xii/blob/master/challenge.md))
 
 ## Prerequisites
 
 Our understanding of the challenge and assumptions description:
 - The system is based on the pull approach
 - The APIs are basic which supports export(filter, limit), save
 - The system is under high load and should respond in a real-time manner with the very new data (completely synced)
 - The system should persist data in local storage
 - The system should be Dockerized
 
 ### Technical decisions
 - The architecture design should comply with the micro-service pattern:
   - Independency in case of failure (High Availability)
   - Independent scalability
   - Independent, easy and frequent deployment
   - Independent technology stacks
   - Independent and parallel development
 - Used language should supports micro-service environment, so we chose Golang:
   - A Native language/high performance
   - Rich support of concurrency
   
 ### Implemented solution
 #### Terms
 - Dashboard
   - HTTP/REST: HTTP multiplexer/router
   - Reporter: Request(get) handler service backed by database and Aggregator service
   - Aggregator: Data aggregator and orchestra service which time-synchronize Requests(Reporter), Providers(Worker), and data-synchronize Local Storage
   - Worker: Query interpreter/dialect-adaptor which also fetches new data from provider
   - Storage: The database service
 - Survey
   - HTTP/REST: HTTP multiplexer/router
   - Reporter: Request(post/get) handler service backed by database
   
 #### Architecture and Flow Description
 It is roughly like this:
 
 <img alt="arch & flow desc" src="https://github.com/rehehe/xii/blob/master/arch_flow_desc.jpg" height="60%" width="60%">
 
 As illustrated above:
 - The user request routes by REST service to Reporter
 - Reporter submits synchronizing request to Aggregator through (Golang) channel
 - Aggregator does these 4 tasks concurrently:
   - Register all Reporter requests(which is a signal-channel) in the callback-list
   - Checks the callback-list whether there is any request to signal Workers or not
   - Watches workers results to aggregate them all, save it in storage, signal all Reporters exist in callback-list and clear callback-list
   - Timeout those workers which didn't finish their task after a specific amount of time and do the previous process to the buffered data 
 - Each Worker corresponds to one Survey server which fetches new data from it
 - Storage is our database which persists data
 - Reporter fetches updated/synchronized data from database after Aggregator sends a signal via callback channel
 - REST service gets the output from Reporter and responds to the request
 
 ## Installation
 
 ### Docker
 ```bash
docker-compose up
```
 
 ### Makefile
 build:
 ```bash
 make
 ```
 run:
 ```bash
 make run-dashboard
 make run-survey
 ```
 or get help to run it :
 ```bash
 ./build/dashboard-linux -h
 ./build/survey-linux -h
 ```
 
 post to survey & get from dashboard
 
 ```bash
 curl -H "Content-type: application/json" \
      -d '{"text":"request: 01 to blue survey","device":"x","location":"y"}' \
      http://localhost:1101/api/v1/reporter
 
 curl -H "Content-type: application/json" \
       -d '{"text":"request: 01 to red survey","device":"x","location":"y"}' \
       http://localhost:1102/api/v1/reporter
            
 curl http://localhost:1100/api/v1/reporter
 ```
 
 run MySQL over docker:
 ```bash
 docker run \
      --name mysql \
      -e MYSQL_ALLOW_EMPTY_PASSWORD="yes" \
      -e MYSQL_DATABASE=devdb \
      -e MYSQL_USER=dbuser \
      -e MYSQL_PASSWORD=dbpassword \
      -p 3306:3306 \
      -d \
      mysql
 ```
  
 ## Intended implementation approach
 - Using advanced message queue for inter-service communications e.g. NSQ, RabbitMQ
 - Using better data serializing structure e.g. Protobuf
 - Using RPC framework e.g. gRPC
 - Using the context concept:
   - Cancel propagation e.g. Golang native context
   - Tracing mechanism e.g. OpenTracing implemented as jaeger in Golang
 - Database Sharding and HA e.g. using MongoDB
 - Using a storage cache layer e.g. Redis
 - Load balance orchestration
 - Central remote config server for all services e.g. etcd
 
 ## Built With
 
 * [GO](https://golang.org)
 * [Docker](https://docker.com)
 * [gORM](https://gorm.io)
 
 ## Authors
 
 * **Rahi Dehzani** - *Initial work* - [Multi-provider Data Aggregation](https://github.com/rehehe/xii)
 
