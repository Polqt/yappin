# Yappin - Real-time Chat Application

A modern, real-time chat application built with a SvelteKit frontend and Go backend, featuring user authentication, chat rooms, achievements, and statistics tracking.

## Features

- **Real-time Messaging**: WebSocket-powered instant messaging in chat rooms
- **User Authentication**: Secure login/signup with JWT tokens
- **Chat Rooms**: Create and join public/private rooms
- **User Profiles**: View user stats, achievements, and daily check-ins
- **Achievements System**: Unlock achievements based on activity
- **Responsive Design**: Mobile-friendly interface with TailwindCSS
- **Admin Panel**: Database management with Adminer

## Tech Stack

- **Frontend**: SvelteKit, TypeScript, TailwindCSS, Socket.io-client
- **Backend**: Go, Chi router, PostgreSQL, Gorilla WebSockets
- **Database**: PostgreSQL with migrations via Goose
- **Deployment**: Docker Compose for local development

## Prerequisites

- Go 1.24+
- Node.js 18+
- Docker and Docker Compose
- Git

## Getting Started

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd yappin
   ```

2. **Set up environment variables**

   Copy the example environment file and configure your settings:
   ```bash
   cp server/.env.example server/.env
   ```

   Edit `server/.env` with your database credentials and other settings.

3. **Start the database**
   ```bash
   docker-compose up -d db adminer
   ```

4. **Run database migrations**
   ```bash
   cd server
   go run db/migrations/migrate.go up
   cd ..
   ```

5. **Start the backend server**
   ```bash
   cd server
   go run main.go
   ```

6. **Start the frontend (in a new terminal)**
   ```bash
   cd client
   npm install
   npm run dev
   ```

7. **Access the application**
   - Frontend: http://localhost:5173
   - Adminer (Database UI): http://localhost:8080

## Project Structure

```
yappin/
├── client/          # SvelteKit frontend
├── server/          # Go backend
├── docker-compose.yml
└── README.md
```

## Contributing

We welcome contributions! This project participates in **Hacktoberfest**. See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

[MIT License](LICENSE) - feel free to use this project for your own purposes.

## Support

If you have questions or need help, please open an issue on GitHub.