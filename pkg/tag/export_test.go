package tag

import (
	"errors"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	tagID      = 1
	noImageID  = 2
	errImageID = 3
)

const tagName = "testTag"

var createTime time.Time = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
var updateTime time.Time = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)

func createTag(id int) models.Tag {
	return models.Tag{
		ID:   id,
		Name: tagName,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createJSONTag(image string) *jsonschema.Tag {
	return &jsonschema.Tag{
		Name: tagName,
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
		Image: image,
	}
}

type testScenario struct {
	tag      models.Tag
	expected *jsonschema.Tag
	err      bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		testScenario{
			createTag(tagID),
			createJSONTag("PHN2ZwogICB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iCiAgIHhtbG5zOmNjPSJodHRwOi8vY3JlYXRpdmVjb21tb25zLm9yZy9ucyMiCiAgIHhtbG5zOnJkZj0iaHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyIKICAgeG1sbnM6c3ZnPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogICB4bWxuczpzb2RpcG9kaT0iaHR0cDovL3NvZGlwb2RpLnNvdXJjZWZvcmdlLm5ldC9EVEQvc29kaXBvZGktMC5kdGQiCiAgIHhtbG5zOmlua3NjYXBlPSJodHRwOi8vd3d3Lmlua3NjYXBlLm9yZy9uYW1lc3BhY2VzL2lua3NjYXBlIgogICB3aWR0aD0iMjAwIgogICBoZWlnaHQ9IjIwMCIKICAgaWQ9InN2ZzIiCiAgIHZlcnNpb249IjEuMSIKICAgaW5rc2NhcGU6dmVyc2lvbj0iMC40OC40IHI5OTM5IgogICBzb2RpcG9kaTpkb2NuYW1lPSJ0YWcuc3ZnIj4KICA8ZGVmcwogICAgIGlkPSJkZWZzNCIgLz4KICA8c29kaXBvZGk6bmFtZWR2aWV3CiAgICAgaWQ9ImJhc2UiCiAgICAgcGFnZWNvbG9yPSIjMDAwMDAwIgogICAgIGJvcmRlcmNvbG9yPSIjNjY2NjY2IgogICAgIGJvcmRlcm9wYWNpdHk9IjEuMCIKICAgICBpbmtzY2FwZTpwYWdlb3BhY2l0eT0iMSIKICAgICBpbmtzY2FwZTpwYWdlc2hhZG93PSIyIgogICAgIGlua3NjYXBlOnpvb209IjEiCiAgICAgaW5rc2NhcGU6Y3g9IjE4MS43Nzc3MSIKICAgICBpbmtzY2FwZTpjeT0iMjc5LjcyMzc2IgogICAgIGlua3NjYXBlOmRvY3VtZW50LXVuaXRzPSJweCIKICAgICBpbmtzY2FwZTpjdXJyZW50LWxheWVyPSJsYXllcjEiCiAgICAgc2hvd2dyaWQ9ImZhbHNlIgogICAgIGZpdC1tYXJnaW4tdG9wPSIwIgogICAgIGZpdC1tYXJnaW4tbGVmdD0iMCIKICAgICBmaXQtbWFyZ2luLXJpZ2h0PSIwIgogICAgIGZpdC1tYXJnaW4tYm90dG9tPSIwIgogICAgIGlua3NjYXBlOndpbmRvdy13aWR0aD0iMTkyMCIKICAgICBpbmtzY2FwZTp3aW5kb3ctaGVpZ2h0PSIxMDE3IgogICAgIGlua3NjYXBlOndpbmRvdy14PSItOCIKICAgICBpbmtzY2FwZTp3aW5kb3cteT0iLTgiCiAgICAgaW5rc2NhcGU6d2luZG93LW1heGltaXplZD0iMSIgLz4KICA8bWV0YWRhdGEKICAgICBpZD0ibWV0YWRhdGE3Ij4KICAgIDxyZGY6UkRGPgogICAgICA8Y2M6V29yawogICAgICAgICByZGY6YWJvdXQ9IiI+CiAgICAgICAgPGRjOmZvcm1hdD5pbWFnZS9zdmcreG1sPC9kYzpmb3JtYXQ+CiAgICAgICAgPGRjOnR5cGUKICAgICAgICAgICByZGY6cmVzb3VyY2U9Imh0dHA6Ly9wdXJsLm9yZy9kYy9kY21pdHlwZS9TdGlsbEltYWdlIiAvPgogICAgICAgIDxkYzp0aXRsZT48L2RjOnRpdGxlPgogICAgICA8L2NjOldvcms+CiAgICA8L3JkZjpSREY+CiAgPC9tZXRhZGF0YT4KICA8ZwogICAgIGlua3NjYXBlOmxhYmVsPSJMYXllciAxIgogICAgIGlua3NjYXBlOmdyb3VwbW9kZT0ibGF5ZXIiCiAgICAgaWQ9ImxheWVyMSIKICAgICB0cmFuc2Zvcm09InRyYW5zbGF0ZSgtMTU3Ljg0MzU4LC01MjQuNjk1MjIpIj4KICAgIDxwYXRoCiAgICAgICBpZD0icGF0aDI5ODciCiAgICAgICBkPSJtIDIyOS45NDMxNCw2NjkuMjY1NDkgLTM2LjA4NDY2LC0zNi4wODQ2NiBjIC00LjY4NjUzLC00LjY4NjUzIC00LjY4NjUzLC0xMi4yODQ2OCAwLC0xNi45NzEyMSBsIDM2LjA4NDY2LC0zNi4wODQ2NyBhIDEyLjAwMDQ1MywxMi4wMDA0NTMgMCAwIDEgOC40ODU2LC0zLjUxNDggbCA3NC45MTQ0MywwIGMgNi42Mjc2MSwwIDEyLjAwMDQxLDUuMzcyOCAxMi4wMDA0MSwxMi4wMDA0MSBsIDAsNzIuMTY5MzMgYyAwLDYuNjI3NjEgLTUuMzcyOCwxMi4wMDA0MSAtMTIuMDAwNDEsMTIuMDAwNDEgbCAtNzQuOTE0NDMsMCBhIDEyLjAwMDQ1MywxMi4wMDA0NTMgMCAwIDEgLTguNDg1NiwtMy41MTQ4MSB6IG0gLTEzLjQ1NjM5LC01My4wNTU4NyBjIC00LjY4NjUzLDQuNjg2NTMgLTQuNjg2NTMsMTIuMjg0NjggMCwxNi45NzEyMSA0LjY4NjUyLDQuNjg2NTIgMTIuMjg0NjcsNC42ODY1MiAxNi45NzEyLDAgNC42ODY1MywtNC42ODY1MyA0LjY4NjUzLC0xMi4yODQ2OCAwLC0xNi45NzEyMSAtNC42ODY1MywtNC42ODY1MiAtMTIuMjg0NjgsLTQuNjg2NTIgLTE2Ljk3MTIsMCB6IgogICAgICAgaW5rc2NhcGU6Y29ubmVjdG9yLWN1cnZhdHVyZT0iMCIKICAgICAgIHN0eWxlPSJmaWxsOiNmZmZmZmY7ZmlsbC1vcGFjaXR5OjEiIC8+CiAgPC9nPgo8L3N2Zz4="),
			false,
		},
		testScenario{
			createTag(noImageID),
			createJSONTag(""),
			false,
		},
		testScenario{
			createTag(errImageID),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockTagReader := &mocks.TagReaderWriter{}

	imageErr := errors.New("error getting image")

	mockTagReader.On("GetTagImage", tagID).Return(models.DefaultTagImage, nil).Once()
	mockTagReader.On("GetTagImage", noImageID).Return(nil, nil).Once()
	mockTagReader.On("GetTagImage", errImageID).Return(nil, imageErr).Once()

	for i, s := range scenarios {
		tag := s.tag
		json, err := ToJSON(mockTagReader, &tag)

		if !s.err && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if s.err && err == nil {
			t.Errorf("[%d] expected error not returned", i)
		} else {
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockTagReader.AssertExpectations(t)
}
