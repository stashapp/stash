//go:generate go run ../testdata/gqlgen.go

package integration

import (
	"context"
	"fmt"
	"time"

	models "github.com/99designs/gqlgen/integration/models-go"
	"github.com/99designs/gqlgen/integration/remote_api"
)

type CustomError struct {
	UserMessage   string
	InternalError string
}

func (e *CustomError) Error() string {
	return e.InternalError
}

type Resolver struct{}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (r *Resolver) Element() ElementResolver {
	return &elementResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type elementResolver struct{ *Resolver }

func (r *elementResolver) Error(ctx context.Context, obj *models.Element) (bool, error) {
	// A silly hack to make the result order stable
	time.Sleep(time.Duration(obj.ID) * 10 * time.Millisecond)

	return false, fmt.Errorf("boom")
}

func (r *elementResolver) Mismatched(ctx context.Context, obj *models.Element) ([]bool, error) {
	return []bool{true}, nil
}

func (r *elementResolver) Child(ctx context.Context, obj *models.Element) (*models.Element, error) {
	return &models.Element{ID: obj.ID * 10}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Error(ctx context.Context, typeArg *models.ErrorType) (bool, error) {
	if *typeArg == models.ErrorTypeCustom {
		return false, &CustomError{"User message", "Internal Message"}
	}

	return false, fmt.Errorf("normal error")
}

func (r *queryResolver) Path(ctx context.Context) ([]*models.Element, error) {
	return []*models.Element{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}, nil
}

func (r *queryResolver) Coercion(ctx context.Context, input []*models.ListCoercion) (bool, error) {
	return true, nil
}

func (r *queryResolver) Date(ctx context.Context, filter models.DateFilter) (bool, error) {
	if filter.Value != "asdf" {
		return false, fmt.Errorf("value must be asdf")
	}

	if *filter.Timezone != "UTC" {
		return false, fmt.Errorf("timezone must be utc")
	}

	if *filter.Op != models.DateFilterOpEq {
		return false, fmt.Errorf("unknown op %s", *filter.Op)
	}

	return true, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*models.Viewer, error) {
	return &models.Viewer{
		User: &remote_api.User{Name: "Bob"},
	}, nil
}

func (r *queryResolver) JSONEncoding(ctx context.Context) (string, error) {
	return "\U000fe4ed", nil
}

func (r *queryResolver) Complexity(ctx context.Context, value int) (bool, error) {
	return true, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Likes(ctx context.Context, obj *remote_api.User) ([]string, error) {
	return obj.Likes, nil
}
