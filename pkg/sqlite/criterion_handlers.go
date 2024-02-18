package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// shared criterion handlers go here

func orientationCriterionHandler(orientation *models.OrientationCriterionInput, heightColumn string, widthColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if orientation != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			var clauses []sqlClause

			for _, v := range orientation.Value {
				// width mod height
				mod := ""
				switch v {
				case models.OrientationPortrait:
					mod = "<"
				case models.OrientationLandscape:
					mod = ">"
				case models.OrientationSquare:
					mod = "="
				}

				if mod != "" {
					clauses = append(clauses, makeClause(fmt.Sprintf("%s %s %s", widthColumn, mod, heightColumn)))
				}
			}

			if len(clauses) > 0 {
				f.whereClauses = append(f.whereClauses, orClauses(clauses...))
			}
		}
	}
}
