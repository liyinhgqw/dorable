package doracle

import (
	"encoding/binary"
)

type DoracleStateMachine struct {
	orc *Oracle
}

func (s *DoracleStateMachine) Save() ([]byte, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(s.orc.maxTs))
	return b, nil
}

func (s *DoracleStateMachine) Recovery(b []byte) error {
	maxts := binary.LittleEndian.Uint64(b)
	s.orc.maxTs = int64(maxts)
	return nil
}
