package apollotracing

import (
	"context"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	Tracer struct{}

	TracingExtension struct {
		mu         sync.Mutex
		Version    int           `json:"version"`
		StartTime  time.Time     `json:"startTime"`
		EndTime    time.Time     `json:"endTime"`
		Duration   time.Duration `json:"duration"`
		Parsing    Span          `json:"parsing"`
		Validation Span          `json:"validation"`
		Execution  struct {
			Resolvers []*ResolverExecution `json:"resolvers"`
		} `json:"execution"`
	}

	Span struct {
		StartOffset time.Duration `json:"startOffset"`
		Duration    time.Duration `json:"duration"`
	}

	ResolverExecution struct {
		Path        ast.Path      `json:"path"`
		ParentType  string        `json:"parentType"`
		FieldName   string        `json:"fieldName"`
		ReturnType  string        `json:"returnType"`
		StartOffset time.Duration `json:"startOffset"`
		Duration    time.Duration `json:"duration"`
	}
)

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Tracer{}

func (Tracer) ExtensionName() string {
	return "ApolloTracing"
}

func (Tracer) Validate(graphql.ExecutableSchema) error {
	return nil
}

func (Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	td, ok := graphql.GetExtension(ctx, "tracing").(*TracingExtension)
	if !ok {
		return next(ctx)
	}

	start := graphql.Now()

	defer func() {
		end := graphql.Now()

		rc := graphql.GetOperationContext(ctx)
		fc := graphql.GetFieldContext(ctx)
		resolver := &ResolverExecution{
			Path:        fc.Path(),
			ParentType:  fc.Object,
			FieldName:   fc.Field.Name,
			ReturnType:  fc.Field.Definition.Type.String(),
			StartOffset: start.Sub(rc.Stats.OperationStart),
			Duration:    end.Sub(start),
		}

		td.mu.Lock()
		td.Execution.Resolvers = append(td.Execution.Resolvers, resolver)
		td.mu.Unlock()
	}()

	return next(ctx)
}

func (Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	rc := graphql.GetOperationContext(ctx)

	start := rc.Stats.OperationStart

	td := &TracingExtension{
		Version:   1,
		StartTime: start,
		Parsing: Span{
			StartOffset: rc.Stats.Parsing.Start.Sub(start),
			Duration:    rc.Stats.Parsing.End.Sub(rc.Stats.Parsing.Start),
		},

		Validation: Span{
			StartOffset: rc.Stats.Validation.Start.Sub(start),
			Duration:    rc.Stats.Validation.End.Sub(rc.Stats.Validation.Start),
		},
	}

	graphql.RegisterExtension(ctx, "tracing", td)
	resp := next(ctx)

	end := graphql.Now()
	td.EndTime = end
	td.Duration = end.Sub(start)

	return resp
}
