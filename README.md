# Go Project Billing Loan

This is a Go project that utilizes various tools and libraries including Go 1.23.1, MySQL, RabbitMQ, Viper, Cobra, Gomock & Mockgen, and Goose.
## Flow
### Database
![Amartha - Billing System](https://github.com/user-attachments/assets/a8cc0377-4380-4eb0-84b2-89b67c343df5)

## Feature
* **API Create Loan**
  - Create Loan
  - Get All Loan
  - Make Payment
  - ![image](https://github.com/user-attachments/assets/a5779a99-491f-4d6e-85e6-e3d1e1609b22)
    - Create Payment and Save to DB as Pending
    - Publish to RabbitMQ for Process Payment
* **Cronjob**
  - Background job that update each **PENDING** loan bills status to **Billed** or **Overdue** every weekly in monday
  - If users has more than 1 **OVERDUE**, will update users to delinquent and wouldn't create loan unless he pays all **OVERDUE** bills
* **Worker** is the worker that listening or as consumer message from rabbitMQ
  ![image](https://github.com/user-attachments/assets/ed001307-4798-4621-90c7-50385603ca07)
  - Subscribe payment message and **PROCESS**
  - Update payment status to process
  - Validation loan, loan bill and amount
  - Update loans status and bill status under Trx
  - Update payment status to **SUCCESS** if success, and **FAILED** if has errors
  - Count total overdue
    - if has less than 2 & user is delinquent, update user to is not delinquent

## Setup & Installation

### 1. Clone the project

```bash
  git clone https://github.com/okiww/billing-loan-system.git
```

### 2. Go to the project directory

```bash
  cd billing-loan-system
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
### 8. Via Docker
The project itself provide **Dockerfile** and **docker-compose**

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
To format code and format import code:
```bash
make format
```
To generate mock:
```bash
make gen-mocks
```
To run unit test:
```bash
make test
```
![image](https://github.com/user-attachments/assets/515fc2f6-5e1b-434b-9f26-11d4b14c46ac)

## API Collection
You can see at project root folder `biling-loan-system/billing-engine.json`

## Tech Stack
**Server:** Go, MySQL, RabbitMQ, Viper, Cobra, Mockgen, Goose
