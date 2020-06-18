package scraper

import (
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/yaml.v2"
)

// adapted from https://www.freeones.com/html/m_links/bio_Mia_Malkova.php
const htmlDoc1 = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en" dir="ltr">
	<head>
		<title>Freeones:  Mia Malkova Biography</title>
	</head>
	<body data-babe="Mia Malkova">
		<div class="ContentBlock Block1">
			<div class="ContentBlockBody" style="padding: 0px;">
				<table id="biographyTable" border="0" cellspacing="0" cellpadding="0" width="100%">
					<tbody>
						<tr>
							<td class="paramname">
								<div><b>Babe Name:</b></div>
							</td>
							<td class="paramvalue">
								<a href="/html/m_links/Mia_Malkova/">Mia Malkova</a>&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<div><b>Profession:</b></div>
							</td>
							<td class="paramvalue">Porn Star
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Ethnicity:</b>
							</td>
							<td class="paramvalue">
								Caucasian&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Country of Origin:</b>
							</td>
							<td class="paramvalue">
								
								<span class="country-us">
								
									United States
								<span>
							</span></span></td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Date of Birth:</b>
							</td>
							<td class="paramvalue">
								July 1, 1992 (27 years old)&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Aliases:</b>
							</td>
							<td class="paramvalue">
								Mia Bliss, Madison Clover, Madison Swan, Mia Mountain, Jessica&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Eye Color:</b>
							</td>
							<td class="paramvalue">
								Hazel&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Hair Color:</b>
							</td>
							<td class="paramvalue">
								Blonde&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Height:</b>
							</td>
							<td class="paramvalue">
								<script type="text/javascript">
									<!--
									heightcm = "171";
									morethenone = 'inch';
									feet = heightcm / 30.48;
									inches = (feet - Math.floor(feet)) * 30.48 / 2.54;
					
									feet = Math.floor(feet);
									inches = inches.toFixed(0);
					
									if (inches > 1) {
										morethenone = 'inches';
									}
					
									if (heightcm == 0) {
										message = 'Unknown';
									} else {
										message = '171 cm - ' + feet + ' feet and ' + inches + ' ' + morethenone;
									}
									document.write(message);
									// -->
								</script>&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Measurements:</b>
							</td>
							<td class="paramvalue">
								34C-26-36
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Fake boobs:</b>
							</td>
							<td class="paramvalue">
								No&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Career Start And End</b>
							</td>
							<td class="paramvalue">
								2012 - 2019
								(7 Years In The Business)
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Tattoos:</b>
							</td>
							<td class="paramvalue">
								None&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<b>Piercings:</b>
							</td>
							<td class="paramvalue">
								None&nbsp;
							</td>
						</tr>
						<tr>
							<td class="paramname">
								<div><b>Social Network Links:</b></div>
							</td>
							<td class="paramvalue">
								<ul id="socialmedia">
									<li class="twitter"><a href="https://twitter.com/MiaMalkova" target="_blank" alt="Mia Malkova Twitter" title="Mia Malkova Twitter">Twitter</a></li>
									<li class="facebook"><a href="https://www.facebook.com/MiaMalcove" target="_blank" alt="Mia Malkova Facebook" title="Mia Malkova Facebook">Facebook</a></li>
									<li class="youtube"><a href="https://www.youtube.com/channel/UCEPR0sZKa_ScMoyhemfB7nA" target="_blank" alt="Mia Malkova YouTube" title="Mia Malkova YouTube">YouTube</a></li>
									<li class="instagram"><a href="https://www.instagram.com/mia_malkova/" target="_blank" alt="Mia Malkova Instagram" title="Mia Malkova Instagram">Instagram</a></li>
								</ul>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</body>
</html>
`

func makeCommonXPath(attr string) string {
	return `//table[@id="biographyTable"]//tr/td[@class="paramname"]//b[text() = '` + attr + `']/ancestor::tr/td[@class="paramvalue"]`
}

func makeReplaceRegex(regex string, with string) map[interface{}]interface{} {
	ret := make(map[interface{}]interface{})

	ret["regex"] = regex
	ret["with"] = with
	return ret
}

