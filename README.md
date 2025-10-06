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

### For Contributors (Hacktoberfest)

1. **Fork the repository**
   - Click the "Fork" button on the top right of this repository
   - This creates a copy of the repository in your GitHub account

2. **Clone your fork**
   ```bash
   git clone https://github.com/your-username/yappin.git
   cd yappin
   ```

3. **Set up upstream remote** (optional, for staying updated)
   ```bash
   git remote add upstream https://github.com/Polqt/yappin.git
   git fetch upstream
   ```

### Local Development Setup

4. **Set up environment variables**

   Copy the example environment file and configure your settings:
   ```bash
   cp server/.env.example server/.env
   ```

   Edit `server/.env` with your database credentials and other settings.

5. **Start the database**
   ```bash
   docker-compose up -d db adminer
   ```

6. **Run database migrations**
   ```bash
   cd server
   go run db/migrations/migrate.go up
   cd ..
   ```

7. **Start the backend server**
   ```bash
   cd server
   go run main.go
   ```

8. **Start the frontend (in a new terminal)**
   ```bash
   cd client
   npm install
   npm run dev
   ```

9. **Access the application**
   - Frontend: http://localhost:5173
   - Adminer (Database UI): http://localhost:8080

## Making Contributions

After setting up your development environment:

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or for bug fixes:
   git checkout -b fix/issue-description
   ```

2. **Make your changes**
   - Follow the coding standards in [CONTRIBUTING.md](CONTRIBUTING.md)
   - Test your changes thoroughly
   - Update documentation if needed

3. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   # Use conventional commit format (see CONTRIBUTING.md)
   ```

4. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request**
   - Go to your fork on GitHub
   - Click "Compare & pull request"
   - Fill out the PR template with details about your changes
   - Submit the PR for review

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