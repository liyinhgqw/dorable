package doracle

import (
	"github.com/goraft/raft"
)

type OracleCommand struct {
	Num int32 `json:"num"`
}

func NewOracleCommand(num int32) *OracleCommand {
	return &OracleCommand{
		Num: num,
	}
}

func (c *OracleCommand) CommandName() string {
	return "getts"
}

func (c *OracleCommand) Apply(server raft.Server) (interface{}, error) {
	orc := server.Context().(*Oracle)
	return orc.GetTimestamp(c.Num), nil
}
