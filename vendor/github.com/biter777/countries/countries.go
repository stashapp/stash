// Package countries - ISO 3166 (ISO3166-1, ISO3166, Digit, Alpha-2 and Alpha-3) countries codes and names (on eng and rus), ISO 4217 currency designators, ITU-T E.164 IDD calling phone codes, countries capitals, UN M.49 regions codes, ccTLD countries domains, IOC/NOC and FIFA letters codes, VERY FAST, NO maps[], NO slices[], NO external links/files/data, NO interface{}, NO specific dependencies, Databases compatible, Emoji countries flags and currencies support, full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts. Full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts.
/*
Package countries - ISO 3166 (ISO3166-1, ISO3166, Digit, Alpha-2 and Alpha-3) countries codes and names (on eng and rus), ISO 4217 currency designators, ITU-T E.164 IDD calling phone codes, countries capitals, UN M.49 regions codes, ccTLD countries domains, IOC/NOC and FIFA letters codes, VERY FAST, NO maps[], NO slices[], NO external links/files/data, NO interface{}, NO specific dependencies, Databases compatible, Emoji countries flags and currencies support, full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts.
Full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts.

Usage

func main() {
	countryJapan := countries.Japan
	fmt.Printf("Country name in english: %v\n", countryJapan)                   // Japan
	fmt.Printf("Country name in russian: %v\n", countryJapan.StringRus())       // Япония
	fmt.Printf("Country ISO-3166 digit code: %d\n", countryJapan)               // 392
	fmt.Printf("Country ISO-3166 Alpha-2 code: %v\n", countryJapan.Alpha2())    // JP
	fmt.Printf("Country ISO-3166 Alpha-3 code: %v\n", countryJapan.Alpha3())    // JPN
	fmt.Printf("Country IOC/NOC code: %v\n", countryJapan.IOC())                // JPN
	fmt.Printf("Country FIFA code: %v\n", countryJapan.FIFA())                  // JPN
	fmt.Printf("Country Capital: %v\n", countryJapan.Capital())                 // Tokyo
	fmt.Printf("Country ITU-T E.164 call code: %v\n", countryJapan.CallCodes()) // +81
	fmt.Printf("Country ccTLD domain: %v\n", countryJapan.Domain())             // .jp
	fmt.Printf("Country UN M.49 region name: %v\n", countryJapan.Region())      // Asia
	fmt.Printf("Country UN M.49 region code: %d\n", countryJapan.Region())      // 142
	fmt.Printf("Country emoji/flag: %v\n\n", countryJapan.Emoji())              // 🇯🇵

	currencyJapan := countryJapan.Currency()
	fmt.Printf("Country ISO-4217 Currency name in english: %v\n", currencyJapan)           // Yen
	fmt.Printf("Country ISO-4217 Currency digit code: %d\n", currencyJapan)                // 392
	fmt.Printf("Country ISO-4217 Currency Alpha code: %v\n", currencyJapan.Alpha())        // JPY
	fmt.Printf("Country Currency emoji: %v\n", currencyJapan.Emoji())                      // 💴
	fmt.Printf("Country of Currency %v: %v\n\n", currencyJapan, currencyJapan.Countries()) // Japan

	// OR you can alternative use:
	japanInfo := countries.Japan.Info()
	fmt.Printf("Country name in english: %v\n", japanInfo.Name)                          // Japan
	fmt.Printf("Country ISO-3166 digit code: %d\n", japanInfo.Code)                      // 392
	fmt.Printf("Country ISO-3166 Alpha-2 code: %v\n", japanInfo.Alpha2)                  // JP
	fmt.Printf("Country ISO-3166 Alpha-3 code: %v\n", japanInfo.Alpha3)                  // JPN
	fmt.Printf("Country IOC/NOC code: %v\n", japanInfo.IOC)                              // JPN
	fmt.Printf("Country FIFA code: %v\n", japanInfo.FIFA)                                // JPN
	fmt.Printf("Country Capital: %v\n", japanInfo.Capital)                               // Tokyo
	fmt.Printf("Country ITU-T E.164 call code: %v\n", japanInfo.CallCodes)               // +81
	fmt.Printf("Country ccTLD domain: %v\n", japanInfo.Domain)                           // .jp
	fmt.Printf("Country UN M.49 region name: %v\n", japanInfo.Region)                    // Asia
	fmt.Printf("Country UN M.49 region code: %d\n", japanInfo.Region)                    // 142
	fmt.Printf("Country emoji/flag: %v\n", japanInfo.Emoji)                              // 🇯🇵
	fmt.Printf("Country ISO-4217 Currency name in english: %v\n", japanInfo.Currency)    // Yen
	fmt.Printf("Country ISO-4217 Currency digit code: %d\n", japanInfo.Currency)         // 392
	fmt.Printf("Country ISO-4217 Currency Alpha code: %v\n", japanInfo.Currency.Alpha()) // JPY

	// Detection usage
	// Detect by name
	country := countries.ByName("angola")
	fmt.Printf("Country name in english: %v\n", country)                // Angola
	fmt.Printf("Country ISO-3166 digit code: %d\n", country)            // 24
	fmt.Printf("Country ISO-3166 Alpha-2 code: %v\n", country.Alpha2()) // AO
	fmt.Printf("Country ISO-3166 Alpha-3 code: %v\n", country.Alpha3()) // AGO
	// Detect by code/numeric
	country = countries.ByNumeric(24)
	fmt.Printf("Country name in english: %v\n", country)                // Angola
	fmt.Printf("Country ISO-3166 digit code: %d\n", country)            // 24
	fmt.Printf("Country ISO-3166 Alpha-2 code: %v\n", country.Alpha2()) // AO
	fmt.Printf("Country ISO-3166 Alpha-3 code: %v\n", country.Alpha3()) // AGO

	// Comparing usage
	// Compare by code/numeric
	if countries.ByName("angola") == countries.AGO {
		fmt.Println("Yes! It's Angola!") // Yes! It's Angola!
	}
	// Compare by name
	if strings.EqualFold("angola", countries.AGO.String()) {
		fmt.Println("Yes! It's Angola!") // Yes! It's Angola!
	}

	// Database usage
	type User struct {
		gorm.Model
		Name     string
		Country  countries.CountryCode
		Currency countries.CurrencyCode
	}
	user := &User{Name: "Helen", Country: countries.Slovenia, Currency: countries.CurrencyEUR}
	db, err := gorm.Open("postgres", 500, "host=127.0.0.2 port=5432 user=usr password=1234567 dbname=db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Create(user)
}

Contributing

 Welcome pull requests, bug fixes and issue reports.
 Before proposing a change, please discuss it first by raising an issue. */
// Package countries - ISO 3166 (ISO3166-1, ISO3166, Digit, Alpha-2 and Alpha-3) countries codes and names (on eng and rus), ISO 4217 currency designators, ITU-T E.164 IDD calling phone codes, countries capitals, UN M.49 regions codes, ccTLD countries domains, IOC/NOC and FIFA letters codes, VERY FAST, NO maps[], NO slices[], NO external links/files/data, NO interface{}, NO specific dependencies, Databases compatible, Emoji countries flags and currencies support, full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts. Full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and ccTLD standarts.
//nolint:misspell
package countries //nolint:misspell

import (
	"encoding/json"
	"fmt"
)

// CountryCode - country code (254 countries). Three codes present, for example Russia == RU == RUS == 643.
type CountryCode int64 // int64 for database/sql/driver.Valuer compatibility

// Country - all info about country
type Country struct {
	Name      string       `json:"name"`
	Alpha2    string       `json:"cca2"`
	Alpha3    string       `json:"cca3"`
	IOC       string       `json:"ioc"`
	FIFA      string       `json:"fifa"`
	Emoji     string       `json:"emoji"`
	Code      CountryCode  `json:"code"`
	Currency  CurrencyCode `json:"currency"`
	Capital   CapitalCode  `json:"capital"`
	CallCodes []CallCode   `json:"callingCode"`
	Domain    DomainCode   `json:"domain"`
	Region    RegionCode   `json:"region"`
}

// Typer - typer interface, provide a name of type
type Typer interface {
	Type() string
}

// Total - returns number of codes in the package, countries.Total() == len(countries.All()) but static value for performance
func Total() int {
	return 252
}

// Emoji - return a country Alpha-2 (ISO2) as Emoji flag (example "RU" as "🇷🇺")
func (c CountryCode) Emoji() string {
	iso2 := c.Alpha2()
	buf := [...]byte{240, 159, 135, 0, 240, 159, 135, 0}
	buf[3] = iso2[0] + (166 - 'A')
	buf[7] = iso2[1] + (166 - 'A')
	return string(buf[:])
}

// Emoji3 - return a country Alpha-3 (ISO3) as Emoji (example "RUS" as "🇷🇺🇸")
func (c CountryCode) Emoji3() string {
	iso3 := c.Alpha3()
	buf := [...]byte{240, 159, 135, 0, 240, 159, 135, 0, 240, 159, 135, 0}
	buf[3] = iso3[0] + (166 - 'A')
	buf[7] = iso3[1] + (166 - 'A')
	buf[11] = iso3[2] + (166 - 'A')
	return string(buf[:])
}

// Type implements Typer interface.
func (c CountryCode) Type() string {
	return TypeCountryCode
}

// String - implements fmt.Stringer, returns a english name of country
//nolint:gocyclo
func (c CountryCode) String() string { //nolint:gocyclo
	switch c {
	case 8:
		return "Albania"
	case 12:
		return "Algeria"
	case 16:
		return "American Samoa"
	case 20:
		return "Andorra"
	case 24:
		return "Angola"
	case 660:
		return "Anguilla"
	case 10:
		return "Antarctica"
	case 28:
		return "Antigua and Barbuda"
	case 32:
		return "Argentina"
	case 51:
		return "Armenia"
	case 533:
		return "Aruba"
	case 36:
		return "Australia"
	case 40:
		return "Austria"
	case 31:
		return "Azerbaijan"
	case 44:
		return "Bahamas"
	case 48:
		return "Bahrain"
	case 50:
		return "Bangladesh"
	case 52:
		return "Barbados"
	case 112:
		return "Belarus"
	case 56:
		return "Belgium"
	case 84:
		return "Belize"
	case 204:
		return "Benin"
	case 60:
		return "Bermuda"
	case 64:
		return "Bhutan"
	case 68:
		return "Bolivia"
	case 70:
		return "Bosnia and Herzegovina"
	case 72:
		return "Botswana"
	case 74:
		return "Bouvet Island"
	case 76:
		return "Brazil"
	case 86:
		return "British Indian Ocean Territory"
	case 96:
		return "Brunei Darussalam"
	case 100:
		return "Bulgaria"
	case 854:
		return "Burkina Faso"
	case 108:
		return "Burundi"
	case 116:
		return "Cambodia"
	case 120:
		return "Cameroon"
	case 124:
		return "Canada"
	case 132:
		return "Cape Verde"
	case 136:
		return "Cayman Islands"
	case 140:
		return "Central African Republic"
	case 148:
		return "Chad"
	case 152:
		return "Chile"
	case 156:
		return "China"
	case 162:
		return "Christmas Island"
	case 166:
		return "Cocos (Keeling) Islands"
	case 170:
		return "Colombia"
	case 174:
		return "Comoros"
	case 178:
		return "Congo"
	case 180:
		return "Democratic Republic of the Congo"
	case 184:
		return "Cook Islands"
	case 188:
		return "Costa Rica"
	case 384:
		return "Cote d`Ivoire" // Ivory Coast
	case 191:
		return "Croatia"
	case 192:
		return "Cuba"
	case 196:
		return "Cyprus"
	case 203:
		return "Czech Republic"
	case 208:
		return "Denmark"
	case 262:
		return "Djibouti"
	case 212:
		return "Dominica"
	case 214:
		return "Dominican Republic"
	case 218:
		return "Ecuador"
	case 818:
		return "Egypt"
	case 222:
		return "El Salvador"
	case 226:
		return "Equatorial Guinea"
	case 232:
		return "Eritrea"
	case 233:
		return "Estonia"
	case 231:
		return "Ethiopia"
	case 234:
		return "Faroe Islands"
	case 238:
		return "Falkland Islands (Malvinas)"
	case 242:
		return "Fiji"
	case 246:
		return "Finland"
	case 250:
		return "France"
	case 254:
		return "French Guiana"
	case 258:
		return "French Polynesia"
	case 260:
		return "French Southern Territories"
	case 266:
		return "Gabon"
	case 270:
		return "Gambia"
	case 268:
		return "Georgia"
	case 276:
		return "Germany"
	case 288:
		return "Ghana"
	case 292:
		return "Gibraltar"
	case 300:
		return "Greece"
	case 304:
		return "Greenland"
	case 308:
		return "Grenada"
	case 312:
		return "Guadeloupe"
	case 316:
		return "Guam"
	case 320:
		return "Guatemala"
	case 324:
		return "Guinea"
	case 624:
		return "Guinea-Bissau"
	case 328:
		return "Guyana"
	case 332:
		return "Haiti"
	case 334:
		return "Heard Island and McDonald Islands"
	case 340:
		return "Honduras"
	case 344:
		return "Hong Kong (Special Administrative Region of China)"
	case 348:
		return "Hungary"
	case 352:
		return "Iceland"
	case 356:
		return "India"
	case 360:
		return "Indonesia"
	case 364:
		return "Iran (Islamic Republic of)"
	case 368:
		return "Iraq"
	case 372:
		return "Ireland"
	case 376:
		return "Israel"
	case 380:
		return "Italy"
	case 388:
		return "Jamaica"
	case 392:
		return "Japan"
	case 400:
		return "Jordan"
	case 398:
		return "Kazakhstan"
	case 404:
		return "Kenya"
	case 296:
		return "Kiribati"
	case 410:
		return "Republic of Korea"
	case 408:
		return "Democratic People`s Republic of Korea"
	case 414:
		return "Kuwait"
	case 417:
		return "Kyrgyzstan"
	case 418:
		return "Lao People`s Democratic Republic"
	case 428:
		return "Latvia"
	case 422:
		return "Lebanon"
	case 426:
		return "Lesotho"
	case 430:
		return "Liberia"
	case 434:
		return "Libyan Arab Jamahiriya"
	case 438:
		return "Liechtenstein"
	case 440:
		return "Lithuania"
	case 442:
		return "Luxembourg"
	case 446:
		return "Macau (Special Administrative Region of China)"
	case 807:
		return "North Macedonia (Republic of North Macedonia)"
	case 450:
		return "Madagascar"
	case 454:
		return "Malawi"
	case 458:
		return "Malaysia"
	case 462:
		return "Maldives"
	case 466:
		return "Mali"
	case 470:
		return "Malta"
	case 584:
		return "Marshall Islands"
	case 474:
		return "Martinique"
	case 478:
		return "Mauritania"
	case 480:
		return "Mauritius"
	case 175:
		return "Mayotte"
	case 484:
		return "Mexico"
	case 583:
		return "Micronesia (Federated States of)"
	case 498:
		return "Moldova (Republic of)"
	case 492:
		return "Monaco"
	case 496:
		return "Mongolia"
	case 500:
		return "Montserrat"
	case 504:
		return "Morocco"
	case 508:
		return "Mozambique"
	case 104:
		return "Myanmar"
	case 516:
		return "Namibia"
	case 520:
		return "Nauru"
	case 524:
		return "Nepal"
	case 528:
		return "Netherlands"
	case 530:
		return "Netherlands Antilles"
	case 540:
		return "New Caledonia"
	case 554:
		return "New Zealand"
	case 558:
		return "Nicaragua"
	case 562:
		return "Niger"
	case 566:
		return "Nigeria"
	case 570:
		return "Niue"
	case 574:
		return "Norfolk Island"
	case 580:
		return "Northern Mariana Islands"
	case 578:
		return "Norway"
	case 512:
		return "Oman"
	case 586:
		return "Pakistan"
	case 585:
		return "Palau"
	case 275:
		return "Palestinian Territory (Occupied)"
	case 591:
		return "Panama"
	case 598:
		return "Papua New Guinea"
	case 600:
		return "Paraguay"
	case 604:
		return "Peru"
	case 608:
		return "Philippines"
	case 612:
		return "Pitcairn"
	case 616:
		return "Poland"
	case 620:
		return "Portugal"
	case 630:
		return "Puerto Rico"
	case 634:
		return "Qatar"
	case 638:
		return "Reunion"
	case 642:
		return "Romania"
	case 643:
		return "Russian Federation"
	case 646:
		return "Rwanda"
	case 654:
		return "Saint Helena"
	case 659:
		return "Saint Kitts and Nevis"
	case 662:
		return "Saint Lucia"
	case 666:
		return "Saint Pierre and Miquelon"
	case 670:
		return "Saint Vincent and the Grenadines"
	case 882:
		return "Samoa"
	case 674:
		return "San Marino"
	case 678:
		return "Sao Tome and Principe"
	case 682:
		return "Saudi Arabia"
	case 686:
		return "Senegal"
	case 690:
		return "Seychelles"
	case 694:
		return "Sierra Leone"
	case 702:
		return "Singapore"
	case 703:
		return "Slovakia"
	case 705:
		return "Slovenia"
	case 90:
		return "Solomon Islands"
	case 706:
		return "Somalia"
	case 710:
		return "South Africa"
	case 239:
		return "South Georgia and The South Sandwich Islands"
	case 724:
		return "Spain"
	case 144:
		return "Sri Lanka"
	case 736:
		return "Sudan"
	case 740:
		return "Suriname"
	case 744:
		return "Svalbard and Jan Mayen Islands"
	case 748:
		return "Swaziland"
	case 752:
		return "Sweden"
	case 756:
		return "Switzerland"
	case 760:
		return "Syrian Arab Republic"
	case 158:
		return "Taiwan (Province of China)"
	case 762:
		return "Tajikistan"
	case 834:
		return "Tanzania (United Republic of)"
	case 764:
		return "Thailand"
	case 626:
		return "Timor-Leste (East Timor)"
	case 768:
		return "Togo"
	case 772:
		return "Tokelau"
	case 776:
		return "Tonga"
	case 780:
		return "Trinidad and Tobago"
	case 788:
		return "Tunisia"
	case 792:
		return "Turkey"
	case 795:
		return "Turkmenistan"
	case 796:
		return "Turks and Caicos Islands"
	case 798:
		return "Tuvalu"
	case 800:
		return "Uganda"
	case 804:
		return "Ukraine"
	case 784:
		return "United Arab Emirates"
	case 826:
		return "United Kingdom"
	case 840:
		return "United States"
	case 581:
		return "United States Minor Outlying Islands"
	case 858:
		return "Uruguay"
	case 860:
		return "Uzbekistan"
	case 548:
		return "Vanuatu"
	case 336:
		return "Holy See (Vatican City State)"
	case 862:
		return "Venezuela"
	case 704:
		return "Vietnam"
	case 92:
		return "Virgin Islands British"
	case 850:
		return "Virgin Islands US"
	case 876:
		return "Wallis and Futuna Islands"
	case 732:
		return "Western Sahara"
	case 887:
		return "Yemen"
	case 891:
		return "Yugoslavia"
	case 894:
		return "Zambia"
	case 716:
		return "Zimbabwe"
	case 4:
		return "Afghanistan"
	case 688:
		return "Serbia"
	case 248:
		return "Aland Islands"
	case 535:
		return "Bonaire, Sint Eustatius And Saba"
	case 831:
		return "Guernsey"
	case 832:
		return "Jersey"
	case 531:
		return "Curacao"
	case 833:
		return "Isle Of Man"
	case 652:
		return "Saint Barthelemy"
	case 663:
		return "Saint Martin French"
	case 534:
		return "Sint Maarten Dutch"
	case 499:
		return "Montenegro"
	case 728:
		return "South Sudan"
	case 900:
		return "Kosovo"
	case 998:
		return "None"
	case 999:
		return "International"
	case 999800:
		return "International Freephone"
	case 999870:
		return "Inmarsat"
	case 999875:
		return "Maritime Mobile service"
	case 999878:
		return "Universal Personal Telecommunications services"
	case 999879:
		return "National non-commercial purposes"
	case 999881:
		return "Global Mobile Satellite System"
	case 999882:
		return "International Networks"
	case 999888:
		return "Disaster Relief"
	case 999979:
		return "International Premium Rate Service"
	case 999991:
		return "International Telecommunications Public Correspondence Service"
	}
	return UnknownMsg
}

