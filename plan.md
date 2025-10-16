## Simple Tasks (doable anywhere)
> Administrative, research, and planning tasks — great for downtime or while traveling.

### Main

- [ ] Brainstorm project name, logo, and branding
- [ ] Research competitors (Docker Hub, Replit, Railway, Podman)
- [ ] Write user flow draft (how users create, run, and share setups)
- [x] Decide on tech stack (backend, frontend, orchestration) 
> Using Go backend and React frontend

### Extra
- [ ] Register hosting provider (Vercel / Fly.io / Hetzner / DigitalOcean)
- [ ] Register domain + set up email 

## Complex Tasks (require focus)
> Architecture, backend, frontend, and testing — schedule these for deep work sessions.

### Architecture

#### Main
- [ ] Define high-level system diagram (frontend ↔ API ↔ runner engine)
- [ ] Research sandboxing methods (Docker)
- [ ] Define API structure and endpoints
- [ ] Plan basic authentication/authorization system

#### Extra
- [ ] Plan extended authentication/authorization system (JWT / OAuth2)

### Backend (API + Engine)
- [x] Setup project structure and basic server
> Basic vite and go server setup
- [x] Implement basic user authentication
> Username and a hashed salted password
- [ ] CRUD endpoints for setups (create, read, update, delete)
- [ ] Implement setup runner (container execution)
- [ ] Add logs and job status tracking
- [ ] Connect to database (PostgreSQL / MongoDB)
- [ ] Add rate limiting and input validation

### Frontend
- [x] Initialize React project (Vite)
> Vite project created, tailwind and react added
- [ ] Implement login/register UI
> Baseic login implemented
- [ ] Create dashboard to list user setups
- [ ] Build setup editor (YAML/JSON configuration)
- [ ] Display logs and container run status
- [ ] Allow cloning/copying other users’ public setups
- [ ] Integrate API calls with backend

### Testing
- [ ] Unit tests for API endpoints
- [ ] Integration tests 
- [ ] Basic end-to-end test 
- [ ] Security sanity checks (input sanitization, token validation)

---

## Necessary but Boring
> Not exciting, but critical for stability, deployment, and user trust.

- [ ] Write documentation (API reference + setup guide)
- [ ] Add Terms of Service + Privacy Policy
- [ ] Optimize Dockerfile (size, cache, multi-stage)
- [ ] Setup staging + production environments
- [ ] Implement CORS, HTTPS, and security headers
- [ ] Add onboarding/help page for new users

---

## Timeline Plan

| Period | Focus | Key Goals |
|---------|--------|-----------|
| **Oct 2025** | Research + Setup | Architecture, repo, and stack finalized |
| **Nov 2025** | Backend MVP | Working API (auth + CRUD + runner) |
| **Dec 2025** | Frontend + Integration | UI connected to backend, basic usability |
| **Jan 2026** | Testing + Polish | CI/CD, docs, and 90% MVP completion |

## Stretch Goals
- [ ] Implement billing/subscriptions
- [ ] Public profiles & community templates
- [ ] Template marketplace
- [ ] Live metrics (CPU, RAM, uptime)
- [ ] CLI tool for setup deployment
- [ ] WebSocket live logs
