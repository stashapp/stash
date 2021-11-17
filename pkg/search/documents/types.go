package documents

import "fmt"

type DocType string

const (
	TypeStudio    DocType = "studio"
	TypePerformer DocType = "performer"
	TypeScene     DocType = "scene"
	TypeTag       DocType = "tag"
)

func NewDocType(in string) DocType {
	switch in {
	case "studio":
		return TypeStudio
	case "performer":
		return TypePerformer
	case "scene":
		return TypeScene
	case "tag":
		return TypeTag
	}

	panic(fmt.Sprintf("unhandled case in NewDocType: %v", in))
}
