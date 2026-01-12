<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React">
  <img src="https://img.shields.io/badge/TypeScript-5.0-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/SQLite-Database-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
</p>

<h1 align="center">PollApp</h1>

<p align="center">
  <strong>A modern, full-stack polling application</strong><br>
  Create polls, gather votes, and see real-time results with voter transparency.
</p>

<p align="center">
  <a href="https://poll-app-frontend-ylqk.onrender.com">View Demo</a>
  &nbsp;&middot;&nbsp;
  <a href="#quick-start">Quick Start</a>
  &nbsp;&middot;&nbsp;
  <a href="#api-reference">API Docs</a>
</p>

---

## Overview

PollApp is a complete polling solution built with **Go** on the backend and **React** on the frontend. It demonstrates modern web development practices including JWT authentication, RESTful API design, and responsive UI/UX.

### Key Features

| Feature | Description |
|---------|-------------|
| **User Authentication** | Secure sign-up and login with JWT tokens |
| **Poll Management** | Create, edit, and delete polls with multiple options |
| **Voting System** | Vote on polls with ability to change or clear your vote anytime |
| **Real-time Updates** | Poll results refresh automatically every 3 seconds |
| **Vote Notifications** | Poll creators get notified when someone changes or clears their vote |
| **Poll Edit Alerts** | Voters see a notification when a poll they voted on is modified |
| **Owner Vote Visibility** | Poll creators can see vote counts without voting themselves |
| **Voter Transparency** | Click on any vote count to see who voted for that option |
| **Responsive Design** | Modern teal/navy theme that works on all devices |

---

## Quick Start

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)

### 1. Clone the Repository

```bash
git clone https://github.com/bPratyush/poll_app.git
cd poll_app
```

### 2. Start the Backend

```bash
cd backend
go mod tidy
go run main.go
```

Server runs at `http://localhost:8080`

### 3. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

App runs at `http://localhost:3000`

### 4. Open in Browser

Navigate to `http://localhost:3000`, create an account, and start polling!

---

## Tech Stack

<table>
<tr>
<td width="50%" valign="top">

### Backend
- **Go** - Fast, compiled language
- **httprouter** - Lightweight HTTP router
- **ent** - Type-safe ORM by Facebook
- **SQLite** - Embedded database
- **JWT** - Stateless authentication
- **bcrypt** - Password hashing

</td>
<td width="50%" valign="top">

### Frontend
- **React 18** - UI component library
- **TypeScript** - Type-safe JavaScript
- **Vite** - Next-gen build tool
- **Axios** - HTTP client
- **React Router v6** - Client-side routing
- **CSS Variables** - Theming system

</td>
</tr>
</table>

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                 React Frontend                      │    │
│  │         (Static files served via CDN)               │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS (REST API)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         Server                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                  Go Backend                         │    │ 
│  │     httprouter │ JWT Auth │ ent ORM │ CORS          │    │
│  └─────────────────────────────────────────────────────┘    │
│                              │                              │
│                              ▼                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                 SQLite Database                     │    │
│  │        Users │ Polls │ Options │ Votes              │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## Database Schema

<details>
<summary><strong>View Entity Relationships</strong></summary>

```
Relationships:
  • User (1) ──────► (N) Poll        : User creates many polls
  • User (1) ──────► (N) Vote        : User casts many votes
  • User (1) ──────► (N) Notification: User receives many notifications
  • Poll (1) ──────► (N) PollOption  : Poll has many options
  • PollOption (1) ► (N) Vote        : Option receives many votes
  • Vote references both User and PollOption (unique constraint)
```

</details>

<details>
<summary><strong>View Table Definitions</strong></summary>

#### Users
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| username | VARCHAR | UNIQUE, NOT NULL |
| email | VARCHAR | UNIQUE, NOT NULL |
| password | VARCHAR | NOT NULL (hashed) |
| created_at | TIMESTAMP | DEFAULT NOW |

#### Polls
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| title | VARCHAR | NOT NULL |
| description | TEXT | NULLABLE |
| creator_id | INTEGER | FOREIGN KEY → users |
| created_at | TIMESTAMP | DEFAULT NOW |
| updated_at | TIMESTAMP | DEFAULT NOW |

#### PollOptions
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| text | VARCHAR | NOT NULL |
| poll_id | INTEGER | FOREIGN KEY → polls (CASCADE) |

