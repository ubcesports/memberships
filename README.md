# UBCEA Memberships Application

## Stack

### Frontend

- Next.js
- TypeScript
- Tailwind CSS

### Backend

- Golang
- PostgreSQL

## Deployment

Production deploys run automatically after CI succeeds on `main`, or manually via the **CD** workflow (`workflow_dispatch`).

### VPS setup

The VPS runs the app with Docker Compose. One-time server setup:

1. Install Docker and Docker Compose.
2. Create a deploy directory (this path becomes `VPS_DEPLOY_PATH`).
3. Add a `.env` file in that directory with the variables used by `deploy/compose.yaml`:
   - `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
   - `BACKEND_PORT` (default `8080`), `FRONTEND_PORT` (default `3000`)
   - Backend settings such as `DATABASE_URL`, `RESEND_API_KEY`, `SENDER_EMAIL`, OAuth credentials, etc. (see `backend/.env.example`)
4. Log in to GHCR on the VPS so it can pull images:
   ```bash
   echo <GITHUB_PAT> | docker login ghcr.io -u <GITHUB_USERNAME> --password-stdin
   ```
5. Add the deploy SSH public key to `~/.ssh/authorized_keys` for the deploy user.

### GitHub configuration

**Secrets** (Settings → Secrets and variables → Actions):

| Secret | Description |
| --- | --- |
| `VPS_SSH_KEY` | Private SSH key used by the CD workflow to connect to the VPS |
| `VPS_HOST` | VPS hostname or IP address |
| `VPS_USER` | SSH user on the VPS |
| `VPS_DEPLOY_PATH` | Absolute path to the deploy directory on the VPS |

**Variables**:

| Variable | Description |
| --- | --- |
| `NEXT_PUBLIC_API_URL` | Public API URL baked into the frontend image at build time |

### Deploy flow

On each deploy, the CD workflow:

1. Builds and pushes backend and frontend images to `ghcr.io/<owner>/<repo>/backend` and `ghcr.io/<owner>/<repo>/frontend`, tagged with the commit SHA and `latest`.
2. Copies `deploy/compose.yaml` to the VPS as a template.
3. Renders the template with the new image tags, pulls images while services are still running, stops services, then starts them. Compose runs migrations as a one-shot container before the backend starts.
