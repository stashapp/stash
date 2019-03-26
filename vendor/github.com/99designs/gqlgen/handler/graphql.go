package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/complexity"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/golang-lru"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

type params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type Config struct {
	cacheSize                       int
	upgrader                        websocket.Upgrader
	recover                         graphql.RecoverFunc
	errorPresenter                  graphql.ErrorPresenterFunc
	resolverHook                    graphql.FieldMiddleware
	requestHook                     graphql.RequestMiddleware
	tracer                          graphql.Tracer
	complexityLimit                 int
	complexityLimitFunc             graphql.ComplexityLimitFunc
	disableIntrospection            bool
	connectionKeepAlivePingInterval time.Duration
}

func (c *Config) newRequestContext(es graphql.ExecutableSchema, doc *ast.QueryDocument, op *ast.OperationDefinition, query string, variables map[string]interface{}) *graphql.RequestContext {
	reqCtx := graphql.NewRequestContext(doc, query, variables)
	reqCtx.DisableIntrospection = c.disableIntrospection

	if hook := c.recover; hook != nil {
		reqCtx.Recover = hook
	}

	if hook := c.errorPresenter; hook != nil {
		reqCtx.ErrorPresenter = hook
	}

	if hook := c.resolverHook; hook != nil {
		reqCtx.ResolverMiddleware = hook
	}

	if hook := c.requestHook; hook != nil {
		reqCtx.RequestMiddleware = hook
	}

	if hook := c.tracer; hook != nil {
		reqCtx.Tracer = hook
	}

	if c.complexityLimit > 0 || c.complexityLimitFunc != nil {
		reqCtx.ComplexityLimit = c.complexityLimit
		operationComplexity := complexity.Calculate(es, op, variables)
		reqCtx.OperationComplexity = operationComplexity
	}

	return reqCtx
}

type Option func(cfg *Config)

func WebsocketUpgrader(upgrader websocket.Upgrader) Option {
	return func(cfg *Config) {
		cfg.upgrader = upgrader
	}
}

func RecoverFunc(recover graphql.RecoverFunc) Option {
	return func(cfg *Config) {
		cfg.recover = recover
	}
}

// ErrorPresenter transforms errors found while resolving into errors that will be returned to the user. It provides
// a good place to add any extra fields, like error.type, that might be desired by your frontend. Check the default
// implementation in graphql.DefaultErrorPresenter for an example.
func ErrorPresenter(f graphql.ErrorPresenterFunc) Option {
	return func(cfg *Config) {
		cfg.errorPresenter = f
	}
}

// IntrospectionEnabled = false will forbid clients from calling introspection endpoints. Can be useful in prod when you dont
// want clients introspecting the full schema.
func IntrospectionEnabled(enabled bool) Option {
	return func(cfg *Config) {
		cfg.disableIntrospection = !enabled
	}
}

// ComplexityLimit sets a maximum query complexity that is allowed to be executed.
// If a query is submitted that exceeds the limit, a 422 status code will be returned.
func ComplexityLimit(limit int) Option {
	return func(cfg *Config) {
		cfg.complexityLimit = limit
	}
}

// ComplexityLimitFunc allows you to define a function to dynamically set the maximum query complexity that is allowed
// to be executed.
// If a query is submitted that exceeds the limit, a 422 status code will be returned.
func ComplexityLimitFunc(complexityLimitFunc graphql.ComplexityLimitFunc) Option {
	return func(cfg *Config) {
		cfg.complexityLimitFunc = complexityLimitFunc
	}
}

// ResolverMiddleware allows you to define a function that will be called around every resolver,
// useful for logging.
func ResolverMiddleware(middleware graphql.FieldMiddleware) Option {
	return func(cfg *Config) {
		if cfg.resolverHook == nil {
			cfg.resolverHook = middleware
			return
		}

		lastResolve := cfg.resolverHook
		cfg.resolverHook = func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			return lastResolve(ctx, func(ctx context.Context) (res interface{}, err error) {
				return middleware(ctx, next)
			})
		}
	}
}

// RequestMiddleware allows you to define a function that will be called around the root request,
// after the query has been parsed. This is useful for logging
func RequestMiddleware(middleware graphql.RequestMiddleware) Option {
	return func(cfg *Config) {
		if cfg.requestHook == nil {
			cfg.requestHook = middleware
			return
		}

		lastResolve := cfg.requestHook
		cfg.requestHook = func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
			return lastResolve(ctx, func(ctx context.Context) []byte {
				return middleware(ctx, next)
			})
		}
	}
}

// Tracer allows you to add a request/resolver tracer that will be called around the root request,
// calling resolver. This is useful for tracing
func Tracer(tracer graphql.Tracer) Option {
	return func(cfg *Config) {
		if cfg.tracer == nil {
			cfg.tracer = tracer

		} else {
			lastResolve := cfg.tracer
			cfg.tracer = &tracerWrapper{
				tracer1: lastResolve,
				tracer2: tracer,
			}
		}

		opt := RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
			ctx = tracer.StartOperationExecution(ctx)
			resp := next(ctx)
			tracer.EndOperationExecution(ctx)

			return resp
		})
		opt(cfg)
	}
}