// StringRus - returns a russian name of country
//nolint:gocyclo
func (c CountryCode) StringRus() string { //nolint:gocyclo
	switch c {
	case 8:
		return "Албания"
	case 12:
		return "Алжир"
	case 16:
		return "Американский Самоа"
	case 20:
		return "Андорра"
	case 24:
		return "Ангола"
	case 660:
		return "Ангилья"
	case 10:
		return "Антарктика"
	case 28:
		return "Антигуа и Барбуда"
	case 32:
		return "Аргентина"
	case 51:
		return "Армения"
	case 533:
		return "Аруба"
	case 36:
		return "Австралия"
	case 40:
		return "Австрия"
	case 31:
		return "Азербайджан"
	case 44:
		return "Багамские острова"
	case 48:
		return "Бахрейн"
	case 50:
		return "Бангладеш"
	case 52:
		return "Барбадос"
	case 112:
		return "Беларусь"
	case 56:
		return "Бельгия"
	case 84:
		return "Белиз"
	case 204:
		return "Бенин"
	case 60:
		return "Бермуды"
	case 64:
		return "Бутан"
	case 68:
		return "Боливия"
	case 70:
		return "Босния и Герцеговина"
	case 72:
		return "Ботсвана"
	case 74:
		return "остров Буве"
	case 76:
		return "Бразилия"
	case 86:
		return "Британские территории Индийского океана"
	case 96:
		return "Бруней"
	case 100:
		return "Болгария"
	case 854:
		return "Буркина Фасо"
	case 108:
		return "Бурунди"
	case 116:
		return "Камбоджа"
	case 120:
		return "Камерун"
	case 124:
		return "Канада"
	case 132:
		return "острова Зеленого Мыса"
	case 136:
		return "Каймановы острова"
	case 140:
		return "Центральная Африканская Республика"
	case 148:
		return "Чад"
	case 152:
		return "Чили"
	case 156:
		return "Китайская Народная Республика"
	case 162:
		return "остров Рождества"
	case 166:
		return "Кокосовые острова"
	case 170:
		return "Колумбия"
	case 174:
		return "Коморские острова"
	case 178:
		return "Конго"
	case 180:
		return "Демократическая республика Конго"
	case 184:
		return "острова Кука"
	case 188:
		return "Коста Рика"
	case 384:
		return "Кот-д`Ивуар"
	case 191:
		return "Хорватия"
	case 192:
		return "Куба"
	case 196:
		return "Кипр"
	case 203:
		return "Чехия"
	case 208:
		return "Дания"
	case 262:
		return "Джибути"
	case 212:
		return "Доминика"
	case 214:
		return "Доминиканская республика"
	case 218:
		return "Эквадор"
	case 818:
		return "Египет"
	case 222:
		return "Сальвадор"
	case 226:
		return "Экваториальная Гвинея"
	case 232:
		return "Эритрея"
	case 233:
		return "Эстония"
	case 231:
		return "Эфиопия"
	case 234:
		return "Фарерские острова"
	case 238:
		return "Фолклендские (Мальвинские) острова"
	case 242:
		return "Фиджи"
	case 246:
		return "Финляндия"
	case 250:
		return "Франция"
	case 254:
		return "Французская Гвиана"
	case 258:
		return "Французская Полинезия"
	case 260:
		return "Французские Южные Территории"
	case 266:
		return "Габон"
	case 270:
		return "Гамбия"
	case 268:
		return "Грузия"
	case 276:
		return "Германия"
	case 288:
		return "Гана"
	case 292:
		return "Гибралтар"
	case 300:
		return "Греция"
	case 304:
		return "Гренландия"
	case 308:
		return "Гренада"
	case 312:
		return "Гваделупа"
	case 316:
		return "Гуам"
	case 320:
		return "Гватемала"
	case 324:
		return "Гвинея"
	case 624:
		return "Гвинея-Бисау"
	case 328:
		return "Гайана"
	case 332:
		return "Гаити"
	case 334:
		return "острова Герда и МакДональда"
	case 340:
		return "Гондурас"
	case 344:
		return "Гонконг (Китай)"
	case 348:
		return "Венгрия"
	case 352:
		return "Исландия"
	case 356:
		return "Индия"
	case 360:
		return "Индонезия"
	case 364:
		return "Иран"
	case 368:
		return "Ирак"
	case 372:
		return "Ирландия"
	case 376:
		return "Израиль"
	case 380:
		return "Италия"
	case 388:
		return "Ямайка"
	case 392:
		return "Япония"
	case 400:
		return "Иордания"
	case 398:
		return "Казахстан"
	case 404:
		return "Кения"
	case 296:
		return "Кирибати"
	case 410:
		return "Корея"
	case 408:
		return "Корейская Народная Демократическая республика"
	case 414:
		return "Кувейт"
	case 417:
		return "Кыргызстан" // Киргизия
	case 418:
		return "Лаос"
	case 428:
		return "Латвия"
	case 422:
		return "Ливан"
	case 426:
		return "Лесото"
	case 430:
		return "Либерия"
	case 434:
		return "Ливия"
	case 438:
		return "Лихтенштейн"
	case 440:
		return "Литва"
	case 442:
		return "Люксембург"
	case 446:
		return "Макао (Китай)"
	case 807:
		return "Македония"
	case 450:
		return "Мадагаскар"
	case 454:
		return "Малави"
	case 458:
		return "Малайзия"
	case 462:
		return "Мальдивские острова"
	case 466:
		return "Мали"
	case 470:
		return "Мальта"
	case 584:
		return "Маршалловы острова"
	case 474:
		return "Мартиника"
	case 478:
		return "Мавритания"
	case 480:
		return "Маврикий"
	case 175:
		return "Майотта"
	case 484:
		return "Мексика"
	case 583:
		return "Микронезия"
	case 498:
		return "Молдова"
	case 492:
		return "Монако"
	case 496:
		return "Монголия"
	case 500:
		return "Монтсеррат"
	case 504:
		return "Марокко"
	case 508:
		return "Мозамбик"
	case 104:
		return "Мьянма"
	case 516:
		return "Намибия"
	case 520:
		return "Науру"
	case 524:
		return "Непал"
	case 528:
		return "Нидерланды"
	case 530:
		return "Антильские острова"
	case 540:
		return "Новая Каледония"
	case 554:
		return "Новая Зеландия"
	case 558:
		return "Никарагуа"
	case 562:
		return "Нигер"
	case 566:
		return "Нигерия"
	case 570:
		return "Ниуэ"
	case 574:
		return "остров Норфолк"
	case 580:
		return "Марианские острова"
	case 578:
		return "Норвегия"
	case 512:
		return "Оман"
	case 586:
		return "Пакистан"
	case 585:
		return "Палау"
	case 275:
		return "Палестина"
	case 591:
		return "Панама"
	case 598:
		return "Папуа - Новая Гвинея"
	case 600:
		return "Парагвай"
	case 604:
		return "Перу"
	case 608:
		return "Филиппины"
	case 612:
		return "остров Питкэрн"
	case 616:
		return "Польша"
	case 620:
		return "Португалия"
	case 630:
		return "Пуэрто-Рико"
	case 634:
		return "Катар"
	case 638:
		return "Реюньон"
	case 642:
		return "Румыния"
	case 643:
		return "Россия"
	case 646:
		return "Руанда"
	case 654:
		return "остров Святой Елены"
	case 659:
		return "Сент-Китс и Невис"
	case 662:
		return "Сент-Люсия"
	case 666:
		return "Сен-Пьер и Микелон"
	case 670:
		return "Сент-Винсент и Гренадины"
	case 882:
		return "острова Самоа"
	case 674:
		return "Сан-Марино"
	case 678:
		return "Сан-Томе и Принсипи"
	case 682:
		return "Саудовская Аравия"
	case 686:
		return "Сенегал"
	case 690:
		return "Сейшельские острова"
	case 694:
		return "Сьерра-Леоне"
	case 702:
		return "Сингапур"
	case 703:
		return "Словакия"
	case 705:
		return "Словения"
	case 90:
		return "Соломоновы острова"
	case 706:
		return "Сомали"
	case 710:
		return "ЮАР"
	case 239:
		return "Южная Георгия и Южные Сандвичевы острова"
	case 724:
		return "Испания"
	case 144:
		return "Шри Ланка"
	case 736:
		return "Судан"
	case 740:
		return "Суринам"
	case 744:
		return "острова Свалбард и Ян Майен"
	case 748:
		return "Свазиленд"
	case 752:
		return "Швеция"
	case 756:
		return "Швейцария"
	case 760:
		return "Сирия"
	case 158:
		return "Тайвань (Республика Китай)"
	case 762:
		return "Таджикистан"
	case 834:
		return "Танзания"
	case 764:
		return "Тайланд"
	case 626:
		return "Восточный Тимор"
	case 768:
		return "Того"
	case 772:
		return "Токелау"
	case 776:
		return "Тонга"
	case 780:
		return "Тринидад и Тобаго"
	case 788:
		return "Тунис"
	case 792:
		return "Турция"
	case 795:
		return "Туркменистан"
	case 796:
		return "острова Туркс и Кайкос"
	case 798:
		return "Тувалу"
	case 800:
		return "Уганда"
	case 804:
		return "Украина"
	case 784:
		return "Арабские Эмираты"
	case 826:
		return "Великобритания"
	case 840:
		return "Соединенные Штаты Америки"
	case 581:
		return "Отдаленные Острова США"
	case 858:
		return "Уругвай"
	case 860:
		return "Узбекистан"
	case 548:
		return "Вануату"
	case 336:
		return "Ватикан"
	case 862:
		return "Венесуэла"
	case 704:
		return "Вьетнам"
	case 92:
		return "Виргинские острова (Британские)"
	case 850:
		return "Виргинские острова (США)"
	case 876:
		return "острова Валлис и Футуна"
	case 732:
		return "Западная Сахара"
	case 887:
		return "Йемен"
	case 891:
		return "Югославия"
	case 894:
		return "Замбия"
	case 716:
		return "Зимбабве"
	case 4:
		return "Афганистан"
	case 688:
		return "Сербия"
	case 248:
		return "Аландские острова"
	case 535:
		return "Бонэйр, Синт-Эстатиус и Саба"
	case 831:
		return "Гернси"
	case 832:
		return "Джерси"
	case 531:
		return "Кюрасао"
	case 833:
		return "Остров Мэн"
	case 652:
		return "Сен-Бартелеми"
	case 663:
		return "Сен-Мартен"
	case 534:
		return "Синт-Мартен"
	case 499:
		return "Черногория"
	case 728:
		return "Южный Судан"
	case 900:
		return "Косово"
	case None:
		return "Отсутствует"
	case International:
		return "International"
	case 999800:
		return "Бесплатный номер"
	case 999870:
		return "Инмарсат"
	case 999875:
		return "Морская подвижная служба"
	case 999878:
		return "Универсальная персональная связь"
	case 999879:
		return "Некоммерческое использование"
	case 999881:
		return "Глобальная мобильная спутниковая система"
	case 999882:
		return "Глобальные телефонные номера"
	case 999888:
		return "Ликвидация последствий катастроф"
	case 999979:
		return "Международная услуга оплаты вызова по повышенному тарифу"
	case 999991:
		return "Служба международной телекоммуникационной корреспонденции"
	}
	return UnknownMsg
}

// Alpha2 - returns a default Alpha (Alpha-2/ISO2, 2 chars) code of country
//nolint:gocyclo
func (c CountryCode) Alpha2() string { //nolint:gocyclo
	switch c {
	case 8:
		return "AL"
	case 12:
		return "DZ"
	case 16:
		return "AS"
	case 20:
		return "AD"
	case 24:
		return "AO"
	case 660:
		return "AI"
	case 10:
		return "AQ"
	case 28:
		return "AG"
	case 32:
		return "AR"
	case 51:
		return "AM"
	case 533:
		return "AW"
	case 36:
		return "AU"
	case 40:
		return "AT"
	case 31:
		return "AZ"
	case 44:
		return "BS"
	case 48:
		return "BH"
	case 50:
		return "BD"
	case 52:
		return "BB"
	case 112:
		return "BY"
	case 56:
		return "BE"
	case 84:
		return "BZ"
	case 204:
		return "BJ"
	case 60:
		return "BM"
	case 64:
		return "BT"
	case 68:
		return "BO"
	case 70:
		return "BA"
	case 72:
		return "BW"
	case 74:
		return "BV"
	case 76:
		return "BR"
	case 86:
		return "IO"
	case 96:
		return "BN"
	case 100:
		return "BG"
	case 854:
		return "BF"
	case 108:
		return "BI"
	case 116:
		return "KH"
	case 120:
		return "CM"
	case 124:
		return "CA"
	case 132:
		return "CV"
	case 136:
		return "KY"
	case 140:
		return "CF"
	case 148:
		return "TD"
	case 152:
		return "CL"
	case 156:
		return "CN"
	case 162:
		return "CX"
	case 166:
		return "CC"
	case 170:
		return "CO"
	case 174:
		return "KM"
	case 178:
		return "CG"
	case 180:
		return "CD"
	case 184:
		return "CK"
	case 188:
		return "CR"
	case 384:
		return "CI"
	case 191:
		return "HR"
	case 192:
		return "CU"
	case 196:
		return "CY"
	case 203:
		return "CZ"
	case 208:
		return "DK"
	case 262:
		return "DJ"
	case 212:
		return "DM"
	case 214:
		return "DO"
	case 218:
		return "EC"
	case 818:
		return "EG"
	case 222:
		return "SV"
	case 226:
		return "GQ"
	case 232:
		return "ER"
	case 233:
		return "EE"
	case 231:
		return "ET"
	case 238:
		return "FK"
	case 242:
		return "FJ"
	case 246:
		return "FI"
	case 234:
		return "FO"
	case 250:
		return "FR"
	case 254:
		return "GF"
	case 258:
		return "PF"
	case 260:
		return "TF"
	case 266:
		return "GA"
	case 270:
		return "GM"
	case 268:
		return "GE"
	case 276:
		return "DE"
	case 288:
		return "GH"
	case 292:
		return "GI"
	case 300:
		return "GR"
	case 304:
		return "GL"
	case 308:
		return "GD"
	case 312:
		return "GP"
	case 316:
		return "GU"
	case 320:
		return "GT"
	case 324:
		return "GN"
	case 624:
		return "GW"
	case 328:
		return "GY"
	case 332:
		return "HT"
	case 334:
		return "HM"
	case 340:
		return "HN"
	case 344:
		return "HK"
	case 348:
		return "HU"
	case 352:
		return "IS"
	case 356:
		return "IN"
	case 360:
		return "ID"
	case 364:
		return "IR"
	case 368:
		return "IQ"
	case 372:
		return "IE"
	case 376:
		return "IL"
	case 380:
		return "IT"
	case 388:
		return "JM"
	case 392:
		return "JP"
	case 400:
		return "JO"
	case 398:
		return "KZ"
	case 404:
		return "KE"
	case 296:
		return "KI"
	case 410:
		return "KR"
	case 408:
		return "KP"
	case 414:
		return "KW"
	case 417:
		return "KG"
	case 418:
		return "LA"
	case 428:
		return "LV"
	case 422:
		return "LB"
	case 426:
		return "LS"
	case 430:
		return "LR"
	case 434:
		return "LY"
	case 438:
		return "LI"
	case 440:
		return "LT"
	case 442:
		return "LU"
	case 446:
		return "MO"
	case 807:
		return "MK"
	case 450:
		return "MG"
	case 454:
		return "MW"
	case 458:
		return "MY"
	case 462:
		return "MV"
	case 466:
		return "ML"
	case 470:
		return "MT"
	case 584:
		return "MH"
	case 474:
		return "MQ"
	case 478:
		return "MR"
	case 480:
		return "MU"
	case 175:
		return "YT"
	case 484:
		return "MX"
	case 583:
		return "FM"
	case 498:
		return "MD"
	case 492:
		return "MC"
	case 496:
		return "MN"
	case 500:
		return "MS"
	case 504:
		return "MA"
	case 508:
		return "MZ"
	case 104:
		return "MM"
	case 516:
		return "NA"
	case 520:
		return "NR"
	case 524:
		return "NP"
	case 528:
		return "NL"
	case 530:
		return "AN"
	case 540:
		return "NC"
	case 554:
		return "NZ"
	case 558:
		return "NI"
	case 562:
		return "NE"
	case 566:
		return "NG"
	case 570:
		return "NU"
	case 574:
		return "NF"
	case 580:
		return "MP"
	case 578:
		return "NO"
	case 512:
		return "OM"
	case 586:
		return "PK"
	case 585:
		return "PW"
	case 275:
		return "PS"
	case 591:
		return "PA"
	case 598:
		return "PG"
	case 600:
		return "PY"
	case 604:
		return "PE"
	case 608:
		return "PH"
	case 612:
		return "PN"
	case 616:
		return "PL"
	case 620:
		return "PT"
	case 630:
		return "PR"
	case 634:
		return "QA"
	case 638:
		return "RE"
	case 642:
		return "RO"
	case 643:
		return "RU"
	case 646:
		return "RW"
	case 654:
		return "SH"
	case 659:
		return "KN"
	case 662:
		return "LC"
	case 666:
		return "PM"
	case 670:
		return "VC"
	case 882:
		return "WS"
	case 674:
		return "SM"
	case 678:
		return "ST"
	case 682:
		return "SA"
	case 686:
		return "SN"
	case 690:
		return "SC"
	case 694:
		return "SL"
	case 702:
		return "SG"
	case 703:
		return "SK"
	case 705:
		return "SI"
	case 90:
		return "SB"
	case 706:
		return "SO"
	case 710:
		return "ZA"
	case 239:
		return "GS"
	case 724:
		return "ES"
	case 144:
		return "LK"
	case 736:
		return "SD"
	case 740:
		return "SR"
	case 744:
		return "SJ"
	case 748:
		return "SZ"
	case 752:
		return "SE"
	case 756:
		return "CH"
	case 760:
		return "SY"
	case 158:
		return "TW"
	case 762:
		return "TJ"
	case 834:
		return "TZ"
	case 764:
		return "TH"
	case 626:
		return "TL"
	case 768:
		return "TG"
	case 772:
		return "TK"
	case 776:
		return "TO"
	case 780:
		return "TT"
	case 788:
		return "TN"
	case 792:
		return "TR"
	case 795:
		return "TM"
	case 796:
		return "TC"
	case 798:
		return "TV"
	case 800:
		return "UG"
	case 804:
		return "UA"
	case 784:
		return "AE"
	case 826:
		return "GB"
	case 840:
		return "US"
	case 581:
		return "UM"
	case 858:
		return "UY"
	case 860:
		return "UZ"
	case 548:
		return "VU"
	case 336:
		return "VA"
	case 862:
		return "VE"
	case 704:
		return "VN"
	case 92:
		return "VG"
	case 850:
		return "VI"
	case 876:
		return "WF"
	case 732:
		return "EH"
	case 887:
		return "YE"
	case 891:
		return "YU"
	case 894:
		return "ZM"
	case 716:
		return "ZW"
	case 4:
		return "AF"
	case 688:
		return "RS"
	case 248:
		return "AX"
	case 535:
		return "BQ"
	case 831:
		return "GG"
	case 832:
		return "JE"
	case 531:
		return "CW"
	case 833:
		return "IM"
	case 652:
		return "BL"
	case 663:
		return "MF"
	case 534:
		return "SX"
	case 499:
		return "ME"
	case 728:
		return "SS"
	case 900:
		return "XK"
	case 998:
		return "None"
	case 999:
		return "International"
	case 999800:
		return "International Freephone"
	case 999870:
		return "Inmarsat"
	case 999875:
		return "Maritime Mobile service"
	case 999878:
		return "Universal Personal Telecommunications services"
	case 999879:
		return "National non-commercial purposes"
	case 999881:
		return "Global Mobile Satellite System"
	case 999882:
		return "International Networks"
	case 999888:
		return "Disaster Relief"
	case 999979:
		return "International Premium Rate Service"
	case 999991:
		return "International Telecommunications Public Correspondence Service"
	}
	return UnknownMsg
}

