# Gallery People draft checkpoint

Updated: 27 October 2026, 12:45 Europe/Budapest.

This is an unfinished design and feasibility checkpoint. It is safe to review as source; it is not a working plugin release. Do not install the prototype into a real Stash library.

## User requirements captured

- Keep the existing Stash application and its normal capabilities.
- Windows is the target machine; the library may contain thousands of images.
- Keep images, recognition, and face data local to the Stash machine.
- Reuse galleries and image selection; analyze **only images the user selects**.
- Automatically produce hair-colour and visible skin-tone tags for that selected batch.
- Suggest groups of the same face, let the user name/correct groups, and support reviewed performer assignments and filtering.
- The final answer to the import-trigger question was **selected images only**. Do not add automatic scanning on import.
- Prepare beginner-friendly Windows setup; no GitHub/source-build knowledge should be required.
- Review upstream compatibility before completing implementation.
- Save draft work to GitHub as progress is made.

## Work saved

- [Compatibility and implementation plan](GALLERY_PEOPLE_PLAN.md): upstream findings, existing plugin comparison, integration points, product flow, Windows packaging, and validation requirements.
- `plugins/gallery-people/`: preliminary Python face-grouping/storage code, launcher, setup script, model manifest, dependencies, and plugin definition.
- The prototype predates the final requirements. It has a scan-all task, no frontend files, no hair/skin tagging, no completed Windows packaging, and no recognition-quality validation. The plan's preliminary-code review lists other corrections required.
- No tracked upstream application files, dependency declarations, or lockfiles have been changed.

## Repository state

- Work target: `ucbbrla9011-jpg/stash232323`, base commit `9e9ad972643320e737c711884f10e4285f0f3a88`. This matched upstream `stashapp/stash` develop HEAD at review time.
- Added repository: `ucbbrla9011-jpg/Stashnewprivate`, cloned at `/workspace/Stashnewprivate`, branch `6221`, commit `9866796a1495534a8037b2d6cc9c455f82032812`. Its checkout was clean and matched `origin/6221`.
- The user selected both repositories. Do not remove either or transfer changes between them without a reason established in the conversation.
- Existing checkout directories are the intended workspaces; do not create a worktree for ordinary cloud tasks.

## Cloud setup evidence at checkpoint time

- Go 1.25.14 downloaded from the official Go source with SHA-256 verification.
- Node 24.19.0, GCC, and FFmpeg 7.1.5 are available.
- UI dependencies installed using repository-pinned pnpm 10.33.0 and the frozen lockfile.
- Go module checksums verified. One artifact was fetched directly from its official Git repository after the Go proxy's storage destination was blocked.
- UI GraphQL generation, TypeScript checking, and the existing frontend production build passed. Upstream deprecation/bundle-size warnings remain.
- Backend generation was stopped at the user's requested handoff, with caches retained. A completed backend build, test suite, server startup, or functional gallery request has **not** yet been demonstrated.
- The face worker's Python dependencies installed in a private local environment. Model-weight download was blocked by a proxy 403. No face model has been loaded or validated.
- Required network additions are `storage.googleapis.com`, `media.githubusercontent.com`, and `api.github.com`; preserve the user's existing `docs.stashapp.cc` and `github.com` rules. Draft persistence and runtime application must be distinguished.
- GitHub branch `codex/gallery-people-draft-20261027` was successfully pushed and verified remotely. Automatic draft-PR creation failed with a proxy denial of `api.github.com`; no PR was created. The root README has manual next steps.
- The cloud environment draft now saves the dependency/frontend `install_script`, explicit work-in-progress resume instructions in `start_skill`, and the required network additions while preserving the existing domain rules. Saving does not execute these instructions or publish an environment. Backend startup remains unverified.

Local environment files are under `/workspace/stash-cloud/`: `env.sh`, tool/dependency caches, logs, and isolated development state. They are not committed to GitHub. Source `/workspace/stash-cloud/env.sh` before using the installed Go/pnpm tooling. Credentials and environment values must never be copied into a draft or recovery archive.

## Resume order

1. Inspect the current process/log state before restarting expensive generation or builds. Do not duplicate a running build.
2. Finish baseline build, UI checks, relevant existing tests, and a functional gallery/image API check in isolated test data.
3. Save tested install/start instructions in the cloud environment draft without changing the saved repository selection.
4. Implement the selected-image plugin seam from the compatibility plan; replace the prototype's scan-all behavior.
5. Evaluate local face and appearance models, redistribution licenses, and Windows runtime compatibility before claiming recognition/tagging readiness.
6. Complete review/correction, tag provenance, resumable jobs, tests, and Windows packaging. Keep this status document current and save meaningful draft checkpoints.

No release has been published, and no application readiness or model accuracy claim is implied by this checkpoint.