type tracerWrapper struct {
	tracer1 graphql.Tracer
	tracer2 graphql.Tracer
}

func (tw *tracerWrapper) StartOperationParsing(ctx context.Context) context.Context {
	ctx = tw.tracer1.StartOperationParsing(ctx)
	ctx = tw.tracer2.StartOperationParsing(ctx)
	return ctx
}

func (tw *tracerWrapper) EndOperationParsing(ctx context.Context) {
	tw.tracer2.EndOperationParsing(ctx)
	tw.tracer1.EndOperationParsing(ctx)
}

func (tw *tracerWrapper) StartOperationValidation(ctx context.Context) context.Context {
	ctx = tw.tracer1.StartOperationValidation(ctx)
	ctx = tw.tracer2.StartOperationValidation(ctx)
	return ctx
}

func (tw *tracerWrapper) EndOperationValidation(ctx context.Context) {
	tw.tracer2.EndOperationValidation(ctx)
	tw.tracer1.EndOperationValidation(ctx)
}

func (tw *tracerWrapper) StartOperationExecution(ctx context.Context) context.Context {
	ctx = tw.tracer1.StartOperationExecution(ctx)
	ctx = tw.tracer2.StartOperationExecution(ctx)
	return ctx
}

func (tw *tracerWrapper) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	ctx = tw.tracer1.StartFieldExecution(ctx, field)
	ctx = tw.tracer2.StartFieldExecution(ctx, field)
	return ctx
}

func (tw *tracerWrapper) StartFieldResolverExecution(ctx context.Context, rc *graphql.ResolverContext) context.Context {
	ctx = tw.tracer1.StartFieldResolverExecution(ctx, rc)
	ctx = tw.tracer2.StartFieldResolverExecution(ctx, rc)
	return ctx
}

func (tw *tracerWrapper) StartFieldChildExecution(ctx context.Context) context.Context {
	ctx = tw.tracer1.StartFieldChildExecution(ctx)
	ctx = tw.tracer2.StartFieldChildExecution(ctx)
	return ctx
}

func (tw *tracerWrapper) EndFieldExecution(ctx context.Context) {
	tw.tracer2.EndFieldExecution(ctx)
	tw.tracer1.EndFieldExecution(ctx)
}

func (tw *tracerWrapper) EndOperationExecution(ctx context.Context) {
	tw.tracer2.EndOperationExecution(ctx)
	tw.tracer1.EndOperationExecution(ctx)
}

// CacheSize sets the maximum size of the query cache.
// If size is less than or equal to 0, the cache is disabled.
func CacheSize(size int) Option {
	return func(cfg *Config) {
		cfg.cacheSize = size
	}
}

// WebsocketKeepAliveDuration allows you to reconfigure the keepalive behavior.
// By default, keepalive is enabled with a DefaultConnectionKeepAlivePingInterval
// duration. Set handler.connectionKeepAlivePingInterval = 0 to disable keepalive
// altogether.
func WebsocketKeepAliveDuration(duration time.Duration) Option {
	return func(cfg *Config) {
		cfg.connectionKeepAlivePingInterval = duration
	}
}

const DefaultCacheSize = 1000
const DefaultConnectionKeepAlivePingInterval = 25 * time.Second

