package followschema

// EmbeddedCase1 model
type EmbeddedCase1 struct {
	Empty
	*ExportedEmbeddedPointerAfterInterface
}

// Empty interface
type Empty interface{}

// ExportedEmbeddedPointerAfterInterface model
type ExportedEmbeddedPointerAfterInterface struct{}

// ExportedEmbeddedPointerExportedMethod method
func (*ExportedEmbeddedPointerAfterInterface) ExportedEmbeddedPointerExportedMethod() string {
	return "ExportedEmbeddedPointerExportedMethodResponse"
}

// EmbeddedCase2 model
type EmbeddedCase2 struct {
	*unexportedEmbeddedPointer
}

type unexportedEmbeddedPointer struct{}

// UnexportedEmbeddedPointerExportedMethod method
func (*unexportedEmbeddedPointer) UnexportedEmbeddedPointerExportedMethod() string {
	return "UnexportedEmbeddedPointerExportedMethodResponse"
}

// EmbeddedCase3 model
type EmbeddedCase3 struct {
	unexportedEmbeddedInterface
}

type unexportedEmbeddedInterface interface {
	nestedInterface
}

type nestedInterface interface {
	UnexportedEmbeddedInterfaceExportedMethod() string
}