func makeXPathConfig() xpathScraperConfig {
	config := make(xpathScraperConfig)

	config["Name"] = makeCommonXPath("Babe Name:") + `/a`
	config["Ethnicity"] = makeCommonXPath("Ethnicity:")
	config["Country"] = makeCommonXPath("Country of Origin:")
	config["Aliases"] = makeCommonXPath("Aliases:")
	config["EyeColor"] = makeCommonXPath("Eye Color:")
	config["Measurements"] = makeCommonXPath("Measurements:")
	config["FakeTits"] = makeCommonXPath("Fake boobs:")
	config["Height"] = makeCommonXPath("Height:")
	config["Tattoos"] = makeCommonXPath("Tattoos:")
	config["Piercings"] = makeCommonXPath("Piercings:")

	// special handling for birthdate
	birthdateAttrConfig := make(map[interface{}]interface{})
	birthdateAttrConfig["selector"] = makeCommonXPath("Date of Birth:")

	var birthdateReplace []interface{}
	birthdateReplace = append(birthdateReplace, makeReplaceRegex(` \(.* years old\)`, ""))

	birthdateAttrConfig["replace"] = birthdateReplace
	birthdateAttrConfig["parseDate"] = "January 2, 2006" // "July 1, 1992 (27 years old)&nbsp;"
	config["Birthdate"] = birthdateAttrConfig

	// special handling for career length
	careerLengthAttrConfig := make(map[interface{}]interface{})
	// no colon in attribute header
	careerLengthAttrConfig["selector"] = makeCommonXPath("Career Start And End")

	var careerLengthReplace []interface{}
	careerLengthReplace = append(careerLengthReplace, makeReplaceRegex(`\s+\(.*\)`, ""))
	careerLengthAttrConfig["replace"] = careerLengthReplace

	config["CareerLength"] = careerLengthAttrConfig

	return config
}

func verifyField(t *testing.T, expected string, actual *string, field string) {
	t.Helper()

	if actual == nil || *actual != expected {
		if actual == nil {
			t.Errorf("Expected %s to be set to %s, instead got nil", field, expected)
		} else {
			t.Errorf("Expected %s to be set to %s, instead got %s", field, expected, *actual)
		}
	}
}

func TestScrapePerformerXPath(t *testing.T) {
	reader := strings.NewReader(htmlDoc1)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	xpathConfig := makeXPathConfig()

	scraper := xpathScraper{
		Performer: xpathConfig,
	}

	performer, err := scraper.scrapePerformer(doc)

	if err != nil {
		t.Errorf("Error scraping performer: %s", err.Error())
		return
	}

	const performerName = "Mia Malkova"
	const ethnicity = "Caucasian"
	const country = "United States"
	const birthdate = "1992-07-01"
	const aliases = "Mia Bliss, Madison Clover, Madison Swan, Mia Mountain, Jessica"
	const eyeColor = "Hazel"
	const measurements = "34C-26-36"
	const fakeTits = "No"
	const careerLength = "2012 - 2019"
	const tattoosPiercings = "None"

	verifyField(t, performerName, performer.Name, "Name")
	verifyField(t, ethnicity, performer.Ethnicity, "Ethnicity")
	verifyField(t, country, performer.Country, "Country")

	verifyField(t, birthdate, performer.Birthdate, "Birthdate")

	verifyField(t, aliases, performer.Aliases, "Aliases")
	verifyField(t, eyeColor, performer.EyeColor, "EyeColor")
	verifyField(t, measurements, performer.Measurements, "Measurements")
	verifyField(t, fakeTits, performer.FakeTits, "FakeTits")

	verifyField(t, careerLength, performer.CareerLength, "CareerLength")

	verifyField(t, tattoosPiercings, performer.Tattoos, "Tattoos")
	verifyField(t, tattoosPiercings, performer.Piercings, "Piercings")
}

func TestConcatXPath(t *testing.T) {
	const firstName = "FirstName"
	const lastName = "LastName"
	const eyeColor = "EyeColor"
	const separator = " "
	const testDoc = `
	<html>
	<div>` + firstName + `</div>
	<div>` + lastName + `</div>
	<span>` + eyeColor + `</span>
	</html>
	`

	reader := strings.NewReader(testDoc)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	xpathConfig := make(xpathScraperConfig)
	nameAttrConfig := make(map[interface{}]interface{})
	nameAttrConfig["selector"] = "//div"
	nameAttrConfig["concat"] = separator
	xpathConfig["Name"] = nameAttrConfig
	xpathConfig["EyeColor"] = "//span"

	scraper := xpathScraper{
		Performer: xpathConfig,
	}

	performer, err := scraper.scrapePerformer(doc)

	if err != nil {
		t.Errorf("Error scraping performer: %s", err.Error())
		return
	}

	const performerName = firstName + separator + lastName

	verifyField(t, performerName, performer.Name, "Name")
	verifyField(t, eyeColor, performer.EyeColor, "EyeColor")
}

