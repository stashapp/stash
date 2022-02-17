package countries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DomainCode - domain code
type DomainCode int64 // int64 for database/sql/driver.Valuer compatibility

// Domain - capital info
type Domain struct {
	Name    string
	Code    DomainCode
	Country CountryCode
}

// Type implements Typer interface
func (c DomainCode) Type() string {
	return TypeDomainCode
}

// String - implements fmt.Stringer, returns a domain (internet ccTDL)
//nolint:gocyclo
func (c DomainCode) String() string { //nolint:gocyclo
	switch c {
	case DomainArpa:
		return ".arpa"
	case DomainCom:
		return ".com"
	case DomainOrg:
		return ".org"
	case DomainNet:
		return ".net"
	case DomainEdu:
		return ".edu"
	case DomainGov:
		return ".gov"
	case DomainMil:
		return ".mil"
	case DomainTest:
		return ".test"
	case DomainBiz:
		return ".biz"
	case DomainInfo:
		return ".info"
	case DomainName:
		return ".name"
	case DomainBV, DomainSJ:
		return ".no"
	case DomainGB:
		return ".uk"
	case DomainXX:
		return ""
	}

	if c >= 999 {
		c = DomainCode(999)
	}

	a2 := CountryCode(c).Alpha2()
	if a2 == UnknownMsg {
		return UnknownMsg
	}

	return "." + strings.ToLower(a2)
}

// IsValid - returns true, if code is correct
func (c DomainCode) IsValid() bool {
	return c.String() != UnknownMsg
}

// Country - returns a country of domain
func (c DomainCode) Country() CountryCode {
	if !c.IsValid() {
		return Unknown
	}
	return CountryCode(c)
}

// Info - returns domain information as Domain
func (c DomainCode) Info() *Domain {
	return &Domain{
		Name:    c.String(),
		Code:    c,
		Country: c.Country(),
	}
}

// Type implements Typer interface
func (c Domain) Type() string {
	return TypeDomain
}

// Value implements database/sql/driver.Valuer
func (c Domain) Value() (Value, error) {
	return json.Marshal(c)
}

// Scan implements database/sql.Scanner
func (c *Domain) Scan(src interface{}) error {
	if c == nil {
		return fmt.Errorf("countries::Scan: Domain scan err: domain == nil")
	}
	switch src := src.(type) {
	case *Domain:
		*c = *src
	case Domain:
		*c = src
	default:
		return fmt.Errorf("countries::Scan: domain scan err: unexpected value of type %T for %T", src, *c)
	}
	return nil
}

// DomainCodeByName - return DomainCode by name, case-insensitive, example: domainAE := DomainCodeByName(".ae") OR capitalAE := domainAE("ae")
func DomainCodeByName(name string) DomainCode {
	country := ByName(name)
	if country == Unknown {
		return DomainUnknown
	}
	return DomainCode(country)
}

