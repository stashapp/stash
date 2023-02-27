package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql"
)

type SSE struct{}

var _ graphql.Transport = SSE{}

func (t SSE) Supports(r *http.Request) bool {
	if !strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}
	return r.Method == http.MethodPost && mediaType == "application/json"
}

func (t SSE) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	ctx := r.Context()
	flusher, ok := w.(http.Flusher)
	if !ok {
		SendErrorf(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}
	defer flusher.Flush()

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
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
		return
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

	ctx = graphql.WithOperationContext(ctx, rc)

	w.Header().Set("Content-Type", "text/event-stream")
	fmt.Fprint(w, ":\n\n")
	flusher.Flush()

	responses, ctx := exec.DispatchOperation(ctx, rc)

	for {
		response := responses(ctx)
		if response == nil {
			break
		}
		writeJsonWithSSE(w, response)
		flusher.Flush()
	}

	fmt.Fprint(w, "event: complete\n\n")
}

func writeJsonWithSSE(w io.Writer, response *graphql.Response) {
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "event: next\ndata: %s\n\n", b)
}