// Alpha3 - returns a Alpha-3 (ISO3, 3 chars) code of country
//nolint:gocyclo
func (c CountryCode) Alpha3() string { //nolint:gocyclo
	switch c {
	case 8:
		return "ALB"
	case 12:
		return "DZA"
	case 16:
		return "ASM"
	case 20:
		return "AND"
	case 24:
		return "AGO"
	case 660:
		return "AIA"
	case 10:
		return "ATA"
	case 28:
		return "ATG"
	case 32:
		return "ARG"
	case 51:
		return "ARM"
	case 533:
		return "ABW"
	case 36:
		return "AUS"
	case 40:
		return "AUT"
	case 31:
		return "AZE"
	case 44:
		return "BHS"
	case 48:
		return "BHR"
	case 50:
		return "BGD"
	case 52:
		return "BRB"
	case 112:
		return "BLR"
	case 56:
		return "BEL"
	case 84:
		return "BLZ"
	case 204:
		return "BEN"
	case 60:
		return "BMU"
	case 64:
		return "BTN"
	case 68:
		return "BOL"
	case 70:
		return "BIH"
	case 72:
		return "BWA"
	case 74:
		return "BVT"
	case 76:
		return "BRA"
	case 86:
		return "IOT"
	case 96:
		return "BRN"
	case 100:
		return "BGR"
	case 854:
		return "BFA"
	case 108:
		return "BDI"
	case 116:
		return "KHM"
	case 120:
		return "CMR"
	case 124:
		return "CAN"
	case 132:
		return "CPV"
	case 136:
		return "CYM"
	case 140:
		return "CAF"
	case 148:
		return "TCD"
	case 152:
		return "CHL"
	case 156:
		return "CHN"
	case 162:
		return "CXR"
	case 166:
		return "CCK"
	case 170:
		return "COL"
	case 174:
		return "COM"
	case 178:
		return "COG"
	case 180:
		return "COD"
	case 184:
		return "COK"
	case 188:
		return "CRI"
	case 384:
		return "CIV"
	case 191:
		return "HRV"
	case 192:
		return "CUB"
	case 196:
		return "CYP"
	case 203:
		return "CZE"
	case 208:
		return "DNK"
	case 262:
		return "DJI"
	case 212:
		return "DMA"
	case 214:
		return "DOM"
	case 218:
		return "ECU"
	case 818:
		return "EGY"
	case 222:
		return "SLV"
	case 226:
		return "GNQ"
	case 232:
		return "ERI"
	case 233:
		return "EST"
	case 231:
		return "ETH"
	case 238:
		return "FLK"
	case 242:
		return "FJI"
	case 246:
		return "FIN"
	case 250:
		return "FRA"
	case 234:
		return "FRO"
	case 254:
		return "GUF"
	case 258:
		return "PYF"
	case 260:
		return "ATF"
	case 266:
		return "GAB"
	case 270:
		return "GMB"
	case 268:
		return "GEO"
	case 276:
		return "DEU"
	case 288:
		return "GHA"
	case 292:
		return "GIB"
	case 300:
		return "GRC"
	case 304:
		return "GRL"
	case 308:
		return "GRD"
	case 312:
		return "GLP"
	case 316:
		return "GUM"
	case 320:
		return "GTM"
	case 324:
		return "GIN"
	case 624:
		return "GNB"
	case 328:
		return "GUY"
	case 332:
		return "HTI"
	case 334:
		return "HMD"
	case 340:
		return "HND"
	case 344:
		return "HKG"
	case 348:
		return "HUN"
	case 352:
		return "ISL"
	case 356:
		return "IND"
	case 360:
		return "IDN"
	case 364:
		return "IRN"
	case 368:
		return "IRQ"
	case 372:
		return "IRL"
	case 376:
		return "ISR"
	case 380:
		return "ITA"
	case 388:
		return "JAM"
	case 392:
		return "JPN"
	case 400:
		return "JOR"
	case 398:
		return "KAZ"
	case 404:
		return "KEN"
	case 296:
		return "KIR"
	case 410:
		return "KOR"
	case 408:
		return "PRK"
	case 414:
		return "KWT"
	case 417:
		return "KGZ"
	case 418:
		return "LAO"
	case 428:
		return "LVA"
	case 422:
		return "LBN"
	case 426:
		return "LSO"
	case 430:
		return "LBR"
	case 434:
		return "LBY"
	case 438:
		return "LIE"
	case 440:
		return "LTU"
	case 442:
		return "LUX"
	case 446:
		return "MAC"
	case 807:
		return "MKD"
	case 450:
		return "MDG"
	case 454:
		return "MWI"
	case 458:
		return "MYS"
	case 462:
		return "MDV"
	case 466:
		return "MLI"
	case 470:
		return "MLT"
	case 584:
		return "MHL"
	case 474:
		return "MTQ"
	case 478:
		return "MRT"
	case 480:
		return "MUS"
	case 175:
		return "MYT"
	case 484:
		return "MEX"
	case 583:
		return "FSM"
	case 498:
		return "MDA"
	case 492:
		return "MCO"
	case 496:
		return "MNG"
	case 500:
		return "MSR"
	case 504:
		return "MAR"
	case 508:
		return "MOZ"
	case 104:
		return "MMR"
	case 516:
		return "NAM"
	case 520:
		return "NRU"
	case 524:
		return "NPL"
	case 528:
		return "NLD"
	case 530:
		return "ANT"
	case 540:
		return "NCL"
	case 554:
		return "NZL"
	case 558:
		return "NIC"
	case 562:
		return "NER"
	case 566:
		return "NGA"
	case 570:
		return "NIU"
	case 574:
		return "NFK"
	case 580:
		return "MNP"
	case 578:
		return "NOR"
	case 512:
		return "OMN"
	case 586:
		return "PAK"
	case 585:
		return "PLW"
	case 275:
		return "PSE"
	case 591:
		return "PAN"
	case 598:
		return "PNG"
	case 600:
		return "PRY"
	case 604:
		return "PER"
	case 608:
		return "PHL"
	case 612:
		return "PCN"
	case 616:
		return "POL"
	case 620:
		return "PRT"
	case 630:
		return "PRI"
	case 634:
		return "QAT"
	case 638:
		return "REU"
	case 642:
		return "ROU"
	case 643:
		return "RUS"
	case 646:
		return "RWA"
	case 654:
		return "SHN"
	case 659:
		return "KNA"
	case 662:
		return "LCA"
	case 666:
		return "SPM"
	case 670:
		return "VCT"
	case 882:
		return "WSM"
	case 674:
		return "SMR"
	case 678:
		return "STP"
	case 682:
		return "SAU"
	case 686:
		return "SEN"
	case 690:
		return "SYC"
	case 694:
		return "SLE"
	case 702:
		return "SGP"
	case 703:
		return "SVK"
	case 705:
		return "SVN"
	case 90:
		return "SLB"
	case 706:
		return "SOM"
	case 710:
		return "ZAF"
	case 239:
		return "SGS"
	case 724:
		return "ESP"
	case 144:
		return "LKA"
	case 736:
		return "SDN"
	case 740:
		return "SUR"
	case 744:
		return "SJM"
	case 748:
		return "SWZ"
	case 752:
		return "SWE"
	case 756:
		return "CHE"
	case 760:
		return "SYR"
	case 158:
		return "TWN"
	case 762:
		return "TJK"
	case 834:
		return "TZA"
	case 764:
		return "THA"
	case 626:
		return "TLS"
	case 768:
		return "TGO"
	case 772:
		return "TKL"
	case 776:
		return "TON"
	case 780:
		return "TTO"
	case 788:
		return "TUN"
	case 792:
		return "TUR"
	case 795:
		return "TKM"
	case 796:
		return "TCA"
	case 798:
		return "TUV"
	case 800:
		return "UGA"
	case 804:
		return "UKR"
	case 784:
		return "ARE"
	case 826:
		return "GBR"
	case 840:
		return "USA"
	case 581:
		return "UMI"
	case 858:
		return "URY"
	case 860:
		return "UZB"
	case 548:
		return "VUT"
	case 336:
		return "VAT"
	case 862:
		return "VEN"
	case 704:
		return "VNM"
	case 92:
		return "VGB"
	case 850:
		return "VIR"
	case 876:
		return "WLF"
	case 732:
		return "ESH"
	case 887:
		return "YEM"
	case 891:
		return "YUG"
	case 894:
		return "ZMB"
	case 716:
		return "ZWE"
	case 4:
		return "AFG"
	case 688:
		return "SRB"
	case 248:
		return "ALA"
	case 535:
		return "BES"
	case 831:
		return "GGY"
	case 832:
		return "JEY"
	case 531:
		return "CUW"
	case 833:
		return "IMN"
	case 652:
		return "BLM"
	case 663:
		return "MAF"
	case 534:
		return "SXM"
	case 499:
		return "MNE"
	case 728:
		return "SSD"
	case 900:
		return "XKX"
	case 998:
		return "None"
	case 999:
		return "International"
	case 999800:
		return "International Freephone"
	case 999870:
		return "Inmarsat"
	case 999875:
		return "Maritime Mobile service"
	case 999878:
		return "Universal Personal Telecommunications services"
	case 999879:
		return "National non-commercial purposes"
	case 999881:
		return "Global Mobile Satellite System"
	case 999882:
		return "International Networks"
	case 999888:
		return "Disaster Relief"
	case 999979:
		return "International Premium Rate Service"
	case 999991:
		return "International Telecommunications Public Correspondence Service"
	}
	return UnknownMsg
}

// FIFA - returns a FIFA (AFC, CAF, CONCACAF, CONMEBOL, OFC and UEFA) three-letter country code
func (c CountryCode) FIFA() string {
	switch c {
	case GBR:
		return "ENG"
	case VNM:
		return "VIE"
	case DEU:
		return "GER"
	case NLD:
		return "NED"
	case PSE:
		return "PLE"
	case TWN:
		return "TPE"
	case HRV:
		return "CRO"
	case CAF:
		return "CTA"
	}
	return c.Alpha3()
}

// IOC - returns The International Olympic Committee (IOC) three-letter abbreviation country codes
//nolint:gocyclo
func (c CountryCode) IOC() string { //nolint:gocyclo
	switch c {
	case DZA:
		return `ALG`
	case ASM:
		return `ASA`
	case VIR:
		return `ISV`
	case AGO:
		return `ANG`
	case ATG:
		return `ANT`
	case GNQ:
		return `GEQ`
	case ABW:
		return `ARU`
	case BHS:
		return `BAH`
	case BHR:
		return `BRN`
	case BGD:
		return `BAN`
	case BRB:
		return `BAR`
	case BLZ:
		return `BIZ`
	case BMU:
		return `BER`
	case BTN:
		return `BHU`
	case BWA:
		return `BOT`
	case VGB:
		return `IVB`
	case BRN:
		return `BRU`
	case BGR:
		return `BUL`
	case BFA:
		return `BUR`
	case CHL:
		return `CHI`
	case CRI:
		return `CRC`
	case DNK:
		return `DEN`
	case DEU:
		return `GER`
	case SLV:
		return `ESA`
	case FJI:
		return `FIJ`
	case GMB:
		return `GAM`
	case GRD:
		return `GRN`
	case GRC:
		return `GRE`
	case GTM:
		return `GUA`
	case GIN:
		return `GUI`
	case GNB:
		return `GBS`
	case HTI:
		return `HAI`
	case HND:
		return `HON`
	case IDN:
		return `INA`
	case IRN:
		return `IRI`
	case CYM:
		return `CAY`
	case KHM:
		return `CAM`
	case COG:
		return `CGO`
	case XKX:
		return `KOS`
	case NON:
		return `NONE`
	case HRV:
		return `CRO`
	case KWT:
		return `KUW`
	case LSO:
		return `LES`
	case LVA:
		return `LAT`
	case LBN:
		return `LIB`
	case LBY:
		return `LBA`
	case MDG:
		return `MAD`
	case MWI:
		return `MAW`
	case MYS:
		return `MAS`
	case MRT:
		return `MTN`
	case MUS:
		return `MRI`
	case MCO:
		return `MON`
	case MNG:
		return `MGL`
	case MMR:
		return `MYA`
	case NPL:
		return `NEP`
	case NIC:
		return `NCA`
	case NLD:
		return `NED`
	case ANT:
		return `AHO`
	case NER:
		return `NIG`
	case NGA:
		return `NGR`
	case OMN:
		return `OMA`
	case PSE:
		return `PLE`
	case PRY:
		return `PAR`
	case PHL:
		return `PHI`
	case PRT:
		return `POR`
	case PRI:
		return `PUR`
	case TWN:
		return `TPE`
	case ROU:
		return `ROU`
	case SLB:
		return `SOL`
	case ZMB:
		return `ZAM`
	case WSM:
		return `SAM`
	case SAU:
		return `KSA`
	case CHE:
		return `SUI`
	case SYC:
		return `SEY`
	case ZWE:
		return `ZIM`
	case SVN:
		return `SLO`
	case LKA:
		return `SRI`
	case KNA:
		return `SKN`
	case VCT:
		return `VIN`
	case ZAF:
		return `RSA`
	case SDN:
		return `SUD`
	case TZA:
		return `TAN`
	case TGO:
		return `TOG`
	case TON:
		return `TGA`
	case TTO:
		return `TRI`
	case TCD:
		return `CHA`
	case URY:
		return `URU`
	case VUT:
		return `VAN`
	case ARE:
		return `UAE`
	case VNM:
		return `VIE`
	}
	return c.Alpha3()
}

