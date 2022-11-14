package common

const (
	NotLogged = iota // 0：未登录  1：会员  2：管理员
	HasLogged
	Admin
)

var RoleName = map[int]string{
	NotLogged: "tourists",
	HasLogged: "member",
	Admin:     "admin",
}

const (
	PgsqlUsername = "postgres"
	PgsqlPassword = "123456"
	PgsqlPort     = 5432
	PgsqlHost     = "127.0.0.1"

	TokenKey = "hwdhy-0426-0125"
)