const sceneHTML = `
<!DOCTYPE html>

<head>
    <title>Test Video - Pornhub.com</title>

    <meta property="og:title" content="Test Video" />
    <meta property="og:description"
        content="Watch Test Video on Pornhub.com, the best hardcore porn site. Pornhub is home to the widest selection of free Babe sex videos full of the hottest pornstars. If you&#039;re craving 3some XXX movies you&#039;ll find them here." />
    <meta property="og:image"
        content="https://di.phncdn.com/videos/201910/13/254476211/thumbs_80/(m=eaAaGwObaaaa)(mh=_V1YEGdMFS1rEYoW)9.jpg" />

    <script type="application/ld+json">
		{
			"@context": "http://schema.org/",
			"@type": "VideoObject",
			"name": "Test Video",
			"embedUrl": "https://www.pornhub.com/embed/ph5da270596459c",
			"duration": "PT00H33M27S",
			"thumbnailUrl": "https://di.phncdn.com/videos/201910/13/254476211/thumbs_80/(m=eaAaGwObaaaa)(mh=_V1YEGdMFS1rEYoW)9.jpg",
			"uploadDate": "2019-10-13T00:33:51+00:00",
			"description": "Watch Test Video on Pornhub&period;com&comma; the best hardcore porn site&period; Pornhub is home to the widest selection of free Babe sex videos full of the hottest pornstars&period; If you&apos;re craving 3some XXX movies you&apos;ll find them here&period;",
				"author" : "Mia Malkova",                "interactionStatistic": [
			{
					"@type": "InteractionCounter",
					"interactionType": "http://schema.org/WatchAction",
					"userInteractionCount": "5,908,861"
			},
			{
					"@type": "InteractionCounter",
					"interactionType": "http://schema.org/LikeAction",
					"userInteractionCount": "22,090"
				}
			]
		}
	</script>
</head>

<body class="logged-out">
    <div class="container  ">


        <div id="main-container" class="clearfix" data-delete-check="1" data-is-private="1" data-is-premium=""
            data-liu="0" data-next-shuffle="ph5da270596459c" data-pkey="" data-platform-pc="1" data-playlist-check="0"
            data-playlist-id-check="0" data-playlist-geo-check="0" data-friend="0" data-playlist-user-check="0"
            data-playlist-video-check="0" data-playlist-shuffle="0" data-shuffle-forward="ph5da270596459c"
            data-shuffle-back="ph5da270596459c" data-min-large="1350"
            data-video-title="Test Video">

            <div id="vpContentContainer">
                <div id="hd-leftColVideoPage">
                    <div class="video-wrapper">
                        <div class="title-container">
                            <i class="isMe tooltipTrig" data-title="Video of verified member"></i>
                            <h1 class="title">
                                <span class="inlineFree">Test Video</span>
                            </h1>
                        </div>

                        <div class="video-actions-container">
                            <div class="video-actions-tabs">
                                <div class="video-action-tab about-tab active">
                                    <div class="video-detailed-info">
                                        <div class="video-info-row">
                                            From:&nbsp;

                                            <div class="usernameWrap clearfix" data-type="channel" data-userid="492538092"
                                                data-liu-user="0"
                                                data-json-url="/user/box?id=492538092&amp;token=MTU3NzA1NTkzNIqATol8v_WrhmNTXkeflvG09C2U7UUT_NyoZUFa7iKq0mlzBkmdgAH1aNHZkJmIOHbbwmho1BehHDoA63K5Wn4."
                                                data-disable-popover="0">

                                                <a rel="" href="/channels/sis-loves-me" class="bolded">Sis Loves Me</a>
                                                <div class="avatarPosition"></div>
                                            </div>

                                            <span class="verified-icon flag tooltipTrig"
                                                data-title="Verified member"></span>
                                            - 87 videos
                                            <span class="subscribers-count">&nbsp;459466</span>
                                        </div>

                                        <div class="video-info-row">
                                            <div class="pornstarsWrapper">
                                                Pornstars:&nbsp;
                                                <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Alex D" data-id="251341" data-login="1"
                                                    href="/pornstar/alex-d">Alex D <span
                                                        class="psbox-link-container display-none"></span>
                                                </a>
                                                , <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Mia Malkova" data-id="10641" data-login="1"
                                                    href="/pornstar/mia-malkova">Mia Malkova <span
                                                        class="psbox-link-container display-none"></span>
                                                </a>
                                                , <a class="pstar-list-btn js-mxp" data-mxptype="Pornstar"
                                                    data-mxptext="Riley Reid" data-id="5343" data-login="1"
                                                    href="/pornstar/riley-reid">Riley Reid <span
                                                        class="psbox-link-container display-none"></span>
                                                </a>
                                                <div class="tooltipTrig suggestBtn" data-title="Add a pornstar">
                                                    <a class="add-btn-small add-pornstar-btn-2">+
                                                        <span>Suggest</span></a>
                                                </div>
                                                <div id="deletePornstarResult" class="suggest-result"></div>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            <div class="categoriesWrapper">
                                                Categories:&nbsp;
                                                <a href="/video?c=3"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Amateur</a>,
                                                <a href="/categories/babe"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Babe</a>,
                                                <a href="/video?c=13"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Blowjob</a>,
                                                <a href="/video?c=115"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Exclusive</a>,
                                                <a href="/hd"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">HD
                                                    Porn</a>, <a href="/categories/pornstar"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Pornstar</a>,
                                                <a href="/video?c=24"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Public</a>,
                                                <a href="/video?c=131"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Pussy
                                                    Licking</a>, <a href="/video?c=65"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Threesome</a>,
                                                <a href="/video?c=139"
                                                    onclick="ga('send', 'event', 'Watch Page', 'click', 'Category');">Verified
                                                    Models</a>
                                                <div class="tooltipTrig suggestBtn" data-title="Suggest Categories">
                                                    <a id="categoryLink" class="add-btn-small ">+
                                                        <span>Suggest</span></a>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            <div class="productionWrapper">
                                                Production:&nbsp;
                                                <a href="/video?p=professional" rel="nofollow"
                                                    class="production">professional</a>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            <div class="tagsWrapper">
                                                Tags:&nbsp;
                                                <a href="/video/search?search=3some">3some</a>, <a
                                                    href="/video?c=9">blonde</a>, <a href="/video?c=59">small tits</a>,
                                                <a href="/video/search?search=butt">butt</a>, <a
                                                    href="/video/search?search=natural+tits">natural tits</a>, <a
                                                    href="/video/search?search=petite">petite</a>, <a
                                                    href="/video?c=24">public</a>, <a
                                                    href="/video/search?search=outside">outside</a>, <a
                                                    href="/video/search?search=car">car</a>, <a
                                                    href="/video/search?search=garage">garage</a>, <a
                                                    href="/video?c=65">threesome</a>, <a
                                                    href="/video/search?search=bgg">bgg</a>, <a
                                                    href="/video/search?search=girlfrien+d">girlfrien d</a>, <a
                                                    href="/video/search?search=parking">parking</a>, <a
                                                    href="/video/search?search=sex">sex</a>, <a
                                                    href="/video/search?search=gagging">gagging</a>, <a
                                                    href="/video?c=13">blowjob</a>, <a
                                                    href="/video/search?search=bj">bj</a>, <a
                                                    href="/video/search?search=double">double</a>, <a
                                                    href="/video/search?search=ass">ass</a>
                                                <div class="tooltipTrig suggestBtn" data-title="Suggest Tags">
                                                    <a id="tagLink" class="add-btn-small">+ <span>Suggest</span></a>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="video-info-row showLess">
                                            Added on: <span class="white">2 months ago</span>
                                        </div>

                                        <div class="video-info-row showLess">
                                            Featured on: <span class="white">1 month ago</span>
                                        </div>
                                    </div>
                                </div>

                                <div class="video-action-tab jump-to-tab">
                                    <div class="title">Jump to your favorite action</div>

                                    <div class="filters mainFilter float-right">
                                        <div class="dropdownTrigger">
                                            <div>
                                                <span class="textFilter" id="tagSort">Sequence</span>
                                                <span class="arrowFilters"></span>
                                            </div>
                                            <ul class="filterListItem dropdownWrapper">
                                                <li class="active"><a class="actionTagSort"
                                                        data-sort="seconds">Sequence</a></li>
                                                <li><a class="actionTagSort" data-sort="tag">Alphabetical</a></li>
                                            </ul>
                                        </div>
                                    </div>

                                    <div class="reset"></div>
                                    <div class="display-grid col-4 gap-row-none sortBy seconds">
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(862), ga('send', 'event', 'Video Page', 'click', 'Jump to Blowjob');">
                                                    Blowjob </a>
                                                &nbsp;
                                                <var>14:22</var>
                                            </li>
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1117), ga('send', 'event', 'Video Page', 'click', 'Jump to Reverse Cowgirl');">
                                                    Reverse Cowgirl </a>
                                                &nbsp;
                                                <var>18:37</var>
                                            </li>
                                        </ul>
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1182), ga('send', 'event', 'Video Page', 'click', 'Jump to Cowgirl');">
                                                    Cowgirl </a>
                                                &nbsp;
                                                <var>19:42</var>
                                            </li>
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1625), ga('send', 'event', 'Video Page', 'click', 'Jump to Cowgirl');">
                                                    Cowgirl </a>
                                                &nbsp;
                                                <var>27:05</var>
                                            </li>
                                        </ul>
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1822), ga('send', 'event', 'Video Page', 'click', 'Jump to Doggystyle');">
                                                    Doggystyle </a>
                                                &nbsp;
                                                <var>30:22</var>
                                            </li>
                                        </ul>

                                    </div>
                                    <div class="display-grid col-4 gap-row-none sortBy tag">
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(862), ga('send', 'event', 'Video Page', 'click', 'Jump to Blowjob');">
                                                    Blowjob </a>
                                                &nbsp;
                                                <var>14:22</var>
                                            </li>
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1117), ga('send', 'event', 'Video Page', 'click', 'Jump to Reverse Cowgirl');">
                                                    Reverse Cowgirl </a>
                                                &nbsp;
                                                <var>18:37</var>
                                            </li>
                                        </ul>
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1182), ga('send', 'event', 'Video Page', 'click', 'Jump to Cowgirl');">
                                                    Cowgirl </a>
                                                &nbsp;
                                                <var>19:42</var>
                                            </li>
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1625), ga('send', 'event', 'Video Page', 'click', 'Jump to Cowgirl');">
                                                    Cowgirl </a>
                                                &nbsp;
                                                <var>27:05</var>
                                            </li>
                                        </ul>
                                        <ul class="actionTagList full-width margin-none">
                                            <li>
                                                <a class="js-triggerJumpCat"
                                                    onclick="jumpToAction(1822), ga('send', 'event', 'Video Page', 'click', 'Jump to Doggystyle');">
                                                    Doggystyle </a>
                                                &nbsp;
                                                <var>30:22</var>
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

func makeSceneXPathConfig() xpathScraper {
	common := make(commonXPathConfig)

	common["$performerElem"] = `//div[@class="pornstarsWrapper"]/a[@data-mxptype="Pornstar"]`
	common["$studioElem"] = `//div[@data-type="channel"]/a`

	config := make(xpathScraperConfig)

	config["Title"] = `//meta[@property="og:title"]/@content`
	// this needs post-processing
	config["Date"] = `//script[@type="application/ld+json"]`

	tagConfig := make(map[interface{}]interface{})
	tagConfig["Name"] = `//div[@class="categoriesWrapper"]//a[not(@class="add-btn-small ")]`
	config["Tags"] = tagConfig

	performerConfig := make(map[interface{}]interface{})
	performerConfig["Name"] = `$performerElem/@data-mxptext`
	performerConfig["URL"] = `$performerElem/@href`
	config["Performers"] = performerConfig

	studioConfig := make(map[interface{}]interface{})
	studioConfig["Name"] = `$studioElem`
	studioConfig["URL"] = `$studioElem/@href`
	config["Studio"] = studioConfig

	const sep = " "
	moviesNameConfig := make(map[interface{}]interface{})
	moviesNameConfig["selector"] = `//i[@class="isMe tooltipTrig"]/@data-title`
	moviesNameConfig["split"] = sep
	moviesConfig := make(map[interface{}]interface{})
	moviesConfig["Name"] = moviesNameConfig
	config["Movies"] = moviesConfig

	scraper := xpathScraper{
		Scene:  config,
		Common: common,
	}

	return scraper
}

