# Yappin Chat Application - Frontend

A modern real-time chat application built with SvelteKit and TailwindCSS.

## Features

- 🔐 User Authentication (Login/Signup)
- 💬 Real-time WebSocket Chat
- 🏆 Leaderboard System
- 👤 User Profiles with Statistics
- 🎯 Achievement System
- 📊 Activity Tracking
- 🎨 Modern UI with TailwindCSS
- 📱 Responsive Design

## Tech Stack

- **Framework**: SvelteKit 2.x
- **Styling**: TailwindCSS 4.x
- **Language**: TypeScript
- **Icons**: Lucide Svelte
- **Real-time**: WebSocket
- **Build Tool**: Vite

## Getting Started

### Prerequisites

- Node.js 18+
- npm, pnpm, or yarn

### Installation

1. Install dependencies:

```bash
npm install
# or
pnpm install
# or
yarn install
```

2. Create environment file:

```bash
cp .env.example .env
```

3. Update `.env` with your API URLs:

```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

### Development

Start the development server:

```bash
npm run dev
# or
npm run dev -- --open  # Opens browser automatically
```

The app will be available at `http://localhost:5173`

### Building for Production

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

## Project Structure

```
client/
├── src/
│   ├── lib/
│   │   ├── components/      # Reusable components
│   │   │   ├── chat/        # Chat-related components
│   │   │   ├── common/      # Common UI components
│   │   │   ├── layout/      # Layout components (Header, Footer)
│   │   │   └── profile/     # Profile-related components
│   │   ├── constants/       # API endpoints and constants
│   │   ├── middleware/      # Route guards and middleware
│   │   ├── types/           # TypeScript type definitions
│   │   └── utils/           # Utility functions
│   ├── routes/              # SvelteKit routes
│   │   ├── dashboard/       # Dashboard pages
│   │   ├── login/           # Login page
│   │   ├── profile/         # Profile page
│   │   ├── room/            # Room chat pages
│   │   └── signup/          # Signup page
│   ├── services/            # API service layers
│   └── stores/              # Svelte stores (state management)
├── static/                  # Static assets
└── package.json
```

## Key Features

### Authentication

- Secure JWT-based authentication
- Cookie-based session management
- Protected routes with middleware

### Real-time Chat

- WebSocket connection for instant messaging
- Room-based chat system
- Message history
- User presence indicators

### User Profiles

- Personal statistics dashboard
- Achievement badges
- Activity graphs
- Leaderboard rankings

### UI Components

- Responsive design for all devices
- Dark/light theme support
- Accessible components
- Loading states and error handling

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run check` - Run Svelte type checking
- `npm run lint` - Lint code with ESLint
- `npm run format` - Format code with Prettier

## Environment Variables

| Variable       | Description     | Default                 |
| -------------- | --------------- | ----------------------- |
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_WS_URL`  | WebSocket URL   | `ws://localhost:8080`   |

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

MIT
