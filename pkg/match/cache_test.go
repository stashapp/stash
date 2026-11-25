package match

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

type fakeTagStore struct {
	tags    []*models.Tag
	aliases map[int][]string
}

func (f *fakeTagStore) Query(_ context.Context, filter *models.TagFilterType, _ *models.FindFilterType) ([]*models.Tag, int, error) {
	var nameRe, aliasRe *regexp.Regexp
	if filter != nil && filter.Name != nil && filter.Name.Modifier == models.CriterionModifierMatchesRegex {
		nameRe = regexp.MustCompile(filter.Name.Value)
	}
	if filter != nil && filter.OperatorFilter.Or != nil && filter.OperatorFilter.Or.Aliases != nil &&
		filter.OperatorFilter.Or.Aliases.Modifier == models.CriterionModifierMatchesRegex {
		aliasRe = regexp.MustCompile(filter.OperatorFilter.Or.Aliases.Value)
	}

	var matched []*models.Tag
	for _, t := range f.tags {
		if nameRe != nil && nameRe.MatchString(t.Name) {
			matched = append(matched, t)
			continue
		}
		if aliasRe != nil {
			for _, a := range f.aliases[t.ID] {
				if aliasRe.MatchString(a) {
					matched = append(matched, t)
					break
				}
			}
		}
	}
	return matched, len(matched), nil
}

func (f *fakeTagStore) GetAliases(_ context.Context, id int) ([]string, error) {
	return f.aliases[id], nil
}

func (f *fakeTagStore) QueryForAutoTag(_ context.Context, words []string) ([]*models.Tag, error) {
	var matched []*models.Tag
	for _, t := range f.tags {
		if t.IgnoreAutoTag {
			continue
		}
		if hasAnyPrefix(t.Name, words) || hasAnyAliasPrefix(f.aliases[t.ID], words) {
			matched = append(matched, t)
		}
	}
	return matched, nil
}

type fakeStudioStore struct {
	studios []*models.Studio
	aliases map[int][]string
}

func (f *fakeStudioStore) Query(_ context.Context, filter *models.StudioFilterType, _ *models.FindFilterType) ([]*models.Studio, int, error) {
	var nameRe *regexp.Regexp
	if filter != nil && filter.Name != nil && filter.Name.Modifier == models.CriterionModifierMatchesRegex {
		nameRe = regexp.MustCompile(filter.Name.Value)
	}
	var matched []*models.Studio
	for _, s := range f.studios {
		if nameRe != nil && nameRe.MatchString(s.Name) {
			matched = append(matched, s)
		}
	}
	return matched, len(matched), nil
}

func (f *fakeStudioStore) QueryCount(ctx context.Context, filter *models.StudioFilterType, findFilter *models.FindFilterType) (int, error) {
	res, count, err := f.Query(ctx, filter, findFilter)
	_ = res
	return count, err
}

func (f *fakeStudioStore) GetAliases(_ context.Context, id int) ([]string, error) {
	return f.aliases[id], nil
}

func (f *fakeStudioStore) QueryForAutoTag(_ context.Context, words []string) ([]*models.Studio, error) {
	var matched []*models.Studio
	for _, s := range f.studios {
		if s.IgnoreAutoTag {
			continue
		}
		if hasAnyPrefix(s.Name, words) || hasAnyAliasPrefix(f.aliases[s.ID], words) {
			matched = append(matched, s)
		}
	}
	return matched, nil
}

type fakePerformerStore struct {
	performers []*models.Performer
}

func (f *fakePerformerStore) Query(_ context.Context, filter *models.PerformerFilterType, _ *models.FindFilterType) ([]*models.Performer, int, error) {
	var nameRe *regexp.Regexp
	if filter != nil && filter.Name != nil && filter.Name.Modifier == models.CriterionModifierMatchesRegex {
		nameRe = regexp.MustCompile(filter.Name.Value)
	}
	var matched []*models.Performer
	for _, p := range f.performers {
		if nameRe != nil && nameRe.MatchString(p.Name) {
			matched = append(matched, p)
		}
	}
	return matched, len(matched), nil
}

func (f *fakePerformerStore) QueryCount(ctx context.Context, filter *models.PerformerFilterType, findFilter *models.FindFilterType) (int, error) {
	res, count, err := f.Query(ctx, filter, findFilter)
	_ = res
	return count, err
}

func (f *fakePerformerStore) GetAliases(_ context.Context, _ int) ([]string, error) {
	return nil, nil
}

func (f *fakePerformerStore) QueryForAutoTag(_ context.Context, words []string) ([]*models.Performer, error) {
	var matched []*models.Performer
	for _, p := range f.performers {
		if p.IgnoreAutoTag {
			continue
		}
		if hasAnyPrefix(p.Name, words) {
			matched = append(matched, p)
		}
	}
	return matched, nil
}