// AllDomains - returns all domains codes
func AllDomains() []DomainCode {
	return []DomainCode{
		DomainArpa,
		DomainCom,
		DomainOrg,
		DomainNet,
		DomainEdu,
		DomainGov,
		DomainMil,
		DomainTest,
		DomainBiz,
		DomainInfo,
		DomainName,
		DomainAU,
		DomainAT,
		DomainAZ,
		DomainAL,
		DomainDZ,
		DomainAS,
		DomainAI,
		DomainAO,
		DomainAD,
		DomainAQ,
		DomainAG,
		DomainAN,
		DomainAE,
		DomainAR,
		DomainAM,
		DomainAW,
		DomainAF,
		DomainBS,
		DomainBD,
		DomainBB,
		DomainBH,
		DomainBY,
		DomainBZ,
		DomainBE,
		DomainBJ,
		DomainBM,
		DomainBG,
		DomainBO,
		DomainBA,
		DomainBW,
		DomainBR,
		DomainIO,
		DomainBN,
		DomainBF,
		DomainBI,
		DomainBT,
		DomainVU,
		DomainVA,
		DomainGB,
		DomainHU,
		DomainVE,
		DomainVG,
		DomainVI,
		DomainTL,
		DomainVN,
		DomainGA,
		DomainHT,
		DomainGY,
		DomainGM,
		DomainGH,
		DomainGP,
		DomainGT,
		DomainGN,
		DomainGW,
		DomainDE,
		DomainGI,
		DomainHN,
		DomainHK,
		DomainGD,
		DomainGL,
		DomainGR,
		DomainGE,
		DomainGU,
		DomainDK,
		DomainCD,
		DomainDJ,
		DomainDM,
		DomainDO,
		DomainEG,
		DomainZM,
		DomainEH,
		DomainZW,
		DomainIL,
		DomainIN,
		DomainID,
		DomainJO,
		DomainIQ,
		DomainIR,
		DomainIE,
		DomainIS,
		DomainES,
		DomainIT,
		DomainYE,
		DomainKZ,
		DomainKY,
		DomainKH,
		DomainCM,
		DomainCA,
		DomainQA,
		DomainKE,
		DomainCY,
		DomainKI,
		DomainCN,
		DomainCC,
		DomainCO,
		DomainKM,
		DomainCG,
		DomainKP,
		DomainKR,
		DomainCR,
		DomainCI,
		DomainCU,
		DomainKW,
		DomainKG,
		DomainLA,
		DomainLV,
		DomainLS,
		DomainLR,
		DomainLB,
		DomainLY,
		DomainLT,
		DomainLI,
		DomainLU,
		DomainMU,
		DomainMR,
		DomainMG,
		DomainYT,
		DomainMO,
		DomainMK,
		DomainMW,
		DomainMY,
		DomainML,
		DomainMV,
		DomainMT,
		DomainMP,
		DomainMA,
		DomainMQ,
		DomainMH,
		DomainMX,
		DomainFM,
		DomainMZ,
		DomainMD,
		DomainMC,
		DomainMN,
		DomainMS,
		DomainMM,
		DomainNA,
		DomainNR,
		DomainNP,
		DomainNE,
		DomainNG,
		DomainNL,
		DomainNI,
		DomainNU,
		DomainNZ,
		DomainNC,
		DomainNO,
		DomainOM,
		DomainBV,
		DomainIM,
		DomainNF,
		DomainPN,
		DomainCX,
		DomainSH,
		DomainWF,
		DomainHM,
		DomainCV,
		DomainCK,
		DomainWS,
		DomainSJ,
		DomainTC,
		DomainUM,
		DomainPK,
		DomainPW,
		DomainPS,
		DomainPA,
		DomainPG,
		DomainPY,
		DomainPE,
		DomainPL,
		DomainPT,
		DomainPR,
		DomainRE,
		DomainRU,
		DomainRW,
		DomainRO,
		DomainSV,
		DomainSM,
		DomainST,
		DomainSA,
		DomainSZ,
		DomainSC,
		DomainSN,
		DomainPM,
		DomainVC,
		DomainKN,
		DomainLC,
		DomainSG,
		DomainSY,
		DomainSK,
		DomainSI,
		DomainUS,
		DomainSB,
		DomainSO,
		DomainSD,
		DomainSR,
		DomainSL,
		DomainTJ,
		DomainTW,
		DomainTH,
		DomainTZ,
		DomainTG,
		DomainTK,
		DomainTO,
		DomainTT,
		DomainTV,
		DomainTN,
		DomainTM,
		DomainTR,
		DomainUG,
		DomainUZ,
		DomainUA,
		DomainUY,
		DomainFO,
		DomainFJ,
		DomainPH,
		DomainFI,
		DomainFK,
		DomainFR,
		DomainGF,
		DomainPF,
		DomainTF,
		DomainHR,
		DomainCF,
		DomainTD,
		DomainCZ,
		DomainCL,
		DomainCH,
		DomainSE,
		DomainXS,
		DomainLK,
		DomainEC,
		DomainGQ,
		DomainER,
		DomainEE,
		DomainET,
		DomainZA,
		DomainYU,
		DomainGS,
		DomainJM,
		DomainME,
		DomainBL,
		DomainSX,
		DomainRS,
		DomainAX,
		DomainBQ,
		DomainGG,
		DomainJE,
		DomainCW,
		DomainMF,
		DomainSS,
		DomainJP,
	}
}

// AllDomainsInfo - return all domains as []*Domain
func AllDomainsInfo() []*Domain {
	all := AllDomains()
	domains := make([]*Domain, 0, len(all))
	for _, v := range all {
		domains = append(domains, v.Info())
	}
	return domains
}

// TotalDomains - returns number of domains in the package, countries.TotalDomains() == len(countries.AllDomains()) but static value for performance
func TotalDomains() int {
	return 263
}
