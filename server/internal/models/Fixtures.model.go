package models

type Fixtures struct {
	Get        string      `json:"get,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Results    int32       `json:"results,omitempty"`
	Paging     struct {
		Current int32 `json:"current,omitempty"`
		Total   int32 `json:"total,omitempty"`
	} `json:"paging,omitempty"`
	Response []struct {
		Fixture struct {
			Id        int32       `json:"id,omitempty"`
			Referee   string      `json:"refereed,omitempty"`
			Timezone  string      `json:"timezone,omitempty"`
			Date      string      `json:"date,omitempty"`
			Timestamp int32       `json:"timestamp,omitempty"`
			Periods   interface{} `json:"periods,omitempty"`
			Venue     interface{} `json:"venue,omitempty"`
			Status    interface{} `json:"status,omitempty"`
		}
		League    interface{} `json:"league,omitempty"`
		Teams     interface{} `json:"teams,omitempty"`
		Goals     interface{} `json:"goals,omitempty"`
		Score     interface{} `json:"score,omitempty"`
		Fulltime  interface{} `json:"fulltime,omitempty"`
		Extratime interface{} `json:"extratime,omitempty"`
		Penalty   interface{} `json:"penalty,omitempty"`
	}
}
