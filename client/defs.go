package client

type CmdType int
type ResposeStatus int

const (
	UndefinedCmd CmdType = iota
	GetCmd
	PutCmd
	DeleteCmd
)

const (
	UndefinedStatus ResposeStatus = iota
	OkStatus
	FailStatus
)

type Arg struct {
	Key   string
	Value []byte `json:",omitempty"`
}

type CommandRequest struct {
	Command CmdType
	Args    []Arg `json:",omitempty"`
}

type CommandResponse struct {
	Status ResposeStatus
	Args   []Arg `json:",omitempty"`
}
