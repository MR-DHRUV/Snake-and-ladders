# Snakes and Ladders Game

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
  - [System Architecture](#system-architecture)
  - [Kubernetes Architecture](#kubernetes-architecture)
- [Frontend](#frontend)
- [Backend](#backend)
- [Game Flow](#game-flow)
- [Deployment](#deployment)
  - [Prerequisites](#prerequisites)
  - [Deployment Steps](#deployment-steps)
- [Development Setup](#development-setup)

## Overview

This is a modern multiplayer Snakes and Ladders game with real-time gameplay. Players can create games, join existing games and chat with other players.

## Architecture

### System Architecture

```mermaid
graph TD
    Client[Client Browser] --> Frontend[Frontend React App]
    Frontend --> |HTTP Requests| BackendREST[Backend REST API]
    Frontend --> |WebSocket| BackendWS[Backend WebSocket]
    BackendREST --> MongoDB
    BackendWS --> MongoDB
    BackendWS --> CM[Connection Manager]
    CM --> |Broadcasts| BackendWS
```

### Kubernetes Architecture

```mermaid
graph TD
    subgraph "Kubernetes Cluster"
        Ingress[Ingress Controller]
        
        subgraph "Frontend Deployment"
            FrontendPod1[Frontend Pod 1]
            FrontendPod2[Frontend Pod 2]
            FrontendService[Frontend Service]
        end
        
        subgraph "Backend Deployment"
            BackendPod1[Backend Pod 1]
            BackendPod2[Backend Pod N]
            BackendService[Backend Service]
            BackendHPA[Horizontal Pod Autoscaler]
        end
        
        subgraph "MongoDB StatefulSet"
            MongoPod[MongoDB Pod]
            MongoService[MongoDB Service]
        end
        
        Ingress --> FrontendService
        Ingress --> BackendService
        FrontendService --> FrontendPod1
        FrontendService --> FrontendPod2
        BackendService --> BackendHPA
        BackendHPA --> BackendPod1
        BackendHPA --> BackendPod2
        BackendPod1 --> MongoService
        BackendPod2 --> MongoService
        MongoService --> MongoPod
    end
    
    User[User] --> Ingress
```

The Kubernetes architecture deploys the application using:
- **Ingress:** Manages external access to services
- **Deployments:** For frontend and backend with multiple replicas
- **Services:** Expose frontend and backend pods
- **StatefulSet:** Manages MongoDB for persistent data storage
- **Horizontal Pod Autoscalers (HPA):** Automatically scales pods based on CPU usage
- **ConfigMaps/Secrets:** Store configuration and sensitive data

## Frontend
The frontend is developed with React and TypeScript, ensuring a robust, type-safe, and modern development experience. It leverages Vite for lightning-fast development and optimized production builds. The game board is rendered using HTML5 Canvas, while WebSockets enable real-time communication with the backend. The implementation follows industry best practices, featuring proactive error handling and reconnection strategies to deliver a seamless and resilient user experience.

### Technology Stack
- **Framework:** React with TypeScript
- **Build Tool:** Vite
- **State Management:** React Context API
- **UI Components:** Shadcn


## Backend
The backend is built with Go, offering exceptional performance and concurrency for real-time gameplay. It runs two dedicated HTTP servers—one serving the RESTful API and another handling WebSocket connections. The REST API manages user authentication, as well as the creation and retrieval of past games, while real-time gameplay, dice rolls, and in-game chat are seamlessly facilitated through WebSockets.

MongoDB is used for data persistence, providing scalability, flexibility, and high availability. The backend efficiently leverages goroutines to enable concurrent handling of multiple game sessions and player interactions, ensuring a smooth experience even under heavy load.

For authentication, the system employs JWT with RSA signing, integrating Google OAuth for secure and streamlined user login.

The application follows a multi-layered architecture—comprising transport, service, model, repository, controller, and handler layers—to ensure clear separation of concerns and maintainability. The design adheres to sound system design principles and SOLID best practices, promoting clean, modular, and testable code.

### Technology Stack
- **Language:** Go
- **Web Framework:** Gorilla Mux for routing
- **WebSockets:** Gorilla WebSocket
- **Database:** MongoDB with official Go driver
- **Authentication:** JWT with RSA signing

### API Endpoints

```mermaid
graph TD
    subgraph "WebSocket Endpoints"
        WSGame[/game - Game WebSocket/]
    end
    
    subgraph "WebSocket Actions"
        JoinGame[joinGame]
        StartGame2[startGame]
        NextTurn[nextTurn]
        ChatMessage[chatMessage]
        RestartGame[restartGame]
    end
    
    WSGame --> JoinGame
    WSGame --> StartGame2
    WSGame --> NextTurn
    WSGame --> ChatMessage
    WSGame --> RestartGame
```

```mermaid
graph TD
    subgraph "REST API"
        Auth[/auth - Google OAuth/]
        User[/user - Get User Info/]
        CreateGame[/game - Create New Game/]
        GetGames[/games - Get Past Games/]
        StartGame[/game/start - Start Game/]
    end
```


### WebSockets

```mermaid
sequenceDiagram
    participant Client
    participant Server
    participant ConnectionManager
    participant GameService
    
    Client->>Server: Connect to WebSocket
    Server->>ConnectionManager: Register connection
    Client->>Server: Send action (joinGame)
    Server->>GameService: Process join game
    GameService->>Server: Return game state
    Server->>ConnectionManager: Broadcast to all game players
    ConnectionManager->>Client: Send updated game state
```

## Game Flow

```mermaid
stateDiagram-v2
    [*] --> Created: User creates game
    Created --> InProgress: Host starts game
    InProgress --> PlayerTurn: Player's turn
    PlayerTurn --> DiceRoll: Roll dice
    DiceRoll --> MovePlayer: Update position
    MovePlayer --> CheckWin: Check win condition
    CheckWin --> InProgress: No winner
    CheckWin --> Finished: Winner(s) found
    Finished --> Created: Restart game
    Created --> [*]: Game abandoned
```

## Deployment

### Prerequisites
- Kubernetes cluster
- kubectl configured
- Ingress controller set up

### Deployment Steps

1. Build and push Docker images:
```bash
# From project root
./deploy.sh
```

This script:
- Starts a local Docker registry in the cluster for securely managing images
- Builds frontend and backend Docker images
- Pushes images to the local docker registry
- Applies Kubernetes manifests:
  - namespace.yml
  - secrets.yml
  - mongo/statefulset.yml, service.yml
  - backend/deployment.yml, service.yml, hpa.yml
  - frontend/deployment.yml, service.yml
  - ingress.yml

## Development Setup

### Prerequisites
- Google Cloud Project with OAuth credentials configured for localhost.
- MongoDB instance running locally or accessible remotely.

Clone the repository and follow the steps below to set up the development environment.

### Backend Setup
1. Create a `.env` file in the `backend` directory from the provided `.env.template` and fill in the required environment variables.

2. Generate RSA keys for JWT authentication:
```bash
cd backend/keys
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
```

3. Start the application
```bash
cd backend
go mod download
go run main.go
```

### Frontend Setup

1. Create a `.env` file in the `frontend` directory from the provided `.env.template` and fill in the required environment variables.

2. Start the application
```bash
cd frontend
pnpm install
pnpm run dev
```

4. **Access the application**
- Frontend: http://localhost:5000
- Backend REST API: http://localhost:8081
- Backend WebSocket: ws://localhost:9999