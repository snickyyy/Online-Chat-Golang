package domain

type DTOInterface interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte, v any) error
}
