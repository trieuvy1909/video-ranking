# System Architecture Diagram

## Overview

The Video Ranking application utilizes a microservices architecture to efficiently manage and rank videos based on user interactions. Below is a high-level architecture diagram illustrating the interaction between different components:

- **API Gateway**: Serves as the entry point for client requests, routing them to appropriate microservices.
- **Video Service**: Handles video-related operations such as upload, retrieval, and ranking.
- **User Service**: Manages user information and authentication.
- **Interaction Service**: Records and processes user interactions with videos.
- **Database**: Stores persistent data related to videos, users, and interactions.
- **Caching System**: Improves performance by storing frequently accessed data.

![Architecture Diagram](architecture-diagram.svg)

## Interaction
The API Gateway routes requests to the appropriate microservice based on the endpoint. Each microservice interacts with the database and caching system to perform its operations efficiently. The architecture ensures scalability and flexibility, allowing for easy integration of additional services in the future.

For detailed documentation, refer to the `docs/document.md` file.

## Using Queue for Sequential Processing

The system employs a queue to sequentially process interaction events, ensuring consistency and high performance. When a new interaction event is created, it is added to the queue via the `EnqueueInteractionEvent` method in `QueueServices`. This queue is managed by `QueueConsumer`, a component responsible for processing events in the order they are added.

### How It Works

1. **Queue Initialization and Management**: The queue is initialized when the application starts and is managed by `QueueServices`. Interaction events are added to this queue for sequential processing.

2. **Adding Events to the Queue**: When a new interaction event is created, it is added to the queue using the `EnqueueInteractionEvent` method. This ensures that all events are processed in the order they are created.

3. **Event Processing**: `QueueConsumer` retrieves each event from the queue and processes it. During this process, functions in `videoService` are called to update data such as scores, likes, comments, and views of the video.

### Benefits

- **Consistency**: Ensures that events are processed in order, avoiding data conflicts.
- **High Performance**: Reduces system load by processing events sequentially and in a controlled manner.

Using a queue helps the system operate smoothly and efficiently, ensuring that data is always updated accurately and promptly.