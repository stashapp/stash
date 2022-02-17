package countries

// TypeRegionCode for Typer interface
const TypeRegionCode string = "countries.RegionCode"

// TypeRegion for Typer interface
const TypeRegion string = "countries.Region"

const (
	// RegionUnknown RegionCode = 0
	RegionUnknown RegionCode = 0
	// RegionAF      RegionCode = 2
	RegionAF RegionCode = 2
	// RegionNA      RegionCode = 3
	RegionNA RegionCode = 3
	// RegionSA      RegionCode = 5
	RegionSA RegionCode = 5
	// RegionOC      RegionCode = 9
	RegionOC RegionCode = 9
	// RegionAN      RegionCode = 998
	RegionAN RegionCode = RegionCode(AN)
	// RegionAS      RegionCode = 142
	RegionAS RegionCode = 142
	// RegionEU      RegionCode = 150
	RegionEU RegionCode = 150
	// RegionNone    RegionCode = RegionCode(None)
	RegionNone RegionCode = RegionCode(None)
)

const (
	// RegionAfrica       RegionCode = 2
	RegionAfrica RegionCode = 2
	// RegionNorthAmerica RegionCode = 3
	RegionNorthAmerica RegionCode = 3
	// RegionSouthAmerica RegionCode = 5
	RegionSouthAmerica RegionCode = 5
	// RegionOceania      RegionCode = 9
	RegionOceania RegionCode = 9
	// RegionAntarctica   RegionCode = 999
	RegionAntarctica RegionCode = 999
	// RegionAsia         RegionCode = 142
	RegionAsia RegionCode = 142
	// RegionEurope       RegionCode = 150
	RegionEurope RegionCode = 150
)