#### Votes
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| user_id | INTEGER | FOREIGN KEY → users |
| option_id | INTEGER | FOREIGN KEY → poll_options (CASCADE) |
| created_at | TIMESTAMP | DEFAULT NOW |
| | | UNIQUE(user_id, option_id) |

#### Notifications
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY |
| user_id | INTEGER | FOREIGN KEY → users |
| message | VARCHAR | NOT NULL |
| type | VARCHAR | DEFAULT 'vote_changed' |
| poll_id | INTEGER | NULLABLE |
| read | BOOLEAN | DEFAULT FALSE |
| created_at | TIMESTAMP | DEFAULT NOW |

</details>

---

## API Reference

Base URL: `https://poll-app-backend-lj26.onrender.com/api`

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/signup` | Register new user |
| `POST` | `/api/auth/login` | Login and get token |
| `GET` | `/api/auth/me` | Get current user |

### Polls

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/polls` | List all polls |
| `POST` | `/api/polls` | Create a poll |
| `GET` | `/api/polls/:id` | Get poll details |
| `PUT` | `/api/polls/:id` | Update a poll |
| `DELETE` | `/api/polls/:id` | Delete a poll |

### Voting

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/polls/:id/vote` | Vote on a poll (or change vote) |
| `DELETE` | `/api/polls/:id/vote` | Clear/remove your vote |
| `GET` | `/api/options/:id/voters` | Get voters for option |

### Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/notifications` | List all notifications |
| `GET` | `/api/notifications/unread-count` | Get unread count |
| `PUT` | `/api/notifications/:id/read` | Mark notification as read |
| `POST` | `/api/notifications/mark-all-read` | Mark all as read |

<details>
<summary><strong>View Request/Response Examples</strong></summary>

#### Sign Up
```http
POST /api/auth/signup
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

#### Create Poll
```http
POST /api/polls
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Favorite Programming Language?",
  "description": "Vote for your preferred language",
  "options": ["Go", "Python", "JavaScript", "Rust"]
}
```

#### Vote
```http
POST /api/polls/1/vote
Authorization: Bearer <token>
Content-Type: application/json

{
  "option_id": 1
}
```

</details>

---

## Project Structure

```
poll_app/
├── backend/
│   ├── main.go              # Application entry point
│   ├── Dockerfile           # Container configuration
│   ├── handlers/
│   │   └── handlers.go      # API route handlers
│   └── ent/
│       └── schema/          # Database models
│           ├── user.go
│           ├── poll.go
│           ├── option.go
│           ├── vote.go
│           └── notification.go
│
├── frontend/
│   ├── src/
│   │   ├── pages/           # Route components
│   │   ├── components/      # Reusable UI components (Navbar with notifications)
│   │   ├── context/         # React context (auth)
│   │   ├── services/        # API client
│   │   └── types/           # TypeScript definitions
│   ├── index.html
│   └── vite.config.ts
│
├── render.yaml              # Deployment configuration
└── README.md
```

---

## Deployment

The application is deployed on [Render](https://render.com) with the following setup:

| Service | Type | URL |
|---------|------|-----|
| Backend | Docker Web Service | [poll-app-backend-lj26.onrender.com](https://poll-app-backend-lj26.onrender.com) |
| Frontend | Static Site | [poll-app-frontend-ylqk.onrender.com](https://poll-app-frontend-ylqk.onrender.com) |

<details>
<summary><strong>Deploy Your Own Instance</strong></summary>

### Using Render Blueprint (Recommended)

1. Fork this repository
2. Go to [Render Dashboard](https://dashboard.render.com)
3. Click **New** → **Blueprint**
4. Connect your forked repository
5. Render auto-detects `render.yaml` and deploys both services

### Environment Variables

#### Backend
| Variable | Description |
|----------|-------------|
| `PORT` | Server port (default: 8080) |
| `DATABASE_PATH` | SQLite file path |
| `JWT_SECRET` | Token signing secret |
| `FRONTEND_URL` | CORS allowed origin |

#### Frontend
| Variable | Description |
|----------|-------------|
| `VITE_API_URL` | Backend API URL |

</details>

---

## Development

### Running Tests

```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm test
```

### Building for Production

```bash
# Backend
cd backend
go build -o poll_app

# Frontend
cd frontend
npm run build
```

---

## License

This project is open source and available under the [MIT License](LICENSE).

---

<p align="center">
  <strong>Built with Go and React</strong><br>
  <sub>Pratyush Bindal</sub>
</p>
