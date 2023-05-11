package models

type FixturePlayer struct {
	Get        string
	Parameters interface{}
	Errors     interface{}
	Results    int32
	Paging     struct {
		Current int32
		Total   int32
	}
	Response []struct {
		Team    interface{}
		Players []struct {
			Player struct {
				Id    int32
				Name  string
				Photo string
			}
			Statistics []struct {
				Games struct {
					Minutes    int32
					Number     int32
					Position   string
					Rating     string
					Captain    bool
					Substitute bool
				}
				Offsides interface{}
				Shots    interface{}
				Goals    interface{}
				Passes   interface{}
				Tackles  interface{}
				Duels    interface{}
				Dribbles interface{}
				Fouls    interface{}
				Cards    interface{}
				Penalty  interface{}
			}
		}
	}
}
