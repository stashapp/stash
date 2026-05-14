# Architecture

This document provides an overview of the Stash codebase architecture for new contributors.

## Project Overview

Stash is a self-hosted web application written in Go that organizes and serves diverse media collections, catering to both SFW and NSFW needs. It gathers information about videos and images from the internet through extensible community-built plugins and scrapers, supports a wide variety of formats, enables tagging and filtering, and provides statistics about performers, tags, studios, and more.

**Core purpose**: Manage local media libraries with automatic metadata scraping, tagging, and organization.

**Key design philosophy**: 
- Backend: Go with GraphQL API and SQLite database
- Frontend: React/TypeScript with Apollo Client
- Extensibility: Plugin and scraper systems for community contributions
- Self-hosted: Single binary deployment with embedded frontend assets

## Repository Structure

```
stash/
├── cmd/              # Application entry points
│   ├── stash/        # Main application (cmd/stash/main.go)
│   └── phasher/      # Perceptual hash utility
├── pkg/              # Reusable Go packages
│   ├── file/         # File system operations and scanning
│   ├── ffmpeg/       # FFmpeg integration for media processing
│   ├── gallery/      # Gallery-specific business logic
│   ├── group/        # Group (movie) business logic
│   ├── hash/         # Hashing utilities (MD5, oshash, phash)
│   ├── image/        # Image-specific business logic
│   ├── job/          # Background job management
│   ├── logger/       # Logging utilities
│   ├── models/       # Interface definitions for data entities
│   ├── performer/    # Performer-specific business logic
│   ├── plugin/       # Plugin system
│   ├── scraper/      # Metadata scraping system
│   ├── scene/        # Scene-specific business logic
│   ├── sqlite/       # SQLite implementations of datalayer interfaces
│   ├── studio/       # Studio-specific business logic
│   ├── tag/          # Tag-specific business logic
│   └── ...           # Other utility packages
├── internal/         # Internal application code
│   ├── api/          # GraphQL API layer (resolvers, server)
│   ├── manager/      # Core application manager and services
│   ├── autotag/      # Auto-tagging functionality
│   ├── identify/     # Scene identification
│   ├── dlna/         # DLNA media server
│   ├── desktop/      # Desktop integration
│   ├── static/       # Static asset serving
│   └── log/          # Implementation of log system
├── ui/               # React/TypeScript frontend
│   ├── login/        # Login page
│   └── v2.5/         # Main frontend application
├── graphql/          # GraphQL schema definitions
│   ├── schema/       # Main schema files
│   └── stash-box/     # Stash-box integration schema
├── docker/           # Docker configuration
│   └── production/   # Production Docker setup
├── docs/             # Documentation
├── scripts/          # Utility scripts
├── go.mod            # Go module definition
├── go.sum            # Go dependency checksums
├── gqlgen.yml        # GraphQL code generation config
└── Makefile          # Build automation
```

## Backend Architecture

### Go Package Organization

The backend follows a layered architecture with clear separation of concerns:

**`pkg/models/` - Interface Layer**
- Defines interfaces for each entity (Scene, Image, Gallery, Performer, Studio, Tag, etc.)
- Each entity has `Reader`, `Writer`, and `ReaderWriter` interfaces
- Contains data model structs and query/filter types
- Example: `repository.go` defines the main `Repository` struct with all entity repositories
- Example: `repository_scene.go` defines `SceneReaderWriter` interface

**`pkg/sqlite/` - Implementation of datalayer interfaces**
- Implements the interfaces defined in `pkg/models/`
- Uses `goqu` query builder for SQL generation
- Contains all database access logic
- Example: `scene.go` implements `SceneStore` with CRUD operations
- Example: `scene_filter.go` implements filtering logic
- Handles transactions and connection pooling

**`internal/api/` - API Layer**
- GraphQL resolvers that implement the schema
- Each resolver method calls repository methods
- Handles authentication, authorization, and validation
- Example: `resolver_query_find_scene.go` implements scene query resolvers
- Example: `resolver_mutation_scene.go` implements scene mutation resolvers
- `server.go` sets up the HTTP server and GraphQL handler

### Layering Pattern

```
GraphQL Query/Mutation
    ↓
Resolver (internal/api/resolver_*.go)
    ↓
Repository Interface (pkg/models/repository_*.go)
    ↓
SQLite Implementation (pkg/sqlite/*.go)
    ↓
SQLite Database
```

Note: `gqlgen.yml` maps GraphQL types to Go structs and controls code generation. Update it when adding new types or fields.

### GraphQL Request Lifecycle

1. **Request**: Frontend sends GraphQL query to `/graphql` endpoint
2. **Routing**: `internal/api/server.go` routes to GraphQL handler (gqlgen)
3. **Parsing**: gqlgen parses the query and validates against schema
4. **Resolver Execution**: Appropriate resolver method in `internal/api/` is called
5. **Transaction**: Resolver wraps operation in read or write transaction via `withReadTxn()` or `withTxn()`
6. **Repository Call**: Resolver calls repository method (e.g., `r.repository.Scene.Find()`)
7. **SQL Execution**: SQLite implementation executes SQL query using goqu
8. **Response**: Data flows back through layers to frontend as JSON


