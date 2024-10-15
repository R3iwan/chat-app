
## Getting Started

### Prerequisites

- Go 1.16+
- Docker
- Node.js (for frontend development)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/chat-application.git
    cd chat-application
    ```

2. Set up the environment variables:
    ```sh
    cp .env.example .env
    ```

3. Start the services using Docker Compose:
    ```sh
    docker-compose up --build
    ```

4. Navigate to the `frontend` directory and install dependencies:
    ```sh
    cd frontend
    npm install
    ```

5. Start the frontend development server:
    ```sh
    npm start
    ```

## Project Modules

### Backend

- **User Service**: Handles user authentication and registration.
  - [cmd/user-service/main.go](cmd/user-service/main.go)
  - [internal/user/login.go](internal/user/login.go)
  - [internal/user/register.go](internal/user/register.go)

- **Message Handling**: Manages sending and receiving messages.
  - [internal/message/message.go](internal/message/message.go)
  - [internal/chat/handlers.go](internal/chat/handlers.go)

- **Chat Rooms**: Manages chat rooms and user connections.
  - [internal/chat/rooms.go](internal/chat/rooms.go)
  - [internal/chat/private.go](internal/chat/private.go)

- **Middleware**: JWT authentication middleware.
  - [internal/middleware/jwt.go](internal/middleware/jwt.go)

- **Database**: PostgreSQL integration.
  - [internal/db/postgres.go](internal/db/postgres.go)

### Frontend

- **Chat Interface**: Implements the chat UI.
  - [frontend/chat.js](frontend/chat.js)
  - [frontend/index.html](frontend/index.html)

## Usage

1. Open your browser and navigate to `http://localhost:8080`.
2. Register or log in with your credentials.
3. Start chatting with other users.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
