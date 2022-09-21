package solax

import (
	"encoding/json"
)

type Response struct {
	Success   bool   `json:"success"`
	Exception string `json:"exception"`
	Result    Result `json:"result"`
}

type Result result

type result struct {
	Error          string  `json:"-"`
	InverterSN     string  `json:"inverterSN"`
	SN             string  `json:"sn"`
	ACPower        float64 `json:"acpower"`
	YieldToday     float64 `json:"yieldtoday"`
	YieldTotal     float64 `json:"yieldtotal"`
	FeedInPower    float64 `json:"feedinpower"`
	FeedInEnergy   float64 `json:"feedinenergy"`
	ConsumeEnergy  float64 `json:"consumeenergy"`
	FeedInPowerM2  float64 `json:"feedinpowerM2"`
	Soc            float64 `json:"soc"`
	Peps1          float64 `json:"peps1"`
	Peps2          float64 `json:"peps2"`
	Peps3          float64 `json:"peps3"`
	InverterType   string  `json:"inverterType"`
	InverterStatus string  `json:"inverterStatus"`
	UploadTime     string  `json:"uploadTime"`
	BatPower       float64 `json:"batPower"`
	PowerDC1       float64 `json:"powerdc1,omitempty"`
	PowerDC2       float64 `json:"powerdc2,omitempty"`
	PowerDC3       float64 `json:"powerdc3,omitempty"`
	PowerDC4       float64 `json:"powerdc4,omitempty"`
	BatStatus      string  `json:"batStatus"`
}

func (r *Result) UnmarshalJSON(data []byte) (err error) {
	var str string
	var res result
	err = json.Unmarshal(data, &str)
	if err == nil {
		r.Error = str
		return nil
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	r.BatPower = res.BatPower
	r.InverterType = res.InverterType
	r.SN = res.SN
	r.BatStatus = res.BatStatus
	r.ACPower = res.ACPower
	r.ConsumeEnergy = res.ConsumeEnergy
	r.InverterSN = res.InverterSN
	r.FeedInEnergy = res.FeedInEnergy
	r.FeedInPower = res.FeedInPower
	r.FeedInPowerM2 = res.FeedInPowerM2
	r.Peps1 = res.Peps1
	r.Peps2 = res.Peps2
	r.Peps3 = res.Peps3
	r.PowerDC1 = res.PowerDC1
	r.PowerDC2 = res.PowerDC2
	r.PowerDC3 = res.PowerDC3
	r.PowerDC4 = res.PowerDC4
	r.Soc = res.Soc
	r.YieldTotal = res.YieldTotal
	r.YieldToday = res.YieldToday
	r.InverterStatus = res.InverterStatus
	r.UploadTime = res.UploadTime
	return nil
}