### Plugin System

**Location**: `pkg/plugin/`

- JavaScript-based plugins that can extend Stash functionality
- Plugins are configured via YAML files in the plugins directory
- Support for hooks that trigger on events (e.g., `Scene.Create.Post`)
- Plugin cache in manager for performance
- RPC communication between Go and JavaScript plugins
- Example hooks: `Scene.Create.Post`, `Scene.Update.Post`, `Scan.Post`

Key files:
- `plugins.go` - Plugin loading and execution
- `hooks.go` - Hook system implementation
- `config.go` - Plugin configuration parsing

### Scraper System

**Location**: `pkg/scraper/`

- YAML-configured scrapers for fetching metadata from websites
- Supports multiple scraper types: XPath, JSON, GraphQL, script-based
- Scrapers can fetch performers, scenes, galleries, studios, tags
- Stash-box integration for crowd-sourced metadata
- Cache for scraper definitions
- Post-processing for transforming scraped data

Key files:
- `cache.go` - Scraper caching
- `definition.go` - Scraper configuration parsing
- `xpath.go` - XPath-based scraping
- `json.go` - JSON-based scraping
- `mapped.go` - Mapping scraped data to Stash models

### Task/Job System

**Location**: `pkg/job/`

- Background job management for long-running operations
- Progress reporting via GraphQL subscriptions
- Task queue with parallel execution
- Cancellation support
- Job types: Scan, Generate, Clean, Auto-tag, Identify, Export, Import

Key files:
- `manager.go` - Job manager implementation
- `job.go` - Job interface and progress tracking
- `subscribe.go` - Subscription support for job updates

Task implementations in `internal/manager/task/`:
- `task_scan.go` - File scanning
- `task_generate.go` - Thumbnail/sprite generation
- `task_clean.go` - Orphaned file cleanup
- `task_autotag.go` - Automatic tagging
- `task_identify.go` - Scene identification

## Frontend Architecture

### React/TypeScript Structure

**Location**: `ui/v2.5/`

```
ui/v2.5/
├── src/
│   ├── core/              # Core services and GraphQL client
│   │   ├── StashService.ts      # Main GraphQL client
│   │   ├── generated-graphql.ts # Auto-generated TypeScript types
│   │   ├── createClient.ts      # Apollo client setup
│   │   ├── config.ts            # Configuration
│   │   ├── scenes.ts            # Scene-specific queries
│   │   ├── performers.ts       # Performer-specific queries
│   │   └── ...
│   ├── components/       # React components
│   │   ├── Scenes/            # Scene-related components
│   │   ├── Performers/        # Performer-related components
│   │   ├── Galleries/         # Gallery-related components
│   │   ├── Images/            # Image-related components
│   │   ├── Studios/           # Studio-related components
│   │   ├── Tags/              # Tag-related components
│   │   ├── Settings/          # Settings components
│   │   ├── Shared/            # Shared/reusable components
│   │   └── ...
│   ├── hooks/            # Custom React hooks
│   │   ├── data.ts           # Data fetching hooks
│   │   ├── LocalForage.ts    # Local storage hooks
│   │   ├── Toast.tsx         # Toast notifications
│   │   └── ...
│   ├── models/           # Frontend data models
│   │   └── list-filter/      # Filter and list models
│   ├── locales/          # i18n translations
│   │   ├── en-GB.json        # English
│   │   ├── de-DE.json        # German
│   │   └── ...
│   ├── utils/            # Utility functions
│   ├── App.tsx           # Main application component
│   └── index.tsx         # Application entry point
├── graphql/              # GraphQL queries and fragments
├── public/               # Static assets
├── package.json          # Dependencies and scripts
├── codegen.ts            # GraphQL codegen configuration
└── vite.config.js        # Vite build configuration
```

### Communication with Backend

**GraphQL via Apollo Client**:
- Frontend uses Apollo Client (`@apollo/client`) for GraphQL communication
- GraphQL queries defined in `graphql/` directory
- Code generation via `@graphql-codegen/cli` generates TypeScript types
- Generated types in `src/core/generated-graphql.ts`
- WebSocket subscriptions for real-time updates (job progress, logging)


**Service Layer** (`src/core/`):
- `StashService.ts` - Main GraphQL client with typed queries/mutations
- Domain-specific files (scenes.ts, performers.ts, etc.) - Organized queries
- `createClient.ts` - Apollo client setup with authentication and uploads

## Database Layer

### SQLite Usage

**Database**: Single SQLite database file (default: `stash-go.sqlite`)
- WAL (Write-Ahead Logging) mode for concurrency
- Connection pooling: 1 write connection, 10 read connections
- 30-second idle connection timeout
- Configurable cache size via `STASH_SQLITE_CACHE_SIZE` environment variable

**Blob Storage**:
- Configurable storage for cover images and other binary data
- Options: Database (BLOB columns) or Filesystem (separate directory)
- Managed via `BlobStore` in `pkg/sqlite/blob.go`

### Migration System

**Location**: `pkg/sqlite/migrations/`

