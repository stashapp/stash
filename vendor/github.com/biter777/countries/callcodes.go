package countries

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// CallCode - calling code of country
type CallCode int64 // int64 for database/sql/driver.Valuer compatibility

// CallCodeInfo - all info about CallCode
type CallCodeInfo struct {
	Code      CallCode
	Countries []CountryCode
}

// String - implements fmt.Stringer, returns a calling phone code in string, example for UK: "+44"
func (c CallCode) String() string {
	return "+" + strconv.FormatInt(int64(c), 10)
}

// Type implements Typer interface
func (c CallCode) Type() string {
	return TypeCallCode
}

// Info - returns CallCodeInfo
func (c CallCode) Info() *CallCodeInfo {
	return &CallCodeInfo{
		Code:      c,
		Countries: c.Countries(),
	}
}

// TotalCallCodes - returns number of call codes in the package, countries.TotalCallCodes() == len(countries.AllCallCodes()), but static value for performance
func TotalCallCodes() int {
	return 264
}

// AllCallCodes - return all countries call phone codes
func AllCallCodes() []CallCode {
	return []CallCode{
		CallCodeUnknown,
		CallCode1,
		CallCode7,
		CallCode20,
		CallCode27,
		CallCode30,
		CallCode31,
		CallCode32,
		CallCode33,
		CallCode34,
		CallCode36,
		CallCode38,
		CallCode39,
		CallCode40,
		CallCode41,
		CallCode43,
		CallCode44,
		CallCode45,
		CallCode46,
		CallCode47,
		CallCode48,
		CallCode49,
		CallCode51,
		CallCode52,
		CallCode53,
		CallCode54,
		CallCode55,
		CallCode56,
		CallCode57,
		CallCode58,
		CallCode60,
		CallCode61,
		CallCode62,
		CallCode63,
		CallCode64,
		CallCode65,
		CallCode66,
		CallCode81,
		CallCode82,
		CallCode84,
		CallCode86,
		CallCode90,
		CallCode91,
		CallCode92,
		CallCode93,
		CallCode94,
		CallCode95,
		CallCode98,
		CallCode211,
		CallCode212,
		CallCode213,
		CallCode216,
		CallCode218,
		CallCode220,
		CallCode221,
		CallCode222,
		CallCode223,
		CallCode224,
		CallCode225,
		CallCode226,
		CallCode227,
		CallCode228,
		CallCode229,
		CallCode230,
		CallCode231,
		CallCode232,
		CallCode233,
		CallCode234,
		CallCode235,
		CallCode236,
		CallCode237,
		CallCode238,
		CallCode239,
		CallCode240,
		CallCode241,
		CallCode242,
		CallCode243,
		CallCode244,
		CallCode245,
		CallCode246,
		CallCode248,
		CallCode249,
		CallCode250,
		CallCode251,
		CallCode252,
		CallCode253,
		CallCode254,
		CallCode255,
		CallCode256,
		CallCode257,
		CallCode258,
		CallCode260,
		CallCode261,
		CallCode262,
		CallCode263,
		CallCode264,
		CallCode265,
		CallCode266,
		CallCode267,
		CallCode268,
		CallCode269,
		CallCode290,
		CallCode291,
		CallCode297,
		CallCode298,
		CallCode299,
		CallCode350,
		CallCode351,
		CallCode352,
		CallCode353,
		CallCode354,
		CallCode355,
		CallCode356,
		CallCode357,
		CallCode358,
		CallCode359,
		CallCode370,
		CallCode371,
		CallCode372,
		CallCode373,
		CallCode374,
		CallCode375,
		CallCode376,
		CallCode377,
		CallCode378,
		CallCode380,
		CallCode381,
		CallCode382,
		CallCode383,
		CallCode385,
		CallCode386,
		CallCode387,
		CallCode389,
		CallCode420,
		CallCode421,
		CallCode423,
		CallCode500,
		CallCode501,
		CallCode502,
		CallCode503,
		CallCode504,
		CallCode505,
		CallCode506,
		CallCode507,
		CallCode508,
		CallCode509,
		CallCode590,
		CallCode591,
		CallCode592,
		CallCode593,
		CallCode594,
		CallCode595,
		CallCode596,
		CallCode597,
		CallCode598,
		CallCode599,
		CallCode670,
		CallCode672,
		CallCode673,
		CallCode674,
		CallCode675,
		CallCode676,
		CallCode677,
		CallCode678,
		CallCode679,
		CallCode680,
		CallCode681,
		CallCode682,
		CallCode683,
		CallCode685,
		CallCode686,
		CallCode687,
		CallCode688,
		CallCode689,
		CallCode690,
		CallCode691,
		CallCode692,
		CallCode800,
		CallCode850,
		CallCode852,
		CallCode853,
		CallCode855,
		CallCode856,
		CallCode870,
		CallCode875,
		CallCode876,
		CallCode877,
		CallCode878,
		CallCode879,
		CallCode880,
		CallCode881,
		CallCode882,
		CallCode883,
		CallCode886,
		CallCode888,
		CallCode960,
		CallCode961,
		CallCode962,
		CallCode963,
		CallCode964,
		CallCode965,
		CallCode966,
		CallCode967,
		CallCode968,
		CallCode970,
		CallCode971,
		CallCode972,
		CallCode973,
		CallCode974,
		CallCode975,
		CallCode976,
		CallCode977,
		CallCode979,
		CallCode991,
		CallCode992,
		CallCode994,
		CallCode995,
		CallCode996,
		CallCode998,
		CallCode1242,
		CallCode1246,
		CallCode1264,
		CallCode1268,
		CallCode1284,
		CallCode1340,
		CallCode1345,
		CallCode1441,
		CallCode1473,
		CallCode1649,
		CallCode1658,
		CallCode1664,
		CallCode1670,
		CallCode1671,
		CallCode1684,
		CallCode1721,
		CallCode1758,
		CallCode1767,
		CallCode1784,
		CallCode1787,
		CallCode1808,
		CallCode1809,
		CallCode1829,
		CallCode1849,
		CallCode1868,
		CallCode1869,
		CallCode1876,
		CallCode1939,
		CallCode4779,
		CallCode5993,
		CallCode5994,
		CallCode5995,
		CallCode5997,
		CallCode5998,
		CallCode5999,
		CallCode993,
		CallCode35818,
		CallCode262269,
		CallCode262639,
		CallCode441481,
		CallCode441534,
		CallCode441624,
		CallCode3906698,
		CallCode6189162,
		CallCode6189164,
	}
}

