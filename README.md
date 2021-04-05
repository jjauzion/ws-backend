---> **WORK IN PROGRESS (missing auth)** <---

# Workstation
The workstation project is machine learning job management system.  

It consists of a task queue where user can create new jobs and one or more worker nodes that will pull
 jobs from this queue, run the algorithm and return the result when it is finished.

The jobs are submitted via a Docker image that shall be available on a public container registry.

The project is made of three repositories:
- the worker node: https://github.com/jjauzion/ws-worker
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
IS_DEV_ENV=true
TOKEN_DURATION_HOURS=24
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
- Now run the following to list all the existing users:
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
  

## Login
Before being able to create user and task you will need to login. As we started the GraphQL server
with bootstrap option, two default users have been created in the DB.  
We will login with the admin user using the GraphQL API.
- open the GraphQL playground: http://localhost:8080/playground
- you can find the doc and schema of our API thanks to the "DOCS" and "SCHEMA" tabs on the right side
  of the screen. This will help you later to build your own request
- copy / paste the following request to log in as the admin user:
```graphql
query login {
  login (id: "admin-user@email.com", pwd: "") {
    ... on Token {
      token
      userId
      username
    }
    ... on Error {
      code
      message
    }
  }
}
```
You should get a response similar to this:
```json
{
  "data": {
    "login": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTczODMzNDksInVzZXJfaWQiOiJkZjljNDYzZC00ZmIwLTRmYzAtYTU5OC00YmQ3NzEzMzg2ZDAifQ.Xj_rUGIB7l90kiXD_U12ni2kf9U-afARaCZKbEao-oU",
      "userId": "df9c463d-4fb0-4fc0-a598-4bd7713386d0",
      "username": "admin-user@email.com"
    }
  }
}
```
As you can see the server successfully authenticated your request and have generated a JWT token 
that you can use for further request to prove that you are authenticated.  

Copy the `token` value and `userId` somewhere as you will need those later.
  
## Create a user
We will now see how to create user and task with the GraphQL API.
- To create a new user, paste the following in the console:
```graphql
mutation tuto_create_user {
  create_user(input:{email:"just-for-test@email.com"}) {
    id
    email
  }
}
```
- if you send the request like this you will get a `403` error because you are not authenticated (WIP at the moment only create task requires authentication)
- you need to pass the token we generated in the chapter before: on the bottom of the console, 
  click on `HTTP HEADERS` and paste the following (replace the token value with yours):
```json
{
  "auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTczODMwNTUsInVzZXJfaWQiOiI0MzVmNTA3OC02NjFlLTRkOGMtODJjZS0zNDJhZTQ1ZTQ4MzcifQ.RTzseF7mSjR8aop-9CCiNt1-IkqFGem9nNWymaJKBRo"
}
```
- you can now send the request. You should have a response like this:
```json
{
  "data": {
    "create_user": {
      "id": "86c776ec-9abe-43a0-93f1-4dac0997ba90",
      "email": "just-for-tes3t@email.com"
    }
  }
}
```
## Create a task
- now let's create a task. Run the following command (replace the user id with yours):
```graphql
mutation createTask {
  create_task(input:{docker_image:"jjauzion/ws-mock-container", dataset:"s3//"}) {
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
- if you got a `403` error, check you didn't forget the `auth` Header in your request (see previous chapter)

Congratulations !! You have created a user and a new jobs :) You can go to the kibana console and run 
the search to see your creations.

## Run a worker node
Now that we have created a new task, it would be nice to have a worker to actually run that task right?  
But before starting a worker node, we need to start the gRPC server:
- Go in the `ws-backend` repository and run: `make grpc`  

Now let's run the worker:  
- Clone the worker repository: `git clone https://github.com/jjauzion/ws-worker.git`
- go in the `ws-worker` repo: `cd ws-worker`  
- Create the `.env` file with the same values as the `.env` file created for the `ws-backend`
  (you can create a simlink to the `.env` of the ws-backend file)
- Start the worker: `make run`  

This will start the worker and it will automatically pull the task you have created in the 
  previous chapter and run it.

You can go to kibana and check your task, you will see the status going from "NOT_STARTED" 
to "RUNNING" and "ENDED"

## Create a Machine Learning task
Let's create a real job: running a ML jobs and tracking your jobs parameters while it is running.

For this we will use wandb (https://wandb.ai/site) so you must create a user and copy your private key.

Then paste the following in the playground console and put your wandb key in the env variable.
Your key will be encrypted on the server and will never be stored in clear (WIP, not done yet)

```graphql
mutation createTask {
  create_task(input:{env:"WANDB_API_KEY=putYourKeyHere", docker_image:"jjauzion/wandb-test", dataset:"s3//"}) {
    id
    user_id
        created_at
        started_at
        ended_at
        status
    job { dataset, docker_image, env }
  }
}
```
Wait until the task status is updated to "RUNNING" (can take up to 30sec), then log to wandb.
You should see your work ongoing.

When you are done, go in the ws-backend repo and run `make down` to stop and the elastic containers

That's it folks :)