func verifyTags(t *testing.T, expectedTagNames []string, actualTags []*models.ScrapedSceneTag) {
	t.Helper()

	i := 0
	for i < len(expectedTagNames) || i < len(actualTags) {
		expectedTag := ""
		actualTag := ""
		if i < len(expectedTagNames) {
			expectedTag = expectedTagNames[i]
		}
		if i < len(actualTags) {
			actualTag = actualTags[i].Name
		}

		if expectedTag != actualTag {
			t.Errorf("Expected tag %s, got %s", expectedTag, actualTag)
		}
		i++
	}
}

func verifyMovies(t *testing.T, expectedMovieNames []string, actualMovies []*models.ScrapedSceneMovie) {
	t.Helper()

	i := 0
	for i < len(expectedMovieNames) || i < len(actualMovies) {
		expectedMovie := ""
		actualMovie := ""
		if i < len(expectedMovieNames) {
			expectedMovie = expectedMovieNames[i]
		}
		if i < len(actualMovies) {
			actualMovie = actualMovies[i].Name
		}

		if expectedMovie != actualMovie {
			t.Errorf("Expected movie %s, got %s", expectedMovie, actualMovie)
		}
		i++
	}
}

func verifyPerformers(t *testing.T, expectedNames []string, expectedURLs []string, actualPerformers []*models.ScrapedScenePerformer) {
	t.Helper()

	i := 0
	for i < len(expectedNames) || i < len(actualPerformers) {
		expectedName := ""
		actualName := ""
		expectedURL := ""
		actualURL := ""
		if i < len(expectedNames) {
			expectedName = expectedNames[i]
		}
		if i < len(expectedURLs) {
			expectedURL = expectedURLs[i]
		}
		if i < len(actualPerformers) {
			actualName = actualPerformers[i].Name
			if actualPerformers[i].URL != nil {
				actualURL = *actualPerformers[i].URL
			}
		}

		if expectedName != actualName {
			t.Errorf("Expected performer name %s, got %s", expectedName, actualName)
		}
		if expectedURL != actualURL {
			t.Errorf("Expected perfromer URL %s, got %s", expectedName, actualName)
		}
		i++
	}
}

