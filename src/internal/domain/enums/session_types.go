package enums

const (
	EMAIL_CONFIRM  = 0
	AUTHORIZE      = 1
	RESET_PASSWORD = 2
)

var SessTypesToLabels map[int]string = map[int]string{
	EMAIL_CONFIRM:  "confirm_email",
	AUTHORIZE:      "authorize",
	RESET_PASSWORD: "reset_password",
}
