# FireWatch
FireWatch is a platform designed to gather statistics about nature. It provides a system for IoT sensors to be placed in a forest (or other environment) and monitor the data in a single place.

## Development
Developement speed will be limited as I'm stil learning GoLang and sadly don't have all the time available to continue working on FireWatch.

As of now I don't accept pull requests, due to the fact that this is meant to be a personal project. However, you are allowed to fork the project, edit the source and deploy it whether it's for personal goals or commercial purposes. Read the LICENSE for more information about what you can and cannot do.

You'll need to have GoLang 1.1 installed. Any dependencies can be automatically retrieved by running ```go mod download```.

## Deployment
Deployment will be simplified by a docker setup. However, this is currently not provided. It'll be added when it's necessarry.

## Dependencies
As any project we rely on others to build our systems. Special thanks to the creators of GoLang and the creates of the following packages.

* [MUX](https://www.github.com/gorilla/mux) HTTP Router
* [PQ](https://www.github.com/lib/pq) Postgres Database
* [BCrypt](https://www.golang.org/x/crypto/bcrypt) Hashing

For instructions on how to install the dependencies see the development section of this readme.