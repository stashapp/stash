// Worker for offloading stash media generation to a Windows + NVIDIA GPU box.
// Separate module from stashapp/stash so build deps don't pollute the main repo.
// See ../docs/llm/EXTERNAL-WORKERS.md for design.
module github.com/Ryokushen/stash/worker

go 1.25

require (
	github.com/corona10/goimagehash v1.1.0
	github.com/disintegration/imaging v1.6.2
	golang.org/x/image v0.18.0
)

require github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
