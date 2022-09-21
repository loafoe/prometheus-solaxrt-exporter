package inverter

import (
	"github.com/loafoe/prometheus-solaxrt-exporter/solax/inverter/fields"
	"github.com/loafoe/prometheus-solaxrt-exporter/solax/inverter/units"
)

type X1BoostAirMini struct {
	Type        string    `json:"type"`
	SN          string    `json:"SN"`
	Ver         string    `json:"ver"`
	Data        []float64 `json:"Data"`
	Information []any     `json:"Information"`
}

type IndexUnit struct {
	Index int
	Unit  string
}

type Decoder map[string]IndexUnit

func (x X1BoostAirMini) Field(field string) float64 {
	f, ok := x.decode()[field]
	if !ok {
		return 0
	}
	return x.Data[f.Index]
}

func (x X1BoostAirMini) decode() Decoder {
	return Decoder{
		fields.PV1_Current:          {0, units.A},
		fields.PV2_Current:          {1, units.A},
		fields.PV1_Voltage:          {2, units.V},
		fields.PV2_Voltage:          {3, units.V},
		fields.Output_Current:       {4, units.A},
		fields.Network_Voltage:      {5, units.V},
		fields.AC_Power:             {6, units.W},
		fields.Inverter_Temperature: {7, units.C},
		fields.Todays_Energy:        {8, units.KWH},
		fields.Total_Energy:         {9, units.KWH},
		fields.Exported_Power:       {10, units.W},
		fields.PV1_Power:            {11, units.W},
		fields.PV2_Power:            {12, units.W},
		fields.Total_FeedIn_Energy:  {41, units.KWH},
		fields.Total_Consumption:    {42, units.KWH},
		fields.Power_Now:            {43, units.W},
		fields.Grid_Frequency:       {50, units.HZ},
	}
}
