# Poll App

A full-stack polling application built with Go (backend) and React/TypeScript (frontend).

## Features

- ✅ User sign up & sign in (JWT-based authentication)
- ✅ Create, read, update, and delete polls
- ✅ Vote on polls (one vote per user per poll)
- ✅ View vote counts after voting
- ✅ Click on vote counts to see who voted for each option
- ✅ Beautiful, responsive UI

## Tech Stack

### Backend
- **Language**: Go
- **Router**: httprouter
- **ORM**: ent (entgo.io)
- **Database**: SQLite
- **Authentication**: JWT

### Frontend
- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **HTTP Client**: Axios
- **Routing**: React Router v6

## Project Structure

```
poll_app/
├── backend/
│   ├── main.go              # Entry point
│   ├── go.mod               # Go dependencies
│   ├── handlers/
│   │   └── handlers.go      # API handlers
│   └── ent/
│       └── schema/          # Database schemas
│           ├── user.go
│           ├── poll.go
│           ├── option.go
│           └── vote.go
└── frontend/
    ├── package.json
    ├── vite.config.ts
    ├── index.html
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── index.css
        ├── types/
        │   └── index.ts
        ├── services/
        │   └── api.ts
        ├── context/
        │   └── AuthContext.tsx
        ├── components/
        │   └── Navbar.tsx
        └── pages/
            ├── Login.tsx
            ├── SignUp.tsx
            ├── Polls.tsx
            ├── CreatePoll.tsx
            ├── EditPoll.tsx
            └── PollDetail.tsx
```

## Database Design

### Users Table
- id (primary key)
- username (unique)
- email (unique)
- password (hashed)
- created_at

### Polls Table
- id (primary key)
- title
- description
- creator_id (foreign key to users)
- created_at
- updated_at

### Options Table
- id (primary key)
- text
- poll_id (foreign key to polls)

### Votes Table
- id (primary key)
- user_id (foreign key to users)
- option_id (foreign key to options)
- created_at
- Unique constraint on (user_id, option_id)

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Download dependencies:
   ```bash
   go mod tidy
   ```

3. Generate ent code:
   ```bash
   go generate ./ent
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

The server will start on http://localhost:8080

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will be available at http://localhost:3000

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Register a new user
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/me` - Get current user info

### Polls
- `GET /api/polls` - List all polls
- `POST /api/polls` - Create a new poll
- `GET /api/polls/:id` - Get a specific poll
- `PUT /api/polls/:id` - Update a poll
- `DELETE /api/polls/:id` - Delete a poll

### Voting
- `POST /api/polls/:id/vote` - Vote on a poll
- `GET /api/options/:id/voters` - Get voters for an option

## Usage

1. Start both the backend and frontend servers
2. Open http://localhost:3000 in your browser
3. Sign up for a new account
4. Create polls with multiple options
5. Vote on polls
6. Click on vote counts to see who voted for each option
