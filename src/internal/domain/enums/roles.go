package enums

const (
	ANONYMOUS = 0
	USER      = 1
	ADMIN     = 2
)

var RolesToLabels map[int]string = map[int]string{
	ANONYMOUS: "anonymous",
	USER:      "user",
	ADMIN:     "admin",
}
