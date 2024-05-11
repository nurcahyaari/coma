## ⭐️ What is Coma?

Coma is a configuration manager that stores, manages, and distributes your configuration across your services. With Coma, you can keep your configuration and seamlessly changing your configuration without updating it on the server and redeploying your server.

## 📚 Table Of Contents

- [Getting Started]()
- [Features]()
- [Roadmaps]

## 🚀 Getting Started

### Manual Installation

- make sure Go has been installed on your local machine
- you must clone this repository
  ```bash
  clone with ssh
  git@github.com:nurcahyaari/Coma.git

  or clone with http
  https://github.com/nurcahyaari/Coma.git
  ```
- Make this code as an executable files
  - using make build
    - type `make build`
    - you'll see coma file under the build folder
    - move the executable file into your path directory
  - using make install
    - type `make install`
    ![alt text](<assets/Screenshot 2024-05-11 at 22.34.44.png>)

    - coma will be installed under your $GOPATH/bin directory
    - you can try to type `ls $GOPATH/bin`
    - you'll see coma has already installed under this folder
- Running coma from your local machine
![alt text](<assets/Screenshot 2024-05-11 at 22.37.52.png>)


### Containerize with Docker

under construction

## Features
- 👨‍💻 Managing user with their access control
- 📦 Simpel to manage stage, application, key, and its on air configuration
- 🚀 Real-time publishing configuration to the listener
- 📦 SDK with many programming language implementation (currently only supports Golang & Node.js)

## Roadmaps

Coma has its biggest roadmap that we want to achieve, here are several planning for coma to make it better:

- [x] Create websocket connection
- [x] Create basic feature of coma (create application, stage, configuration)
- [x] Create authentication and application key
- [x] Distribute data through websocket within application key
- [x] Create local publisher/ subscriber pattern for handling queuing data through Golang Channel
- [x] Create SDK
  - [x] Nodejs
  - [x] Golang
  - [ ] etc
- [ ] Create coma as a server apps, installable and make as a binary distribution
  - [ ] Change how coma to communicate and store its configuration
- [ ] Create coma as installable through Package manager such as brew
- [ ] Create a concise beauty documentation
- [ ] Create UI for managing its data
- [ ] Create coma as a distributed system
- [ ] etc


# How to use

This is the general case of how to use coma, 
- first you need to run the service
- Open swagger http://localhost:YOUR_PORT_SETTING/swagger/index.html
- Create user root [POST /v1/users/root]
- Create stage [POST /v1/stages]
- Create your application [POST /v1/applications]
- Create your application's key [POST /v1/keys]
- Set the configuration [POST /v1/configuration/upsert]




