package telebotai

type Response struct {
	Id      string            `json:"id"`
	Object  string            `json:"object"`
	Created uint              `json:"created"`
	Model   string            `json:"model"`
	Choices map[string]string `json:"choices"`
	Usage   map[string]int    `json:"usage"`
}

type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Max_tokens  int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
	Top_p       int    `json:"top_p"`
}
