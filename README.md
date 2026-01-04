# im-chat-golang-react

An Instant messaging application written in Go lang and React with WebSocket support, real-time messaging, and channel-based communication.

## Features

✅ **WebSocket Communication**: Real-time bidirectional communication between clients and server  
✅ **User Management**: Create and manage users with UUID-based identification  
✅ **Direct Messaging**: Send messages directly between users  
✅ **Channel Support**: Create and subscribe to channels for group communication  
✅ **Automatic Subscriber Management**: Channel creators are automatically added as subscribers  
✅ **PostgreSQL Database**: Persistent storage with GORM ORM  
✅ **Transaction Safety**: Database operations use transactions for data consistency  
✅ **RESTful API**: Clean API endpoints for all operations  

## Tech Stack

### Backend
- **Go** - Primary programming language
- **Gin** - HTTP web framework
- **Gorilla WebSocket** - WebSocket implementation
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **Docker** - Containerization

### Database Schema

#### Users Table
- `id` (primary key)
- `created_at`, `updated_at`, `deleted_at`
- `username` (unique, not null)
- `user_uuid` (UUID, unique, not null)
- `email` (optional, unique)
- `first_name`, `last_name` (optional)

#### Channels Table
- `id` (primary key)
- `created_at`, `updated_at`, `deleted_at`
- `channel_uuid` (UUID, unique, indexed)
- `name` (required)
- `description` (optional)
- `created_by` (user UUID)

#### Channel Subscribers (Join Table)
- `channel_id` (foreign key to channels.id)
- `user_id` (foreign key to users.id)

#### Chat Messages Table
- `id` (primary key)
- `sender_id` (UUID, not null, indexed)
- `recipient_id` (UUID, not null, indexed)
- `content` (text, not null)
- `message_type` (default: 'direct')

## Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL
- Docker (optional)

### Environment Variables

Create a `.env` file in the `server` directory:

```env
PG_HOST=localhost
PG_PORT=5432
PG_USER=your_username
PG_PASSWORD=your_password
PG_DB=chat_db
```

### Running with Docker

```bash
cd server
docker-compose up
```

### Running Locally

1. Install dependencies:
```bash
cd server
go mod download
```

2. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### User Endpoints

#### Create a User
**POST** `/user`

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response:**
```json
{
  "message": "User john_doe created successfully"
}
```

#### Get All Users
**GET** `/users`

**Response:**
```json
{
  "users": [
    {
      "ID": 1,
      "username": "john_doe",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe"
    }
  ]
}
```

### Channel Endpoints

#### Create a Channel
**POST** `/channel`

Creates a new channel and automatically adds the creator as a subscriber.

