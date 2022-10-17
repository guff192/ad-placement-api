package placement

type Tile struct {
	Id    uint    `json:"id"`
	Width uint    `json:"width"`
	Ratio float32 `json:"ratio"`
}

type Context struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type PlacementRequest struct {
	Id      string `json:"id"`
	Tiles   []Tile `json:"tiles"`
	Context `json:"context"`
}
