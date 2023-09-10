package canal_handler

import "go-canal/libs/canal/canal_entity"

func GetHeaderMap() THeaderMap {
	return make(THeaderMap)
}

func PutHeaderMap(hvm THeaderMap) {
	hvm = nil
}

func RowHeaderValueMap(entityRow canal_entity.TRow, headerValueMap THeaderMap, eventRow TEventRow) {
	for header, index := range headerValueMap {
		if index >= len(eventRow) {
			if index-1 < len(eventRow) {
				index = index - 1
			} else {
				continue
			}
		}
		entityRow[header] = eventRow[index]
	}
}