func GraphQL(exec graphql.ExecutableSchema, options ...Option) http.HandlerFunc {
	cfg := &Config{
		cacheSize:                       DefaultCacheSize,
		connectionKeepAlivePingInterval: DefaultConnectionKeepAlivePingInterval,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	for _, option := range options {
		option(cfg)
	}

	var cache *lru.Cache
	if cfg.cacheSize > 0 {
		var err error
		cache, err = lru.New(cfg.cacheSize)
		if err != nil {
			// An error is only returned for non-positive cache size
			// and we already checked for that.
			panic("unexpected error creating cache: " + err.Error())
		}
	}
	if cfg.tracer == nil {
		cfg.tracer = &graphql.NopTracer{}
	}

	handler := &graphqlHandler{
		cfg:   cfg,
		cache: cache,
		exec:  exec,
	}

	return handler.ServeHTTP
}

var _ http.Handler = (*graphqlHandler)(nil)

type graphqlHandler struct {
	cfg   *Config
	cache *lru.Cache
	exec  graphql.ExecutableSchema
}

func (gh *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "OPTIONS, GET, POST")
		w.WriteHeader(http.StatusOK)
		return
	}

	if strings.Contains(r.Header.Get("Upgrade"), "websocket") {
		connectWs(gh.exec, w, r, gh.cfg, gh.cache)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var reqParams params
	switch r.Method {
	case http.MethodGet:
		reqParams.Query = r.URL.Query().Get("query")
		reqParams.OperationName = r.URL.Query().Get("operationName")

		if variables := r.URL.Query().Get("variables"); variables != "" {
			if err := jsonDecode(strings.NewReader(variables), &reqParams.Variables); err != nil {
				sendErrorf(w, http.StatusBadRequest, "variables could not be decoded")
				return
			}
		}
	case http.MethodPost:
		if err := jsonDecode(r.Body, &reqParams); err != nil {
			sendErrorf(w, http.StatusBadRequest, "json body could not be decoded: "+err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var doc *ast.QueryDocument
	var cacheHit bool
	if gh.cache != nil {
		val, ok := gh.cache.Get(reqParams.Query)
		if ok {
			doc = val.(*ast.QueryDocument)
			cacheHit = true
		}
	}

	ctx, doc, gqlErr := gh.parseOperation(ctx, &parseOperationArgs{
		Query:     reqParams.Query,
		CachedDoc: doc,
	})
	if gqlErr != nil {
		sendError(w, http.StatusUnprocessableEntity, gqlErr)
		return
	}

	ctx, op, vars, listErr := gh.validateOperation(ctx, &validateOperationArgs{
		Doc:           doc,
		OperationName: reqParams.OperationName,
		CacheHit:      cacheHit,
		R:             r,
		Variables:     reqParams.Variables,
	})
	if len(listErr) != 0 {
		sendError(w, http.StatusUnprocessableEntity, listErr...)
		return
	}

	if gh.cache != nil && !cacheHit {
		gh.cache.Add(reqParams.Query, doc)
	}

	reqCtx := gh.cfg.newRequestContext(gh.exec, doc, op, reqParams.Query, vars)
	ctx = graphql.WithRequestContext(ctx, reqCtx)

	defer func() {
		if err := recover(); err != nil {
			userErr := reqCtx.Recover(ctx, err)
			sendErrorf(w, http.StatusUnprocessableEntity, userErr.Error())
		}
	}()

	if gh.cfg.complexityLimitFunc != nil {
		reqCtx.ComplexityLimit = gh.cfg.complexityLimitFunc(ctx)
	}

	if reqCtx.ComplexityLimit > 0 && reqCtx.OperationComplexity > reqCtx.ComplexityLimit {
		sendErrorf(w, http.StatusUnprocessableEntity, "operation has complexity %d, which exceeds the limit of %d", reqCtx.OperationComplexity, reqCtx.ComplexityLimit)
		return
	}

	switch op.Operation {
	case ast.Query:
		b, err := json.Marshal(gh.exec.Query(ctx, op))
		if err != nil {
			panic(err)
		}
		w.Write(b)
	case ast.Mutation:
		b, err := json.Marshal(gh.exec.Mutation(ctx, op))
		if err != nil {
			panic(err)
		}
		w.Write(b)
	default:
		sendErrorf(w, http.StatusBadRequest, "unsupported operation type")
	}
}

type parseOperationArgs struct {
	Query     string
	CachedDoc *ast.QueryDocument
}

func (gh *graphqlHandler) parseOperation(ctx context.Context, args *parseOperationArgs) (context.Context, *ast.QueryDocument, *gqlerror.Error) {
	ctx = gh.cfg.tracer.StartOperationParsing(ctx)
	defer func() { gh.cfg.tracer.EndOperationParsing(ctx) }()

	if args.CachedDoc != nil {
		return ctx, args.CachedDoc, nil
	}

	doc, gqlErr := parser.ParseQuery(&ast.Source{Input: args.Query})
	if gqlErr != nil {
		return ctx, nil, gqlErr
	}

	return ctx, doc, nil
}

type validateOperationArgs struct {
	Doc           *ast.QueryDocument
	OperationName string
	CacheHit      bool
	R             *http.Request
	Variables     map[string]interface{}
}

func (gh *graphqlHandler) validateOperation(ctx context.Context, args *validateOperationArgs) (context.Context, *ast.OperationDefinition, map[string]interface{}, gqlerror.List) {
	ctx = gh.cfg.tracer.StartOperationValidation(ctx)
	defer func() { gh.cfg.tracer.EndOperationValidation(ctx) }()

	if !args.CacheHit {
		listErr := validator.Validate(gh.exec.Schema(), args.Doc)
		if len(listErr) != 0 {
			return ctx, nil, nil, listErr
		}
	}

	op := args.Doc.Operations.ForName(args.OperationName)
	if op == nil {
		return ctx, nil, nil, gqlerror.List{gqlerror.Errorf("operation %s not found", args.OperationName)}
	}

	if op.Operation != ast.Query && args.R.Method == http.MethodGet {
		return ctx, nil, nil, gqlerror.List{gqlerror.Errorf("GET requests only allow query operations")}
	}

	vars, err := validator.VariableValues(gh.exec.Schema(), op, args.Variables)
	if err != nil {
		return ctx, nil, nil, gqlerror.List{err}
	}

	return ctx, op, vars, nil
}

func jsonDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func sendError(w http.ResponseWriter, code int, errors ...*gqlerror.Error) {
	w.WriteHeader(code)
	b, err := json.Marshal(&graphql.Response{Errors: errors})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func sendErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	sendError(w, code, &gqlerror.Error{Message: fmt.Sprintf(format, args...)})
}
