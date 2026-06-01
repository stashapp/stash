// Worker for offloading stash media generation to a Windows + NVIDIA GPU box.
// Separate module from stashapp/stash so build deps don't pollute the main repo.
// See ../docs/llm/EXTERNAL-WORKERS.md for design.
module github.com/Ryokushen/stash/worker

go 1.25
