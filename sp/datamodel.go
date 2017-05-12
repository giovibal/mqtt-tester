package sp

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Snapshot struct {
	Ref            string        `json:"ref,omitempty"`
	Address        string        `json:"address,omitempty"`
	Cat            string        `json:"cat,omitempty"`
	Ctx            string        `json:"ctx,omitempty"`
	Type           string        `json:"type,omitempty"`
	T              int64         `json:"t,omitempty"`
	Tz             string        `json:"tz,omitempty"`
	Msgcode        string        `json:"msgcode,omitempty"`
	Level          int32         `json:"level,omitempty"`
	Cause          string        `json:"cause,omitempty"`
	Msg            string        `json:"msg,omitempty"`
	Is             *bool         `json:"is,omitempty"`
	Weight         int32         `json:"weight,omitempty"`
	Correlation_id string        `json:"correlation_id,omitempty"`
	M              []Measurement `json:"m,omitempty"`
	R              []Relation    `json:"r,omitempty"`
}

type Measurement struct {
	K  string  `json:"k,omitempty"`
	T  int64   `json:"t,omitempty"`
	V  float64 `json:"v,omitempty"`
	U  string  `json:"u,omitempty"`
	Tz string  `json:"tz,omitempty"`
	X  string  `json:"x,omitempty"`
}

type Relation struct {
	Measurement
	C string `json:"c,omitempty"`
	S string `json:"s,omitempty"`
	D string `json:"d,omitempty"`
}

func (m Measurement) ToString() string {
	s := fmt.Sprintf("%s = %g", m.K, m.V)
	return s
}

func (r Relation) ToString() string {
	s := fmt.Sprintf("(%s)---[%s:%s]--->(%s) = %g", r.S, r.C, r.K, r.D, r.V)
	return s
}

func FromJson(b []byte) (*Snapshot, error) {
	// json to snapshot
	var s Snapshot
	err := json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func ToJsonPretty(s Snapshot) []byte {
	output, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return output
}
func ToJson(s Snapshot) []byte {
	output, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return output
}

func (s *Snapshot) SetTimestamp(ts time.Time) {
	s.T = ts.UnixNano() / int64(time.Millisecond)
	s.Tz = ts.Format("2006-01-02T15:04:05.000Z07:00")
	for i := 0; i < len(s.M); i++ {
		s.M[i].T = s.T
		s.M[i].Tz = s.Tz
	}
	for i := 0; i < len(s.R); i++ {
		s.R[i].T = s.T
		s.R[i].Tz = s.Tz
	}
}