**Request:**
```json
{
  "name": "General Discussion",
  "description": "A channel for general discussions",
  "createdBy": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (200 OK):**
```json
{
  "message": "Channel General Discussion created successfully",
  "channel": {
    "ID": 1,
    "CreatedAt": "2026-01-03T10:00:00Z",
    "UpdatedAt": "2026-01-03T10:00:00Z",
    "channelId": "660e8400-e29b-41d4-a716-446655440000",
    "name": "General Discussion",
    "description": "A channel for general discussions",
    "createdBy": "550e8400-e29b-41d4-a716-446655440000",
    "subscribers": [
      {
        "ID": 1,
        "username": "john_doe",
        "userId": "550e8400-e29b-41d4-a716-446655440000"
      }
    ]
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or missing createdBy field
- `500 Internal Server Error` - Failed to create channel (e.g., user not found)

**Notes:**
- The `createdBy` user must exist in the database before creating a channel
- Channel creator is automatically added to the subscribers list
- Channel creation is performed in a database transaction for atomicity

#### Get a Channel
**GET** `/channel/:channelId`

Retrieves a channel by its UUID along with all its subscribers.

**Response:**
```json
{
  "channel": {
    "ID": 1,
    "channelId": "660e8400-e29b-41d4-a716-446655440000",
    "name": "General Discussion",
    "description": "A channel for general discussions",
    "subscribers": [...]
  }
}
```

#### Get User's Channels
**GET** `/user/:userId/channels`

Retrieves all channels that a user is subscribed to.

**Response:**
```json
{
  "channels": [
    {
      "ID": 1,
      "channelId": "660e8400-e29b-41d4-a716-446655440000",
      "name": "General Discussion",
      "subscribers": [...]
    }
  ]
}
```

### Message Endpoints

#### Get Messages Between Users
**GET** `/messages/:userId/:recipientId`

Retrieves all messages exchanged between two users, ordered by creation time.

**Response:**
```json
{
  "messages": [
    {
      "senderId": "550e8400-e29b-41d4-a716-446655440000",
      "recipientId": "660e8400-e29b-41d4-a716-446655440000",
      "content": "Hello!",
      "messageType": "direct"
    }
  ]
}
```

### WebSocket Endpoint

#### Connect to WebSocket
**GET** `/connect/:clientId`

Upgrades HTTP connection to WebSocket for real-time communication.

**WebSocket Message Format:**
```json
{
  "type": "message",
  "content": "Hello, World!",
  "channelId": "channel-uuid",
  "recipientId": "user-uuid",
  "senderId": "sender-uuid"
}
```

### Health Check

#### Health Check
**GET** `/health`

**Response:**
```json
{
  "status": "UP"
}
```

## Example Usage with cURL

### Create a User
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Create a Channel
```bash
curl -X POST http://localhost:8080/channel \
  -H "Content-Type: application/json" \
  -d '{
    "name": "General Discussion",
    "description": "A channel for general discussions",
    "createdBy": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

### Get a Channel
```bash
curl http://localhost:8080/channel/660e8400-e29b-41d4-a716-446655440000
```

### Get User's Channels
```bash
curl http://localhost:8080/user/550e8400-e29b-41d4-a716-446655440000/channels
```

### Get All Users
```bash
curl http://localhost:8080/users
```

### Get Messages Between Users
```bash
curl http://localhost:8080/messages/550e8400-e29b-41d4-a716-446655440000/660e8400-e29b-41d4-a716-446655440000
```

## Project Structure

```
im-chat-golang-react/
├── server/
│   ├── db/
│   │   └── postgres/
│   │       ├── pg_conn.go      # Database connection and migrations
│   │       └── pg_db.go        # Database query methods
│   ├── models/
│   │   └── models.go           # Data models (User, Channel, Message)
│   ├── websocket/
│   │   ├── ws_server.go        # WebSocket server and HTTP handlers
│   │   ├── ws_client.go        # WebSocket client management
│   │   └── ws_hub.go           # WebSocket hub for broadcasting
│   ├── main.go                 # Application entry point
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── docker-compose.yaml
│   └── .env
└── README.md
```

## Architecture

### WebSocket Server
- Handles WebSocket connections and upgrades
- Manages client connections through a Hub
- Routes messages to appropriate recipients
- Persists messages to the database

### Database Layer
- Uses GORM for object-relational mapping
- Implements repository pattern for data access
- Supports transactions for complex operations
- Auto-migrates database schema on startup

### Models
- **User**: User account information
- **Channel**: Group communication channels with many-to-many relationship to users
- **ChatMessage**: Persistent message storage
- **Message**: Wire format for WebSocket communication

## Development

### Building the Application
```bash
cd server
go build -o chat-server
```

### Running Tests
```bash
go test ./...
```

### Code Verification
The codebase has been verified to compile successfully with no errors.

## Future Enhancements

- [ ] Add endpoint to add/remove subscribers from a channel
- [ ] Add endpoint to list all channels
- [ ] Add pagination for channel and message lists
- [ ] Add search/filter functionality for channels and users
- [ ] Add channel permissions/roles (admin, moderator, member)
- [ ] Add WebSocket events for channel subscriptions
- [ ] Implement message read receipts
- [ ] Add file/image sharing capabilities
- [ ] Implement user presence (online/offline status)
- [ ] Add authentication and authorization (JWT)
- [ ] Build React frontend

## License

This project is licensed under the terms specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
