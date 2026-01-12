# PollHub - Full-Stack Polling Application

A comprehensive polling application designed to demonstrate proficiency in modern web development technologies, specifically Go for backend services and React for frontend interfaces.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Technology Stack](#technology-stack)
4. [System Architecture](#system-architecture)
5. [Database Design](#database-design)
6. [API Documentation](#api-documentation)
7. [Installation Guide](#installation-guide)
8. [Deployment](#deployment)
9. [Project Structure](#project-structure)

---

## Project Overview

PollHub is a full-stack web application that enables users to create, manage, and participate in polls. The application implements secure JWT-based authentication, real-time vote tracking, and voter transparency features.

### Objectives

- Demonstrate proficiency in Go programming language and RESTful API design
- Implement a responsive React frontend with TypeScript
- Utilize an ORM (ent) for database operations
- Apply industry-standard authentication mechanisms (JWT)
- Deploy a production-ready application

---

## Features

### Authentication
- User registration with secure password hashing (bcrypt)
- JWT-based session management
- Protected routes and API endpoints

### Poll Management
- Create polls with multiple options
- Edit existing polls (restricted to poll creators)
- Delete polls (restricted to poll creators)
- List all available polls with pagination support

### Voting System
- One vote per user per poll (enforced at database level)
- Real-time vote count display after participation
- Voter transparency: view all users who selected each option

### User Interface
- Responsive design for desktop and mobile devices
- Modern, accessible component architecture
- Loading states and error handling

---

## Technology Stack

### Backend

| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.21+ | Server-side logic |
| Router | httprouter | HTTP request routing |
| ORM | ent (entgo.io) | Database operations |
| Database | SQLite | Data persistence |
| Authentication | golang-jwt/v5 | JWT token management |
| Password Hashing | bcrypt | Secure credential storage |
| CORS | rs/cors | Cross-origin resource sharing |

### Frontend

| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | React 18 | UI component library |
| Language | TypeScript | Type-safe JavaScript |
| Build Tool | Vite 5 | Development and bundling |
| HTTP Client | Axios | API communication |
| Routing | React Router v6 | Client-side navigation |
| Styling | CSS3 | Custom styling with CSS variables |

---

## System Architecture

```
+-------------------+         HTTPS         +-------------------+
|                   |  <------------------> |                   |
|   React Frontend  |                       |    Go Backend     |
|   (Static Site)   |                       |   (Docker/API)    |
|                   |                       |                   |
+-------------------+                       +-------------------+
        |                                           |
        | Served via CDN                            | SQLite
        |                                           |
+-------------------+                       +-------------------+
|   User Browser    |                       |     Database      |
+-------------------+                       +-------------------+
```

### Request Flow

1. User interacts with React frontend in browser
2. Frontend makes authenticated API requests to backend
3. Backend validates JWT token and processes request
4. Backend queries/updates SQLite database via ent ORM
5. Response returned to frontend for rendering

---

## Database Design

### Entity-Relationship Diagram

```
+-------------+       +-------------+       +---------------+
|    User     |       |    Poll     |       |  PollOption   |
+-------------+       +-------------+       +---------------+
| id (PK)     |<----->| id (PK)     |<----->| id (PK)       |
| username    |   1:N | title       |   1:N | text          |
| email       |       | description |       | poll_id (FK)  |
| password    |       | creator_id  |       +---------------+
| created_at  |       | created_at  |               |
+-------------+       | updated_at  |               |
      |               +-------------+               |
      |                                             |
      |               +-------------+               |
      +-------------->|    Vote     |<--------------+
                  1:N +-------------+ N:1
                      | id (PK)     |
                      | user_id(FK) |
                      | option_id   |
                      | created_at  |
                      +-------------+
```

### Table Specifications

#### Users
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY, AUTO INCREMENT |
| username | VARCHAR(255) | UNIQUE, NOT NULL |
| email | VARCHAR(255) | UNIQUE, NOT NULL |
| password | VARCHAR(255) | NOT NULL |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

#### Polls
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY, AUTO INCREMENT |
| title | VARCHAR(255) | NOT NULL |
| description | TEXT | NULLABLE |
| creator_id | INTEGER | FOREIGN KEY (users.id) |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

#### PollOptions
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY, AUTO INCREMENT |
| text | VARCHAR(255) | NOT NULL |
| poll_id | INTEGER | FOREIGN KEY (polls.id), ON DELETE CASCADE |

#### Votes
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY, AUTO INCREMENT |
| user_id | INTEGER | FOREIGN KEY (users.id) |
| option_id | INTEGER | FOREIGN KEY (poll_options.id), ON DELETE CASCADE |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| | | UNIQUE(user_id, option_id) |

---

## API Documentation

### Base URL
- Development: `http://localhost:8080/api`
- Production: `https://poll-app-backend-lj26.onrender.com/api`

### Authentication Endpoints

#### Register User
```
POST /api/auth/signup
Content-Type: application/json

Request:
{
  "username": "string",
  "email": "string",
  "password": "string"
}

Response: 201 Created
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string"
  }
}
```

#### Login
```
POST /api/auth/login
Content-Type: application/json

Request:
{
  "email": "string",
  "password": "string"
}

Response: 200 OK
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string"
  }
}
```

#### Get Current User
```
GET /api/auth/me
Authorization: Bearer <token>

Response: 200 OK
{
  "id": 1,
  "username": "string",
  "email": "string"
}
```

### Poll Endpoints

#### List Polls
```
GET /api/polls
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "title": "string",
    "description": "string",
    "creator": { "id": 1, "username": "string" },
    "options": [
      { "id": 1, "text": "string", "vote_count": 0 }
    ],
    "user_voted_option_id": null,
    "created_at": "2026-01-12T00:00:00Z"
  }
]
```

#### Create Poll
```
POST /api/polls
Authorization: Bearer <token>
Content-Type: application/json

Request:
{
  "title": "string",
  "description": "string",
  "options": ["Option 1", "Option 2"]
}

Response: 201 Created
```

#### Get Poll
```
GET /api/polls/:id
Authorization: Bearer <token>

Response: 200 OK
```

#### Update Poll
```
PUT /api/polls/:id
Authorization: Bearer <token>
Content-Type: application/json

Request:
{
  "title": "string",
  "description": "string",
  "options": [
    { "id": 1, "text": "Updated Option" },
    { "text": "New Option" }
  ]
}

Response: 200 OK
```

#### Delete Poll
```
DELETE /api/polls/:id
Authorization: Bearer <token>

Response: 204 No Content
```

### Voting Endpoints

#### Submit Vote
```
POST /api/polls/:id/vote
Authorization: Bearer <token>
Content-Type: application/json

Request:
{
  "option_id": 1
}

Response: 200 OK
```

#### Get Voters for Option
```
GET /api/options/:id/voters
Authorization: Bearer <token>

Response: 200 OK
[
  { "id": 1, "username": "string", "email": "string" }
]
```

---

## Installation Guide

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm 9 or higher
- Git

### Backend Setup

```bash
# Navigate to backend directory
cd backend

# Download Go dependencies
go mod tidy

# Generate ent schema code
go generate ./ent

# Start the development server
go run main.go
```

The backend server will start at `http://localhost:8080`

### Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install npm dependencies
npm install

# Start the development server
npm run dev
```

The frontend application will be available at `http://localhost:3000`

### Environment Variables

#### Backend
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DATABASE_PATH | SQLite file path | poll_app.db |
| JWT_SECRET | JWT signing key | (required in production) |
| FRONTEND_URL | Allowed CORS origin | http://localhost:3000 |

#### Frontend
| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_URL | Backend API URL | http://localhost:8080 |

---

## Deployment

### Production URLs

- Frontend: https://poll-app-frontend-ylqk.onrender.com
- Backend: https://poll-app-backend-lj26.onrender.com

### Render Platform Deployment

The application is configured for deployment on Render using the included `render.yaml` blueprint file.

#### Automated Deployment (Recommended)

1. Push code to GitHub repository
2. Navigate to Render Dashboard
3. Select "New" followed by "Blueprint"
4. Connect the GitHub repository
5. Render will automatically detect `render.yaml` and configure both services

#### Manual Deployment

##### Backend Service
1. Create new Web Service on Render
2. Configure as Docker runtime
3. Set root directory to `backend`
4. Configure environment variables as specified above

##### Frontend Service
1. Create new Static Site on Render
2. Set root directory to `frontend`
3. Build command: `npm ci && npx tsc && npx vite build`
4. Publish directory: `dist`
5. Add rewrite rule: `/*` to `/index.html`

### Deployment Considerations

- Free tier services enter sleep mode after 15 minutes of inactivity
- Initial request after sleep may experience 30-60 second delay
- SQLite database persists within the container but resets on redeployment
- For production use, consider migrating to PostgreSQL

---

## Project Structure

```
poll_app/
├── README.md
├── render.yaml                 # Render deployment configuration
├── backend/
│   ├── Dockerfile              # Container configuration
│   ├── go.mod                  # Go module definition
│   ├── go.sum                  # Dependency checksums
│   ├── main.go                 # Application entry point
│   ├── handlers/
│   │   └── handlers.go         # HTTP request handlers
│   └── ent/
│       ├── schema/             # Database entity definitions
│       │   ├── user.go
│       │   ├── poll.go
│       │   ├── option.go
│       │   └── vote.go
│       └── [generated files]   # Auto-generated ORM code
└── frontend/
    ├── package.json            # Node.js dependencies
    ├── tsconfig.json           # TypeScript configuration
    ├── vite.config.ts          # Vite build configuration
    ├── index.html              # HTML entry point
    └── src/
        ├── main.tsx            # React entry point
        ├── App.tsx             # Root component
        ├── index.css           # Global styles
        ├── types/
        │   └── index.ts        # TypeScript type definitions
        ├── services/
        │   └── api.ts          # API client configuration
        ├── context/
        │   └── AuthContext.tsx # Authentication state management
        ├── components/
        │   └── Navbar.tsx      # Navigation component
        └── pages/
            ├── Login.tsx       # Login page
            ├── SignUp.tsx      # Registration page
            ├── Polls.tsx       # Poll listing page
            ├── CreatePoll.tsx  # Poll creation page
            ├── EditPoll.tsx    # Poll editing page
            └── PollDetail.tsx  # Poll detail and voting page
```

---

## License

This project was developed as part of a technical exercise for demonstrating full-stack development capabilities.

---

## Author

Pratyush Bindal
