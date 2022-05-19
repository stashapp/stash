package singlefile

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	introspection1 "github.com/99designs/gqlgen/codegen/testserver/singlefile/introspection"
	invalid_packagename "github.com/99designs/gqlgen/codegen/testserver/singlefile/invalid-packagename"
	"github.com/99designs/gqlgen/codegen/testserver/singlefile/otherpkg"
)

type Resolver struct{}

func (r *backedByInterfaceResolver) ID(ctx context.Context, obj BackedByInterface) (string, error) {
	panic("not implemented")
}

func (r *errorsResolver) A(ctx context.Context, obj *Errors) (*Error, error) {
	panic("not implemented")
}

func (r *errorsResolver) B(ctx context.Context, obj *Errors) (*Error, error) {
	panic("not implemented")
}

func (r *errorsResolver) C(ctx context.Context, obj *Errors) (*Error, error) {
	panic("not implemented")
}

func (r *errorsResolver) D(ctx context.Context, obj *Errors) (*Error, error) {
	panic("not implemented")
}

func (r *errorsResolver) E(ctx context.Context, obj *Errors) (*Error, error) {
	panic("not implemented")
}

func (r *forcedResolverResolver) Field(ctx context.Context, obj *ForcedResolver) (*Circle, error) {
	panic("not implemented")
}

