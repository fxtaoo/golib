package mail

type Smtp struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Address string `json:"address"`
	PassWd  string `json:"passWd"`
	Name    string `json:"user"`
}
