package golf

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type bound struct {
	min  int
	max  int
	step int
}

var typeMap = map[string]bound{
	"uint":    {min: 1, max: 256, step: 8},
	"int":     {min: 1, max: 256, step: 8},
	"bytes":   {min: 1, max: 32, step: 1},
	"address": {},
	"string":  {},
	"bool":    {},
}

func areValidtypes(types []string) error {
	for _, t := range types {
		err := checkValidType(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkValidType(t string) error {
	var typeName strings.Builder
	for _, c := range t {
		if unicode.IsLetter(c) {
			typeName.WriteRune(c)
		} else {
			break
		}
	}

	typeStr := typeName.String()

	if _, exists := typeMap[typeStr]; !exists || (typeStr == "") {
		return fmt.Errorf("invalid function arg type %s", t)
	}

	b := typeMap[typeStr]
	if b == (bound{}) {
		if len(typeStr) == len(t) {
			return nil
		} else {
			return fmt.Errorf("invalid function arg type %s", t)
		}
	}

	typeVal := strings.TrimPrefix(t, typeStr)
	typeInt, typeErr := strconv.Atoi(typeVal)
	if typeErr != nil {
		return fmt.Errorf("invalid function arg suffix %s", typeVal)
	}

	if typeInt > b.min && typeInt <= b.max && typeInt%b.step == 0 {
		return nil
	}

	return fmt.Errorf("invalid function arg %s, outside suffix bounds min:%d max:%d step:%d", t, b.min, b.max, b.step)
}