func (r *modelMethodsResolver) ResolverField(ctx context.Context, obj *ModelMethods) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DefaultInput(ctx context.Context, input DefaultInput) (*DefaultParametersMirror, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateSomething(ctx context.Context, input SpecialInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePtrToPtr(ctx context.Context, input UpdatePtrToPtrOuter) (*PtrToPtrOuter, error) {
	panic("not implemented")
}

func (r *overlappingFieldsResolver) OldFoo(ctx context.Context, obj *OverlappingFields) (int, error) {
	panic("not implemented")
}

func (r *panicsResolver) FieldScalarMarshal(ctx context.Context, obj *Panics) ([]MarshalPanic, error) {
	panic("not implemented")
}

func (r *panicsResolver) ArgUnmarshal(ctx context.Context, obj *Panics, u []MarshalPanic) (bool, error) {
	panic("not implemented")
}

func (r *primitiveResolver) Value(ctx context.Context, obj *Primitive) (int, error) {
	panic("not implemented")
}

func (r *primitiveStringResolver) Value(ctx context.Context, obj *PrimitiveString) (string, error) {
	panic("not implemented")
}

func (r *primitiveStringResolver) Len(ctx context.Context, obj *PrimitiveString) (int, error) {
	panic("not implemented")
}

func (r *queryResolver) InvalidIdentifier(ctx context.Context) (*invalid_packagename.InvalidIdentifier, error) {
	panic("not implemented")
}

func (r *queryResolver) Collision(ctx context.Context) (*introspection1.It, error) {
	panic("not implemented")
}

func (r *queryResolver) MapInput(ctx context.Context, input map[string]interface{}) (*bool, error) {
	panic("not implemented")
}

func (r *queryResolver) Recursive(ctx context.Context, input *RecursiveInputSlice) (*bool, error) {
	panic("not implemented")
}

func (r *queryResolver) NestedInputs(ctx context.Context, input [][]*OuterInput) (*bool, error) {
	panic("not implemented")
}

func (r *queryResolver) NestedOutputs(ctx context.Context) ([][]*OuterObject, error) {
	panic("not implemented")
}

func (r *queryResolver) ModelMethods(ctx context.Context) (*ModelMethods, error) {
	panic("not implemented")
}

func (r *queryResolver) User(ctx context.Context, id int) (*User, error) {
	panic("not implemented")
}

func (r *queryResolver) NullableArg(ctx context.Context, arg *int) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) InputSlice(ctx context.Context, arg []string) (bool, error) {
	panic("not implemented")
}

func (r *queryResolver) InputNullableSlice(ctx context.Context, arg []string) (bool, error) {
	panic("not implemented")
}

func (r *queryResolver) ShapeUnion(ctx context.Context) (ShapeUnion, error) {
	panic("not implemented")
}

func (r *queryResolver) Autobind(ctx context.Context) (*Autobind, error) {
	panic("not implemented")
}

func (r *queryResolver) DeprecatedField(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) Overlapping(ctx context.Context) (*OverlappingFields, error) {
	panic("not implemented")
}

func (r *queryResolver) DefaultParameters(ctx context.Context, falsyBoolean *bool, truthyBoolean *bool) (*DefaultParametersMirror, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveArg(ctx context.Context, arg string) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveNullableArg(ctx context.Context, arg *int, arg2 *int, arg3 *string) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveInputNullable(ctx context.Context, arg *InputDirectives) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveInput(ctx context.Context, arg InputDirectives) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveInputType(ctx context.Context, arg InnerInput) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveObject(ctx context.Context) (*ObjectDirectives, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveObjectWithCustomGoModel(ctx context.Context) (*ObjectDirectivesWithCustomGoModel, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveFieldDef(ctx context.Context, ret string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveField(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveDouble(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) DirectiveUnimplemented(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) EmbeddedCase1(ctx context.Context) (*EmbeddedCase1, error) {
	panic("not implemented")
}

func (r *queryResolver) EmbeddedCase2(ctx context.Context) (*EmbeddedCase2, error) {
	panic("not implemented")
}

func (r *queryResolver) EmbeddedCase3(ctx context.Context) (*EmbeddedCase3, error) {
	panic("not implemented")
}

func (r *queryResolver) EnumInInput(ctx context.Context, input *InputWithEnumValue) (EnumTest, error) {
	panic("not implemented")
}

func (r *queryResolver) Shapes(ctx context.Context) ([]Shape, error) {
	panic("not implemented")
}

func (r *queryResolver) NoShape(ctx context.Context) (Shape, error) {
	panic("not implemented")
}

func (r *queryResolver) Node(ctx context.Context) (Node, error) {
	panic("not implemented")
}

func (r *queryResolver) NoShapeTypedNil(ctx context.Context) (Shape, error) {
	panic("not implemented")
}

func (r *queryResolver) Animal(ctx context.Context) (Animal, error) {
	panic("not implemented")
}

func (r *queryResolver) NotAnInterface(ctx context.Context) (BackedByInterface, error) {
	panic("not implemented")
}

func (r *queryResolver) Issue896a(ctx context.Context) ([]*CheckIssue896, error) {
	panic("not implemented")
}

func (r *queryResolver) MapStringInterface(ctx context.Context, in map[string]interface{}) (map[string]interface{}, error) {
	panic("not implemented")
}

func (r *queryResolver) MapNestedStringInterface(ctx context.Context, in *NestedMapInput) (map[string]interface{}, error) {
	panic("not implemented")
}

func (r *queryResolver) ErrorBubble(ctx context.Context) (*Error, error) {
	panic("not implemented")
}

func (r *queryResolver) ErrorBubbleList(ctx context.Context) ([]*Error, error) {
	panic("not implemented")
}

func (r *queryResolver) ErrorList(ctx context.Context) ([]*Error, error) {
	panic("not implemented")
}

func (r *queryResolver) Errors(ctx context.Context) (*Errors, error) {
	panic("not implemented")
}

func (r *queryResolver) Valid(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) Panics(ctx context.Context) (*Panics, error) {
	panic("not implemented")
}

func (r *queryResolver) PrimitiveObject(ctx context.Context) ([]Primitive, error) {
	panic("not implemented")
}

func (r *queryResolver) PrimitiveStringObject(ctx context.Context) ([]PrimitiveString, error) {
	panic("not implemented")
}

func (r *queryResolver) PtrToSliceContainer(ctx context.Context) (*PtrToSliceContainer, error) {
	panic("not implemented")
}

func (r *queryResolver) Infinity(ctx context.Context) (float64, error) {
	panic("not implemented")
}

func (r *queryResolver) StringFromContextInterface(ctx context.Context) (*StringFromContextInterface, error) {
	panic("not implemented")
}

func (r *queryResolver) StringFromContextFunction(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) DefaultScalar(ctx context.Context, arg string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) Slices(ctx context.Context) (*Slices, error) {
	panic("not implemented")
}

func (r *queryResolver) ScalarSlice(ctx context.Context) ([]byte, error) {
	panic("not implemented")
}

func (r *queryResolver) Fallback(ctx context.Context, arg FallbackToStringEncoding) (FallbackToStringEncoding, error) {
	panic("not implemented")
}

func (r *queryResolver) OptionalUnion(ctx context.Context) (TestUnion, error) {
	panic("not implemented")
}

func (r *queryResolver) VOkCaseValue(ctx context.Context) (*VOkCaseValue, error) {
	panic("not implemented")
}

func (r *queryResolver) VOkCaseNil(ctx context.Context) (*VOkCaseNil, error) {
	panic("not implemented")
}

func (r *queryResolver) ValidType(ctx context.Context) (*ValidType, error) {
	panic("not implemented")
}

func (r *queryResolver) WrappedStruct(ctx context.Context) (*WrappedStruct, error) {
	panic("not implemented")
}

func (r *queryResolver) WrappedScalar(ctx context.Context) (otherpkg.Scalar, error) {
	panic("not implemented")
}

func (r *queryResolver) WrappedMap(ctx context.Context) (WrappedMap, error) {
	panic("not implemented")
}

func (r *queryResolver) WrappedSlice(ctx context.Context) (WrappedSlice, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) Updated(ctx context.Context) (<-chan string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) InitPayload(ctx context.Context) (<-chan string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) DirectiveArg(ctx context.Context, arg string) (<-chan *string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) DirectiveNullableArg(ctx context.Context, arg *int, arg2 *int, arg3 *string) (<-chan *string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) DirectiveDouble(ctx context.Context) (<-chan *string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) DirectiveUnimplemented(ctx context.Context) (<-chan *string, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) Issue896b(ctx context.Context) (<-chan []*CheckIssue896, error) {
	panic("not implemented")
}

func (r *userResolver) Friends(ctx context.Context, obj *User) ([]*User, error) {
	panic("not implemented")
}

func (r *wrappedMapResolver) Get(ctx context.Context, obj WrappedMap, key string) (string, error) {
	panic("not implemented")
}

func (r *wrappedSliceResolver) Get(ctx context.Context, obj WrappedSlice, idx int) (string, error) {
	panic("not implemented")
}

// BackedByInterface returns BackedByInterfaceResolver implementation.
func (r *Resolver) BackedByInterface() BackedByInterfaceResolver {
	return &backedByInterfaceResolver{r}
}

// Errors returns ErrorsResolver implementation.
func (r *Resolver) Errors() ErrorsResolver { return &errorsResolver{r} }

// ForcedResolver returns ForcedResolverResolver implementation.
func (r *Resolver) ForcedResolver() ForcedResolverResolver { return &forcedResolverResolver{r} }

// ModelMethods returns ModelMethodsResolver implementation.
func (r *Resolver) ModelMethods() ModelMethodsResolver { return &modelMethodsResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// OverlappingFields returns OverlappingFieldsResolver implementation.
func (r *Resolver) OverlappingFields() OverlappingFieldsResolver {
	return &overlappingFieldsResolver{r}
}

// Panics returns PanicsResolver implementation.
func (r *Resolver) Panics() PanicsResolver { return &panicsResolver{r} }

// Primitive returns PrimitiveResolver implementation.
func (r *Resolver) Primitive() PrimitiveResolver { return &primitiveResolver{r} }

// PrimitiveString returns PrimitiveStringResolver implementation.
func (r *Resolver) PrimitiveString() PrimitiveStringResolver { return &primitiveStringResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

// WrappedMap returns WrappedMapResolver implementation.
func (r *Resolver) WrappedMap() WrappedMapResolver { return &wrappedMapResolver{r} }

// WrappedSlice returns WrappedSliceResolver implementation.
func (r *Resolver) WrappedSlice() WrappedSliceResolver { return &wrappedSliceResolver{r} }

type backedByInterfaceResolver struct{ *Resolver }
type errorsResolver struct{ *Resolver }
type forcedResolverResolver struct{ *Resolver }
type modelMethodsResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type overlappingFieldsResolver struct{ *Resolver }
type panicsResolver struct{ *Resolver }
type primitiveResolver struct{ *Resolver }
type primitiveStringResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type wrappedMapResolver struct{ *Resolver }
type wrappedSliceResolver struct{ *Resolver }