// Currency - returns a currency of the country
//nolint:gocyclo
func (c CountryCode) Currency() CurrencyCode { //nolint:gocyclo
	switch c {
	case AUS:
		return CurrencyAUD
	case AUT, AND, MAF:
		return CurrencyEUR
	case AZE:
		return CurrencyAZN
	case ALB:
		return CurrencyALL
	case DZA:
		return CurrencyDZD
	case ASM, BES:
		return CurrencyUSD
	case AIA:
		return CurrencyXCD
	case AGO:
		return CurrencyAOA
	case ATG:
		return CurrencyXCD
	case ANT, CUW:
		return CurrencyANG
	case ARE:
		return CurrencyAED
	case ARG:
		return CurrencyARS
	case ARM:
		return CurrencyAMD
	case ABW:
		return CurrencyAWG
	case AFG:
		return CurrencyAFN
	case BHS:
		return CurrencyBSD
	case BGD:
		return CurrencyBDT
	case BRB:
		return CurrencyBBD
	case BHR:
		return CurrencyBHD
	case BLR:
		return CurrencyBYR
	case BLZ:
		return CurrencyBZD
	case BEL:
		return CurrencyEUR
	case BEN:
		return CurrencyXOF
	case BMU:
		return CurrencyBMD
	case BGR:
		return CurrencyBGN
	case BOL:
		return CurrencyBOB
	case BIH:
		return CurrencyBAM
	case BWA:
		return CurrencyBWP
	case BRA:
		return CurrencyBRL
	case IOT:
		return CurrencyUSD
	case BRN:
		return CurrencyBND
	case BFA:
		return CurrencyXOF
	case BDI:
		return CurrencyBIF
	case BTN:
		return CurrencyBTN
	case VUT:
		return CurrencyVUV
	case VAT:
		return CurrencyEUR
	case GBR, GGY, JEY, IMN:
		return CurrencyGBP
	case HUN:
		return CurrencyHUF
	case VEN:
		return CurrencyVEF
	case VGB:
		return CurrencyUSD
	case VIR:
		return CurrencyUSD
	case TLS:
		return CurrencyUSD
	case VNM:
		return CurrencyVND
	case GAB:
		return CurrencyXAF
	case HTI:
		return CurrencyHTG
	case GUY:
		return CurrencyGYD
	case GMB:
		return CurrencyGMD
	case GHA:
		return CurrencyGHS
	case GLP:
		return CurrencyEUR
	case GTM:
		return CurrencyGTQ
	case GIN:
		return CurrencyGNF
	case GNB:
		return CurrencyXOF
	case DEU:
		return CurrencyEUR
	case GIB:
		return CurrencyGIP
	case HND:
		return CurrencyHNL
	case HKG:
		return CurrencyHKD
	case GRD:
		return CurrencyXCD
	case GRL:
		return CurrencyDKK
	case GRC:
		return CurrencyEUR
	case GEO:
		return CurrencyGEL
	case GUM:
		return CurrencyUSD
	case DNK:
		return CurrencyDKK
	case COD:
		return CurrencyCDF
	case DJI:
		return CurrencyDJF
	case DMA:
		return CurrencyXCD
	case DOM:
		return CurrencyDOP
	case EGY:
		return CurrencyEGP
	case ZMB:
		return CurrencyZMW
	case ESH:
		return CurrencyMAD
	case ZWE:
		return CurrencyZWL
	case ISR:
		return CurrencyILS
	case IND:
		return CurrencyINR
	case IDN:
		return CurrencyIDR
	case JOR:
		return CurrencyJOD
	case IRQ:
		return CurrencyIQD
	case IRN:
		return CurrencyIRR
	case IRL:
		return CurrencyEUR
	case ISL:
		return CurrencyISK
	case ESP:
		return CurrencyEUR
	case ITA:
		return CurrencyEUR
	case YEM:
		return CurrencyYER
	case KAZ:
		return CurrencyKZT
	case CYM:
		return CurrencyKYD
	case KHM:
		return CurrencyKHR
	case CMR:
		return CurrencyXAF
	case CAN:
		return CurrencyCAD
	case QAT:
		return CurrencyQAR
	case KEN:
		return CurrencyKES
	case CYP:
		return CurrencyEUR
	case KIR:
		return CurrencyAUD
	case CHN:
		return CurrencyCNY
	case CCK:
		return CurrencyAUD
	case COL:
		return CurrencyCOP
	case COM:
		return CurrencyKMF
	case COG:
		return CurrencyXAF
	case PRK:
		return CurrencyKPW
	case KOR:
		return CurrencyKRW
	case CRI:
		return CurrencyCRC
	case CIV:
		return CurrencyXOF
	case CUB:
		return CurrencyCUC
	case KWT:
		return CurrencyKWD
	case KGZ:
		return CurrencyKGS
	case LAO:
		return CurrencyLAK
	case LVA:
		return CurrencyEUR
	case LSO:
		return CurrencyLSL
	case LBR:
		return CurrencyLRD
	case LBN:
		return CurrencyLBP
	case LBY:
		return CurrencyLYD
	case LTU:
		return CurrencyEUR
	case LIE:
		return CurrencyCHF
	case LUX:
		return CurrencyEUR
	case MUS:
		return CurrencyMUR
	case MRT:
		return CurrencyMRU
	case MDG:
		return CurrencyMGA
	case MYT:
		return CurrencyEUR
	case MAC:
		return CurrencyMOP
	case MKD:
		return CurrencyMKD
	case MWI:
		return CurrencyMWK
	case MYS:
		return CurrencyMYR
	case MLI:
		return CurrencyXOF
	case MDV:
		return CurrencyMVR
	case MLT:
		return CurrencyEUR
	case MNP:
		return CurrencyUSD
	case MAR:
		return CurrencyMAD
	case MTQ:
		return CurrencyEUR
	case MHL:
		return CurrencyUSD
	case MEX:
		return CurrencyMXN
	case FSM:
		return CurrencyUSD
	case MOZ:
		return CurrencyMZN
	case MDA:
		return CurrencyMDL
	case MCO:
		return CurrencyEUR
	case MNG:
		return CurrencyMNT
	case MSR:
		return CurrencyXCD
	case MMR:
		return CurrencyMMK
	case NAM:
		return CurrencyNAD
	case NRU:
		return CurrencyAUD
	case NPL:
		return CurrencyNPR
	case NER:
		return CurrencyXOF
	case NGA:
		return CurrencyNGN
	case NLD:
		return CurrencyEUR
	case NIC:
		return CurrencyNIO
	case NIU:
		return CurrencyNZD
	case NZL:
		return CurrencyNZD
	case NCL:
		return CurrencyXPF
	case NOR:
		return CurrencyNOK
	case ChannelIslands:
		return CurrencyEUR
	case OMN:
		return CurrencyOMR
	case BVT:
		return CurrencyNOK
	case NFK:
		return CurrencyAUD
	case PCN:
		return CurrencyNZD
	case CXR:
		return CurrencyAUD
	case SHN:
		return CurrencySHP
	case WLF:
		return CurrencyXPF
	case HMD:
		return CurrencyAUD
	case CPV:
		return CurrencyCVE
	case COK:
		return CurrencyNZD
	case WSM:
		return CurrencyWST
	case SJM:
		return CurrencyNOK
	case TCA:
		return CurrencyUSD
	case UMI:
		return CurrencyUSD
	case PAK:
		return CurrencyPKR
	case PLW:
		return CurrencyUSD
	case PSE:
		return CurrencyILS
	case PAN:
		return CurrencyPAB
	case PNG:
		return CurrencyPGK
	case PRY:
		return CurrencyPYG
	case PER:
		return CurrencyPEN
	case POL:
		return CurrencyPLN
	case PRT:
		return CurrencyEUR
	case PRI:
		return CurrencyUSD
	case REU:
		return CurrencyEUR
	case RUS:
		return CurrencyRUB
	case RWA:
		return CurrencyRWF
	case ROU:
		return CurrencyRON
	case SLV:
		return CurrencySVC
	case SMR:
		return CurrencyEUR
	case STP:
		return CurrencySTN
	case SAU:
		return CurrencySAR
	case SWZ:
		return CurrencySZL
	case SYC:
		return CurrencySCR
	case SEN:
		return CurrencyXOF
	case SPM:
		return CurrencyEUR
	case VCT:
		return CurrencyXCD
	case KNA:
		return CurrencyXCD
	case LCA:
		return CurrencyXCD
	case SGP:
		return CurrencySGD
	case SYR:
		return CurrencySYP
	case SVK:
		return CurrencyEUR
	case SVN:
		return CurrencyEUR
	case USA:
		return CurrencyUSD
	case SLB:
		return CurrencySBD
	case SOM:
		return CurrencySOS
	case SDN:
		return CurrencySDG
	case SUR:
		return CurrencySRD
	case SLE:
		return CurrencySLL
	case TJK:
		return CurrencyTJS
	case TWN:
		return CurrencyTWD
	case THA:
		return CurrencyTHB
	case TZA:
		return CurrencyTZS
	case TGO:
		return CurrencyXOF
	case TKL:
		return CurrencyNZD
	case TON:
		return CurrencyTOP
	case TTO:
		return CurrencyTTD
	case TUV:
		return CurrencyAUD
	case TUN:
		return CurrencyTND
	case TKM:
		return CurrencyTMT
	case TUR:
		return CurrencyTRY
	case UGA:
		return CurrencyUGX
	case UZB:
		return CurrencyUZS
	case UKR:
		return CurrencyUAH
	case URY:
		return CurrencyUYI
	case FRO:
		return CurrencyDKK
	case FJI:
		return CurrencyFJD
	case PHL:
		return CurrencyPHP
	case FIN:
		return CurrencyEUR
	case FLK:
		return CurrencyFKP
	case FRA:
		return CurrencyEUR
	case GUF:
		return CurrencyEUR
	case PYF:
		return CurrencyXPF
	case ATF:
		return CurrencyEUR
	case HRV:
		return CurrencyHRK
	case CAF:
		return CurrencyXAF
	case TCD:
		return CurrencyXAF
	case CZE:
		return CurrencyCZK
	case CHL:
		return CurrencyCLF
	case CHE:
		return CurrencyCHF
	case SWE:
		return CurrencySEK
	case LKA:
		return CurrencyLKR
	case ECU:
		return CurrencyUSD
	case GNQ:
		return CurrencyXAF
	case ERI:
		return CurrencyERN
	case EST:
		return CurrencyEUR
	case ETH:
		return CurrencyETB
	case ZAF:
		return CurrencyZAR
	case YUG:
		return CurrencyYUD
	case SGS:
		return CurrencyGBP
	case JAM:
		return CurrencyJMD
	case JPN:
		return CurrencyJPY
	case BLM, MNE, ALA:
		return CurrencyEUR
	case SXM:
		return CurrencyANG
	case SRB:
		return CurrencyRSD
	case SSD:
		return CurrencySSP
	case XKX:
		return CurrencyEUR
	case NON, International, NonCountryInternationalFreephone, NonCountryInmarsat, NonCountryMaritimeMobileService,
		NonCountryUniversalPersonalTelecommunicationsServices, NonCountryNationalNonCommercialPurposes, NonCountryGlobalMobileSatelliteSystem,
		NonCountryInternationalNetworks, NonCountryDisasterRelief, NonCountryInternationalPremiumRateService,
		NonCountryInternationalTelecommunicationsCorrespondenceService:
		return CurrencyNone
	}
	return CurrencyUnknown
}

// All - return all country codes
func All() []CountryCode {
	return []CountryCode{
		AUS,
		AUT,
		AZE,
		ALB,
		DZA,
		ASM,
		AIA,
		AGO,
		AND,
		ATA,
		ATG,
		ANT,
		ARE,
		ARG,
		ARM,
		ABW,
		AFG,
		BHS,
		BGD,
		BRB,
		BHR,
		BLR,
		BLZ,
		BEL,
		BEN,
		BMU,
		BGR,
		BOL,
		BIH,
		BWA,
		BRA,
		IOT,
		BRN,
		BFA,
		BDI,
		BTN,
		VUT,
		VAT,
		GBR,
		HUN,
		VEN,
		VGB,
		VIR,
		TLS,
		VNM,
		GAB,
		HTI,
		GUY,
		GMB,
		GHA,
		GLP,
		GTM,
		GIN,
		GNB,
		DEU,
		GIB,
		HND,
		HKG,
		GRD,
		GRL,
		GRC,
		GEO,
		GUM,
		DNK,
		COD,
		DJI,
		DMA,
		DOM,
		EGY,
		ZMB,
		ESH,
		ZWE,
		ISR,
		IND,
		IDN,
		JOR,
		IRQ,
		IRN,
		IRL,
		ISL,
		ESP,
		ITA,
		YEM,
		KAZ,
		CYM,
		KHM,
		CMR,
		CAN,
		QAT,
		KEN,
		CYP,
		KIR,
		CHN,
		CCK,
		COL,
		COM,
		COG,
		PRK,
		KOR,
		CRI,
		CIV,
		CUB,
		KWT,
		KGZ,
		LAO,
		LVA,
		LSO,
		LBR,
		LBN,
		LBY,
		LTU,
		LIE,
		LUX,
		MUS,
		MRT,
		MDG,
		MYT,
		MAC,
		MKD,
		MWI,
		MYS,
		MLI,
		MDV,
		MLT,
		MNP,
		MAR,
		MTQ,
		MHL,
		MEX,
		FSM,
		MOZ,
		MDA,
		MCO,
		MNG,
		MSR,
		MMR,
		NAM,
		NRU,
		NPL,
		NER,
		NGA,
		NLD,
		NIC,
		NIU,
		NZL,
		NCL,
		NOR,
		OMN,
		BVT,
		IMN,
		NFK,
		PCN,
		CXR,
		SHN,
		WLF,
		HMD,
		CPV,
		COK,
		WSM,
		SJM,
		TCA,
		UMI,
		PAK,
		PLW,
		PSE,
		PAN,
		PNG,
		PRY,
		PER,
		POL,
		PRT,
		PRI,
		REU,
		RUS,
		RWA,
		ROU,
		SLV,
		SMR,
		STP,
		SAU,
		SWZ,
		SYC,
		SEN,
		SPM,
		VCT,
		KNA,
		LCA,
		SGP,
		SYR,
		SVK,
		SVN,
		USA,
		SLB,
		SOM,
		SDN,
		SUR,
		SLE,
		TJK,
		TWN,
		THA,
		TZA,
		TGO,
		TKL,
		TON,
		TTO,
		TUV,
		TUN,
		TKM,
		TUR,
		UGA,
		UZB,
		UKR,
		URY,
		// XWA, // ignored for All(), part of GB
		FRO,
		FJI,
		PHL,
		FIN,
		FLK,
		FRA,
		GUF,
		PYF,
		ATF,
		HRV,
		CAF,
		TCD,
		CZE,
		CHL,
		CHE,
		SWE,
		// XSC, // ignored for All(), part of GB
		LKA,
		ECU,
		GNQ,
		ERI,
		EST,
		ETH,
		ZAF,
		YUG,
		SGS,
		JAM,
		MNE,
		BLM,
		SXM,
		SRB,
		ALA,
		BES,
		GGY,
		JEY,
		CUW,
		MAF,
		SSD,
		JPN,
		XKX,
	}
}

// AllNonCountries - return all non-country codes
func AllNonCountries() []CountryCode {
	return []CountryCode{
		NonCountryInternationalFreephone,
		NonCountryInmarsat,
		NonCountryMaritimeMobileService,
		NonCountryUniversalPersonalTelecommunicationsServices,
		NonCountryNationalNonCommercialPurposes,
		NonCountryGlobalMobileSatelliteSystem,
		NonCountryInternationalNetworks,
		NonCountryDisasterRelief,
		NonCountryInternationalPremiumRateService,
		NonCountryInternationalTelecommunicationsCorrespondenceService,
	}
}

