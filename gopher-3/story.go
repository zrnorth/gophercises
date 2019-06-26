package cyoa

// Story : self-explanatory
type Story map[string]Arc

// Arc : A path through the choose-your-own-adventure story
type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option for each Arc
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
