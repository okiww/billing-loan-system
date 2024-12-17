
# Go Project Billing Loan

This is a Go project that utilizes various tools and libraries including Go 1.23.1, MySQL, RabbitMQ, Viper, Cobra, Gomock & Mockgen, and Goose.
## Flow
### Database

## Feature
* **API Create Loan**
    - Create Loan
    - Get All Loan
    - Make Payment
* **Cronjob** is the background job that generate bills payment every weekly in monday
* **Worker** is the worker that listening or as consumer message from rabbitMQ

## Setup & Installation

### 1. Clone the project

```bash
  git clone https://github.com/poscompany/gprc-contract
```

### 2. Go to the project directory

```bash
  cd <your-project-directory>
```

### 3. Install dependencies
Install the necessary Go dependencies:
```bash
  go mod tidy
```

### 4. Setup database
Make sure your MySQL server is running and create a database:
```bash
  CREATE DATABASE your_database_name;
```

### 5. Setup local env
Make sure env is Setup
```bash
  cp configs/env.yml configs/env-local.yml
```
Setup MySQL DSN and RabbitMQ DSN

### 6. Set up RabbitMQ
```bash
  docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

### 7. Run Migrations with Goose
To manage your database schema, you can use Goose for migrations. First, ensure that your goose binary is available:
```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
```
Then, create migration files and apply them to the database:
```bash
 goose -dir "./db/migrations" mysql "<user>:<password>@tcp(localhost:3306)/<dbname>" up
 ```

## Make Commands

This project uses a Makefile to simplify various tasks.

To run the HTTP server for the application:
```bash
make serve-http
```
To run the Cronjob server for the application:
```bash
make run-cron
```
To run the Worker server for RabbitMQ Subscriber:
```bash
make serve-http
```
## Tech Stack

**Server:** Go, MySQL, RabbitMQ, Viper, Cobra, Mockgen, Goose