// CallCodes - return calling code of country
//nolint:gocyclo
func (c CountryCode) CallCodes() []CallCode { //nolint:gocyclo
	switch c {
	case AUS:
		return []CallCode{CallCode(61)}
	case AUT:
		return []CallCode{CallCode(43)}
	case AZE:
		return []CallCode{CallCode(994)}
	case ALB:
		return []CallCode{CallCode(355)}
	case DZA:
		return []CallCode{CallCode(213)}
	case ASM:
		return []CallCode{CallCode(1684)}
	case AIA:
		return []CallCode{CallCode(1264)}
	case AGO:
		return []CallCode{CallCode(244)}
	case AND:
		return []CallCode{CallCode(376)}
	case ATA:
		return []CallCode{CallCode(672)}
	case ATG:
		return []CallCode{CallCode(1268)}
	case ANT:
		return []CallCode{CallCode(599)}
	case ARE:
		return []CallCode{CallCode(971)}
	case ARG:
		return []CallCode{CallCode(54)}
	case ARM:
		return []CallCode{CallCode(374)}
	case ABW:
		return []CallCode{CallCode(297), CallCode(5998)}
	case AFG:
		return []CallCode{CallCode(93)}
	case BHS:
		return []CallCode{CallCode(1242)}
	case BGD:
		return []CallCode{CallCode(880)}
	case BRB:
		return []CallCode{CallCode(1246)}
	case BHR:
		return []CallCode{CallCode(973)}
	case BLR:
		return []CallCode{CallCode(375)}
	case BLZ:
		return []CallCode{CallCode(501)}
	case BEL:
		return []CallCode{CallCode(32)}
	case BEN:
		return []CallCode{CallCode(229)}
	case BMU:
		return []CallCode{CallCode(1441)}
	case BGR:
		return []CallCode{CallCode(359)}
	case BOL:
		return []CallCode{CallCode(591)}
	case BIH:
		return []CallCode{CallCode(387)}
	case BWA:
		return []CallCode{CallCode(267)}
	case BRA:
		return []CallCode{CallCode(55)}
	case IOT:
		return []CallCode{CallCode(246)}
	case BRN:
		return []CallCode{CallCode(673)}
	case BFA:
		return []CallCode{CallCode(226)}
	case BDI:
		return []CallCode{CallCode(257)}
	case BTN:
		return []CallCode{CallCode(975)}
	case VUT:
		return []CallCode{CallCode(678)}
	case VAT:
		return []CallCode{CallCode(3906698)}
	case GBR:
		return []CallCode{CallCode(44)}
	case HUN:
		return []CallCode{CallCode(36)}
	case VEN:
		return []CallCode{CallCode(58)}
	case VGB:
		return []CallCode{CallCode(1284)}
	case VIR:
		return []CallCode{CallCode(1340)}
	case TLS:
		return []CallCode{CallCode(670)}
	case VNM:
		return []CallCode{CallCode(84)}
	case GAB:
		return []CallCode{CallCode(241)}
	case HTI:
		return []CallCode{CallCode(509)}
	case GUY:
		return []CallCode{CallCode(592)}
	case GMB:
		return []CallCode{CallCode(220)}
	case GHA:
		return []CallCode{CallCode(233)}
	case GLP:
		return []CallCode{CallCode(590)}
	case GTM:
		return []CallCode{CallCode(502)}
	case GIN:
		return []CallCode{CallCode(224)}
	case GNB:
		return []CallCode{CallCode(245)}
	case DEU:
		return []CallCode{CallCode(49)}
	case GIB:
		return []CallCode{CallCode(350)}
	case HND:
		return []CallCode{CallCode(504)}
	case HKG:
		return []CallCode{CallCode(852)}
	case GRD:
		return []CallCode{CallCode(1473)}
	case GRL:
		return []CallCode{CallCode(299)}
	case GRC:
		return []CallCode{CallCode(30)}
	case GEO:
		return []CallCode{CallCode(995)}
	case GUM:
		return []CallCode{CallCode(1671)}
	case DNK:
		return []CallCode{CallCode(45)}
	case COD:
		return []CallCode{CallCode(243)}
	case DJI:
		return []CallCode{CallCode(253)}
	case DMA:
		return []CallCode{CallCode(1767)}
	case DOM:
		return []CallCode{CallCode(1809), CallCode(1829), CallCode(1849)}
	case EGY:
		return []CallCode{CallCode(20)}
	case ZMB:
		return []CallCode{CallCode(260)}
	case ESH:
		return []CallCode{CallCode(212)}
	case ZWE:
		return []CallCode{CallCode(263)}
	case ISR:
		return []CallCode{CallCode(972)}
	case IND:
		return []CallCode{CallCode(91)}
	case IDN:
		return []CallCode{CallCode(62)}
	case JOR:
		return []CallCode{CallCode(962)}
	case IRQ:
		return []CallCode{CallCode(964)}
	case IRN:
		return []CallCode{CallCode(98)}
	case IRL:
		return []CallCode{CallCode(353)}
	case ISL:
		return []CallCode{CallCode(354)}
	case ESP:
		return []CallCode{CallCode(34)}
	case ITA:
		return []CallCode{CallCode(39)}
	case YEM:
		return []CallCode{CallCode(967)}
	case KAZ:
		return []CallCode{CallCode(7)}
	case CYM:
		return []CallCode{CallCode(1345)}
	case KHM:
		return []CallCode{CallCode(855)}
	case CMR:
		return []CallCode{CallCode(237)}
	case CAN:
		return []CallCode{CallCode(1)}
	case QAT:
		return []CallCode{CallCode(974)}
	case KEN:
		return []CallCode{CallCode(254)}
	case CYP:
		return []CallCode{CallCode(357)}
	case KIR:
		return []CallCode{CallCode(686)}
	case CHN:
		return []CallCode{CallCode(86)}
	case CCK:
		return []CallCode{CallCode(672), CallCode(6189162)}
	case COL:
		return []CallCode{CallCode(57)}
	case COM:
		return []CallCode{CallCode(269)}
	case COG:
		return []CallCode{CallCode(242)}
	case PRK:
		return []CallCode{CallCode(850)}
	case KOR:
		return []CallCode{CallCode(82)}
	case CRI:
		return []CallCode{CallCode(506)}
	case CIV:
		return []CallCode{CallCode(225)}
	case CUB:
		return []CallCode{CallCode(53)}
	case KWT:
		return []CallCode{CallCode(965)}
	case KGZ:
		return []CallCode{CallCode(996)}
	case LAO:
		return []CallCode{CallCode(856)}
	case LVA:
		return []CallCode{CallCode(371)}
	case LSO:
		return []CallCode{CallCode(266)}
	case LBR:
		return []CallCode{CallCode(231)}
	case LBN:
		return []CallCode{CallCode(961)}
	case LBY:
		return []CallCode{CallCode(218)}
	case LTU:
		return []CallCode{CallCode(370)}
	case LIE:
		return []CallCode{CallCode(423)}
	case LUX:
		return []CallCode{CallCode(352)}
	case MUS:
		return []CallCode{CallCode(230)}
	case MRT:
		return []CallCode{CallCode(222)}
	case MDG:
		return []CallCode{CallCode(261)}
	case MYT:
		return []CallCode{CallCode(CallCode262269), CallCode(CallCode262639)}
	case MAC:
		return []CallCode{CallCode(853)}
	case MKD:
		return []CallCode{CallCode(389)}
	case MWI:
		return []CallCode{CallCode(265)}
	case MYS:
		return []CallCode{CallCode(60)}
	case MLI:
		return []CallCode{CallCode(223)}
	case MDV:
		return []CallCode{CallCode(960)}
	case MLT:
		return []CallCode{CallCode(356)}
	case MNP:
		return []CallCode{CallCode(1670)}
	case MAR:
		return []CallCode{CallCode(212)}
	case MTQ:
		return []CallCode{CallCode(596)}
	case MHL:
		return []CallCode{CallCode(692)}
	case MEX:
		return []CallCode{CallCode(52)}
	case FSM:
		return []CallCode{CallCode(691)}
	case MOZ:
		return []CallCode{CallCode(258)}
	case MDA:
		return []CallCode{CallCode(373)}
	case MCO:
		return []CallCode{CallCode(377)}
	case MNG:
		return []CallCode{CallCode(976)}
	case MSR:
		return []CallCode{CallCode(1664)}
	case MMR:
		return []CallCode{CallCode(95)}
	case NAM:
		return []CallCode{CallCode(264)}
	case NRU:
		return []CallCode{CallCode(674)}
	case NPL:
		return []CallCode{CallCode(977)}
	case NER:
		return []CallCode{CallCode(227)}
	case NGA:
		return []CallCode{CallCode(234)}
	case NLD:
		return []CallCode{CallCode(31)}
	case NIC:
		return []CallCode{CallCode(505)}
	case NIU:
		return []CallCode{CallCode(683)}
	case NZL:
		return []CallCode{CallCode(64)}
	case NCL:
		return []CallCode{CallCode(687)}
	case NOR:
		return []CallCode{CallCode(47)}
	case OMN:
		return []CallCode{CallCode(968)}
	case BVT:
		return []CallCode{CallCode(47)}
	case IMN:
		return []CallCode{CallCode(441624)}
	case NFK:
		return []CallCode{CallCode(672)}
	case PCN:
		return []CallCode{CallCode(64)}
	case CXR:
		return []CallCode{CallCode(6189164)}
	case SHN:
		return []CallCode{CallCode(290)}
	case WLF:
		return []CallCode{CallCode(681)}
	case HMD:
		return []CallCode{CallCode(61)}
	case CPV:
		return []CallCode{CallCode(238)}
	case COK:
		return []CallCode{CallCode(682)}
	case WSM:
		return []CallCode{CallCode(685)}
	case SJM:
		return []CallCode{CallCode(4779)}
	case TCA:
		return []CallCode{CallCode(1649)}
	case UMI:
		return []CallCode{CallCode(1)}
	case PAK:
		return []CallCode{CallCode(92)}
	case PLW:
		return []CallCode{CallCode(680)}
	case PSE:
		return []CallCode{CallCode(970)}
	case PAN:
		return []CallCode{CallCode(507)}
	case PNG:
		return []CallCode{CallCode(675)}
	case PRY:
		return []CallCode{CallCode(595)}
	case PER:
		return []CallCode{CallCode(51)}
	case POL:
		return []CallCode{CallCode(48)}
	case PRT:
		return []CallCode{CallCode(351)}
	case PRI:
		return []CallCode{CallCode(1787), CallCode(1939)}
	case REU:
		return []CallCode{CallCode(262)}
	case RUS:
		return []CallCode{CallCode(7)}
	case RWA:
		return []CallCode{CallCode(250)}
	case ROU:
		return []CallCode{CallCode(40)}
	case SLV:
		return []CallCode{CallCode(503)}
	case SMR:
		return []CallCode{CallCode(378)}
	case STP:
		return []CallCode{CallCode(239)}
	case SAU:
		return []CallCode{CallCode(966)}
	case SWZ:
		return []CallCode{CallCode(268)}
	case SYC:
		return []CallCode{CallCode(248)}
	case SEN:
		return []CallCode{CallCode(221)}
	case SPM:
		return []CallCode{CallCode(508)}
	case VCT:
		return []CallCode{CallCode(1784)}
	case KNA:
		return []CallCode{CallCode(1869)}
	case LCA:
		return []CallCode{CallCode(1758)}
	case SGP:
		return []CallCode{CallCode(65)}
	case SYR:
		return []CallCode{CallCode(963)}
	case SVK:
		return []CallCode{CallCode(421)}
	case SVN:
		return []CallCode{CallCode(386)}
	case USA:
		return []CallCode{CallCode(1)}
	case SLB:
		return []CallCode{CallCode(677)}
	case SOM:
		return []CallCode{CallCode(252)}
	case SDN:
		return []CallCode{CallCode(249)}
	case SUR:
		return []CallCode{CallCode(597)}
	case SLE:
		return []CallCode{CallCode(232)}
	case TJK:
		return []CallCode{CallCode(992)}
	case TWN:
		return []CallCode{CallCode(886)}
	case THA:
		return []CallCode{CallCode(66)}
	case TZA:
		return []CallCode{CallCode(255)}
	case TGO:
		return []CallCode{CallCode(228)}
	case TKL:
		return []CallCode{CallCode(690)}
	case TON:
		return []CallCode{CallCode(676)}
	case TTO:
		return []CallCode{CallCode(1868)}
	case TUV:
		return []CallCode{CallCode(688)}
	case TUN:
		return []CallCode{CallCode(216)}
	case TKM:
		return []CallCode{CallCode(993)}
	case TUR:
		return []CallCode{CallCode(90)}
	case UGA:
		return []CallCode{CallCode(256)}
	case UZB:
		return []CallCode{CallCode(998)}
	case UKR:
		return []CallCode{CallCode(380)}
	case URY:
		return []CallCode{CallCode(598)}
	case FRO:
		return []CallCode{CallCode(298)}
	case FJI:
		return []CallCode{CallCode(679)}
	case PHL:
		return []CallCode{CallCode(63)}
	case FIN:
		return []CallCode{CallCode(358)}
	case FLK:
		return []CallCode{CallCode(500)}
	case FRA:
		return []CallCode{CallCode(33)}
	case GUF:
		return []CallCode{CallCode(594)}
	case PYF:
		return []CallCode{CallCode(689)}
	case ATF:
		return []CallCode{CallCode(1)}
	case HRV:
		return []CallCode{CallCode(385)}
	case CAF:
		return []CallCode{CallCode(236)}
	case TCD:
		return []CallCode{CallCode(235)}
	case CZE:
		return []CallCode{CallCode(420)}
	case CHL:
		return []CallCode{CallCode(56)}
	case CHE:
		return []CallCode{CallCode(41)}
	case SWE:
		return []CallCode{CallCode(46)}
	case LKA:
		return []CallCode{CallCode(94)}
	case ECU:
		return []CallCode{CallCode(593)}
	case GNQ:
		return []CallCode{CallCode(240)}
	case ERI:
		return []CallCode{CallCode(291)}
	case EST:
		return []CallCode{CallCode(372)}
	case ETH:
		return []CallCode{CallCode(251)}
	case ZAF:
		return []CallCode{CallCode(27)}
	case YUG:
		return []CallCode{CallCode(38)}
	case SGS:
		return []CallCode{CallCode(500)}
	case JAM:
		return []CallCode{CallCode(1876), CallCode(1658)}
	case MNE:
		return []CallCode{CallCode(382)}
	case BLM:
		return []CallCode{CallCode(590)}
	case SXM:
		return []CallCode{CallCode(1721)}
	case SRB:
		return []CallCode{CallCode(381)}
	case ALA:
		return []CallCode{CallCode(35818)}
	case BES:
		return []CallCode{CallCode(5993), CallCode(5994)}
	case GGY:
		return []CallCode{CallCode(441481)}
	case JEY:
		return []CallCode{CallCode(441534)}
	case CUW:
		return []CallCode{CallCode(5999)}
	case MAF:
		return []CallCode{CallCode(590)}
	case SSD:
		return []CallCode{CallCode(211)}
	case XKX:
		return []CallCode{CallCode(383)}
	case NON:
		return []CallCode{}
	case International:
		return []CallCode{CallCode(800), CallCode(870), CallCode(875), CallCode(876), CallCode(877), CallCode(878), CallCode(879), CallCode(881),
			CallCode(882), CallCode(883), CallCode(888), CallCode(979), CallCode(991)}
	case JPN:
		return []CallCode{CallCode(81)}
	case NonCountryInternationalFreephone:
		return []CallCode{CallCode(800)}
	case NonCountryInmarsat:
		return []CallCode{CallCode(870)}
	case NonCountryMaritimeMobileService:
		return []CallCode{CallCode(875), CallCode(876), CallCode(877)}
	case NonCountryUniversalPersonalTelecommunicationsServices:
		return []CallCode{CallCode(878)}
	case NonCountryNationalNonCommercialPurposes:
		return []CallCode{CallCode(879)}
	case NonCountryGlobalMobileSatelliteSystem:
		return []CallCode{CallCode(881)}
	case NonCountryInternationalNetworks:
		return []CallCode{CallCode(882), CallCode(883)}
	case NonCountryDisasterRelief:
		return []CallCode{CallCode(888)}
	case NonCountryInternationalPremiumRateService:
		return []CallCode{CallCode(979)}
	case NonCountryInternationalTelecommunicationsCorrespondenceService:
		return []CallCode{CallCode(991)}
	}
	return []CallCode{CallCode(0)}
}

// Domain - return domain code of country
func (c CountryCode) Domain() DomainCode {
	domain := DomainCode(c)
	if domain.IsValid() {
		return domain
	}
	return DomainUnknown
}

// Region - return Region code ot the country
//nolint:gocyclo
func (c CountryCode) Region() RegionCode { //nolint:gocyclo
	switch c {
	case AUS:
		return RegionOC
	case AUT:
		return RegionEU
	case AZE:
		return RegionAS
	case ALB:
		return RegionEU
	case DZA:
		return RegionAF
	case ASM:
		return RegionOC
	case AIA:
		return RegionNA
	case AGO:
		return RegionAF
	case AND:
		return RegionEU
	case ATA:
		return RegionAN
	case ATG:
		return RegionNA
	case ANT:
		return RegionNA
	case ARE:
		return RegionAS
	case ARG:
		return RegionSA
	case ARM:
		return RegionAS
	case ABW:
		return RegionNA
	case AFG:
		return RegionAS
	case BHS:
		return RegionNA
	case BGD:
		return RegionAS
	case BRB:
		return RegionNA
	case BHR:
		return RegionAS
	case BLR:
		return RegionEU
	case BLZ:
		return RegionNA
	case BEL:
		return RegionEU
	case BEN:
		return RegionAF
	case BMU:
		return RegionNA
	case BGR:
		return RegionEU
	case BOL:
		return RegionSA
	case BIH:
		return RegionEU
	case BWA:
		return RegionAF
	case BRA:
		return RegionSA
	case IOT:
		return RegionAS
	case BRN:
		return RegionAS
	case BFA:
		return RegionAF
	case BDI:
		return RegionAF
	case BTN:
		return RegionAS
	case VUT:
		return RegionOC
	case VAT:
		return RegionEU
	case GBR:
		return RegionEU
	case HUN:
		return RegionEU
	case VEN:
		return RegionSA
	case VGB:
		return RegionNA
	case VIR:
		return RegionNA
	case TLS:
		return RegionAS
	case VNM:
		return RegionAS
	case GAB:
		return RegionAF
	case HTI:
		return RegionNA
	case GUY:
		return RegionSA
	case GMB:
		return RegionAF
	case GHA:
		return RegionAF
	case GLP:
		return RegionNA
	case GTM:
		return RegionNA
	case GIN:
		return RegionAF
	case GNB:
		return RegionAF
	case DEU:
		return RegionEU
	case GIB:
		return RegionEU
	case HND:
		return RegionNA
	case HKG:
		return RegionAS
	case GRD:
		return RegionNA
	case GRL:
		return RegionNA
	case GRC:
		return RegionEU
	case GEO:
		return RegionAS
	case GUM:
		return RegionOC
	case DNK:
		return RegionEU
	case COD:
		return RegionAF
	case DJI:
		return RegionAF
	case DMA:
		return RegionNA
	case DOM:
		return RegionNA
	case EGY:
		return RegionAF
	case ZMB:
		return RegionAF
	case ESH:
		return RegionAF
	case ZWE:
		return RegionAF
	case ISR:
		return RegionAS
	case IND:
		return RegionAS
	case IDN:
		return RegionAS
	case JOR:
		return RegionAS
	case IRQ:
		return RegionAS
	case IRN:
		return RegionAS
	case IRL:
		return RegionEU
	case ISL:
		return RegionEU
	case ESP:
		return RegionEU
	case ITA:
		return RegionEU
	case YEM:
		return RegionAS
	case KAZ:
		return RegionAS
	case CYM:
		return RegionNA
	case KHM:
		return RegionAS
	case CMR:
		return RegionAF
	case CAN:
		return RegionNA
	case QAT:
		return RegionAS
	case KEN:
		return RegionAF
	case CYP:
		return RegionAS
	case KIR:
		return RegionOC
	case CHN:
		return RegionAS
	case CCK:
		return RegionAS
	case COL:
		return RegionSA
	case COM:
		return RegionAF
	case COG:
		return RegionAF
	case PRK:
		return RegionAS
	case KOR:
		return RegionAS
	case CRI:
		return RegionNA
	case CIV:
		return RegionAF
	case CUB:
		return RegionNA
	case KWT:
		return RegionAS
	case KGZ:
		return RegionAS
	case LAO:
		return RegionAS
	case LVA:
		return RegionEU
	case LSO:
		return RegionAF
	case LBR:
		return RegionAF
	case LBN:
		return RegionAS
	case LBY:
		return RegionAF
	case LTU:
		return RegionEU
	case LIE:
		return RegionEU
	case LUX:
		return RegionEU
	case MUS:
		return RegionAF
	case MRT:
		return RegionAF
	case MDG:
		return RegionAF
	case MYT:
		return RegionAF
	case MAC:
		return RegionAS
	case MKD:
		return RegionEU
	case MWI:
		return RegionAF
	case MYS:
		return RegionAS
	case MLI:
		return RegionAF
	case MDV:
		return RegionAS
	case MLT:
		return RegionEU
	case MNP:
		return RegionOC
	case MAR:
		return RegionAF
	case MTQ:
		return RegionNA
	case MHL:
		return RegionOC
	case MEX:
		return RegionNA
	case FSM:
		return RegionOC
	case MOZ:
		return RegionAF
	case MDA:
		return RegionEU
	case MCO:
		return RegionEU
	case MNG:
		return RegionAS
	case MSR:
		return RegionNA
	case MMR:
		return RegionAS
	case NAM:
		return RegionAF
	case NRU:
		return RegionOC
	case NPL:
		return RegionAS
	case NER:
		return RegionAF
	case NGA:
		return RegionAF
	case NLD:
		return RegionEU
	case NIC:
		return RegionNA
	case NIU:
		return RegionOC
	case NZL:
		return RegionOC
	case NCL:
		return RegionOC
	case NOR:
		return RegionEU
	case OMN:
		return RegionAS
	case BVT:
		return RegionAN
	case IMN:
		return RegionEU
	case NFK:
		return RegionOC
	case PCN:
		return RegionOC
	case CXR:
		return RegionAS
	case SHN:
		return RegionAF
	case WLF:
		return RegionOC
	case HMD:
		return RegionAN
	case CPV:
		return RegionAF
	case COK:
		return RegionOC
	case WSM:
		return RegionOC
	case SJM:
		return RegionEU
	case TCA:
		return RegionNA
	case UMI:
		return RegionOC
	case PAK:
		return RegionAS
	case PLW:
		return RegionOC
	case PSE:
		return RegionAS
	case PAN:
		return RegionNA
	case PNG:
		return RegionOC
	case PRY:
		return RegionSA
	case PER:
		return RegionSA
	case POL:
		return RegionEU
	case PRT:
		return RegionEU
	case PRI:
		return RegionNA
	case REU:
		return RegionAF
	case RUS:
		return RegionEU
	case RWA:
		return RegionAF
	case ROU:
		return RegionEU
	case SLV:
		return RegionNA
	case SMR:
		return RegionEU
	case STP:
		return RegionAF
	case SAU:
		return RegionAS
	case SWZ:
		return RegionAF
	case SYC:
		return RegionAF
	case SEN:
		return RegionAF
	case SPM:
		return RegionNA
	case VCT:
		return RegionNA
	case KNA:
		return RegionNA
	case LCA:
		return RegionNA
	case SGP:
		return RegionAS
	case SYR:
		return RegionAS
	case SVK:
		return RegionEU
	case SVN:
		return RegionEU
	case USA:
		return RegionNA
	case SLB:
		return RegionOC
	case SOM:
		return RegionAF
	case SDN:
		return RegionAF
	case SUR:
		return RegionSA
	case SLE:
		return RegionAF
	case TJK:
		return RegionAS
	case TWN:
		return RegionAS
	case THA:
		return RegionAS
	case TZA:
		return RegionAF
	case TGO:
		return RegionAF
	case TKL:
		return RegionOC
	case TON:
		return RegionOC
	case TTO:
		return RegionNA
	case TUV:
		return RegionOC
	case TUN:
		return RegionAF
	case TKM:
		return RegionAS
	case TUR:
		return RegionEU
	case UGA:
		return RegionAF
	case UZB:
		return RegionAS
	case UKR:
		return RegionEU
	case URY:
		return RegionSA
	case FRO:
		return RegionEU
	case FJI:
		return RegionOC
	case PHL:
		return RegionAS
	case FIN:
		return RegionEU
	case FLK:
		return RegionSA
	case FRA:
		return RegionEU
	case GUF:
		return RegionSA
	case PYF:
		return RegionOC
	case ATF:
		return RegionAN
	case HRV:
		return RegionEU
	case CAF:
		return RegionAF
	case TCD:
		return RegionAF
	case CZE:
		return RegionEU
	case CHL:
		return RegionSA
	case CHE:
		return RegionEU
	case SWE:
		return RegionEU
	case LKA:
		return RegionAS
	case ECU:
		return RegionSA
	case GNQ:
		return RegionAF
	case ERI:
		return RegionAF
	case EST:
		return RegionEU
	case ETH:
		return RegionAF
	case ZAF:
		return RegionAF
	case YUG:
		return RegionEU
	case SGS:
		return RegionAN
	case JAM:
		return RegionNA
	case MNE:
		return RegionEU
	case BLM:
		return RegionNA
	case SXM:
		return RegionNA
	case SRB:
		return RegionEU
	case ALA:
		return RegionEU
	case BES:
		return RegionNA
	case GGY:
		return RegionEU
	case JEY:
		return RegionEU
	case CUW:
		return RegionOC
	case MAF:
		return RegionNA
	case SSD:
		return RegionAF
	case XKX:
		return RegionEU
	case NON, International, NonCountryInternationalFreephone, NonCountryInmarsat, NonCountryMaritimeMobileService,
		NonCountryUniversalPersonalTelecommunicationsServices, NonCountryNationalNonCommercialPurposes, NonCountryGlobalMobileSatelliteSystem,
		NonCountryInternationalNetworks, NonCountryDisasterRelief, NonCountryInternationalPremiumRateService,
		NonCountryInternationalTelecommunicationsCorrespondenceService:
		return RegionNone
	case JPN:
		return RegionAS
	}
	return RegionUnknown
}

