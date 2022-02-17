## countries

Countries - ISO 3166 (ISO3166-1, ISO3166, Digit, Alpha-2, Alpha-3) countries codes and names (on eng and rus), ISO 4217 currency designators, ITU-T E.164 IDD calling phone codes, countries capitals, UN M.49 regions codes, IANA ccTLD countries domains, IOC/NOC and FIFA codes, VERY FAST, NO maps[], NO slices[], NO init() funcs, NO external links/files/data, NO interface{}, NO specific dependencies, Databases/JSON/GOB/XML/CSV compatible, Emoji countries flags and currencies support, full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and IANA ccTLD standarts.<br/><br/>Full support ISO-3166-1, ISO-4217, ITU-T E.164, Unicode CLDR and IANA ccTLD standarts.<br/><br/>
[![GoDoc](https://godoc.org/github.com/biter777/countries?status.svg)](https://godoc.org/github.com/biter777/countries)
[![GoDev](https://img.shields.io/badge/godev-reference-5b77b3)](https://pkg.go.dev/github.com/biter777/countries?tab=doc)
[![Go Walker](https://img.shields.io/badge/gowalker-reference-5b77b3)](https://gowalker.org/github.com/biter777/countries)
[![Documentation Status](https://readthedocs.org/projects/countries/badge/?version=latest)](https://countries.readthedocs.io/en/latest/?badge=latest)
[![DOI](https://zenodo.org/badge/182808313.svg)](https://zenodo.org/badge/latestdoi/182808313)
[![codeclimate](https://codeclimate.com/github/biter777/countries/badges/gpa.svg)](https://codeclimate.com/github/biter777/countries)
[![GolangCI](https://golangci.com/badges/github.com/biter777/countries.svg?style=flat)](https://golangci.com/r/github.com/biter777/countries)
[![GoReport](https://goreportcard.com/badge/github.com/biter777/countries)](https://goreportcard.com/report/github.com/biter777/countries)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/08eb1d2ff62e465091b3a288ae078a96)](https://www.codacy.com/manual/biter777/countries?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=biter777/countries&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/biter777/countries/branch/master/graph/badge.svg)](https://codecov.io/gh/biter777/countries)
[![Coverage Status](https://coveralls.io/repos/github/biter777/countries/badge.svg?branch=master)](https://coveralls.io/github/biter777/countries?branch=master)
[![Coverage](https://img.shields.io/badge/coverage-gocover.io-brightgreen)](https://gocover.io/github.com/biter777/countries)
[![ISO](https://img.shields.io/badge/powered%20by-ISO-brightgreen)](https://www.iso.org/)
[![ITU](https://img.shields.io/badge/powered%20by-ITU-brightgreen)](https://www.itu.int/)
[![IANA](https://img.shields.io/badge/powered%20by-IANA-brightgreen)](http://www.iana.org/)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-brightgreen.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbiter777%2Fcountries.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbiter777%2Fcountries?ref=badge_shield)
[![Build Status](https://travis-ci.org/biter777/countries.svg?branch=master)](https://travis-ci.org/biter777/countries)
[![Build status](https://ci.appveyor.com/api/projects/status/t9lpor9o8tpacpmr/branch/master?svg=true)](https://ci.appveyor.com/project/biter777/countries/branch/master)
[![Circle CI](https://circleci.com/gh/biter777/countries/tree/master.svg?style=shield)](https://circleci.com/gh/biter777/countries/tree/master)
[![Semaphore Status](https://biter777.semaphoreci.com/badges/countries.svg?style=shields)](https://biter777.semaphoreci.com/projects/countries)
[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![Codeship Status](https://codeship.com/projects/00db4400-1803-0138-1132-7ab932dd1523/status?branch=master)](https://app.codeship.com/projects/381056) 
<a href="//www.dmca.com/Protection/Status.aspx?ID=7a019cc5-ec73-464b-9707-4b33726f348f" title="DMCA.com Protection Status" class="dmca-badge"> <img src ="https://img.shields.io/badge/DMCA-protected-brightgreen"  alt="DMCA.com Protection Status" /></a> 
[![Dependencies Free](https://img.shields.io/badge/dependencies-free-brightgreen)](https://pkg.go.dev/github.com/biter777/countries?tab=imports)
[![Gluten Free](https://img.shields.io/badge/gluten-free-brightgreen)](https://www.scsglobalservices.com/services/gluten-free-certification)
[![DepShield Badge](https://depshield.sonatype.org/badges/biter777/countries/depshield.svg)](https://depshield.github.io)
[![Stars](https://img.shields.io/github/stars/biter777/countries?label=Please%20like%20us&style=social)](https://github.com/biter777/countries/stargazers)
<br/>

### installation

    go get github.com/biter777/countries

### usage

```go
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
```

### options

You can take a counties names in russian language, use StringRus(). For Emoji use Emoji(). Enjoy!

```go
import "github.com/biter777/countries"
```

For more complex options, consult the [documentation](http://godoc.org/github.com/biter777/countries).

### contributing

1) <b>Welcome pull requests, bug fixes and issue reports.</b><br/>
Before proposing a change, please discuss it first by raising an issue.<br/>
<small>Contributors list: <a href="https://github.com/biter777">@biter777</a>, <a href="https://github.com/gavincarr">@gavincarr</a>, <a href="https://github.com/benja-M-1">@benja-M-1</a>.</small><br/>

2) <b>Donate</b>. A donation isn't necessary, but it's welcome.<br/>
<noscript><a href="https://liberapay.com/biter777/donate"><img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg"></a></noscript> 
[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I2I61D1XZ) <a href="https://pay.cloudtips.ru/p/94fc4268" target="_blank"><img height="30" src="https://usa.visa.com/dam/VCOM/regional/lac/ENG/Default/Partner%20With%20Us/Payment%20Technology/visapos/full-color-800x450.jpg"></a> <a href="https://pay.cloudtips.ru/p/94fc4268" target="_blank"><img height="30" src="https://brand.mastercard.com/content/dam/mccom/brandcenter/thumbnails/mastercard_debit_sym_decal_web_105px.png"></a> <a href="https://pay.cloudtips.ru/p/94fc4268" target="_blank"><img height="30" src="https://developer.apple.com/assets/elements/icons/apple-pay/apple-pay.svg"></a> <a href="https://pay.cloudtips.ru/p/94fc4268" target="_blank"><img height="30" src="https://developers.google.com/pay/api/images/brand-guidelines/google-pay-mark.png"></a> <a href="https://money.yandex.ru/to/4100164702007" target="_blank"><img width="125" height="25" src="https://yastatic.net/q/logoaas/v1/Yandex%20Money.svg"></a><br/>

3) <b>Star us</b>. Give us a star, please, if it's not against your religion :)