**Migration Files**:
- Numbered `.up.sql` files (e.g., `32_files.up.sql`)
- Current schema version is defined in `pkg/sqlite/database.go`
- Migrations embedded via `//go:embed migrations/*.sql`
- Uses `golang-migrate/migrate` library

**Custom Migrations**:
- Pre-migration Go files (e.g., `32_premigrate.go`) - Run before SQL
- Post-migration Go files (e.g., `32_postmigrate.go`) - Run after SQL
- Used for data transformations that SQL cannot handle

**Migration Process** (`pkg/sqlite/migrate.go`):
- The migrator runs pre-migration Go code, executes the SQL migration, then runs post-migration Go code for each version increment.

**Key Migrations**:
- `32_files.up.sql` - Introduced file/folder abstraction
- `45_blobs.up.sql` - Blob storage system
- `71_custom_fields.up.sql` - Custom fields support

### Query Patterns

**Repository Pattern**:
- All database access goes through repository interfaces
- SQLite implementations use goqu query builder
- Transactions managed via `txn.Manager`

**Example Query** (`pkg/sqlite/scene.go`):
```go
func (qb *SceneStore) Find(ctx context.Context, id int) (*models.Scene, error) {
    var scene models.Scene
    err := qb.repository.queryStruct(ctx, qb.sceneQuery(), []interface{}{id}, &scene)
    if err != nil {
        return nil, err
    }
    return &scene, nil
}
```

**Filtering**:
- Complex filtering via `criterion_handlers.go`
- Supports hierarchical filters (tags, studios)
- Joins handled via query builder

## Key Data Flows

### Example 1: GraphQL Query (findScene)

**Flow**:
1. Frontend sends GraphQL query requesting scene data by ID

2. Request hits `internal/api/server.go` at `/graphql` endpoint

3. gqlgen routes to the appropriate resolver in `resolver_query_find_scene.go`

4. Resolver wraps the operation in a read transaction using `withReadTxn()` to ensure consistent database access

5. Repository calls the SQLite implementation in `pkg/sqlite/scene.go` to execute the query

6. SQLite executes a parameterized SQL query via goqu to fetch the scene record

7. Scene object flows back through layers: SQLite → Repository → Resolver → GraphQL → Frontend

8. Frontend receives JSON response with the requested scene data

### Example 2: Scanning a File

**Flow**:
1. User triggers scan via UI (Settings → Metadata → Scan)

2. Frontend sends GraphQL mutation to start the scan job

3. Mutation resolver in `internal/api/resolver_mutation_metadata.go` creates a background job

4. Job manager queues `ScanJob` from `internal/manager/task_scan.go`

5. `ScanJob.Execute()` runs the scan operation with progress tracking

6. Filesystem walk traverses configured paths using `file.SymWalk`, queues files for processing, and filters based on modification time and .stashignore

7. File handlers process each file type: videos become Scenes, images become Images, zip files become Galleries, and folders get Folder records

8. For each video file, the system calculates checksums (MD5, oshash, phash), extracts metadata via FFmpeg, creates File and Scene records, and generates thumbnails, sprites, previews, and interactive heatmaps

9. Progress updates flow via GraphQL subscription with real-time updates on files processed

10. Scan completes with updated statistics, subscription notifies completion, and UI refreshes with new content

**Key Files**:
- `internal/manager/task_scan.go` - Main scan logic
- `pkg/file/` - File system operations
- `pkg/scene/scan.go` - Scene-specific scan logic
- `pkg/image/scan.go` - Image-specific scan logic
- `pkg/gallery/scan.go` - Gallery-specific scan logic

## Development Workflow

### Adding a New GraphQL Field

1. Define field in `graphql/schema/schema.graphql`
2. Run `make generate-backend` to regenerate types
3. Implement resolver in `internal/api/resolver_*.go`
4. If query requires new repository method:
   - Add interface to `pkg/models/repository_*.go`
   - Implement in `pkg/sqlite/*.go`
5. Add frontend query in `ui/v2.5/graphql/`
6. Run `make generate-ui` to regenerate frontend types
<<<<<<< HEAD
- Frontend type checking runs in CI — you do not need to run `tsc` locally.
=======

### Adding a Database Migration

1. Create new migration file: `pkg/sqlite/migrations/{version}_description.up.sql`
2. If needed, create `{version}_premigrate.go` for pre-migration logic
3. If needed, create `{version}_postmigrate.go` for post-migration logic
4. Update `appSchemaVersion` in `pkg/sqlite/database.go`
5. Test migration on development database

### Running Tests

```bash
# Backend test
make it
```

### Building

```bash
# Build backend
make build

# Build frontend
cd ui/v2.5
npm run build

# Full build
make build-ci
```

## Additional Resources

- **Development Guide**: See `docs/DEVELOPMENT.md`
- **Contributing**: See `docs/CONTRIBUTING.md`
- **GraphQL Schema**: `graphql/schema/schema.graphql`
- **API Documentation**: Available in-app via Shift+?
- **Community**: [Discord](https://discord.gg/2TsNFKt) and [Discourse](https://discourse.stashapp.cc)
