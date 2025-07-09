# Fitness Framework API

A Go-based backend API for managing fitness exercises, equipment, and muscle groups, using MongoDB as the database. The API initializes with hardcoded data if the database is empty.

## Features
- Stores and retrieves exercises, equipment, and muscle groups
- Populates MongoDB with initial data on first run
- Simple setup using Docker for MongoDB

## Prerequisites
- Go 1.20 or newer
- Docker (for running MongoDB)

## Running MongoDB with Docker

You can run MongoDB in a Docker container for easy setup and teardown. Below are the most common commands:

### Start MongoDB Container
```bash
docker run --name fitness-framework-db -d -p 27017:27017 mongo:6
```
- This starts MongoDB in the background, accessible at `localhost:27017`.

### Stop the Container
```bash
docker stop fitness-framework-db
```

### Remove the Container
```bash
docker rm fitness-framework-db
```

### (Optional) Persist Data with a Volume
To keep your MongoDB data between container restarts, use a Docker volume:
```bash
docker run --name fitness-framework-db -d -p 27017:27017 -v fitness-framework-db-data:/data/db mongo:6
```
- This will store MongoDB data in a Docker-managed volume named `fitness-framework-db-data`.

## Getting Started

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd fitness-framework-api
```

### 2. Start MongoDB with Docker
See the section above for detailed Docker instructions.

### 3. Configure the API (Optional)
By default, the API connects to MongoDB at `mongodb://localhost:27017` and uses the database name you provide at startup. If you need to change the connection string, update the `MongoURI` constant in `internal/mongodb/mongodb.go`.

### 4. Build and Run the API

```bash
go build -o fitness-framework-api
./fitness-framework-api
```

- The API will connect to MongoDB and populate the database with initial data if the collections are empty.

### 5. API Usage

The API exposes endpoints for retrieving exercises, equipment, and muscle groups. (See your API documentation or code for available endpoints.)

## Initial Data Population
- On first run, if the `exercises` collection is empty, the API will populate it (and related collections) with hardcoded data from Go constants.

## Troubleshooting
- **MongoDB connection errors:** Ensure the Docker container is running and accessible at `localhost:27017`.
- **Port conflicts:** Make sure nothing else is using port 27017.
- **Data not populating:** The API only populates data if the collections are empty. Drop the collections manually if you want to re-initialize.
