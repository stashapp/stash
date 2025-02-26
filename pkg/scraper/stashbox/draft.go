package stashbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/Yamashou/gqlgenc/graphqljson"
)

func (c *Client) submitDraft(ctx context.Context, query string, input interface{}, image io.Reader, ret interface{}) error {
	vars := map[string]interface{}{
		"input": input,
	}

	r := &clientv2.Request{
		Query:         query,
		Variables:     vars,
		OperationName: "",
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("operations", string(requestBody)); err != nil {
		return err
	}

	if image != nil {
		if err := writer.WriteField("map", "{ \"0\": [\"variables.input.image\"] }"); err != nil {
			return err
		}
		part, _ := writer.CreateFormFile("0", "draft")
		if _, err := io.Copy(part, image); err != nil {
			return err
		}
	} else if err := writer.WriteField("map", "{}"); err != nil {
		return err
	}

	writer.Close()

	req, _ := http.NewRequestWithContext(ctx, "POST", c.box.Endpoint, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("ApiKey", c.box.APIKey)

	httpClient := c.client.Client.Client
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type response struct {
		Data   json.RawMessage `json:"data"`
		Errors json.RawMessage `json:"errors"`
	}

	var respGQL response

	if err := json.Unmarshal(responseBytes, &respGQL); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(responseBytes), err)
	}

	if len(respGQL.Errors) > 0 {
		// try to parse standard graphql error
		errors := &clientv2.GqlErrorList{}
		if e := json.Unmarshal(responseBytes, errors); e != nil {
			return fmt.Errorf("failed to parse graphql errors. Response content %s - %w ", string(responseBytes), e)
		}

		return errors
	}

	if err := graphqljson.UnmarshalData(respGQL.Data, ret); err != nil {
		return err
	}

	return err
}

// we can't currently use this due to https://github.com/Yamashou/gqlgenc/issues/109
// func uploadImage(image io.Reader) client.HTTPRequestOption {
// 	return func(req *http.Request) {
// 		if image == nil {
// 			// return without changing anything
// 			return
// 		}

// 		// we can't handle errors in here, so if one happens, just return
// 		// without changing anything.

// 		// repackage the request to include the image
// 		bodyBytes, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			return
// 		}

// 		newBody := &bytes.Buffer{}
// 		writer := multipart.NewWriter(newBody)
// 		_ = writer.WriteField("operations", string(bodyBytes))

// 		if err := writer.WriteField("map", "{ \"0\": [\"variables.input.image\"] }"); err != nil {
// 			return
// 		}
// 		part, _ := writer.CreateFormFile("0", "draft")
// 		if _, err := io.Copy(part, image); err != nil {
// 			return
// 		}

// 		writer.Close()

// 		// now set the request body to this new body
// 		req.Body = io.NopCloser(newBody)
// 		req.ContentLength = int64(newBody.Len())
// 		req.Header.Set("Content-Type", writer.FormDataContentType())
// 	}
// }
