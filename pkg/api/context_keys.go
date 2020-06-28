package api

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint

type key int

const (
	galleryKey   key = 0
	performerKey key = 1
	sceneKey     key = 2
	studioKey    key = 3
	movieKey     key = 4
	ContextUser  key = 5
	tagKey       key = 6
)
