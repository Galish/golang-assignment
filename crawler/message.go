package crawler

type Message struct {
	ID     int    `json:"id"`
	Index  int    `json:"index"`
	Link   string `json:"link"`
	Avatar string `json:"avatar"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	HTML   string `json:"html"`
	Text   string `json:"text"`
}
