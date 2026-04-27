package urlbuilders

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type AudioURLBuilder struct {
	BaseURL   string
	AudioID   string
	UpdatedAt string
}

func NewAudioURLBuilder(baseURL string, audio *models.Audio) AudioURLBuilder {
	return AudioURLBuilder{
		BaseURL:   baseURL,
		AudioID:   strconv.Itoa(audio.ID),
		UpdatedAt: strconv.FormatInt(audio.UpdatedAt.Unix(), 10),
	}
}

func (b AudioURLBuilder) GetStreamURL(apiKey string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s/audio/%s/stream", b.BaseURL, b.AudioID))
	if err != nil {
		// shouldn't happen
		panic(err)
	}

	if apiKey != "" {
		v := u.Query()
		v.Set("apikey", apiKey)
		u.RawQuery = v.Encode()
	}
	return u
}

func (b AudioURLBuilder) GetCaptionURL() string {
	return b.BaseURL + "/audio/" + b.AudioID + "/caption"
}
