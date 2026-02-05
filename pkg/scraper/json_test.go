package scraper

import (
	"context"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestJsonPerformerScraper(t *testing.T) {
	const yamlStr = `name: Test
jsonScrapers:
  performerScraper:
    common:
      $extras: data.extras
    performer:
      Name: data.name
      Gender: $extras.gender
      Birthdate: $extras.birthday
      Ethnicity: $extras.ethnicity
      Height: $extras.height
      Measurements: $extras.measurements
      Tattoos: $extras.tattoos
      Piercings: $extras.piercings
      Aliases: data.aliases
      Image: data.image
      Details: data.bio
      HairColor: $extras.hair_colour
      Weight: $extras.weight
`

	const json = `
{
	"data": {
        "id": "2cd4146b-637d-49b1-8ff9-19d4a06947bb",
        "name": "Mia Malkova",
        "bio": "Some girls are so damn hot that they can get you bent out of shape, and you will not even be mad at them for doing so. Well, tawny blonde Mia Malkova can bend her body into any shape she pleases, and that’s sure to satisfy all of the horny cocks and wet pussies out there. This girl has acrobatic and contortionist abilities that could even twist a pretzel into a new knot, which can be very helpful in the ... arrow_drop_down Some girls are so damn hot that they can get you bent out of shape, and you will not even be mad at them for doing so. Well, tawny blonde Mia Malkova can bend her body into any shape she pleases, and that’s sure to satisfy all of the horny cocks and wet pussies out there. This girl has acrobatic and contortionist abilities that could even twist a pretzel into a new knot, which can be very helpful in the VR Porn movies – trust us. Ankles behind her neck and feet over her back so she can kiss her toes, turned, twisted and gyrating, she can fuck any which way she wants (and that ass!), will surely make you fall in love with this hot Virtual Reality Porn slut, as she is one of the finest of them all. Talking about perfection, maybe it’s all the acrobatic work that keeps it in such gorgeous shape? Who cares really, because you just want to take a big bite out of it and never let go. But it’s not all about the body. Mia’s also got a great smile, which might not sound kinky, but believe us, it is a smile that will heat up your innards and drop your pants. Is it her golden skin, her innocent pink lips or that heart-shaped face? There is just too much good stuff going on with Mia Malkova, which is maybe why these past few years have heaped awards upon awards on this Southern California native. Mia came to VR Bangers for her first VR Porn video, so you know she’s only going for top-notch scenes with top-game performers, men, and women. Better hit up that yoga studio if you ever dream of being able to bang a flexible and talented chick like lady Malkova. arrow_drop_up",
        "extras": {
            "gender": "Female",
            "birthday": "1992-07-01",
            "birthday_timestamp": 709948800,
            "birthplace": "Palm Springs, California, United States",
            "active": 1,
            "astrology": "Cancer (Jun 21 - Jul 22)",
            "ethnicity": "Caucasian",
            "nationality": "United States",
            "hair_colour": "Blonde",
            "weight": 57,
            "height": "5'6\" (or 167 cm)",
            "measurements": "34-26-36",
            "cupsize": "34C (75C)",
            "tattoos": "None",
            "piercings": "Navel",
            "first_seen": null
        },
        "aliases": [
            "Mia Bliss",
            "Madison Clover",
            "Madison Swan",
            "Mia Mountain",
            "Mia M.",
            "Mia Malvoka",
            "Mia Molkova",
            "Mia Thomas"
        ],
		"image": "https:\/\/thumb.metadataapi.net\/unsafe\/1000x1500\/smart\/filters:sharpen():upscale()\/https%3A%2F%2Fcdn.metadataapi.net%2Fperformer%2F49%2F05%2F30%2Fade2255dc065032a89ebb23f0e038fa%2Fposter%2Fmia-malkova.jpg%3Fid1582610531"
	}
}
`

	c := &Definition{}
	err := yaml.Unmarshal([]byte(yamlStr), &c)

	if err != nil {
		t.Fatalf("Error loading yaml: %s", err.Error())
	}

	// perform scrape using json string
	performerScraper := c.JsonScrapers["performerScraper"]

	q := &jsonQuery{
		doc: json,
	}

	scrapedPerformer, err := performerScraper.scrapePerformer(context.Background(), q)
	if err != nil {
		t.Fatalf("Error scraping performer: %s", err.Error())
	}

	verifyField(t, "Mia Malkova", scrapedPerformer.Name, "Name")
	verifyField(t, "Female", scrapedPerformer.Gender, "Gender")
	verifyField(t, "1992-07-01", scrapedPerformer.Birthdate, "Birthdate")
	verifyField(t, "Caucasian", scrapedPerformer.Ethnicity, "Ethnicity")
	verifyField(t, "5'6\" (or 167 cm)", scrapedPerformer.Height, "Height")
	verifyField(t, "None", scrapedPerformer.Tattoos, "Tattoos")
	verifyField(t, "Navel", scrapedPerformer.Piercings, "Piercings")
	verifyField(t, "Some girls are so damn hot that they can get you bent out of shape, and you will not even be mad at them for doing so. Well, tawny blonde Mia Malkova can bend her body into any shape she pleases, and that’s sure to satisfy all of the horny cocks and wet pussies out there. This girl has acrobatic and contortionist abilities that could even twist a pretzel into a new knot, which can be very helpful in the ... arrow_drop_down Some girls are so damn hot that they can get you bent out of shape, and you will not even be mad at them for doing so. Well, tawny blonde Mia Malkova can bend her body into any shape she pleases, and that’s sure to satisfy all of the horny cocks and wet pussies out there. This girl has acrobatic and contortionist abilities that could even twist a pretzel into a new knot, which can be very helpful in the VR Porn movies – trust us. Ankles behind her neck and feet over her back so she can kiss her toes, turned, twisted and gyrating, she can fuck any which way she wants (and that ass!), will surely make you fall in love with this hot Virtual Reality Porn slut, as she is one of the finest of them all. Talking about perfection, maybe it’s all the acrobatic work that keeps it in such gorgeous shape? Who cares really, because you just want to take a big bite out of it and never let go. But it’s not all about the body. Mia’s also got a great smile, which might not sound kinky, but believe us, it is a smile that will heat up your innards and drop your pants. Is it her golden skin, her innocent pink lips or that heart-shaped face? There is just too much good stuff going on with Mia Malkova, which is maybe why these past few years have heaped awards upon awards on this Southern California native. Mia came to VR Bangers for her first VR Porn video, so you know she’s only going for top-notch scenes with top-game performers, men, and women. Better hit up that yoga studio if you ever dream of being able to bang a flexible and talented chick like lady Malkova. arrow_drop_up", scrapedPerformer.Details, "Details")
	verifyField(t, "Blonde", scrapedPerformer.HairColor, "HairColor")
	verifyField(t, "57", scrapedPerformer.Weight, "Weight")

	notFoundJson := `
{
    "data": null
}`

	q = &jsonQuery{
		doc: notFoundJson,
	}

	scrapedPerformer, err = performerScraper.scrapePerformer(context.Background(), q)
	if err != nil {
		t.Fatalf("Error scraping performer: %s", err.Error())
	}

	if scrapedPerformer != nil {
		t.Errorf("expected nil scraped performer when not found, got %v", scrapedPerformer)
	}
}
