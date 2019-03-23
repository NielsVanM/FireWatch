# FireWatch
FireWatch is a platform designed to gather statistics about nature. It provides a platform for IoT sensors to be placed in a forest (or other environment) and monitor the data in a single place.

## Development
Developement speed will be limited as I'm stil learning GoLang & Vue and sadly don't have all the time available to continue working on FireWatch.

For the backend you'll need to have GoLang 1.1 installed. Any dependencies can be automatically retrieved by running ```go mod download``` in the backend directory.
For the frontend you'll need npm installed, any dependencies can be retrieved using ```npm install``` in the frontend directory

## Deployment
Deployment will be simplified by a docker setup. However, this is currently not provided. It'll be added when it's necessarry.

## Dependencies
As any project I rely on others to build our systems.

Backend:
* [MUX](https://www.github.com/gorilla/mux) HTTP Router
* [Context](https://www.github.com/gorilla/context) Context passing
* [PQ](https://www.github.com/lib/pq) Postgres Database
* [BCrypt](https://www.golang.org/x/crypto/bcrypt) Hashing
* [LogRus](https://www.github.com/sirupsen/logrus) Logging

Frontend:
* [Vue](https://www.vuejs.org) Frontend Framework
* [Vuetify](https://www.vuetifyjs.com) Material Design Framework
* [Vuex](https://vuex.vuejs.org) State management
* [Vue-Router](https://router.vuejs.org) HTTP Routing
* [Axios](https://github.com/axios/axios) HTTP request library

For instructions on how to install the dependencies see the development section of this readme.