// Capital - return a capital of country
//nolint:gocyclo
func (c CountryCode) Capital() CapitalCode { //nolint:gocyclo
	switch c {
	case AUS:
		return CapitalAU
	case AUT:
		return CapitalAT
	case AZE:
		return CapitalAZ
	case ALB:
		return CapitalAL
	case DZA:
		return CapitalDZ
	case ASM:
		return CapitalAS
	case AIA:
		return CapitalAI
	case AGO:
		return CapitalAO
	case AND:
		return CapitalAD
	case ATA:
		return CapitalAQ
	case ATG:
		return CapitalAG
	case ANT:
		return CapitalAN
	case ARE:
		return CapitalAE
	case ARG:
		return CapitalAR
	case ARM:
		return CapitalAM
	case ABW:
		return CapitalAW
	case AFG:
		return CapitalAF
	case BHS:
		return CapitalBS
	case BGD:
		return CapitalBD
	case BRB:
		return CapitalBB
	case BHR:
		return CapitalBH
	case BLR:
		return CapitalBY
	case BLZ:
		return CapitalBZ
	case BEL:
		return CapitalBE
	case BEN:
		return CapitalBJ
	case BMU:
		return CapitalBM
	case BGR:
		return CapitalBG
	case BOL:
		return CapitalBO
	case BIH:
		return CapitalBA
	case BWA:
		return CapitalBW
	case BRA:
		return CapitalBR
	case IOT:
		return CapitalIO
	case BRN:
		return CapitalBN
	case BFA:
		return CapitalBF
	case BDI:
		return CapitalBI
	case BTN:
		return CapitalBT
	case VUT:
		return CapitalVU
	case VAT:
		return CapitalVA
	case GBR:
		return CapitalGB
	case HUN:
		return CapitalHU
	case VEN:
		return CapitalVE
	case VGB:
		return CapitalVG
	case VIR:
		return CapitalVI
	case TLS:
		return CapitalTL
	case VNM:
		return CapitalVN
	case GAB:
		return CapitalGA
	case HTI:
		return CapitalHT
	case GUY:
		return CapitalGY
	case GMB:
		return CapitalGM
	case GHA:
		return CapitalGH
	case GLP:
		return CapitalGP
	case GTM:
		return CapitalGT
	case GIN:
		return CapitalGN
	case GNB:
		return CapitalGW
	case DEU:
		return CapitalDE
	case GIB:
		return CapitalGI
	case HND:
		return CapitalHN
	case HKG:
		return CapitalHK
	case GRD:
		return CapitalGD
	case GRL:
		return CapitalGL
	case GRC:
		return CapitalGR
	case GEO:
		return CapitalGE
	case GUM:
		return CapitalGU
	case DNK:
		return CapitalDK
	case COD:
		return CapitalCD
	case DJI:
		return CapitalDJ
	case DMA:
		return CapitalDM
	case DOM:
		return CapitalDO
	case EGY:
		return CapitalEG
	case ZMB:
		return CapitalZM
	case ESH:
		return CapitalEH
	case ZWE:
		return CapitalZW
	case ISR:
		return CapitalIL
	case IND:
		return CapitalIN
	case IDN:
		return CapitalID
	case JOR:
		return CapitalJO
	case IRQ:
		return CapitalIQ
	case IRN:
		return CapitalIR
	case IRL:
		return CapitalIE
	case ISL:
		return CapitalIS
	case ESP:
		return CapitalES
	case ITA:
		return CapitalIT
	case YEM:
		return CapitalYE
	case KAZ:
		return CapitalKZ
	case CYM:
		return CapitalKY
	case KHM:
		return CapitalKH
	case CMR:
		return CapitalCM
	case CAN:
		return CapitalCA
	case QAT:
		return CapitalQA
	case KEN:
		return CapitalKE
	case CYP:
		return CapitalCY
	case KIR:
		return CapitalKI
	case CHN:
		return CapitalCN
	case CCK:
		return CapitalCC
	case COL:
		return CapitalCO
	case COM:
		return CapitalKM
	case COG:
		return CapitalCG
	case PRK:
		return CapitalKP
	case KOR:
		return CapitalKR
	case CRI:
		return CapitalCR
	case CIV:
		return CapitalCI
	case CUB:
		return CapitalCU
	case KWT:
		return CapitalKW
	case KGZ:
		return CapitalKG
	case LAO:
		return CapitalLA
	case LVA:
		return CapitalLV
	case LSO:
		return CapitalLS
	case LBR:
		return CapitalLR
	case LBN:
		return CapitalLB
	case LBY:
		return CapitalLY
	case LTU:
		return CapitalLT
	case LIE:
		return CapitalLI
	case LUX:
		return CapitalLU
	case MUS:
		return CapitalMU
	case MRT:
		return CapitalMR
	case MDG:
		return CapitalMG
	case MYT:
		return CapitalYT
	case MAC:
		return CapitalMO
	case MKD:
		return CapitalMK
	case MWI:
		return CapitalMW
	case MYS:
		return CapitalMY
	case MLI:
		return CapitalML
	case MDV:
		return CapitalMV
	case MLT:
		return CapitalMT
	case MNP:
		return CapitalMP
	case MAR:
		return CapitalMA
	case MTQ:
		return CapitalMQ
	case MHL:
		return CapitalMH
	case MEX:
		return CapitalMX
	case FSM:
		return CapitalFM
	case MOZ:
		return CapitalMZ
	case MDA:
		return CapitalMD
	case MCO:
		return CapitalMC
	case MNG:
		return CapitalMN
	case MSR:
		return CapitalMS
	case MMR:
		return CapitalMM
	case NAM:
		return CapitalNA
	case NRU:
		return CapitalNR
	case NPL:
		return CapitalNP
	case NER:
		return CapitalNE
	case NGA:
		return CapitalNG
	case NLD:
		return CapitalNL
	case NIC:
		return CapitalNI
	case NIU:
		return CapitalNU
	case NZL:
		return CapitalNZ
	case NCL:
		return CapitalNC
	case NOR:
		return CapitalNO
	case OMN:
		return CapitalOM
	case BVT:
		return CapitalBV
	case IMN:
		return CapitalIM
	case NFK:
		return CapitalNF
	case PCN:
		return CapitalPN
	case CXR:
		return CapitalCX
	case SHN:
		return CapitalSH
	case WLF:
		return CapitalWF
	case HMD:
		return CapitalHM
	case CPV:
		return CapitalCV
	case COK:
		return CapitalCK
	case WSM:
		return CapitalWS
	case SJM:
		return CapitalSJ
	case TCA:
		return CapitalTC
	case UMI:
		return CapitalUM
	case PAK:
		return CapitalPK
	case PLW:
		return CapitalPW
	case PSE:
		return CapitalPS
	case PAN:
		return CapitalPA
	case PNG:
		return CapitalPG
	case PRY:
		return CapitalPY
	case PER:
		return CapitalPE
	case POL:
		return CapitalPL
	case PRT:
		return CapitalPT
	case PRI:
		return CapitalPR
	case REU:
		return CapitalRE
	case RUS:
		return CapitalRU
	case RWA:
		return CapitalRW
	case ROU:
		return CapitalRO
	case SLV:
		return CapitalSV
	case SMR:
		return CapitalSM
	case STP:
		return CapitalST
	case SAU:
		return CapitalSA
	case SWZ:
		return CapitalSZ
	case SYC:
		return CapitalSC
	case SEN:
		return CapitalSN
	case SPM:
		return CapitalPM
	case VCT:
		return CapitalVC
	case KNA:
		return CapitalKN
	case LCA:
		return CapitalLC
	case SGP:
		return CapitalSG
	case SYR:
		return CapitalSY
	case SVK:
		return CapitalSK
	case SVN:
		return CapitalSI
	case USA:
		return CapitalUS
	case SLB:
		return CapitalSB
	case SOM:
		return CapitalSO
	case SDN:
		return CapitalSD
	case SUR:
		return CapitalSR
	case SLE:
		return CapitalSL
	case TJK:
		return CapitalTJ
	case TWN:
		return CapitalTW
	case THA:
		return CapitalTH
	case TZA:
		return CapitalTZ
	case TGO:
		return CapitalTG
	case TKL:
		return CapitalTK
	case TON:
		return CapitalTO
	case TTO:
		return CapitalTT
	case TUV:
		return CapitalTV
	case TUN:
		return CapitalTN
	case TKM:
		return CapitalTM
	case TUR:
		return CapitalTR
	case UGA:
		return CapitalUG
	case UZB:
		return CapitalUZ
	case UKR:
		return CapitalUA
	case URY:
		return CapitalUY
	case FRO:
		return CapitalFO
	case FJI:
		return CapitalFJ
	case PHL:
		return CapitalPH
	case FIN:
		return CapitalFI
	case FLK:
		return CapitalFK
	case FRA:
		return CapitalFR
	case GUF:
		return CapitalGF
	case PYF:
		return CapitalPF
	case ATF:
		return CapitalTF
	case HRV:
		return CapitalHR
	case CAF:
		return CapitalCF
	case TCD:
		return CapitalTD
	case CZE:
		return CapitalCZ
	case CHL:
		return CapitalCL
	case CHE:
		return CapitalCH
	case SWE:
		return CapitalSE
	case LKA:
		return CapitalLK
	case ECU:
		return CapitalEC
	case GNQ:
		return CapitalGQ
	case ERI:
		return CapitalER
	case EST:
		return CapitalEE
	case ETH:
		return CapitalET
	case ZAF:
		return CapitalZA
	case YUG:
		return CapitalYU
	case SGS:
		return CapitalGS
	case JAM:
		return CapitalJM
	case MNE:
		return CapitalME
	case BLM:
		return CapitalBL
	case SXM:
		return CapitalSX
	case SRB:
		return CapitalRS
	case ALA:
		return CapitalAX
	case BES:
		return CapitalBQ
	case GGY:
		return CapitalGG
	case JEY:
		return CapitalJE
	case CUW:
		return CapitalCW
	case MAF:
		return CapitalMF
	case SSD:
		return CapitalSS
	case XKX:
		return CapitalXK
	case JPN:
		return CapitalJP
	case NON, International, NonCountryInternationalFreephone, NonCountryInmarsat, NonCountryMaritimeMobileService,
		NonCountryUniversalPersonalTelecommunicationsServices, NonCountryNationalNonCommercialPurposes, NonCountryGlobalMobileSatelliteSystem,
		NonCountryInternationalNetworks, NonCountryDisasterRelief, NonCountryInternationalPremiumRateService,
		NonCountryInternationalTelecommunicationsCorrespondenceService:
		return CapitalXX
	}
	return CapitalUnknown
}

// Info - return all info about country as Country struct
func (c CountryCode) Info() *Country {
	return &Country{
		Name:      c.String(),
		Alpha2:    c.Alpha2(),
		Alpha3:    c.Alpha3(),
		IOC:       c.IOC(),
		FIFA:      c.FIFA(),
		Emoji:     c.Emoji(),
		Code:      c,
		Capital:   c.Capital(),
		Currency:  c.Currency(),
		CallCodes: c.CallCodes(),
		Domain:    c.Domain(),
		Region:    c.Region(),
	}
}

// Type implements Typer interface.
func (country *Country) Type() string {
	return TypeCountry
}

// Value implements database/sql/driver.Valuer
func (country Country) Value() (Value, error) {
	return json.Marshal(country)
}

// Scan implements database/sql.Scanner
func (country *Country) Scan(src interface{}) error {
	if country == nil {
		return fmt.Errorf("countries::Scan: Country scan err: country == nil")
	}
	switch src := src.(type) {
	case *Country:
		*country = *src
	case Country:
		*country = src
	default:
		return fmt.Errorf("countries::Scan: Country scan err: unexpected value of type %T for %T", src, *country)
	}
	return nil
}

//
// AllInfo - return all currencies as []Currency
func AllInfo() []*Country {
	all := All()
	countries := make([]*Country, 0, len(all))
	for _, v := range all {
		countries = append(countries, v.Info())
	}
	return countries
}

