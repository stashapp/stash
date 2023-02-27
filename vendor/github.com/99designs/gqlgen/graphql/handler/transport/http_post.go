package transport

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql"
)

// POST implements the POST side of the default HTTP transport
// defined in https://github.com/APIs-guru/graphql-over-http#post
type POST struct{}

var _ graphql.Transport = POST{}

func (h POST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}

	return r.Method == "POST" && mediaType == "application/json"
}

func getRequestBody(r *http.Request) (string, error) {
	if r == nil || r.Body == nil {
		return "", nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("unable to get Request Body %w", err)
	}
	return string(body), nil
}

func (h POST) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	params := &graphql.RawParams{}
	start := graphql.Now()
	params.Headers = r.Header
	params.ReadTime = graphql.TraceTiming{
		Start: start,
		End:   graphql.Now(),
	}

	bodyString, err := getRequestBody(r)
	if err != nil {
		gqlErr := gqlerror.Errorf("could not get json request body: %+v", err)
		resp := exec.DispatchError(ctx, gqlerror.List{gqlErr})
		log.Printf("could not get json request body: %+v", err.Error())
		writeJson(w, resp)
	}

	bodyReader := io.NopCloser(strings.NewReader(bodyString))
	if err = jsonDecode(bodyReader, &params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		gqlErr := gqlerror.Errorf(
			"json request body could not be decoded: %+v body:%s",
			err,
			bodyString,
		)
		resp := exec.DispatchError(ctx, gqlerror.List{gqlErr})
		log.Printf("decoding error: %+v body:%s", err.Error(), bodyString)
		writeJson(w, resp)
		return
	}

	rc, OpErr := exec.CreateOperationContext(ctx, params)
	if OpErr != nil {
		w.WriteHeader(statusFor(OpErr))
		resp := exec.DispatchError(graphql.WithOperationContext(ctx, rc), OpErr)
		writeJson(w, resp)
		return
	}

	var responses graphql.ResponseHandler
	responses, ctx = exec.DispatchOperation(ctx, rc)
	writeJson(w, responses(ctx))
}
