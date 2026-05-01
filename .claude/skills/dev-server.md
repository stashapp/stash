# Start Development Server

Set up and run the Stash development environment with hot-reloading.

## When to use
- Starting local development
- Testing backend/frontend changes interactively

## Steps

### First-time setup

1. Install UI dependencies:
   ```bash
   make pre-ui
   ```

2. Generate code:
   ```bash
   make generate
   ```

### Running

1. Start the backend server (terminal 1):
   ```bash
   make server-start
   ```
   This runs from `.local/` directory with `STASH_CONFIG_FILE=config.yml`

2. Start the frontend dev server (terminal 2):
   ```bash
   make ui-start
   ```
   UI runs on port 3000 by default, connects to backend at `http://localhost:9999`

3. Open browser: `http://localhost:3000/`

### First launch wizard

1. Choose a directory with media files
2. Accept default database/generated content locations
3. Navigate to "Tasks -> Library -> Scan" to scan your library

### Reset development environment

```bash
make server-clean   # Removes .local/ directory
make server-start   # Start fresh
```

## Tips
- Backend changes require server restart (Ctrl-C then `make server-start` again)
- Frontend changes hot-reload automatically
- To change the backend URL the UI connects to, set `VITE_APP_PLATFORM_URL` env var
- Auth cookies don't work cross-origin, so the UI dev server can't use auth
