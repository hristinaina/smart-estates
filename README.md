# Smart Home

Project for "Rich Internet Applications" university course.

This project is a software solution that enables tracking smart houses within smart cities.

## Technologies

- [Go](https://golang.org/): High-performance programming language used for server-side logic.
- [Gin](https://gin-gonic.com/): Web framework for Go, enabling easy and efficient handling of HTTP requests.
- [MySQL](https://www.mysql.com/): Relational database management system.
- [InfluxDB](https://www.influxdata.com/): Time-series database for handling metrics and events.
- [React](https://reactjs.org/): JavaScript library for building user interfaces.
- [Material Design](https://material.io/): Design system for creating visually appealing and consistent UIs.
- [Nginx](https://www.nginx.com/): Web server for handling HTTP and reverse proxy requests.
- [Docker](https://www.docker.com/): Platform for developing, shipping, and running applications in containers.


## Installation

Before you begin, ensure you have the latest version of Go installed. You can download it from the official [Go website](https://golang.org/).

Install **InfluxDB** by following the instructions on the [official website](https://docs.influxdata.com/influxdb/v2.0/get-started/). In the *smarthome-back/config* folder add **config.json** file with data to access the influx database.
 
In the same folder you will find script **database.sql** to add data to MySQL database. 


## How to Run

### InfluxDB

Position in the installation folder of the influxdb and run the installation. Influxdb will run on port 8086

### MQTT

Start MQTT broker by positioning in `simulation/config/broker` folder and running the next command:

`docker compose up`

### Backend

In the `smart-home/smarhome-back` folder, run the following command to start the backend server:

`go run main.go`

It will run on port 8081. **After the initial startup, be sure to check the console.**

### Client 

In the `smart-home/smarthome-front` folder, install dependencies and start the frontend:

`npm install`

`npm start`

### Simulation

In the `smart-home/simulation` folder, run next command to start the simulation:

`go run main.go`


## Contributors
- [Anastasija Savić](https://gitlab.com/savic-a)
- [Katarina Vučić](https://gitlab.com/kaca01)
- [Hristina Adamović](https://gitlab.com/hristinaina)
