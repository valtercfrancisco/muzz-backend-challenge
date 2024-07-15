# muzz-backend-challenge

Implementation of take home coding challenge for Muzz.

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Project Structure](#project-structure)
4. [Architecture and Design](#architecture-and-design)
6. [Testing](#testing)
7. [Performance](#performance)
9. [Scalability](#scalability)
10. [Future Improvements](#future-improvements)

## Introduction

### Project Overview

This repository contains the implementation for the Explore Service described in the take-home exercise document. The 
project is a backend application developed in Go that manages user decisions and interactions within a social platform. 
It includes features for recording user decisions (likes or dislikes) on other users and checking for mutual likes
between users.


### Key Features


- `ListLikedYou(ListLikedYouRequest) returns (ListLikedYouResponse)`; // List all users who liked the recipient
- `ListNewLikedYou(ListLikedYouRequest) returns (ListLikedYouResponse)`; // List all users who liked the recipient excluding those who have been liked in return
- `CountLikedYou(CountLikedYouRequest) returns (CountLikedYouResponse)`; // Count the number of users who liked the recipient
- `PutDecision(PutDecisionRequest) returns (PutDecisionResponse)`; // Record the decision of the actor to like or pass the recipient

### Technologies Used

- Go: Main programming language for the backend application.
- PostgreSQL: Relational database used for storing user interactions and decisions.
- Docker: Containerization for easy deployment and development environment setup.
- Docker Compose: Orchestrates multiple Docker containers for development and testing environments.
- gRPC: Remote procedure call framework for efficient communication between services.
- Viper: Library for configuration management in Go, used for handling configuration files and environment variables.
- Testify: Testing framework and utilities for writing and running tests in Go, ensuring comprehensive test coverage.
- Makefile: Automation for common tasks like building, testing, and deploying.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed on your system:

[Docker: Install Docker](https://www.docker.com/)
[Docker Compose: Install Docker Compose](https://docs.docker.com/compose/install/)

For local running make sure you have Go installed.

### Installation Instructions

- Clone the Repository

```
git clone https://github.com/valtercfrancisco/muzz-backend-challenge.git
cd <project_directory>/muzz-backend-challenge
```
### How to Run

#### Running with Docker-Compose

Run the following command to start the application using Docker Compose:

```
docker-compose up --build app postgres 
```

This command will build and start the PostgreSQL database and the main application. Running it this way ensures that the 
main app will retry if the database is not ready yet.

#### Running with Locally

First, ensure you have a database running with the following configurations:

```
POSTGRES_USER=admin
POSTGRES_PASSWORD=password
POSTGRES_DB=muzzdb
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

You can use the PostgreSQL service defined in the docker-compose file with the following command:

```
docker-compose up --build postgres 
```

Then, you can run the project with:

```
go run cmd/server/main.go
```

#### Requesting the API

The server listens on port 8089 and can be accessed at the following URL: http://localhost:8089.

To make requests to the server, you can use tools such as [BloomRPC](https://github.com/bloomrpc/bloomrpc), 
[Postman](https://www.postman.com/) or [grpcurl](https://github.com/fullstorydev/grpcurl) using the following payloads:

**ListLikedYou**

```
{
  "recipient_user_id": "00000000-0000-0000-0000-000000000001",
  "pagination_token": ""
}
```

**ListNewLikedYou**

```
{
  "recipient_user_id": "00000000-0000-0000-0000-000000000001",
  "pagination_token": ""
}
```

**CountLikedYou**

```
{
  "recipient_user_id": "00000000-0000-0000-0000-000000000001"
}
```

**ListLikedYou**

```
{
  "actor_user_id": "00000000-0000-0000-0000-000000000001",
  "recipient_user_id": "00000000-0000-0000-0000-000000000004",
  "liked_recipient": true
}
```

## Project Structure

### Explanation of the Directory Structure

```
.
├── Dockerfile
├── Dockerfile.test
├── Makefile
├── README.md
├── README2.md
├── cmd
│   └── server
│       └── main.go
├── config.yaml
├── db-variables.env
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   └── db
│       ├── migrations
│       │   ├── 000001_create_users_table.down.sql
│       │   ├── 000001_create_users_table.up.sql
│       │   ├── 000002_create_likes_table.down.sql
│       │   ├── 000002_create_likes_table.up.sql
│       │   ├── 000003_create_decisions_table.down.sql
│       │   └── 000003_create_decisions_table.up.sql
│       ├── mock
│       │   ├── insert_mock_data.sql
│       │   └── insert_mock_data_random.sql
│       └── setup.go
└── pkg
    ├── proto
    │   ├── explore-service.pb.go
    │   ├── explore-service.proto
    │   └── explore-service_grpc.pb.go
    ├── repository
    │   ├── explore-repository.go
    │   ├── explore-repository_integration_test.go
    │   └── explore-repository_mock.go
    └── service
        ├── explore-service.go
        └── explore-service_test.go
```

### Overview of Important Files

ate development and to separate concerns. It follows Go best practices and guidelines, and it focuses on managing user
interactions and decisions within a social platform. Below is a breakdown of The project structure is organized to facilite 
the directory and file organization:

- **Dockerfile**: Defines the instructions to build the Docker image for the main application.
- **Dockerfile.test**: Specifies the Dockerfile for building the test environment.
- **Makefile**: Contains a quick and easy way to generate grpc codes. Just run: ` make generate_grpc_code`
- **README.md**: Primary README file containing general project information and setup instructions.
- **config.yaml**: Configuration file in YAML format for application settings.
- **db-variables.env**: Environment variables file used by Docker Compose to configure the PostgreSQL database.
- **docker-compose.yml**: Docker Compose configuration file defining services, networks, and volumes for development environment orchestration.
- **go.mod and go.sum**: Go modules files that manage dependencies for the project.

### Key Directories:

- **cmd/server/main.go**: Entry point for the application.
- **internal/config/config.go**: Configuration setup and management.
- **internal/db/**: Database migration scripts and setup logic.
- **pkg/proto/**: Protobuf files and generated Go code for gRPC service and message formats.
- **pkg/repository/**: Data access and persistence logic.
- **pkg/service/**: Business logic for managing user interactions and decisions.


## Architecture and Design

### Overview

The architecture of the project follows a layered approach, separating concerns into distinct modules and components to
ensure modularity, scalability, and maintainability. It leverages Go's concurrency model, gRPC for efficient communication, 
and PostgreSQL as the relational database.

### Key Components

```
    +---------------------+
    |     Presentation    |
    |        Layer        |
    |   (gRPC Service)    |
    +----------+----------+
               |
               | Protobuf
               |
    +----------v----------+
    |    Business Logic   |
    |        Layer        |
    |      (Service)      |
    +----------+----------+
               |
               | Interface
               |
    +----------v----------+
    |     Data Access     |
    |        Layer        |
    |    (Repository)     |
    +----------+----------+
               |
               | Database
               |
    +----------v----------+
    |     PostgreSQL      |
    |   (Data Storage)    |
    +---------------------+

```

#### Presentation Layer

**gRPC Service**: Handles incoming requests and serves as the interface for client-server communication.
Proto Files: Define service methods and message formats (pkg/proto/explore-service.proto).

#### Business Logic Layer

**Service Layer (pkg/service/explore-service.go)**: Implements business rules and logic for managing user decisions and interactions.
Orchestrates data flow between the repository layer and external services.

#### Data Access Layer

**Repository Layer (pkg/repository/explore-repository.go)**: Encapsulates database operations (CRUD) and abstracts away 
database-specific details. Provides interfaces for data manipulation and retrieval.
Includes mock implementations (explore-repository_mock.go) for testing purposes.

#### Configuration and Setup

**Configuration Handling (internal/config/config.go)**: Manages application configuration using Viper, allowing for easy integration of environment variables and configuration files (config.yaml).

#### Database Management

**Database Migrations (internal/db/migrations/):** SQL scripts for database schema management.
**Setup (internal/db/setup.go)**: Initializes database connection and performs setup operations during application startup.


### Data model and Database

One of the biggest decisions for the project was the data model and the database to use for the likes and decisions.
Two ways were considered, a singular likes tables for both types of decisions from likes to passes and a second approach,
using two separate tables, likes and decisions. The likes table should focus solely on recording likes between users.
The decisions table should record all user decisions (like, pass, block, etc.). I ended up choosing the separate tables approach for a few reasons:

- Detailed Tracking: A decisions table can store not only likes but also other interactions such as passes, blocks, or reports, 
providing a more comprehensive view of user interactions.
- Explicit User Decisions: Separating likes from other decisions (e.g., pass, block) can help in clearly understanding user behavior and preferences.
- Indexing: Separate tables allow for more focused indexing strategies, which can improve the performance of read and write operations.
- Scalability: As user interactions grow, having specialized tables helps in scaling the database efficiently.
- Historical Analysis: Storing detailed decision data allows for better historical analysis of user interactions, which can be useful for understanding trends and improving the app’s matching algorithm.
- Future Proofing: It provides a more flexible schema that can easily accommodate new types of user decisions without major refactoring.

Like everything in software engineering this comes with some trade-off:
- Increased Schema Complexity
- More Tables: Managing multiple tables adds complexity to the database schema and the application logic.
- Implementation Effort: More effort is needed to implement and maintain the logic for multiple tables.
- Synchronization: Keeping the data in sync between the likes and decisions tables requires careful management.

#### Data model

Users: 
```
CREATE TABLE IF NOT EXISTS users (
user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
username VARCHAR(255) UNIQUE NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Likes:
```
CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    actor_user_id UUID NOT NULL,
    recipient_user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(actor_user_id, recipient_user_id),
    FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
    FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
);
```

Decisions:
```
CREATE TABLE IF NOT EXISTS decisions (
actor_user_id UUID NOT NULL,
recipient_user_id UUID NOT NULL,
liked_recipient BOOLEAN NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
UNIQUE(actor_user_id, recipient_user_id),
FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
);
```

#### Database choice
For this challenge two DBs were considered, one SQL one noSQL, specifically Postgres or MongoDB.
While both are good approaches for this project I ended up going with a PostgresSQL database for the following reasons:
a straightforward implementation, better performance for complex queries using joins and transactional consistency.

To help with performance when handling large volumes of data, indexes were added for columns frequently used 
(recipient_user_id, created_at):
```
CREATE INDEX IF NOT EXISTS idx_likes_recipient_user_id ON likes(recipient_user_id);
CREATE INDEX IF NOT EXISTS idx_likes_actor_recipient ON likes(actor_user_id, recipient_user_id);
```

## Testing

### Testing Framework Used

This project contains unit tests for the service class and integration tests for the repository class, all leveraged by testify

### How to Run Tests

You can run them with the following command:

```
docker-compose up --build tests
```

## Documentation

Godoc wand code commentaries were used to generate documentation for this project. You can run godoc with the following command:
```
godoc -http=:6060
```

You can then find the doc in the following URL: ´http://localhost:6060/pkg/muzz-backend-challenge/pkg/service/

## Performance

### Performance Considerations

Indexes were added to improve performance on ListLikedYou and ListNewLikedYou but that still leaves PutDecisions.
I considered this endpoint the biggest potential bottleneck for this application since it needs multiple queries on multiple 
tables to fully do its job. When accounting for scale and concurrency, problems can arise when handling multiple concurrent 
requests which can lead to data consistency issues, especially if the same users are involved in multiple decisions at the same time. 

### Optimization Techniques Used

To address the data consistency issues, I used Transactions to ensure consistency and atomicity.
- Data consistency: Due to its nature of inserting data into multiple tables, and it might need to handle many concurrent
requests, we have to be aware of potential race conditions. With this in mind, we can use transactions to keep data consistency:
```
func (service ExploreService) PutDecision(
	_ context.Context,
	request *explore.PutDecisionRequest,
) (*explore.PutDecisionResponse, error) {
	// Start a transaction
	tx, err := service.repository.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the decision into the decision database
	err = service.repository.InsertDecisionTx(tx, request.ActorUserId, request.RecipientUserId, request.LikedRecipient)
	if err != nil {
		return nil, err
	}

	mutualLikes := false

	// If the user liked the recipient, record the like
	if request.LikedRecipient {
		// Insert the like into the like database
		err = service.repository.InsertLikeTx(tx, request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}

		// Check if the recipient also liked the actor
		mutualLikes, err = service.repository.CheckMutualLikeTx(tx, request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}
	} else {
		// Delete the like if the actor passes on the recipient (unmatched)
		err = service.repository.DeleteLikeTx(tx, request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &explore.PutDecisionResponse{MutualLikes: mutualLikes}, nil
}
```

## Scalability

In this section, I want to talk about scalability specifically the question posed in the document for this challenge:

```
Consider how your solution would scale and perform for users who have given
or received a lot of decisions
- We have some users who have been on the platform for many years
and have decided on hundreds of thousands of users.
```

With this in mind I want to take some time to touch on a few techniques we could use to handle hundreds of thousands of users:

- **Indexing & Query Optimization**: I already use a bit of indexing to improve the performance of ListLikedYou and ListNewLikedYou
but we could take this further and even optimize these queries if possible (We could use tools liken EXPLAIN to understand the performance)
- **Partitioning**: A strategy to reduce the amount the data scanned by queries for users with hundreds of interactions 
is to partition our tables based on user ID or time (yearly for example). This would help lighting the load.
- **Asynchronous Processing**: For operations that can be performed asynchronously, such as processing large sets of decisions or saving likes, 
we can use background jobs and batch processing to distribute the load and avoiding.
  - Message Queues: We can use message queues to handle and process tasks asynchronously to process large volumes of user 
interactions.
- **Horizontal Scaling**: We can have multiple instances of the Explore Service running (If need be we can split 
the decision logic to its own microservice). We can further use load balancers to distribute incoming traffic across 
multiple instances.
- **Caching**: We can use in-memory caching solutions like Redis or Memcached to store frequently accessed data.
- **Monitoring and Alerting**: Monitor system performance using tools like Prometheus, Micrometer and Grafana. This way we can 
tweak and make changes to the system as things change.

## Future Improvements

- PutDecision Asynchronously
As mentioned before PutDecisions has the potential to be a big bottleneck especially when taking into account concurrency. 
There are multiple ways handle this from batch operations to caching which are both good solutions,
but I want to focus on one an improvement that, asynchronous processing.
  - Asynchronous Processing: The most important part of the PutDecisions for me is checking if there's a mutual like, as this
    powers another feature, matches (assumption). This is an important a feature that is extremely crucial and user facing,
    therefore it needs to be performant. With this in mind we can see that the main thing we need to do first with endpoint is
    compute if there's a match, regardless of the outcome the rest of the tasks, updating decisions and likes tables, can be
    done in the background in another thread. Here's a potential implementation of this using goroutines:
```
func (service ExploreService) PutDecision(
    ctx context.Context,
    request *explore.PutDecisionRequest,
) (*explore.PutDecisionResponse, error) {
    // Asynchronously insert the decision into the decision database
    go func() {
        err := service.repository.InsertDecision(request.ActorUserId, request.RecipientUserId, request.LikedRecipient)
        if err != nil {
            // Log the error if needed
            log.Printf("failed to insert decision: %v", err)
        }
    }()

    // Asynchronously handle like operations
    if request.LikedRecipient {
        go func() {
            err := service.repository.InsertLike(request.ActorUserId, request.RecipientUserId)
            if err != nil {
                // Log the error if needed
                log.Printf("failed to insert like: %v", err)
            }
        }()
    } else {
        go func() {
            err := service.repository.DeleteLike(request.ActorUserId, request.RecipientUserId)
            if err != nil {
                // Log the error if needed
                log.Printf("failed to delete like: %v", err)
            }
        }()
    }

    // Synchronously check mutual likes
    mutualLikes, err := service.repository.CheckMutualLike(request.ActorUserId, request.RecipientUserId)
    if err != nil {
        return nil, err
    }

    return &explore.PutDecisionResponse{MutualLikes: mutualLikes}, nil
}
```
Here we see that CheckMutualLike is the only method we must do synchronously. With this approach we can return the response
to the mutual like immediately ensuring the performance of matching is fast, and we handle saving the like and logging
the decision later. This can of course be tweaked further. On an even further improvement, we could separate saving the like
and decision to another service and use a queue to ensure everything is saved later. That way we can further decouple these 
actions. If saving the like is imperative we can also ensure that action is done synchronously.
- Consider separating likes and decision into different databases: As the system grows in  user we might want to separate the likes 
and decisions and separate databases and not just tables. With this we could have decisions use a NoSQL as they are
optimized for high write throughput, which is beneficial when you have a high volume of user interactions and decisions 
being recorded continuously.
- Benchmarks
- Load Tests
- CI/CD




