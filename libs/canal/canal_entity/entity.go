package canal_entity

import (
	"sync"
)

type TRow = map[string]interface{}

type TRowEntity struct {
	Action    string
	Schema    string
	Table     string
	Row       TRow
	OldRow    TRow
	Timestamp uint32
}

type TPosEntity struct {
	Name  string
	Pos   uint32
	Force bool
}

var rowEntityPool = sync.Pool{
	New: func() interface{} {
		entity := new(TRowEntity)
		entity.Row = make(TRow)
		entity.OldRow = make(TRow)
		return entity
	},
}

func GetRowEntity() *TRowEntity {
	return rowEntityPool.Get().(*TRowEntity)
}

func PutRowEntity(re *TRowEntity) {
	re.Row = make(TRow)
	re.OldRow = make(TRow)
	rowEntityPool.Put(re)
}
