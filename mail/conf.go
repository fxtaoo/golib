package mail

type Smtp struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	UserPW string `json:"userpw"`
}
