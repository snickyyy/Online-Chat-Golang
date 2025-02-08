package domain

import domain "libs/src/internal/domain/interfaces"

func Repr(i domain.Stringer) string {
	return i.String()
}
