# Video Ranking Application

## Introduction
The Video Ranking application is designed to manage and rank videos based on user interactions. It leverages a microservices architecture with separate services for API, database, and caching.

## Prerequisites
- Docker and Docker Compose installed on your machine.
- Go programming language installed.

## Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   ```

2. **Navigate to the project directory**
   ```bash
   cd video-ranking
   ```

3. **Start the services using Docker Compose**
   ```bash
   docker-compose up
   ```

4. **Access the application**
   Open your browser and go to `http://localhost:8080/swagger/index.html` to view the API documentation.

## Running Tests
To ensure all components are working correctly, run the test suite:
```bash
go test ./...
```

## Conclusion
This README provides a quick overview of setting up and running the Video Ranking application. For detailed documentation, refer to the `docs/document.md` file.

## System Architecture Diagram

For a high-level understanding of how different components interact within the microservice architecture, refer to the `docs/document.md` file, which includes a detailed architecture diagram.