// AllCallCodesInfo - return all countries call phone codes as []CallCodeInfo
func AllCallCodesInfo() []*CallCodeInfo {
	all := AllCallCodes()
	codes := make([]*CallCodeInfo, 0, len(all))
	for _, v := range all {
		codes = append(codes, v.Info())
	}
	return codes
}

// Countries - returns CountryCodes of CallCode
//nolint:gocyclo
func (c CallCode) Countries() []CountryCode { //nolint:gocyclo
	switch c {
	case CallCode1:
		return []CountryCode{ATF, CAN, UMI, USA}
	case CallCode1808:
		return []CountryCode{USA}
	case CallCode1242:
		return []CountryCode{BHS}
	case CallCode1246:
		return []CountryCode{BRB}
	case CallCode1264:
		return []CountryCode{AIA}
	case CallCode1268:
		return []CountryCode{ATG}
	case CallCode1284:
		return []CountryCode{VGB}
	case CallCode1340:
		return []CountryCode{VIR}
	case CallCode1345:
		return []CountryCode{CYM}
	case CallCode1441:
		return []CountryCode{BMU}
	case CallCode1473:
		return []CountryCode{GRD}
	case CallCode1649:
		return []CountryCode{TCA}
	case CallCode1664:
		return []CountryCode{MSR}
	case CallCode1670:
		return []CountryCode{MNP}
	case CallCode1671:
		return []CountryCode{GUM}
	case CallCode1684:
		return []CountryCode{ASM}
	case CallCode1758:
		return []CountryCode{LCA}
	case CallCode1767:
		return []CountryCode{DMA}
	case CallCode1784:
		return []CountryCode{VCT}
	case CallCode1787, CallCode1939:
		return []CountryCode{PRI}
	case CallCode1809, CallCode1829, CallCode1849:
		return []CountryCode{DOM}
	case CallCode1868:
		return []CountryCode{TTO}
	case CallCode1869:
		return []CountryCode{KNA}
	case CallCode1876, CallCode1658:
		return []CountryCode{JAM}
	case CallCode20:
		return []CountryCode{EGY}
	case CallCode211:
		return []CountryCode{SSD}
	case CallCode212:
		return []CountryCode{ESH, MAR}
	case CallCode213:
		return []CountryCode{DZA}
	case CallCode216:
		return []CountryCode{TUN}
	case CallCode218:
		return []CountryCode{LBY}
	case CallCode220:
		return []CountryCode{GMB}
	case CallCode221:
		return []CountryCode{SEN}
	case CallCode222:
		return []CountryCode{MRT}
	case CallCode223:
		return []CountryCode{MLI}
	case CallCode224:
		return []CountryCode{GIN}
	case CallCode225:
		return []CountryCode{CIV}
	case CallCode226:
		return []CountryCode{BFA}
	case CallCode227:
		return []CountryCode{NER}
	case CallCode228:
		return []CountryCode{TGO}
	case CallCode229:
		return []CountryCode{BEN}
	case CallCode230:
		return []CountryCode{MUS}
	case CallCode231:
		return []CountryCode{LBR}
	case CallCode232:
		return []CountryCode{SLE}
	case CallCode233:
		return []CountryCode{GHA}
	case CallCode234:
		return []CountryCode{NGA}
	case CallCode235:
		return []CountryCode{TCD}
	case CallCode236:
		return []CountryCode{CAF}
	case CallCode237:
		return []CountryCode{CMR}
	case CallCode238:
		return []CountryCode{CPV}
	case CallCode239:
		return []CountryCode{STP}
	case CallCode240:
		return []CountryCode{GNQ}
	case CallCode241:
		return []CountryCode{GAB}
	case CallCode242:
		return []CountryCode{COG}
	case CallCode243:
		return []CountryCode{COD}
	case CallCode244:
		return []CountryCode{AGO}
	case CallCode245:
		return []CountryCode{GNB}
	case CallCode246:
		return []CountryCode{IOT}
	case CallCode248:
		return []CountryCode{SYC}
	case CallCode249:
		return []CountryCode{SDN}
	case CallCode250:
		return []CountryCode{RWA}
	case CallCode251:
		return []CountryCode{ETH}
	case CallCode252:
		return []CountryCode{SOM}
	case CallCode253:
		return []CountryCode{DJI}
	case CallCode254:
		return []CountryCode{KEN}
	case CallCode255:
		return []CountryCode{TZA}
	case CallCode256:
		return []CountryCode{UGA}
	case CallCode257:
		return []CountryCode{BDI}
	case CallCode258:
		return []CountryCode{MOZ}
	case CallCode260:
		return []CountryCode{ZMB}
	case CallCode261:
		return []CountryCode{MDG}
	case CallCode262:
		return []CountryCode{MYT, REU}
	case CallCode262269, CallCode262639:
		return []CountryCode{MYT}
	case CallCode263:
		return []CountryCode{ZWE}
	case CallCode264:
		return []CountryCode{NAM}
	case CallCode265:
		return []CountryCode{MWI}
	case CallCode266:
		return []CountryCode{LSO}
	case CallCode267:
		return []CountryCode{BWA}
	case CallCode268:
		return []CountryCode{SWZ}
	case CallCode269:
		return []CountryCode{COM}
	case CallCode27:
		return []CountryCode{ZAF}
	case CallCode290:
		return []CountryCode{SHN}
	case CallCode291:
		return []CountryCode{ERI}
	case CallCode297, CallCode5998:
		return []CountryCode{ABW}
	case CallCode298:
		return []CountryCode{FRO}
	case CallCode299:
		return []CountryCode{GRL}
	case CallCode30:
		return []CountryCode{GRC}
	case CallCode31:
		return []CountryCode{NLD}
	case CallCode32:
		return []CountryCode{BEL}
	case CallCode33:
		return []CountryCode{FRA}
	case CallCode34:
		return []CountryCode{ESP}
	case CallCode350:
		return []CountryCode{GIB}
	case CallCode351:
		return []CountryCode{PRT}
	case CallCode352:
		return []CountryCode{LUX}
	case CallCode353:
		return []CountryCode{IRL}
	case CallCode354:
		return []CountryCode{ISL}
	case CallCode355:
		return []CountryCode{ALB}
	case CallCode356:
		return []CountryCode{MLT}
	case CallCode357:
		return []CountryCode{CYP}
	case CallCode358:
		return []CountryCode{ALA, FIN}
	case CallCode359:
		return []CountryCode{BGR}
	case CallCode36:
		return []CountryCode{HUN}
	case CallCode370:
		return []CountryCode{LTU}
	case CallCode371:
		return []CountryCode{LVA}
	case CallCode372:
		return []CountryCode{EST}
	case CallCode373:
		return []CountryCode{MDA}
	case CallCode374:
		return []CountryCode{ARM}
	case CallCode375:
		return []CountryCode{BLR}
	case CallCode376:
		return []CountryCode{AND}
	case CallCode377:
		return []CountryCode{MCO}
	case CallCode378:
		return []CountryCode{SMR}
	case CallCode38:
		return []CountryCode{YUG}
	case CallCode380:
		return []CountryCode{UKR}
	case CallCode381:
		return []CountryCode{SRB}
	case CallCode382:
		return []CountryCode{MNE}
	case CallCode383:
		return []CountryCode{XKX}
	case CallCode385:
		return []CountryCode{HRV}
	case CallCode386:
		return []CountryCode{SVN}
	case CallCode387:
		return []CountryCode{BIH}
	case CallCode389:
		return []CountryCode{MKD}
	case CallCode39:
		return []CountryCode{ITA, VAT}
	case CallCode40:
		return []CountryCode{ROU}
	case CallCode41:
		return []CountryCode{CHE}
	case CallCode420:
		return []CountryCode{CZE}
	case CallCode421:
		return []CountryCode{SVK}
	case CallCode423:
		return []CountryCode{LIE}
	case CallCode43:
		return []CountryCode{AUT}
	case CallCode44:
		return []CountryCode{GBR, GGY, IMN, JEY, XSC, XWA}
	case CallCode45:
		return []CountryCode{DNK}
	case CallCode46:
		return []CountryCode{SWE}
	case CallCode47:
		return []CountryCode{BVT, NOR, SJM}
	case CallCode48:
		return []CountryCode{POL}
	case CallCode49:
		return []CountryCode{DEU}
	case CallCode500:
		return []CountryCode{FLK, SGS}
	case CallCode501:
		return []CountryCode{BLZ}
	case CallCode502:
		return []CountryCode{GTM}
	case CallCode503:
		return []CountryCode{SLV}
	case CallCode504:
		return []CountryCode{HND}
	case CallCode505:
		return []CountryCode{NIC}
	case CallCode506:
		return []CountryCode{CRI}
	case CallCode507:
		return []CountryCode{PAN}
	case CallCode508:
		return []CountryCode{SPM}
	case CallCode509:
		return []CountryCode{HTI}
	case CallCode51:
		return []CountryCode{PER}
	case CallCode52:
		return []CountryCode{MEX}
	case CallCode53:
		return []CountryCode{CUB}
	case CallCode54:
		return []CountryCode{ARG}
	case CallCode55:
		return []CountryCode{BRA}
	case CallCode56:
		return []CountryCode{CHL}
	case CallCode57:
		return []CountryCode{COL}
	case CallCode58:
		return []CountryCode{VEN}
	case CallCode590:
		return []CountryCode{BLM, GLP, MAF}
	case CallCode591:
		return []CountryCode{BOL}
	case CallCode592:
		return []CountryCode{GUY}
	case CallCode593:
		return []CountryCode{ECU}
	case CallCode594:
		return []CountryCode{GUF}
	case CallCode595:
		return []CountryCode{PRY}
	case CallCode596:
		return []CountryCode{MTQ}
	case CallCode597:
		return []CountryCode{SUR}
	case CallCode598:
		return []CountryCode{URY}
	case CallCode599:
		return []CountryCode{ANT, BES, CUW}
	case CallCode60:
		return []CountryCode{MYS}
	case CallCode61:
		return []CountryCode{AUS, CXR, CCK}
	case CallCode62:
		return []CountryCode{IDN}
	case CallCode63:
		return []CountryCode{PHL}
	case CallCode64:
		return []CountryCode{NZL, PCN}
	case CallCode65:
		return []CountryCode{SGP}
	case CallCode66:
		return []CountryCode{THA}
	case CallCode670:
		return []CountryCode{TLS}
	case CallCode672:
		return []CountryCode{ATA, CCK, NFK}
	case CallCode673:
		return []CountryCode{BRN}
	case CallCode674:
		return []CountryCode{NRU}
	case CallCode675:
		return []CountryCode{PNG}
	case CallCode676:
		return []CountryCode{TON}
	case CallCode677:
		return []CountryCode{SLB}
	case CallCode678:
		return []CountryCode{VUT}
	case CallCode679:
		return []CountryCode{FJI}
	case CallCode680:
		return []CountryCode{PLW}
	case CallCode681:
		return []CountryCode{WLF}
	case CallCode682:
		return []CountryCode{COK}
	case CallCode683:
		return []CountryCode{NIU}
	case CallCode685:
		return []CountryCode{WSM}
	case CallCode686:
		return []CountryCode{KIR}
	case CallCode687:
		return []CountryCode{NCL}
	case CallCode688:
		return []CountryCode{TUV}
	case CallCode689:
		return []CountryCode{PYF}
	case CallCode690:
		return []CountryCode{TKL}
	case CallCode691:
		return []CountryCode{FSM}
	case CallCode692:
		return []CountryCode{MHL}
	case CallCode7:
		return []CountryCode{KAZ, RUS}
	case CallCode1721, CallCode5995:
		return []CountryCode{SXM}
	case CallCode4779:
		return []CountryCode{SJM}
	case CallCode5993, CallCode5994, CallCode5997:
		return []CountryCode{BES}
	case CallCode993:
		return []CountryCode{TKM}
	case CallCode81:
		return []CountryCode{JPN}
	case CallCode82:
		return []CountryCode{KOR}
	case CallCode84:
		return []CountryCode{VNM}
	case CallCode850:
		return []CountryCode{PRK}
	case CallCode852:
		return []CountryCode{HKG}
	case CallCode853:
		return []CountryCode{MAC}
	case CallCode855:
		return []CountryCode{KHM}
	case CallCode856:
		return []CountryCode{LAO}
	case CallCode86:
		return []CountryCode{CHN}
	case CallCode800:
		return []CountryCode{NonCountryInternationalFreephone}
	case CallCode870:
		return []CountryCode{NonCountryInmarsat}
	case CallCode875, CallCode876, CallCode877:
		return []CountryCode{NonCountryMaritimeMobileService}
	case CallCode878:
		return []CountryCode{NonCountryUniversalPersonalTelecommunicationsServices}
	case CallCode879:
		return []CountryCode{NonCountryNationalNonCommercialPurposes}
	case CallCode880:
		return []CountryCode{BGD}
	case CallCode881:
		return []CountryCode{NonCountryGlobalMobileSatelliteSystem}
	case CallCode882, CallCode883:
		return []CountryCode{NonCountryInternationalNetworks}
	case CallCode886:
		return []CountryCode{TWN}
	case CallCode90:
		return []CountryCode{TUR}
	case CallCode91:
		return []CountryCode{IND}
	case CallCode92:
		return []CountryCode{PAK}
	case CallCode93:
		return []CountryCode{AFG}
	case CallCode94:
		return []CountryCode{LKA}
	case CallCode95:
		return []CountryCode{MMR}
	case CallCode888:
		return []CountryCode{NonCountryDisasterRelief}
	case CallCode960:
		return []CountryCode{MDV}
	case CallCode961:
		return []CountryCode{LBN}
	case CallCode962:
		return []CountryCode{JOR}
	case CallCode963:
		return []CountryCode{SYR}
	case CallCode964:
		return []CountryCode{IRQ}
	case CallCode965:
		return []CountryCode{KWT}
	case CallCode966:
		return []CountryCode{SAU}
	case CallCode967:
		return []CountryCode{YEM}
	case CallCode968:
		return []CountryCode{OMN}
	case CallCode970:
		return []CountryCode{PSE}
	case CallCode971:
		return []CountryCode{ARE}
	case CallCode972:
		return []CountryCode{PSE}
	case CallCode973:
		return []CountryCode{BHR}
	case CallCode974:
		return []CountryCode{QAT}
	case CallCode975:
		return []CountryCode{BTN}
	case CallCode976:
		return []CountryCode{MNG}
	case CallCode977:
		return []CountryCode{NPL}
	case CallCode98:
		return []CountryCode{IRN}
	case CallCode979:
		return []CountryCode{NonCountryInternationalPremiumRateService}
	case CallCode991:
		return []CountryCode{NonCountryInternationalTelecommunicationsCorrespondenceService}
	case CallCode992:
		return []CountryCode{TJK}
	case CallCode994:
		return []CountryCode{AZE}
	case CallCode995:
		return []CountryCode{GEO}
	case CallCode996:
		return []CountryCode{KGZ}
	case CallCode998:
		return []CountryCode{UZB}
	case CallCode5999:
		return []CountryCode{CUW}
	case CallCode35818:
		return []CountryCode{ALA}
	case CallCode441481:
		return []CountryCode{GGY}
	case CallCode441534:
		return []CountryCode{JEY}
	case CallCode441624:
		return []CountryCode{IMN}
	case CallCode3906698:
		return []CountryCode{VAT}
	case CallCode6189162:
		return []CountryCode{CCK}
	case CallCode6189164:
		return []CountryCode{CXR}
	}
	return []CountryCode{Unknown}
}

// IsValid - returns true, if code is correct
func (c CallCode) IsValid() bool {
	return c.Countries()[0] != Unknown
}

// Type implements Typer interface
func (c *CallCodeInfo) Type() string {
	return TypeCallCodeInfo
}

// Value implements database/sql/driver.Valuer
func (c CallCodeInfo) Value() (Value, error) {
	return json.Marshal(c)
}

// Scan implements database/sql.Scanner
func (c *CallCodeInfo) Scan(src interface{}) error {
	if c == nil {
		return fmt.Errorf("countries::Scan: CallCodeInfo scan err: callCodeInfo == nil")
	}
	switch src := src.(type) {
	case *CallCodeInfo:
		*c = *src
	case CallCodeInfo:
		*c = src
	default:
		return fmt.Errorf("countries::Scan: CallCodeInfo scan err: unexpected value of type %T for %T", src, *c)
	}
	return nil
}