func hasAnyPrefix(name string, words []string) bool {
	lcName := strings.ToLower(name)
	for _, w := range words {
		if strings.HasPrefix(lcName, strings.ToLower(w)) {
			return true
		}
	}
	return false
}

func hasAnyAliasPrefix(aliases, words []string) bool {
	for _, a := range aliases {
		if hasAnyPrefix(a, words) {
			return true
		}
	}
	return false
}

func TestSingleFirstCharacterRegex(t *testing.T) {
	re, err := regexp.Compile(singleFirstCharacterRegex)
	if err != nil {
		t.Fatalf("failed to compile singleFirstCharacterRegex: %v", err)
	}

	tests := []struct {
		name string
		want bool
	}{
		{"1 - Tag [XYZ]", true},
		{"1-Tag", true},
		{"1.Tag", true},
		{"1_Tag", true},
		{"1: Tag", true},
		{"1", true},
		{"A - Tag", true},
		{"J.R.", true},
		{"a_b", true},
		{"X", true},
		{"１ - Tag", true},
		{"01 - Tag [XYZ]", false},
		{"10 - Tag", false},
		{"4K", false},
		{"ab", false},
		{"Tag", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := re.MatchString(tt.name); got != tt.want {
				t.Errorf("regexp(%q).MatchString(%q) = %v, want %v", singleFirstCharacterRegex, tt.name, got, tt.want)
			}
		})
	}
}

func TestPathToSingleDigitPrefixEntities(t *testing.T) {
	ctx := context.Background()

	tagStore := &fakeTagStore{
		tags: []*models.Tag{
			{ID: 1, Name: "1 - Tag [XYZ]"},
			{ID: 2, Name: "01 - Tag [XYZ]"},
			{ID: 3, Name: "Distinct Name"},
			{ID: 4, Name: "9"},
		},
		aliases: map[int][]string{
			3: {"2 - Alias Tag"},
		},
	}

	tagPaths := []struct {
		name    string
		path    string
		trimExt bool
		wantIDs []int
	}{
		{
			name:    "single digit prefix posix path",
			path:    "/library/1 - Tag [XYZ]/video.mp4",
			wantIDs: []int{1},
		},
		{
			name:    "single digit prefix windows path",
			path:    `C:\library\1 - Tag [XYZ]\video.mp4`,
			wantIDs: []int{1},
		},
		{
			name:    "zero-padded prefix control",
			path:    "/library/01 - Tag [XYZ]/video.mp4",
			wantIDs: []int{2},
		},
		{
			name:    "single digit prefix tag alias",
			path:    "/library/2 - Alias Tag/video.mp4",
			wantIDs: []int{3},
		},
		{
			name:    "single digit standalone name",
			path:    "9/a.mp4",
			trimExt: true,
			wantIDs: []int{4},
		},
	}

	for _, tt := range tagPaths {
		t.Run("tag/"+tt.name, func(t *testing.T) {
			cache := &Cache{}
			got, err := PathToTags(ctx, tt.path, tagStore, cache, tt.trimExt)
			if err != nil {
				t.Fatalf("PathToTags(%q) error: %v", tt.path, err)
			}
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("PathToTags(%q) returned %d tags (%+v), want %v", tt.path, len(got), got, tt.wantIDs)
			}
			for i, wantID := range tt.wantIDs {
				if got[i].ID != wantID {
					t.Errorf("PathToTags(%q)[%d].ID = %d, want %d", tt.path, i, got[i].ID, wantID)
				}
			}
		})
	}

	studioStore := &fakeStudioStore{
		studios: []*models.Studio{
			{ID: 10, Name: "1 - Studio"},
		},
	}
	gotStudio, err := PathToStudio(ctx, "/media/1 - Studio/clip.mp4", studioStore, &Cache{}, false)
	if err != nil {
		t.Fatalf("PathToStudio error: %v", err)
	}
	if gotStudio == nil || gotStudio.ID != 10 {
		t.Fatalf("PathToStudio returned %+v, want studio ID 10", gotStudio)
	}

	performerStore := &fakePerformerStore{
		performers: []*models.Performer{
			{ID: 20, Name: "1 - Performer"},
		},
	}
	gotPerformers, err := PathToPerformers(ctx, "/media/1 - Performer/clip.mp4", performerStore, &Cache{}, false)
	if err != nil {
		t.Fatalf("PathToPerformers error: %v", err)
	}
	if len(gotPerformers) != 1 || gotPerformers[0].ID != 20 {
		t.Fatalf("PathToPerformers returned %+v, want performer ID 20", gotPerformers)
	}
}
