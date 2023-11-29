package messagetype

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
