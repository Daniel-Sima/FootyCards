package models

type RequestJSON struct {
	Get        string      `json:"get,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Results    int32       `json:"results,omitempty"`
	Paging     struct {
		Current int32 `json:"current,omitempty"`
		Total   int32 `json:"total,omitempty"`
	} `json:"paging,omitempty"`
	Response []struct {
		Player struct {
			Id        int32  `json:"id,omitempty"`
			Name      string `json:"name,omitempty"`
			Firstname string `json:"firstname,omitempty"`
			Lastname  string `json:"lastname,omitempty"`
			Age       int32  `json:"age,omitempty"`
			Birth     struct {
				Date    string `json:"date,omitempty"`
				Place   string `json:"place,omitempty"`
				Country string `json:"country,omitempty"`
			} `json:"birth,omitempty"`
			Nationality string `json:"nationality,omitempty"`
			Height      string `json:"height,omitempty"`
			Weight      string `json:"weight,omitempty"`
			Injured     bool   `json:"injured,omitempty"`
			Photo       string `json:"photo,omitempty"`
		} `json:"player,omitempty"`
		Statistics []struct {
			Team struct {
				Id   int32  `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
				Logo string `json:"logo,omitempty"`
			} `json:"team,omitempty"`
			League struct {
				Id      int32  `json:"id,omitempty"`
				Name    string `json:"name,omitempty"`
				Country string `json:"country,omitempty"`
				Logo    string `json:"logo,omitempty"`
				Flag    string `json:"flag,omitempty"`
				Season  int32  `json:"season,omitempty"`
			} `json:"league,omitempty"`
			Games struct {
				Appearances int32  `json:"appearances,omitempty"`
				Lineups     int32  `json:"lineups,omitempty"`
				Minutes     int32  `json:"minutes,omitempty"`
				Number      int32  `json:"number,omitempty"`
				Position    string `json:"position,omitempty"`
				Rating      string `json:"rating,omitempty"`
				Captain     bool   `json:"captain,omitempty"`
			} `json:"games,omitempty"`
			Substitutes struct {
				In    int32 `json:"in,omitempty"`
				Out   int32 `json:"out,omitempty"`
				Bench int32 `json:"bench,omitempty"`
			} `json:"substitutes,omitempty"`
			Shots struct {
				Total int32 `json:"total,omitempty"`
				On    int32 `json:"on,omitempty"`
			} `json:"shots,omitempty"`
			Goals struct {
				Total    int32 `json:"total,omitempty"`
				Conceded int32 `json:"conceded,omitempty"`
				Assists  int32 `json:"assists,omitempty"`
				Saves    int32 `json:"saves,omitempty"`
			} `json:"goals,omitempty"`
			Passes struct {
				Total     int32 `json:"total,omitempty"`
				Key       int32 `json:"key,omitempty"`
				Accurency int32 `json:"accurency,omitempty"`
			} `json:"passes,omitempty"`
			Tackles struct {
				Total         int32 `json:"total,omitempty"`
				Blocks        int32 `json:"blocks,omitempty"`
				Interceptions int32 `json:"interceptions,omitempty"`
			} `json:"tackles,omitempty"`
			Duels struct {
				Total int32 `json:"total,omitempty"`
				Won   int32 `json:"won,omitempty"`
			} `json:"duels,omitempty"`
			Dribbles struct {
				Attempts int32 `json:"attempts,omitempty"`
				Success  int32 `json:"success,omitempty"`
				Past     int32 `json:"past,omitempty"`
			} `json:"dribbles,omitempty"`
			Fouls struct {
				Drawn    int32 `json:"drawn,omitempty"`
				Commited int32 `json:"commited,omitempty"`
			} `json:"fouls,omitempty"`
			Cards struct {
				Yellow    int32 `json:"yellow,omitempty"`
				Yellowred int32 `json:"yellowred,omitempty"`
				Red       int32 `json:"red,omitempty"`
			} `json:"cards,omitempty"`
			Penalty struct {
				Won      int32 `json:"won,omitempty"`
				Commited int32 `json:"commited,omitempty"`
				Scored   int32 `json:"scored,omitempty"`
				Missed   int32 `json:"missed,omitempty"`
				Saved    int32 `json:"saved,omitempty"`
			} `json:"penalty,omitempty"`
		}
	} `json:"response,omitempty"`
}
