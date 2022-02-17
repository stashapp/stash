package countries

import (
	"encoding/json"
	"fmt"
)

// RegionCode - Region code (UN M.49 code standart)
type RegionCode int64 // int64 for database/sql/driver.Valuer compatibility

// Region - all info about region
type Region struct {
	Name string
	Code RegionCode
}

// Type implements Typer interface
func (c RegionCode) Type() string {
	return TypeRegionCode
}

// String - implements fmt.Stringer, returns a Region name in english
func (c RegionCode) String() string {
	switch c {
	case RegionNone:
		return "None"
	case RegionAF:
		return "Africa"
	case RegionNA:
		return "North America"
	case RegionOC:
		return "Oceania"
	case RegionAN:
		return "Antarctica"
	case RegionAS:
		return "Asia"
	case RegionEU:
		return "Europe"
	case RegionSA:
		return "South America"
	}
	return UnknownMsg
}

// StringRus - returns a Region name in russian
func (c RegionCode) StringRus() string {
	switch c {
	case RegionNone:
		return "Отсутствует"
	case RegionAF:
		return "Африка"
	case RegionNA:
		return "Северная Америка"
	case RegionOC:
		return "Океания"
	case RegionAN:
		return "Антарктика"
	case RegionAS:
		return "Азия"
	case RegionEU:
		return "Европа"
	case RegionSA:
		return "Южная Америка"
	}
	return UnknownMsg
}

// IsValid - returns true, if code is correct
func (c RegionCode) IsValid() bool {
	return c.String() != UnknownMsg
}

// TotalRegions - returns number of Regions codes in the package
func TotalRegions() int {
	return 7
}

// Info - return a RegionCode as Region info
func (c RegionCode) Info() *Region {
	return &Region{
		Name: c.String(),
		Code: c,
	}
}

// Type implements Typer interface
func (r *Region) Type() string {
	return TypeRegion
}

// Value implements database/sql/driver.Valuer
func (r Region) Value() (Value, error) {
	return json.Marshal(r)
}

// Scan implements database/sql.Scanner
func (r *Region) Scan(src interface{}) error {
	if r == nil {
		return fmt.Errorf("countries::Scan: Region scan err: region == nil")
	}
	switch src := src.(type) {
	case *Region:
		*r = *src
	case Region:
		*r = src
	default:
		return fmt.Errorf("countries::Scan: Region scan err: unexpected value of type %T for %T", src, *r)
	}
	return nil
}

// AllRegions - returns all regions codes
func AllRegions() []RegionCode {
	return []RegionCode{
		RegionAF,
		RegionNA,
		RegionOC,
		RegionAN,
		RegionAS,
		RegionEU,
		RegionSA,
	}
}

// AllRegionsInfo - return all regions as []*Region
func AllRegionsInfo() []*Region {
	all := AllRegions()
	regions := make([]*Region, 0, len(all))
	for _, v := range all {
		regions = append(regions, v.Info())
	}
	return regions
}

// RegionCodeByName - return RegionCode by region name, case-insensitive, example: regionEU := RegionCodeByName("eu") OR regionEU := RegionCodeByName("europe")
func RegionCodeByName(name string) RegionCode {
	switch textPrepare(name) {
	case "AF", "AFRICA", "AFRIKA":
		return RegionAF
	case "NA", "NORTHAMERICA", "NORTHAMERIC":
		return RegionNA
	case "SA", "SOUTHAMERICA", "SOUTHAMERIC":
		return RegionSA
	case "OC", "OCEANIA", "OKEANIA", "OCEANIYA", "OKEANIYA":
		return RegionOC
	case "AN", "ANTARCTICA", "ANTARCTIC", "ANTARKTICA", "ANTARKTIC":
		return RegionAN
	case "AS", "ASIA":
		return RegionAS
	case "EU", "EUROPE", "EUROPA", "EVROPA":
		return RegionEU
	case "NONE", "XX", "NON":
		return RegionNone
	}
	return RegionUnknown
}
