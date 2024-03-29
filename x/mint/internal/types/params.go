package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyMintDenom           = []byte("MintDenom")
	KeyInflationRateChange = []byte("InflationRateChange")
	KeyInflationMax        = []byte("InflationMax")
	KeyInflationMin        = []byte("InflationMin")
	KeyGoalBonded          = []byte("GoalBonded")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
)

// mint parameters
type Params struct {
	MintDenom           string  `json:"mint_denom" yaml:"mint_denom"`                       // type of coin to mint
	InflationRateChange sdk.Dec `json:"inflation_rate_change" yaml:"inflation_rate_change"` // maximum annual change in inflation rate
	InflationMax        sdk.Dec `json:"inflation_max" yaml:"inflation_max"`                 // maximum inflation rate
	InflationMin        sdk.Dec `json:"inflation_min" yaml:"inflation_min"`                 // minimum inflation rate
	GoalBonded          sdk.Dec `json:"goal_bonded" yaml:"goal_bonded"`                     // goal of percent bonded atoms
	BlocksPerYear       uint64  `json:"blocks_per_year" yaml:"blocks_per_year"`             // expected blocks per year
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflationRateChange, inflationMax,
inflationMin, goalBonded sdk.Dec, blocksPerYear uint64) Params {

	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom,
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}
}

// validate params
func ValidateParams(params Params) error {
	if params.GoalBonded.IsNegative() {
		return fmt.Errorf("mint parameter GoalBonded should be positive, is %s ", params.GoalBonded.String())
	}
	if params.GoalBonded.GT(sdk.OneDec()) {
		return fmt.Errorf("mint parameter GoalBonded must be <= 1, is %s", params.GoalBonded.String())
	}
	if params.InflationMax.LT(params.InflationMin) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Inflation Rate Change:  %s
  Inflation Max:          %s
  Inflation Min:          %s
  Goal Bonded:            %s
  Blocks Per Year:        %d
`,
		p.MintDenom, p.InflationRateChange, p.InflationMax,
		p.InflationMin, p.GoalBonded, p.BlocksPerYear,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMintDenom, Value: &p.MintDenom},
		{Key: KeyInflationRateChange, Value: &p.InflationRateChange},
		{Key: KeyInflationMax, Value: &p.InflationMax},
		{Key: KeyInflationMin, Value: &p.InflationMin},
		{Key: KeyGoalBonded, Value: &p.GoalBonded},
		{Key: KeyBlocksPerYear, Value: &p.BlocksPerYear},
	}
}

//map[year]inflation
var FixupInflation = map[int]sdk.Dec{
	0: sdk.NewDecWithPrec(300, 4),
	1: sdk.NewDecWithPrec(250, 4),
	2: sdk.NewDecWithPrec(150, 4),
	3: sdk.NewDecWithPrec(125, 4),
	4: sdk.NewDecWithPrec(100, 4),
	5: sdk.NewDecWithPrec(75, 4),
	6: sdk.NewDecWithPrec(0, 4),
}

//unit: umdu
var FixupTotalTokenNumber = sdk.NewInt(200000000 * 1000000)