// ByName - return CountryCode by country Alpha-2 / Alpha-3 / name, case-insensitive, example: rus := ByName("Ru") OR rus := ByName("russia"),
// returns countries.Unknown, if country name not found or not valid
//nolint:misspell,gocyclo
func ByName(name string) CountryCode { //nolint:misspell,gocyclo
	switch textPrepare(name) {
	case "AU", "AUS", "AUSTRALIA", "AVSTRALIA", "AVSTRALIYA", "AUSTRALIYA", "AUSTRALIEN":
		return AUS
	case "AT", "AUT", "AUSTRIA", "AVSTRIA", "AUSTRIYA", "AVSTRIYA", "ÖSTERREICH", "OESTERREICH":
		return AUT
	case "AZ", "AZE", "AZERBAIJAN", "AYZERBAIJAN", "AZERBAIDJAN", "AYZERBAIDJAN", "ASERBAIDSCHAN":
		return AZE
	case "AL", "ALB", "ALBANIA", "ALBANIYA", "ALBANIEN":
		return ALB
	case "DZ", "DZA", "ALGERIA", "ALGERIYA", "ALGERIEN":
		return DZA
	case "AS", "ASM", "AMERICANSAMOA", "AMERICASAMOA", "SAMOAAMERICAN", "SAMOAMERICAN", "SAMOAMERICA", "AMERIKANISCHSAMOA":
		return ASM
	case "AI", "AIA", "ANGUILLA", "ANGUILA":
		return AIA
	case "XEN", "ENG", "ENGLAND", "INGLAND":
		return GBR
	case "AO", "AGO", "ANGOLA", "ANGOLIA":
		return AGO
	case "AD", "AND", "ANDORRA", "ANDORA":
		return AND
	case "AQ", "ATA", "NQ", "ATB", "ATN", "BQAQ", "NQAQ", "ANTARCTICA", "ANTARKTICA", "ANTARCTIKA", "ANTARKTIKA", "ANTARCTIC", "ANTARKTIC", "ANTARCTIK", "ANTARKTIK", "ANTARKTIS":
		return ATA
	case "AG", "ATG", "ANTIGUAANDBARBUDA", "ANTIGUABARBUDA", "ANTIGUA", "ANTIGUAUNDBARBUDA":
		return ATG
	case "AN", "ANT", "AHO", "ANHH", "NETHERLANDSANTILLES", "NETHERLSANTILLES", "NETHERLANDSANTILES", "NETHERLSANTILES", "NIEDERLAENDISCHEANTILLEN", "NIEDERLÄNDISCHANTILLEN":
		return ANT
	case "AE", "ARE", "UAE", "UNITEDARABEMIRATES", "ARABEMIRATES", "UNITEDEMIRATES", "VEREINIGTEARABISCHEEMIRATE":
		return ARE
	case "AR", "ARG", "ARGENTINA", "ARGENTIN", "RA", "ARGENTINIEN":
		return ARG
	case "AM", "ARM", "ARMENIA", "ARMENIYA", "ARMENIAN", "ARMENIEN":
		return ARM
	case "AW", "ABW", "ARUBA":
		return ABW
	case "AF", "AFG", "AFGHANISTAN", "AFHANISTAN", "AFGANISTAN", "AFGHANIAN", "AFGANIAN", "AFGHAN", "AFGHANI":
		return AFG
	case "BS", "BHS", "BAHAMAS", "BAGHAMAS", "BAGAMAS", "BAHAMIAN", "BAGAMIAN":
		return BHS
	case "BD", "BGD", "BANGLADESH", "BANGLADEH", "BANHGLADESH", "BANHLADESH", "BANHLADEH":
		return BGD
	case "BB", "BRB", "BAR", "BDS", "BARBADOS", "BARBODOS":
		return BRB
	case "BH", "BHR", "BAHRAIN", "BAGHRAIN":
		return BHR
	case "BY", "BLR", "BYS", "BYAA", "BELARUS", "BELORUS", "BELLARUSSIA", "BELARUSSIA", "BELLORUSSIA", "BELORUSSIA", "BELLARUSSIAN", "BELARUSSIAN", "BELLORUSSIAN", "BELORUSSIAN", "BYELORUSSIAN", "BYELORUSSIA", "BYELORUSSIYA", "WEISSRUSSLAND":
		return BLR
	case "BZ", "BLZ", "BIZ", "BELIZE":
		return BLZ
	case "BE", "BEL", "BELGIUM", "BELGUM", "BELGIEN":
		return BEL
	case "BJ", "BEN", "DHY", "BENIN", "DY", "DYBJ":
		return BEN
	case "BM", "BMU", "BERMUDA", "BERMUDS", "BERMUD":
		return BMU
	case "BG", "BGR", "BULGARIA", "BULGARIYA", "BULGARY", "BOLGARIA", "BOLGARIYA", "BULGARIEN":
		return BGR
	case "BO", "BOL", "BOLIVIA", "BOLIVIYA", "BOLIVIAN", "BOLIVIAPLURINATIONALSTATEOF", "BOLIVIAPLURINATIONALSTATE", "BOLIVIEN":
		return BOL
	case "BA", "BIH", "BOSNIAANDHERZEGOVINA", "BOSNIAHERZEGOVINA", "BOSNIA", "BOSNIEN", "BOSNIENUNDHERZEGOWINA":
		return BIH
	case "BW", "BWA", "BOTSWANA", "BOTSWANNA", "BOTSVANA", "BOTSVANNA":
		return BWA
	case "BR", "BRA", "BRAZIL", "BRAZILIA", "BRAZILIAN", "BRASILIEN":
		return BRA
	case "IO", "IOT", "BRITISHINDIANOCEANTERRITORY", "BRITISHINDIANTERRITORY", "BRITISCHESTERRITORIUM", "BRITISCHESTERRITORIUMIMINDISCHENOZEAN":
		return IOT
	case "BN", "BRN", "BRU", "BRUNEI", "BRUNEY", "BRUNEIDARUSSALAM":
		return BRN
	case "BF", "BFA", "HV", "HVO", "BURKINAFASO", "BURKINAANDFASO", "BURCINAFASO", "BURCINAANDFASO", "HVBF":
		return BFA
	case "BI", "BDI", "BURUNDI":
		return BDI
	case "BT", "BTN", "BHUTAN", "BGHUTAN":
		return BTN
	case "VU", "VUT", "NHB", "VANUATU", "NH", "NHVU":
		return VUT
	case "VA", "VAT", "HOLYSEEVATICAN", "HOLYSEE", "VATICAN", "VATICANCITYSTATE", "VATICANSTATE", "HOLYSEEVATIKAN", "VATIKAN", "VATIKANCITYSTATE", "VATIKANSTATE", "HOLYSEEVATIKANCITYSTATE", "VATIKANSTADT":
		return VAT
	case "GB", "DG", "GBR", "ADN", "DGA", "UNITEDKINGDOM", "UNITEDKINDOM", "UK", "GREATBRITAN", "GREATBRITAIN", "NORTHERNIRELAND", "BRITAN", "BRITAIN", "GROSSBRITANNIEN", "VEREINIGTESKÖNIGREICH", "VEREINIGTESKOENIGREICH": //nolint
		return GBR
	case "HU", "HUN", "HUNGARY", "HUNGAR", "HUNGARI", "VENGRIYA", "VENGRIA", "UNGARN":
		return HUN
	case "VE", "VEN", "VENEZUELA", "VENEZUELLA", "VENECUELA", "VENECUELLA", "YV": //nolint
		return VEN
	case "VG", "VGB", "IVB", "VIRGINISLANDSBRITISH", "VIRGINISLANDSBRITIH", "VIRGINISLSBRITIH", "VIRGINISLSBRITISH", "VIRGINISLANDSGB", "VIRGINISLANDSUK", "BRITISCHEJUNGFERNINSELN":
		return VGB
	case "VI", "VIR", "ISV", "VIRGINISLANDSUS", "USVIRGINISLANDS", "USVI", "AMERIKANISCHEJUNGFERNINSELN":
		return VIR
	case "TL", "TP", "TLS", "TMP", "TPTL", "TIMORLESTE", "EASTTIMOR", "TIMOR", "TIMORELESTE", "EASTTIMORE", "TIMORE", "TIMORLESTEEASTTIMORE", "OSTTIMOR":
		return TLS
	case "VN", "VNM", "VIE", "VDR", "VD", "VIETNAM", "VETNAM", "VIETNAME", "VETNAME", "VDVN", "VIỆTNAM", "CỘNGHÒAXÃHỘICHỦNGHĨAVIỆTNAM", "CHỦNGHĨAVIỆTNAM", "NGHĨAVIỆTNAM":
		return VNM
	case "GA", "GAB", "GABON", "GABUN":
		return GAB
	case "HT", "HTI", "HAITI", "GAITI":
		return HTI
	case "GY", "GUY", "GUYANA":
		return GUY
	case "GM", "GMB", "WAG", "GAMBIA", "GAMBIYA":
		return GMB
	case "GH", "GHA", "GHANA", "HANA":
		return GHA
	case "GP", "GLP", "GUADELOUPE", "GUADELUPE", "GUADELOOPE", "GUADELOUPA", "GUADELUPA", "GUADELOOPA":
		return GLP
	case "GT", "GTM", "GCA", "GUATEMALA":
		return GTM
	case "GN", "GIN", "GUINEA", "GUINEYA":
		return GIN
	case "GW", "GNB", "GBS", "GUINEABISSAU":
		return GNB
	case "DE", "DEU", "DD", "DDR", "GER", "GERMANY", "GERMANIYA", "DEUTSCHLAND", "DEUTSCH", "DDDE":
		return DEU
	case "GI", "GIB", "GBZ", "GIBRALTAR", "HIBRALTAR":
		return GIB
	case "HN", "HND", "HONDURAS", "GONDURAS":
		return HND
	case "HK", "HKG", "HONGKONG", "HONKONG":
		return HKG
	case "GD", "GRD", "GRENADA", "GRINADA", "WG":
		return GRD
	case "GL", "GRL", "GREENLAND", "GRÖNLAND", "GROENLAND":
		return GRL
	case "GR", "GRC", "GREECE", "GRECE", "GRIECHENLAND", "GRECIYA":
		return GRC
	case "GE", "GEO", "GEORGIA", "GEORGIYA", "GEORGIEN", "GRUZIYA":
		return GEO
	case "GU", "GUM", "GUAM":
		return GUM
	case "DK", "DNK", "DENMARK", "DANMARK", "DÄNEMARK", "DAENEMARK", "KONGERIGETDANMARK", "DANMARKKONGERIGET", "DANIYA":
		return DNK
	case "CD", "COD", "ZRE", "ZAR", "ZR", "ZRCD", "CONGODEMOCRATICREPUBLIC", "DEMOCRATICREPUBLICOFTHECONGO", "CONGODEMOCRATICREP", "CONGODEMOCRATIC", "CONGOTHEDEMOCRATICREPUBLICOF", "KONGODEMOCRACTICREPUBLIC", "KONGODEMOCRATICREP", "KONGODEMOCRATIC", "KONGOTHEDEMOCRATICREPUBLICOF", "ZAIRE", "ZAIR", "DEMOKRATISCHEREPUBLIKKONGO", "CONGOREPUBLIC", "KONGOREPUBLIC", "REPUBLICOFCONGO", "REPUBLICOFKONGO":
		return COD
	case "DJ", "DJI", "AFI", "DJIBOUTI", "AIDJ", "DSCHIBUTI":
		return DJI
	case "DM", "DMA", "DOMINICA", "DOMINIKA":
		return DMA
	case "DO", "DOM", "DOMINICANREPUBLIC", "DOMINICANA", "DOMINIKANA", "DOMINIKANISCHEREPUBLIK":
		return DOM
	case "EG", "EGY", "EGYPT", "ÄGYPTEN", "AEGYPTEN":
		return EGY
	case "ZM", "ZMB", "RNR", "ZAMBIA", "SAMBIA":
		return ZMB
	case "EH", "ESH", "WESTERNSAHARA", "WESTSAHARA":
		return ESH
	case "ZW", "ZWE", "ZIM", "RHO", "RSR", "ZIMBABWE", "ZIMBABVE", "RH", "RHZW", "SIMBABWE":
		return ZWE
	case "IL", "ISR", "ISRAEL", "IZRAIL":
		return ISR
	case "IN", "IND", "INDIA", "INDIAN", "INDIYA", "SKM", "SKIN", "INDIEN":
		return IND
	case "ID", "IDN", "INA", "INDONESIA", "REPUBLICOFINDONESIA", "RI", "INDONESIEN":
		return IDN
	case "JO", "JOR", "HKJ", "JORDAN", "JORDANIEN":
		return JOR
	case "IQ", "IRQ", "IRAQ", "IRAK":
		return IRQ
	case "IR", "IRN", "IRI", "IRAN":
		return IRN
	case "IE", "IRL", "IRELAND", "IRLAND":
		return IRL
	case "IS", "ISL", "ICELAND", "ISLAND":
		return ISL
	case "ES", "EA", "IC", "ESP", "SPAIN", "SPANIEN", "ISPANIA":
		return ESP
	case "IT", "ITA", "ITALY", "ITALIYA", "ITALIEN":
		return ITA
	case "YE", "YEM", "YMD", "YEMEN", "IEMEN", "YD", "YDYE", "JEMEN":
		return YEM
	case "KZ", "KAZ", "KAZAKHSTAN", "KAZAHSTAN", "KASACHSTAN":
		return KAZ
	case "KY", "CYM", "CAYMANISLANDS", "KAYMANISLANDS", "KAIMANINSELN":
		return CYM
	case "KH", "KHM", "CAMBODIA", "KAMBODSCHA":
		return KHM
	case "CM", "CMR", "CAMEROON", "KAMERUN":
		return CMR
	case "CA", "CAN", "CDN", "CANADA", "KANADA":
		return CAN
	case "QA", "QAT", "QATAR", "KATAR":
		return QAT
	case "KE", "KEN", "EAK", "KENYA":
		return KEN
	case "CY", "CYP", "CYPRUS", "CIPRUS", "ZYPERN", "REPUBLIKZYPERN":
		return CYP
	case "KI", "KIR", "CT", "CTE", "CTKI", "KIRIBATI", "CIRIBATI", "KIRIBATY", "CIRIBATY":
		return KIR
	case "CN", "CHN", "CHINA", "CHINESE", "RC", "KITAY":
		return CHN
	case "CC", "CCK", "KEELING", "COCOS", "COCOSKEELINGISLANDS", "COCOSISLANDS", "KOKOSISLANDS", "KOKOSINSELN":
		return CCK
	case "CO", "COL", "COLOMBIA", "KOLUMBIEN":
		return COL
	case "KM", "COM", "COMOROS", "KOMOREN":
		return COM
	case "CG", "COG", "RCB", "CONGO", "KONGO":
		return COG
	case "KP", "PRK", "DEMOCRATICPEOPLESREPUBLICOFKOREA", "KOREADEMOCRATICPEOPLESREPUBLICOF", "KOREADEMOCRATICPEOPLESREPUBLIC", "KOREANORTH", "NORTHKOREA", "NORDKOREA":
		return PRK
	case "KR", "KOR", "ROK", "KOREA", "KOREYA", "SOUTHKOREA", "KOREAREPUBLICOF", "REPUBLICOFKOREA", "KOREAREPOF", "SÜDKOREA", "SUEDKOREA":
		return KOR
	case "CR", "CRI", "COSTARICA", "KOSTARIKA", "KOSTARICA", "COSTARIKA":
		return CRI
	case "CI", "CIV", "COTEDIVOIRE", "CÔTEDIVOIRE", "IVORYCOAST", "ELFENBEINKÜSTE", "ELFENBEINKUESTE":
		return CIV
	case "CU", "CUB", "CUBA", "CUBAREPUBLIC", "REPUBLICCUBA", "KUBA":
		return CUB
	case "KW", "KWT", "KUWAIT":
		return KWT
	case "KG", "KGZ", "KYRGYZSTAN", "KIRGISISTAN":
		return KGZ
	case "LA", "LAO", "LAOS", "LAODEMOCRATICPEOPLESREPUBLIC", "LAOSDEMOCRATICPEOPLESREPUBLIC", "LAOPEOPLESDEMOCRATICREPUBLIC":
		return LAO
	case "LV", "LVA", "LAT", "LATVIA", "LATVIYA", "LETTLAND":
		return LVA
	case "LS", "LSO", "LESOTHO":
		return LSO
	case "LR", "LBR", "LIBERIA":
		return LBR
	case "LB", "LBN", "LEBANON", "RL", "LIBANON":
		return LBN
	case "LY", "LBY", "LBA", "LIBYA", "LIVIA", "LIVIYA", "LIBYAN", "LIBYANARABJAMAHIRIYA", "LF", "LIBYEN":
		return LBY
	case "LT", "LTU", "LITHUANIA", "LITAUEN", "LITVA":
		return LTU
	case "LI", "LIE", "LIECHTENSTEIN", "LIEHTENSTEIN", "FL":
		return LIE
	case "LU", "LUX", "LUXEMBOURG", "LUXEMBURG":
		return LUX
	case "MU", "MUS", "MAURITIUS":
		return MUS
	case "MR", "MRT", "MAURITANIA", "MAURETANIEN":
		return MRT
	case "MG", "MDG", "MADAGASCAR", "RM", "MADAGASKAR":
		return MDG
	case "YT", "MYT", "MAYOTTE":
		return MYT
	case "MO", "MAC", "MACAUCHINA", "MACAU", "MACAO", "MACAUSAR", "MACAOSAR":
		return MAC
	case "MK", "MKD", "MACEDONIA", "MACEDONIAFYRO", "MACEDONIATHEFORMERYUGOSLAVREPUBLICOF", "REPUBLICOFNORTHMACEDONIA", "REPUBLICOFMACEDONIA", "NORTHMACEDONIA", "MACEDONIANORTH", "NORDMAZEDONIEN":
		return MKD
	case "MW", "MWI", "MAW", "MALAWI", "MALAVI":
		return MWI
	case "MY", "MYS", "MAL", "MALAYSIA", "MALAYSIYA":
		return MYS
	case "ML", "MLI", "RMM", "MALI":
		return MLI
	case "MV", "MDV", "MALDIVES", "MALEDIVEN":
		return MDV
	case "MT", "MLT", "MALTA":
		return MLT
	case "MP", "MNP", "NORTHERNMARIANAISLANDS", "NORTHERNMARIANAIS", "MARIANAISLANDS", "NÖRDLICHEMARIANEN", "NOERDLICHEMARIANEN":
		return MNP
	case "MA", "MAR", "MOROCCO", "MOROCO", "MOROKO", "MAROKKO":
		return MAR
	case "MQ", "MTQ", "MARTINIQUE":
		return MTQ
	case "MH", "MHL", "MARSHALLISLANDS", "MARSHALL", "REPUBLICOFTHEMARSHALLISLANDS", "MARSHALLINSELN":
		return MHL
	case "MX", "MEX", "MEXICO", "MEXIKO":
		return MEX
	case "FM", "FSM", "MICRONESIA", "MICRONESIAFEDERATEDSTATESOF", "MICRONESIAFEDST", "MIKRONESIEN":
		return FSM
	case "MZ", "MOZ", "MOZAMBIQUE", "MOZAMBIQ", "MOSAMBIK":
		return MOZ
	case "MD", "MDA", "MOLDOVA", "MOLDAVIA", "REPUBLIKMOLDOVA":
		return MDA
	case "MC", "MCO", "MONACO", "MONAKO":
		return MCO
	case "MN", "MNG", "MONGOLIA", "MONGOLEI":
		return MNG
	case "MS", "MSR", "MONTSERRAT":
		return MSR
	case "MM", "BU", "MMR", "BUMM", "MYANMAR", "BURMA":
		return MMR
	case "NA", "NAM", "NAMIBIA", "NAMIBIE":
		return NAM
	case "NR", "NRU", "NAURU":
		return NRU
	case "NP", "NPL", "NEPAL", "NEPALI":
		return NPL
	case "NE", "NER", "NIGER", "NIGGER", "RN":
		return NER
	case "NG", "NGA", "NGR", "WAN", "NIGERIA", "NIGERIYA", "NIGGERIA", "NIGGERIYA":
		return NGA
	case "NL", "NLD", "NED", "NETHERLANDS", "NETHERLAND", "HOLLAND", "HOLLANDIA", "HOLLANDIYA", "NIEDERLANDE":
		return NLD
	case "NI", "NIC", "NICARAGUA":
		return NIC
	case "NU", "NIU", "NIUE":
		return NIU
	case "NZ", "NZL", "NEWZEALAND", "NEWZELANDIA", "NEWZELAND", "NEUSEELAND":
		return NZL
	case "NC", "NCL", "NEWCALEDONIA", "NEWCALEDONIYA", "NEUKALEDONIEN":
		return NCL
	case "NO", "NOR", "NORWAY", "NORWEGEN":
		return NOR
	case "OM", "OMN", "OMAN":
		return OMN
	case "BV", "BVT", "BOUVET", "BOUVETE", "BOUVETISLAND", "ISLANDOFBOUVET", "BOUVETINSEL":
		return BVT
	case "IM", "IMN", "GBM", "ISLEOFMAN":
		return IMN
	case "NF", "NFK", "NORFOLKISLAND", "NORFOLK", "NORFOLCISLAND", "NORFOLC", "NORFOLKINSEL":
		return NFK
	case "PN", "PCN", "PITCAIRN", "THEPITCAIRN", "PITCAIRNISLANDS", "THEPITCAIRNISLANDS", "DUCIEANDOENOISLANDS", "DUCIEANDOENO", "PITCAIRNINSELN":
		return PCN
	case "CX", "CXR", "CHRISTMASISLAND", "TERRITORYOFCHRISTMASISLAND", "WEIHNACHTSINSEL":
		return CXR
	case "SH", "TA", "SHN", "TAA", "ASC", "SAINTHELENA", "SAINTELENA", "STHELENA", "STELENA", "TRISTAN", "ASCENSIONANDTRISTANDACUNHA", "ASCENSIONTRISTANDACUNHA", "TRISTANDACUNHA", "SANKTHELENA":
		return SHN
	case "WF", "WLF", "WALLISANDFUTUNAISLANDS", "WALLISFUTUNAISLANDS", "WALLISANDFUTUNA", "WALLISFUTUNA", "WALLISUNDFUTUNA":
		return WLF
	case "HM", "HMD", "HEARDISLANDANDMCDONALDISLANDS", "HEARDISLAND", "HEARDUNDMCDONALDINSELN":
		return HMD
	case "CV", "CPV", "CAPEVERDE", "KAPVERDE":
		return CPV
	case "CK", "COK", "COOKISLANDS", "COOKINSELN":
		return COK
	case "WS", "WSM", "SAMOA":
		return WSM
	case "SJ", "SJM", "SVALBARDANDJANMAYENISLANDS", "SVALBARD", "SVALBARDUNDJANMAYEN":
		return SJM
	case "TC", "TCA", "TURKSANDCAICOSISLANDS", "TURKSANDCAICOSIS", "CAICOSISLANDS", "CACOSISLANDS", "TURKSUNDCACIOINSELN":
		return TCA
	case "UM", "UMI", "UNITEDSTATESMINOROUTLYINGISLANDS", "MINOROUTLYINGISLANDS", "MINOROUTLYING", "USMI", "JT", "JTN", "JTUM", "MI", "MID", "MIUM", "PU", "PUS", "PUUM", "WK", "WAK", "WKUM", "KLEINEINSELBESITZUNGENDERVEREINIGTENSTAATEN":
		return UMI
	case "PK", "PAK", "PAKISTAN", "PACISTAN":
		return PAK
	case "PW", "PLW", "PALAU":
		return PLW
	case "PS", "PSE", "PLE", "PALESTINE", "PALESTINA", "PALESTINIAN", "PALESTINIANTERRITORY", "PALÄSTINA", "PALAESTINA":
		return PSE
	case "PA", "PAN", "PCZ", "PANAMA", "PANAMIAN", "PANAM", "PZ", "PZPA":
		return PAN
	case "PG", "PNG", "PAPUANEWGUINEA", "PAPUA", "PAPUANEUGUINEA", "NEWGUINEA", "NEUGUINEA":
		return PNG
	case "PY", "PRY", "PARAGUAY":
		return PRY
	case "PE", "PER", "PERU":
		return PER
	case "PL", "POL", "POLAND", "POLSKI", "POLSHA", "POLEN":
		return POL
	case "PT", "PRT", "PORTUGAL", "PORTUGALIAN", "PORTUGALIYA":
		return PRT
	case "PR", "PRI", "PUERTORICO", "PUERTORIKO":
		return PRI
	case "RE", "REU", "REUNION", "RÉUNION":
		return REU
	case "RU", "RUS", "SUN", "RUSSIA", "RUSSO", "RUSSISH", "RUSSLAND", "RUSLAND", "RUSIA", "ROSSIA", "ROSSIYA", "RUSSIAN", "RUSSIANFEDERATION", "USSR":
		return RUS
	case "RW", "RWA", "RWANDA", "RUANDA", "RUWANDA":
		return RWA
	case "RO", "ROU", "ROM", "ROMANIA", "RUMINIA", "RUMINIYA", "RUMÄNIEN", "RUMAENIEN":
		return ROU
	case "SV", "SLV", "ESA", "ELSALVADOR":
		return SLV
	case "SM", "SMR", "RSM", "SANMARINO":
		return SMR
	case "ST", "STP", "SAOTOMEANDPRINCIPE", "SAOTOME", "SAOTOMEUNDPRINCIPE":
		return STP
	case "SA", "SAU", "SAUDIARABIA", "SAUDI", "SAUDIARABIEN":
		return SAU
	case "SZ", "SWZ", "SWAZILAND", "SWASILAND":
		return SWZ
	case "SC", "SYC", "SEYCHELLES", "SEYCHELLEN":
		return SYC
	case "SN", "SEN", "SENEGAL":
		return SEN
	case "PM", "SPM", "SAINTPIERREANDMIQUELON", "SAINTPIERRE", "STPIERREANDMIQUELON", "STPIERRE", "SANKTPIERRE", "SANKTPIERREUNDMIQUELON":
		return SPM
	case "VC", "VCT", "SAINTVINCENTANDTHEGRENADINES", "SAINTVINCENT", "STVINCENTANDTHEGRENADINES", "STVINCENT", "WV", "STVINCENTUNDDIEGRENADINEN":
		return VCT
	case "KN", "KNA", "SAINTKITTSANDNEVIS", "SAINTKITTSNEVIS", "SAINTKITTS", "STKITTSANDNEVIS", "STKITTSNEVIS", "STKITTS", "SANKTKITTSUNDNEVIS":
		return KNA
	case "LC", "LCA", "SAINTLUCIA", "STLUCIA", "LUCIA", "WL":
		return LCA
	case "SG", "SGP", "SINGAPORE", "SINGPAORE", "SINGAPORECITY", "SINGAPOUR", "SINGAPURA", "SINGAPUR": //nolint
		return SGP
	case "SY", "SYR", "SYRIA", "SYRIAN", "SYRIANARABREPUBLIC", "SYRIEN":
		return SYR
	case "SK", "SVK", "CSHH", "SLOVAKIA", "SLOVAK", "SLOVAKIYA", "SLOVACIA", "SLOVAC", "SLOVACIYA", "SLOWAKEI":
		return SVK
	case "SI", "SVN", "SLO", "SLOVENIA", "SLOVENIYA", "SLOWENIEN":
		return SVN
	case "US", "USA", "UNITEDSTATES", "UNITEDSTATESOFAMERICA", "USOFAMERICA", "USAMERICA", "VEREINIGTESTAATENVONAMERIKA":
		return USA
	case "SB", "SLB", "SOLOMONISLANDS", "SOLOMON", "SALOMONEN":
		return SLB
	case "SO", "SOM", "SOMALIA", "SOMALI":
		return SOM
	case "SD", "SDN", "SUDAN", "SUDANE", "UMHŪRIYYATASSŪDĀN", "SŪDĀN", "جمهوريةالسودان", "السودان":
		return SDN
	case "SR", "SUR", "SME", "SURINAME", "SURINAM":
		return SUR
	case "SL", "SLE", "WAL", "SIERRALEONE", "SIERRALEON", "SIERALEONE", "SIERALEON":
		return SLE
	case "TJ", "TJK", "TAJIKISTAN", "TADJIKISTAN", "TADSCHIKISTAN":
		return TJK
	case "TW", "TWN", "TPE", "TAIWAN":
		return TWN
	case "TH", "THA", "THAILAND", "TAILAND", "THAI":
		return THA
	case "TZ", "TZA", "EAT", "EAZ", "TANZANIA", "TANZANIYA", "TANSANIA":
		return TZA
	case "TG", "TGO", "TOGO":
		return TGO
	case "TK", "TKL", "TOKELAU":
		return TKL
	case "TO", "TON", "TONGA":
		return TON
	case "TT", "TTO", "TRINIDADANDTOBAGO", "TRINIDAD", "TRINADUNDTOBAGO":
		return TTO
	case "TV", "TUV", "TUVALU":
		return TUV
	case "TN", "TUN", "TUNISIA", "TUNESIEN":
		return TUN
	case "TM", "TKM", "TMN", "TURKMENISTAN", "TURKMENISTON", "TURKMENI", "TURKMENIA", "TURKMENIYA":
		return TKM
	case "TR", "TUR", "TURKEY", "TURCIA", "TURKISH", "TÜRKEI", "TUERKEI":
		return TUR
	case "UG", "UGA", "EAU", "UGANDA":
		return UGA
	case "UZ", "UZB", "UZBEKISTAN", "UZBEKISTON":
		return UZB
	case "UA", "UKR", "UKRAINE", "UKRAINA": //nolint
		return UKR
	case "UY", "URY", "URUGUAY":
		return URY
	case "XW", "XWA", "WALES":
		return XWA
	case "FO", "FRO", "FAROEISLANDS", "FAROE", "FÄRÖER", "FAEROERER":
		return FRO
	case "FJ", "FJI", "FIJI", "FIDSCHI":
		return FJI
	case "PH", "PHL", "PHI", "PHILIPPINES", "PHILIPINES", "PI", "RP", "PHILIPPINEN": //nolint
		return PHL
	case "FI", "SF", "FIN", "FINLAND", "FINNISH", "FINNLAND":
		return FIN
	case "FK", "FLK", "FALKLANDISLANDSMALVINAS", "MALVINAS", "FALKLANDISLANDS", "FALKLAND", "FALKLANDINSELN":
		return FLK
	case "FR", "CP", "FX", "FRA", "FXX", "CPT", "FXFR", "FRANCE", "FRENCH", "FRANKREICH":
		return FRA
	case "GF", "GUF", "FRENCHGUIANA", "GUIANA", "FRANZÖSISCHGUYANA", "FRANZOESISCHGUYANA":
		return GUF
	case "PF", "PYF", "FRENCHPOLYNESIA", "POLYNESIA", "FRANZÖSISCHPOLYNESIEN", "FRANZOESISCHPOLYNESIEN":
		return PYF
	case "TF", "ATF", "FRENCHSOUTHERNTERRITORIES", "SOUTHERNTERRITORIESFRENCH", "FRANZÖSISCHESÜDUNDANTARKTISGEBIETE", "FRANZOESISCHESUEDUNDANTARKTISGEBIETE":
		return ATF
	case "HR", "HRV", "CRO", "CROATIA", "KROATIA", "KROATIEN":
		return HRV
	case "CF", "CAF", "CTA", "RCA", "CENTRALAFRICANREPUBLIC", "CENTRALAFRICANREP", "CENTRALAFRICAN", "ZENTRALAFRIKA":
		return CAF
	case "TD", "TCD", "CHAD", "TSCHAD":
		return TCD
	case "CZ", "CZE", "CZECHREPUBLIC", "CZECH", "TSCHECHIEN":
		return CZE
	case "CL", "CHL", "RCH", "CHILE":
		return CHL
	case "CH", "CHE", "SWITZERLAND", "SWISS", "SCHWEIZ", "SUISSE", "SVIZZERA", "SVIZRA", "HELVETIA":
		return CHE
	case "SE", "SWE", "SWEDEN", "SCHWEDEN":
		return SWE
	case "XS", "XSC", "SCOTLAND", "SCHOTTLAND":
		return XSC
	case "LK", "LKA", "SRILANKA":
		return LKA
	case "EC", "ECU", "ECUADOR":
		return ECU
	case "GQ", "GNQ", "EQG", "GEQ", "EQUATORIALGUINEA", "ÄQUATORIALGUINEA", "AEQUATORIALGUINEA":
		return GNQ
	case "ER", "ERI", "ERITREA":
		return ERI
	case "EE", "EST", "ESTONIA", "EW", "ESTLAND":
		return EST
	case "ET", "ETH", "ETHIOPIA", "ÄTHOPIEN", "AETHOPIEN":
		return ETH
	case "ZA", "ZAF", "SOUTHAFRICA", "SÜDAFRIKA", "SUEDAFRIKA":
		return ZAF
	case "YU", "YUG", "YUGOSLAVIA", "UGOSLAVIA", "YUGOSLAVIYA", "UGOSLAVIYA", "SERBIAANDMONTENEGRO", "CS", "SCG", "JUGOSLAWIEN":
		return YUG
	case "GS", "SGS", "SOUTHGEORGIAANDTHESOUTHSANDWICHISLANDS", "SOUTHGEORGIAANDTHESOUTHSANDWICH", "SOUTHGEORGIATHESOUTHSWICHISLANDS", "SOUTHGEORGIA", "SÜDGEORGIEN", "SUEDGEORGIEN":
		return SGS
	case "JM", "JAM", "JAMAICA", "JAMAIKA", "YAMAICA", "YAMAIKA", "JA":
		return JAM
	case "ME", "MNE", "MONTENEGRO":
		return MNE
	case "BL", "BLM", "SAINTBARTHELEMY", "STBARTHELEMY":
		return BLM
	case "SX", "SXM", "SINTMAARTENDUTCH", "SAINTMAARTEN", "SINTMAARTEN", "STMAARTEN":
		return SXM
	case "RS", "SRB", "CSXX", "SERBIA", "SERBIYA", "SERBIEN":
		return SRB
	case "AX", "ALA", "ALANDISLANDS", "ISLANDSALAND", "ALAND":
		return ALA
	case "BQ", "BES", "BONAIRE", "BONAIR", "BONEIRU", "BONAIRESINTEUSTATIUSANDSABA", "BONAIRESINTEUSTATIUSSABA", "BONAIRESTEUSTANDSABA", "BONAIRESTEUSTSABA", "SINTEUSTATIUSANDSABA", "SINTEUSTATIUS":
		return BES
	case "GG", "GGY", "GBA", "GBG", "GUERNSEY":
		return GGY
	case "JE", "JEY", "GBJ", "JERSEY", "JERSIEY":
		return JEY
	case "CW", "CUW", "CURACAO", "CURAÇAO", "CURAQAO", "CURAKAO", "KURACAO", "KURAKAO":
		return CUW
	case "MF", "MAF", "SAINTMARTINFRENCH", "STMARTINFRENCH", "SANKTMARTIN":
		return MAF
	case "SS", "SSD", "SOUTHSUDAN", "SOUTHSUDANE", "REPUBLICOFSOUTHSUDAN", "SOUTHSUDANREPUBLICOF", "PAGUOTTHUDÄN", "SÜDSUDAN", "SUEDSUDAN":
		return SSD
	case "JP", "JPN", "JAPAN":
		return JPN
	case "XK", "XKX", "KOS", "KOSOVO", "COSOVO", "КОСОВО", "KOSOVËS", "РЕПУБЛИКАКОСОВО", "REPUBLIKAKOSOVO", "REPUBLIKACOSOVO", "REPUBLIKAKOSOVËS", "REPUBLICAKOSOVO", "REPUBLICACOSOVO", "REPUBLICAKOSOVËS", "KOSOVOREPUBLIC", "COSOVOREPUBLIC", "KOSOVËSREPUBLIC":
		return XKX
	case "XX", "NONE", "NON", "NICHT", "NICHTS":
		return None
	case "INTERNATIONAL":
		return International
	case "UIFN", "INTERNATIONALFREEPHONE", "TOLLFREEPHONE":
		return NonCountryInternationalFreephone
	case "INMARSAT":
		return NonCountryInmarsat
	case "MMS", "MARITIMEMOBILESERVICE", "MARITIMEMOBILESERVICES", "MARITIMEMOBILE", "MARITIME":
		return NonCountryMaritimeMobileService
	case "UNIVERSALPERSONALTELECOMMUNICATIONSSERVICES", "UNIVERSALPERSONALTELECOMMUNICATIONSSERVICE", "UNIVERSALPERSONALTELECOMMUNICATIONS", "UNIVERSALPERSONALTELECOMMUNICATION":
		return NonCountryUniversalPersonalTelecommunicationsServices
	case "NCP", "NATIONALNONCOMMERCIALPURPOSES", "NONCOMMERCIALPURPOSES", "NATIONALNONCOMMERCIAL", "NONCOMMERCIAL":
		return NonCountryNationalNonCommercialPurposes
	case "GMSS", "GLOBALMOBILESATELLITESYSTEM", "GLOBALMOBILESATELITESYSTEM", "GLOBALMOBILESATELLITE", "GLOBALMOBILESATELITE":
		return NonCountryGlobalMobileSatelliteSystem
	case "INTERNATIONALNETWORKS", "INTERNATIONALNETWORKSSERVICE", "INTERNATIONALNETWORKSSERVICES":
		return NonCountryInternationalNetworks
	case "DISASTERRELIEF", "DISASTER":
		return NonCountryDisasterRelief
	case "IPRS", "INTERNATIONALPREMIUMRATESERVICE", "PREMIUMRATESERVICE", "INTERNATIONALPREMIUMRATESERVICES", "PREMIUMRATESERVICES":
		return NonCountryInternationalPremiumRateService
	case "ITPCS", "INTERNATIONALTELECOMMUNICATIONSPUBLICCORRESPONDENCESERVICETRIAL", "INTERNATIONALTELECOMMUNICATIONSPUBLICCORRESPONDENCESERVICE", "InternationalTELECOMMUNICATIONSPUBLICCORRESPONDENCESERVICES", "InternationalTELECOMMUNICATIONSCORRESPONDENCESERVICE", "InternationalTELECOMMUNICATIONSCORRESPONDENCESERVICES":
		return NonCountryInternationalTelecommunicationsCorrespondenceService
	}
	return Unknown
}

// ByNumeric - return CountryCode by country Alpha-2 / Alpha-3 / numeric code, example: rus := ByNumeric(643),
// returns countries.Unknown, if country code not found or not valid
func ByNumeric(numeric int) CountryCode {
	if code := CountryCode(numeric); code.IsValid() {
		return code
	}
	return Unknown
}

// IsValid - returns true, if code is correct
func (c CountryCode) IsValid() bool {
	return c.Alpha2() != UnknownMsg
}
