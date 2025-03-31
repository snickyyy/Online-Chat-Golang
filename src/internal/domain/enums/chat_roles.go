package enums

const (
	MEMBER     = 0
	CHAT_ADMIN = 1
	OWNER      = 2
)

var ChatRolesToLabels map[int]string = map[int]string{
	MEMBER:     "member",
	CHAT_ADMIN: "admin",
	OWNER:      "owner",
}

var ChatLabelsToRoles map[string]int = map[string]int{
	"member": MEMBER,
	"admin":  CHAT_ADMIN,
	"owner":  OWNER,
}