func TestApplySceneXPathConfig(t *testing.T) {
	reader := strings.NewReader(sceneHTML)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	scraper := makeSceneXPathConfig()

	scene, err := scraper.scrapeScene(doc)

	if err != nil {
		t.Errorf("Error scraping scene: %s", err.Error())
		return
	}

	const title = "Test Video"

	verifyField(t, title, scene.Title, "Title")

	// verify tags
	expectedTags := []string{
		"Amateur",
		"Babe",
		"Blowjob",
		"Exclusive",
		"HD Porn",
		"Pornstar",
		"Public",
		"Pussy Licking",
		"Threesome",
		"Verified Models",
	}
	verifyTags(t, expectedTags, scene.Tags)

	// verify movies
	expectedMovies := []string{
		"Video",
		"of",
		"verified",
		"member",
	}
	verifyMovies(t, expectedMovies, scene.Movies)

	expectedPerformerNames := []string{
		"Alex D",
		"Mia Malkova",
		"Riley Reid",
	}

	expectedPerformerURLs := []string{
		"/pornstar/alex-d",
		"/pornstar/mia-malkova",
		"/pornstar/riley-reid",
	}

	verifyPerformers(t, expectedPerformerNames, expectedPerformerURLs, scene.Performers)

	const expectedStudioName = "Sis Loves Me"
	const expectedStudioURL = "/channels/sis-loves-me"

	verifyField(t, expectedStudioName, &scene.Studio.Name, "Studio.Name")
	verifyField(t, expectedStudioURL, scene.Studio.URL, "Studio.URL")
}

func TestLoadXPathScraperFromYAML(t *testing.T) {
	const yamlStr = `name: Test
performerByURL:
  - action: scrapeXPath
    url: 
      - test.com
    scraper: performerScraper
xPathScrapers:
  performerScraper:
    performer:
      name: //h1[@itemprop="name"]
`

	config := &scraperConfig{}
	err := yaml.Unmarshal([]byte(yamlStr), &config)

	if err != nil {
		t.Errorf("Error loading yaml: %s", err.Error())
		return
	}
}

func TestLoadInvalidXPath(t *testing.T) {
	config := make(xpathScraperConfig)

	config["Name"] = `//a[id=']/span`

	reader := strings.NewReader(htmlDoc1)
	doc, err := htmlquery.Parse(reader)

	if err != nil {
		t.Errorf("Error loading document: %s", err.Error())
		return
	}

	common := make(commonXPathConfig)
	config.process(doc, common)
}
