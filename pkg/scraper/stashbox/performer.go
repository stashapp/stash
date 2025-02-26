package stashbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// QueryStashBoxPerformer queries stash-box for performers using a query string.
func (c Client) QueryStashBoxPerformer(ctx context.Context, queryStr string) ([]*StashBoxPerformerQueryResult, error) {
	performers, err := c.queryStashBoxPerformer(ctx, queryStr)

	res := []*StashBoxPerformerQueryResult{
		{
			Query:   queryStr,
			Results: performers,
		},
	}

	// set the deprecated image field
	for _, p := range res[0].Results {
		if len(p.Images) > 0 {
			p.Image = &p.Images[0]
		}
	}

	return res, err
}

func (c Client) queryStashBoxPerformer(ctx context.Context, queryStr string) ([]*models.ScrapedPerformer, error) {
	performers, err := c.client.SearchPerformer(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	performerFragments := performers.SearchPerformer

	var ret []*models.ScrapedPerformer
	for _, fragment := range performerFragments {
		performer := performerFragmentToScrapedPerformer(*fragment)
		ret = append(ret, performer)
	}

	return ret, nil
}

// FindStashBoxPerformersByNames queries stash-box for performers by name
func (c Client) FindStashBoxPerformersByNames(ctx context.Context, performerIDs []string) ([]*StashBoxPerformerQueryResult, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name != "" {
				performers = append(performers, performer)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return c.findStashBoxPerformersByNames(ctx, performers)
}

func (c Client) FindStashBoxPerformersByPerformerNames(ctx context.Context, performerIDs []string) ([][]*models.ScrapedPerformer, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name != "" {
				performers = append(performers, performer)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	results, err := c.findStashBoxPerformersByNames(ctx, performers)
	if err != nil {
		return nil, err
	}

	var ret [][]*models.ScrapedPerformer
	for _, r := range results {
		ret = append(ret, r.Results)
	}

	return ret, nil
}

func (c Client) findStashBoxPerformersByNames(ctx context.Context, performers []*models.Performer) ([]*StashBoxPerformerQueryResult, error) {
	var ret []*StashBoxPerformerQueryResult
	for _, performer := range performers {
		if performer.Name != "" {
			performerResults, err := c.queryStashBoxPerformer(ctx, performer.Name)
			if err != nil {
				return nil, err
			}

			result := StashBoxPerformerQueryResult{
				Query:   strconv.Itoa(performer.ID),
				Results: performerResults,
			}

			ret = append(ret, &result)
		}
	}

	return ret, nil
}

func findURL(urls []*graphql.URLFragment, urlType string) *string {
	for _, u := range urls {
		if u.Type == urlType {
			ret := u.URL
			return &ret
		}
	}

	return nil
}

func enumToStringPtr(e fmt.Stringer, titleCase bool) *string {
	if e != nil {
		ret := strings.ReplaceAll(e.String(), "_", " ")
		if titleCase {
			c := cases.Title(language.Und)
			ret = c.String(strings.ToLower(ret))
		}
		return &ret
	}

	return nil
}

func translateGender(gender *graphql.GenderEnum) *string {
	var res models.GenderEnum
	switch *gender {
	case graphql.GenderEnumMale:
		res = models.GenderEnumMale
	case graphql.GenderEnumFemale:
		res = models.GenderEnumFemale
	case graphql.GenderEnumIntersex:
		res = models.GenderEnumIntersex
	case graphql.GenderEnumTransgenderFemale:
		res = models.GenderEnumTransgenderFemale
	case graphql.GenderEnumTransgenderMale:
		res = models.GenderEnumTransgenderMale
	case graphql.GenderEnumNonBinary:
		res = models.GenderEnumNonBinary
	}

	if res != "" {
		strVal := res.String()
		return &strVal
	}
	return nil
}

func formatMeasurements(m graphql.MeasurementsFragment) *string {
	if m.BandSize != nil && m.CupSize != nil && m.Hip != nil && m.Waist != nil {
		ret := fmt.Sprintf("%d%s-%d-%d", *m.BandSize, *m.CupSize, *m.Waist, *m.Hip)
		return &ret
	}

	return nil
}

func formatCareerLength(start, end *int) *string {
	if start == nil && end == nil {
		return nil
	}

	var ret string
	switch {
	case end == nil:
		ret = fmt.Sprintf("%d -", *start)
	case start == nil:
		ret = fmt.Sprintf("- %d", *end)
	default:
		ret = fmt.Sprintf("%d - %d", *start, *end)
	}

	return &ret
}

func formatBodyModifications(m []*graphql.BodyModificationFragment) *string {
	if len(m) == 0 {
		return nil
	}

	var retSlice []string
	for _, f := range m {
		if f.Description == nil {
			retSlice = append(retSlice, f.Location)
		} else {
			retSlice = append(retSlice, fmt.Sprintf("%s, %s", f.Location, *f.Description))
		}
	}

	ret := strings.Join(retSlice, "; ")
	return &ret
}

func fetchImage(ctx context.Context, client *http.Client, url string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// determine the image type and set the base64 type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	img := "data:" + contentType + ";base64," + utils.GetBase64StringFromData(body)
	return &img, nil
}

func performerFragmentToScrapedPerformer(p graphql.PerformerFragment) *models.ScrapedPerformer {
	images := []string{}
	for _, image := range p.Images {
		images = append(images, image.URL)
	}

	sp := &models.ScrapedPerformer{
		Name:           &p.Name,
		Disambiguation: p.Disambiguation,
		Country:        p.Country,
		Measurements:   formatMeasurements(*p.Measurements),
		CareerLength:   formatCareerLength(p.CareerStartYear, p.CareerEndYear),
		Tattoos:        formatBodyModifications(p.Tattoos),
		Piercings:      formatBodyModifications(p.Piercings),
		Twitter:        findURL(p.Urls, "TWITTER"),
		RemoteSiteID:   &p.ID,
		Images:         images,
		// TODO - tags not currently supported
		// graphql schema change to accommodate this. Leave off for now.
	}

	if len(sp.Images) > 0 {
		sp.Image = &sp.Images[0]
	}

	if p.Height != nil && *p.Height > 0 {
		hs := strconv.Itoa(*p.Height)
		sp.Height = &hs
	}

	if p.BirthDate != nil {
		sp.Birthdate = padFuzzyDate(p.BirthDate)
	}

	if p.DeathDate != nil {
		sp.DeathDate = padFuzzyDate(p.DeathDate)
	}

	if p.Gender != nil {
		sp.Gender = translateGender(p.Gender)
	}

	if p.Ethnicity != nil {
		sp.Ethnicity = enumToStringPtr(p.Ethnicity, true)
	}

	if p.EyeColor != nil {
		sp.EyeColor = enumToStringPtr(p.EyeColor, true)
	}

	if p.HairColor != nil {
		sp.HairColor = enumToStringPtr(p.HairColor, true)
	}

	if p.BreastType != nil {
		sp.FakeTits = enumToStringPtr(p.BreastType, true)
	}

	if len(p.Aliases) > 0 {
		// #4437 - stash-box may return aliases that are equal to the performer name
		// filter these out
		p.Aliases = sliceutil.Filter(p.Aliases, func(s string) bool {
			return !strings.EqualFold(s, p.Name)
		})

		// #4596 - stash-box may return duplicate aliases. Filter these out
		p.Aliases = stringslice.UniqueFold(p.Aliases)

		alias := strings.Join(p.Aliases, ", ")
		sp.Aliases = &alias
	}

	for _, u := range p.Urls {
		sp.URLs = append(sp.URLs, u.URL)
	}

	return sp
}

func padFuzzyDate(date *string) *string {
	if date == nil {
		return nil
	}

	var paddedDate string
	switch len(*date) {
	case 10:
		paddedDate = *date
	case 7:
		paddedDate = fmt.Sprintf("%s-01", *date)
	case 4:
		paddedDate = fmt.Sprintf("%s-01-01", *date)
	}
	return &paddedDate
}

func (c Client) FindStashBoxPerformerByID(ctx context.Context, id string) (*models.ScrapedPerformer, error) {
	performer, err := c.client.FindPerformerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if performer.FindPerformer == nil {
		return nil, nil
	}

	ret := performerFragmentToScrapedPerformer(*performer.FindPerformer)

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		err := match.ScrapedPerformer(ctx, r.Performer, ret, &c.box.Endpoint)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Client) FindStashBoxPerformerByName(ctx context.Context, name string) (*models.ScrapedPerformer, error) {
	performers, err := c.client.SearchPerformer(ctx, name)
	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedPerformer
	for _, performer := range performers.SearchPerformer {
		if strings.EqualFold(performer.Name, name) {
			ret = performerFragmentToScrapedPerformer(*performer)
		}
	}

	if ret == nil {
		return nil, nil
	}

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		err := match.ScrapedPerformer(ctx, r.Performer, ret, &c.box.Endpoint)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Client) SubmitPerformerDraft(ctx context.Context, performer *models.Performer) (*string, error) {
	draft := graphql.PerformerDraftInput{}
	var image io.Reader
	pqb := c.repository.Performer
	endpoint := c.box.Endpoint

	if err := performer.LoadAliases(ctx, pqb); err != nil {
		return nil, err
	}

	if err := performer.LoadURLs(ctx, pqb); err != nil {
		return nil, err
	}

	if err := performer.LoadStashIDs(ctx, pqb); err != nil {
		return nil, err
	}

	img, _ := pqb.GetImage(ctx, performer.ID)
	if img != nil {
		image = bytes.NewReader(img)
	}

	if performer.Name != "" {
		draft.Name = performer.Name
	}
	if performer.Disambiguation != "" {
		draft.Disambiguation = &performer.Disambiguation
	}
	if performer.Birthdate != nil {
		d := performer.Birthdate.String()
		draft.Birthdate = &d
	}
	if performer.Country != "" {
		draft.Country = &performer.Country
	}
	if performer.Ethnicity != "" {
		draft.Ethnicity = &performer.Ethnicity
	}
	if performer.EyeColor != "" {
		draft.EyeColor = &performer.EyeColor
	}
	if performer.FakeTits != "" {
		draft.BreastType = &performer.FakeTits
	}
	if performer.Gender != nil && performer.Gender.IsValid() {
		v := performer.Gender.String()
		draft.Gender = &v
	}
	if performer.HairColor != "" {
		draft.HairColor = &performer.HairColor
	}
	if performer.Height != nil {
		v := strconv.Itoa(*performer.Height)
		draft.Height = &v
	}
	if performer.Measurements != "" {
		draft.Measurements = &performer.Measurements
	}
	if performer.Piercings != "" {
		draft.Piercings = &performer.Piercings
	}
	if performer.Tattoos != "" {
		draft.Tattoos = &performer.Tattoos
	}
	if len(performer.Aliases.List()) > 0 {
		aliases := strings.Join(performer.Aliases.List(), ",")
		draft.Aliases = &aliases
	}
	if performer.CareerLength != "" {
		var career = strings.Split(performer.CareerLength, "-")
		if i, err := strconv.Atoi(strings.TrimSpace(career[0])); err == nil {
			draft.CareerStartYear = &i
		}
		if len(career) == 2 {
			if y, err := strconv.Atoi(strings.TrimSpace(career[1])); err == nil {
				draft.CareerEndYear = &y
			}
		}
	}

	if len(performer.URLs.List()) > 0 {
		draft.Urls = performer.URLs.List()
	}

	stashIDs, err := pqb.GetStashIDs(ctx, performer.ID)
	if err != nil {
		return nil, err
	}
	var stashID *string
	for _, v := range stashIDs {
		c := v
		if v.Endpoint == endpoint {
			stashID = &c.StashID
			break
		}
	}
	draft.ID = stashID

	var id *string
	var ret graphql.SubmitPerformerDraft
	err = c.submitDraft(ctx, graphql.SubmitPerformerDraftDocument, draft, image, &ret)
	id = ret.SubmitPerformerDraft.ID

	return id, err

	// ret, err := c.client.SubmitPerformerDraft(ctx, draft, uploadImage(image))
	// if err != nil {
	// 	return nil, err
	// }

	// id := ret.SubmitPerformerDraft.ID
	// return id, nil
}
