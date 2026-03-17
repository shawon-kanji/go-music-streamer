# Project: Go Music Streamer (Spotify-like Backend)

## Goal
To learn and practice Golang by building a music streaming backend.

## Phase 1: Core Foundation & Basic Metadata
Focus: Setting up the server, database connection, and basic CRUD operations.

### Features
1.  **Project Setup**
    - [x] Initialize Go module
    - [x] Basic Gin server setup
    - [x] Health check endpoint (`GET /health`)

2.  **User Management (Basic)**
    - [ ] User registration (`POST /register`)
    - [ ] User login (basic auth or simple token) (`POST /login`)
    - [ ] Get user profile (`GET /me`)

3.  **Song Metadata**
    - [ ] CRUD for Songs (Title, Artist, Album, Duration)
    - [ ] List all songs (`GET /songs`)
    - [ ] Get song details (`GET /songs/:id`)
    - [ ] Upload song metadata (admin/mock) (`POST /songs`)

4.  **Playlists (Simple)**
    - [ ] Create playlist (`POST /playlists`)
    - [ ] Add song to playlist (`POST /playlists/:id/songs`)
    - [ ] Get user playlists (`GET /playlists`)

## Future Phases (Preview)
-   **Phase 2:** Audio Streaming (serving static files, range requests).
-   **Phase 3:** Search & Filtering (PostgreSQL search or simple string matching).
-   **Phase 4:** Real-time features (WebSockets for "now playing").

---

## Custom Agent: PR Description Generator

Use this agent whenever you want to generate a Pull Request description.

### Agent Name
`PR Description Agent`

### Agent Task
Create a clean PR description from the current branch changes using the template below.

### Rules
1. Keep language concise and developer-focused.
2. Use bullet points for lists.
3. Do not invent features not present in the diff.
4. If a section has no relevant content, write `N/A`.
5. Mention breaking changes clearly.

### Input
- Git diff / changed files
- Commit messages (optional)
- Related issue/ticket (optional)

### Output Template
```md
## Summary
-

## What Changed
-

## Why
-

## Testing
- [ ] Unit tests added/updated
- [ ] Manually tested
- [ ] Build passes

## Risks / Impact
-

## Breaking Changes
- None

## Related
- Closes #
```

### Usage Prompt
"Run `/agent` to generate the PR description using `agents/pr-description.md`."
