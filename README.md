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

## Deploy to Render (Free)

### Option 1: One-Click Deploy with Blueprint

1. Push your code to GitHub:
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/YOUR_USERNAME/poll-app.git
   git push -u origin main
   ```

2. Go to [Render Dashboard](https://dashboard.render.com/)

3. Click **New** → **Blueprint**

4. Connect your GitHub repository

5. Render will detect `render.yaml` and deploy both services automatically

### Option 2: Manual Deploy

#### Deploy Backend:

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **New** → **Web Service**
3. Connect your GitHub repo
4. Configure:
   - **Name**: `poll-app-backend`
   - **Root Directory**: `backend`
   - **Runtime**: Docker
   - **Plan**: Free
5. Add Environment Variables:
   - `PORT`: `8080`
   - `JWT_SECRET`: (generate a random string)
   - `DATABASE_PATH`: `./poll_app.db`
   - `FRONTEND_URL`: (your frontend URL after deploying it)
6. Click **Create Web Service**

#### Deploy Frontend:

1. Click **New** → **Static Site**
2. Connect the same GitHub repo
3. Configure:
   - **Name**: `poll-app-frontend`
   - **Root Directory**: `frontend`
   - **Build Command**: `npm install && npm run build`
   - **Publish Directory**: `dist`
4. Add Environment Variable:
   - `VITE_API_URL`: `https://poll-app-backend.onrender.com/api`
5. Add Rewrite Rule:
   - **Source**: `/*`
   - **Destination**: `/index.html`
   - **Action**: Rewrite
6. Click **Create Static Site**

### After Deployment

1. Copy your backend URL (e.g., `https://poll-app-backend.onrender.com`)
2. Update the frontend's `VITE_API_URL` environment variable
3. Update the backend's `FRONTEND_URL` environment variable with frontend URL
4. Redeploy both services

### Important Notes

- **Free tier limitations**: Services spin down after 15 minutes of inactivity
- **Cold starts**: First request after idle may take 30-60 seconds
- **Database**: SQLite file is stored on the service, consider upgrading to PostgreSQL for production
