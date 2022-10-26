package common

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"log"
)

// InitAdapter 权限初始化
func InitAdapter(permission []map[string]int) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapter("postgres", fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		PgsqlUsername, PgsqlPassword, "casbin", PgsqlPort))
	if err != nil {
		logrus.Fatal("new gorm adapter err:", err)
	}

	e, err := casbin.NewEnforcer("./auth_model.conf", adapter)
	if err != nil {
		logrus.Fatal("new enforcer err:", err)
	}

	err = e.LoadPolicy()
	if err != nil {
		logrus.Fatal("load policy err:", err)
	}

	_, err = e.AddPolicy("admin", "/*", "*")
	if err != nil {
		logrus.Fatal(err)
	}

	for _, perm := range permission {
		for requestUrl := range perm {
			_, err = e.RemoveFilteredPolicy(0, "", requestUrl, "*")
			if err != nil {
				logrus.Fatalf("remove all policy err: %v", err)
			}
		}
	}

	for _, perm := range permission {
		// 遍历接口增加权限控制
		for requestUrl, p := range perm {
			for roleID, roleName := range RoleName {
				if roleID >= p {
					_, err = e.AddPolicy(roleName, requestUrl, "*")
					if err != nil {
						log.Fatalf("add policy err: %v", err)
					}
				}
			}
		}
	}
	return e
}
