package placement

type Imp struct {
	Id     uint   `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Price  string `json:"price"`
}

type ImpRequest struct {
	Id        uint `json:"id"`
	Minwidth  uint `json:"minwidth"`
	Minheight uint `json:"minheight"`
}
