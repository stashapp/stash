// +build ignore

package stashbox

import (
	"context"

	"github.com/shurcooL/graphql"
)

type findSceneByFingerprintScene struct {
}

type findSceneByFingerprintReturn []findSceneByFingerprintScene

type fingerprintQueryInput struct {
	Hash      graphql.String `graphql:"hash" json:"hash"`
	Algorithm graphql.String `graphql:"algorithm" json:"algorithm"`
}

// findSceneByFingerprint(fingerprint: FingerprintQueryInput!): [Scene!]!
func (c *Client) findSceneByFingerprint(input fingerprintQueryInput) (findSceneByFingerprintReturn, error) {
	var q struct {
		QueryReturn findSceneByFingerprintReturn `graphql:"findSceneByFingerprint(input: $i)"`
	}

	vars := map[string]interface{}{
		"i": input,
	}

	err := c.client.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, err
	}

	return q.QueryReturn, nil
}

// queryScenes(scene_filter: SceneFilterType, filter: QuerySpec): QueryScenesResultType!
