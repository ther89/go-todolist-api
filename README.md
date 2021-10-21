# Go Todo List Backend API

This is a learning project to get familiar with Golang, Visual Studio Code and developer containers. 

The inspiration and the Go basics are coming from [github.com/sdil](https://github.com/sdil)'s [tutorial](https://www.fadhil-blog.dev/blog/golang-todolist/). Kudos for that!

The project can be run inside a *Go* docker developer container. This way you don't need to install anything else in order to start the project, everything that is needed for the development is already described in the `.devcontainer` folder. Your terminal commands are running inside the container.
 
## Prerequisites

- [Visual Studio Code](https://code.visualstudio.com/) 
- [Docker](https://www.docker.com/)
- [Remote - Containers](vscode:extension/ms-vscode-remote.remote-containers) VS Code extension (If you open the workspace you will be prompted to install this. More information about it can be found [here](https://code.visualstudio.com/docs/remote/containers-tutorial))

## Opening the project

Visual Studio Code: You can open the project by navigating to the root folder and open `go-todolist-fe.code-workspace`.

## Project architecture and key components

### Containers
- A MySQL server (listening on port 3306)
- A MySQL webgui *Adminer* for Database management (listening on port 3080)
- A Go developer container hosting the API (listening on port 3010)

### .devcontainer folder
- `devcontainer.json` - used by Visual Studio Code to describe commands and containers to use for remote development. See [devcointainer.json reference](https://code.visualstudio.com/docs/remote/devcontainerjson-reference) for more.
- `docker-compose.yml` & `Dockerfile` - to describe docker containers.
- `postCreateCommand.sh` - extracted set of commands to run once the developer container is created.

### main.go

This project supplies a HTTP api with CRUD operations and a health check. The endpoints accepts and returns json payload.
It is using `github.com/gorilla/mux` package for routing, `github.com/rs/cors` for cors handling. For storing it is using an [ORM](https://en.wikipedia.org/wiki/Object%E2%80%93relational_mapping) tool `gorm.io` it is handling the data definition and manipulation as well based on our models defined. For logging, `github.com/sirupsen/logrus` package is used.

## Setup before first start

Once your project is opened in Visual Studio Code (or your docker containers are up and running) you need to initiate an empty MySQL database called `todolist`
- Navigate to Adminer (http://localhost:3080/)
- Input these: Server: mysql, Username: root, Password: root
- Once logged in, click *Create database* and name it `todolist`

## Start 

`launch.json` is set up so you can simply start the API by hitting F5 (Start debugging) in Visual Studio Code, or with `go run main.go` command in the developer container terminal.

## Frontend

I have created a frontend to be used with this API. The source code and running insturctions for that can be found at https://github.com/tomhudak/ng-todolist-fe.


## Known issues / blind spots

For some reason, despite the packages are acquired with `go get` in the `postCreateCommand.sh`, the import statements shown error as the packages are not there. If you edit one line it detects that the packages are there and removes the error.
