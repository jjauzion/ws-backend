---> **WORK IN PROGRESS (missing auth)** <---

# Workstation
The workstation project is machine learning job management system.  

It consists of a task queue where user can create new jobs and one or more worker nodes that will pull
 jobs from this queue, run the algorithm and return the result when it is finished.

The jobs are submitted via a Docker image that shall be available on a public container registry.

The project is made of three repositories:
- the worker node: *the current repo*
- the backend server: https://github.com/jjauzion/ws-backend
- the frontend: *not yet implemented*

# How to run

This tutorial will guide you through running the whole project.

## Pre requisite
- macOS or linux
- Docker installed
- Docker-compose installed
- go installed
- makefile support
- clone the project repositories to your machine

## Start the backend
Here we will run the entire project on your local machine from scratch, including the database.
 The database will be boostrapped with default users.
- Clone the backend repository: `git clone https://github.com/jjauzion/ws-backend`
- Open the repo: `cd ws-backend`
- Create the `.env` file. For a dev environment use this:
```dotenv
WS_ES_HOST=http://localhost
WS_ES_PORT=9200
WS_KIBANA_PORT=5601
WS_API_PORT=8080
WS_GRPC_HOST=localhost
WS_GRPC_PORT=8090
```
- Start the elastic and kibana cluster: `make elastic`
- Check kibana container logs: `docker logs ws-backend_kibana_1 -f`   
- Wait until you see:
```dockerfile
{"type":"log","@timestamp":"2021-03-28T15:11:50+00:00","tags":["listening","info"],"pid":7,"message":"Server running at http://0:5601"}
{"type":"log","@timestamp":"2021-03-28T15:11:51+00:00","tags":["info","http","server","Kibana"],"pid":7,"message":"http server running at http://0:5601"}
{"type":"log","@timestamp":"2021-03-28T15:11:54+00:00","tags":["warning","plugins","reporting"],"pid":7,"message":"Enabling the Chromium sandbox provides an additional layer of protection."}
```
- Start the GraphQL server `make gql FLAG="--bootstrap"`  
  The bootstrap option initialise the DB by creating the required index and indexing default users
- Open a new terminal in the same repo 
- Start the gRPC server: `make grpc`
  
At this point you have started the database, the graphQL server that interact with the frontend
 and the gRPC server that interact with the worker nodes.  

## Kibana
Before starting the worker node we will learn how to interact with the backend. First, lets check the
database:
- Open Kibana: http://localhost:5601  
- Click on the burger menu in the top left corner and go to the `Dev Tools`
- Copy / Paste the following in the console and run it: `GET _cat/indices?v`  
  This list all the index existing in the DB. You should see an index called `task` and one
  called `user`. Index starting with a dot `.` are system index.
- Now run the following to list all the users existing:
```
GET user/_search
{
  "query": {
    "match_all": {}
  }
}
```
- To get all the task, replace in the previous query `user` by `task`
- You can use this console for debug purpose if you need to check the content of your database.
  You could also create or delete task and user manually from here but it is better to use the
  GraphQL API.
  
## Create a user and a task
We will now see how to create user and task with the GraphQL API.
- Open the GraphQL playground: http://localhost:8080/playground
- You can find the doc and schema of our API thanks to the "DOCS" and "SCHEMA" tabs on the right side
  of the screen
- To create a new user, paste the following in the console:
```graphql
mutation tuto_create_user {
  create_user(input:{email:"just-for-test@email.com"}) {
    id
    email
  }
}
```
You should have a response like this:
```graphql
{
  "data": {
    "create_user": {
      "id": "86c776ec-9abe-43a0-93f1-4dac0997ba90",
      "email": "just-for-tes3t@email.com"
    }
  }
}
```
- Copy the id of the user you've just created
- Now let's create a task. Run the following command (replace the id with yours):
```graphql
mutation createTask {
  create_task(input:{user_id:"65941391-733a-430c-a3bd-2bdd853af7be", docker_image:"jjauzion/ws-mock-container", dataset:"s3//"}) {
    id
    user_id
  	created_at
  	started_at
  	ended_at
  	status
    job { dataset, docker_image }
  }
}
```
Congratulations !! You have created a user and a new jobs :) You can go to the kibana console and run 
the search to see your creations.

## Run a worker node
Now that we have created a new task, it would be nice to have a worker to actually run that task right?  
But before starting a worker node, we need to start the gRPC server:
- Go in the `ws-backend` repository and run: `make grpc`  

Now let's run the worker:  
- Clone the worker repository: `git clone https://github.com/jjauzion/ws-worker.git`
- go in the `ws-worker` repo: `cd ws-worker`  
- Create the `.env` file. For a dev environment use this:
```dotenv
WS_ES_HOST=http://localhost
WS_ES_PORT=9200
WS_KIBANA_PORT=5601
WS_API_PORT=8080
WS_GRPC_HOST=localhost
WS_GRPC_PORT=8090
```
- Start the worker: `make run`  

This will start the worker and it will automatically pull the task you have created in the 
  previous chapter and run it.

You can go to kibana and check your task, you will see the status going from "NOT_STARTED" 
to "RUNNING" and "ENDED"

When you are done, go in the ws-backend repo and run `make down` to stop and the elastic containers

That's it folks :)
