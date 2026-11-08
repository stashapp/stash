package identify

import (
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

func setOptionalString(target *models.OptionalString, existing string, scraped *string, options *FieldOptions) {
	if scraped == nil || existing == *scraped {
		return
	}

	if !shouldSetSingleValueField(options, existing != "") {
		return
	}

	*target = models.NewOptionalString(*scraped)
}

// setPartialDate sets target from the scraped date value, if the field strategy
// allows it and the scraped value differs from the existing one. Parse failures
// are logged and target is left unset, rather than failing the whole identify
// operation over one bad field.
func setPartialDate(target *models.OptionalDate, scrapedDate *string, existing *models.Date, options *FieldOptions, fieldName string) {
	if scrapedDate == nil || (existing != nil && existing.String() == *scrapedDate) {
		return
	}

	if !shouldSetSingleValueField(options, existing != nil) {
		return
	}

	date, err := models.ParseDate(*scrapedDate)
	if err != nil {
		logger.Warnf("Ignoring scraped %s %q: %v", fieldName, *scrapedDate, err)
		return
	}

	*target = models.NewOptionalDate(date)
}

func setOptionalURLs(target **models.UpdateStrings, existing []string, scraped []string, options *FieldOptions) {
	if len(scraped) == 0 || !shouldSetSingleValueField(options, false) {
		return
	}

	switch getFieldStrategy(options) {
	case FieldStrategyOverwrite:
		if !sliceutil.SliceSame(scraped, existing) {
			*target = &models.UpdateStrings{
				Values: scraped,
				Mode:   models.RelationshipUpdateModeSet,
			}
		}
	case FieldStrategyMerge:
		urls := sliceutil.AppendUniques(existing, scraped)
		if len(urls) != len(existing) {
			*target = &models.UpdateStrings{
				Values: urls,
				Mode:   models.RelationshipUpdateModeSet,
			}
		}
	}
}
