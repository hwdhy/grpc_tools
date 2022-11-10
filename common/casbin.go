package common

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
)

// InitAdapter 项目启动重新加载接口权限
func InitAdapter(permission []map[string]int) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapter("postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			PgsqlHost, PgsqlUsername, PgsqlPassword, "casbin", PgsqlPort))
	if err != nil {
		logrus.Fatal("init gorm adapter err:", err)
	}

	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`)

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logrus.Fatal("new enforcer err:", err)
	}

	err = e.LoadPolicy()
	if err != nil {
		logrus.Fatal("load policy err:", err)
	}
	// 加载默认admin权限， 不存在则创建
	_, err = e.AddPolicy("admin", "/*", "*")
	if err != nil {
		logrus.Fatal(err)
	}

	// 遍历当前注册对象，删除对应接口权限
	for _, perm := range permission {
		for requestUrl := range perm {
			_, err = e.RemoveFilteredPolicy(0, "", requestUrl, "*")
			if err != nil {
				logrus.Fatalf("remove all policy err: %v", err)
			}
		}
	}

	// 遍历接口增加权限控制
	for _, perm := range permission {
		for requestUrl, p := range perm {
			for roleID, roleName := range RoleName {
				if roleID >= p {
					_, err = e.AddPolicy(roleName, requestUrl, "*")
					if err != nil {
						logrus.Fatalf("add policy err: %v", err)
					}
				}
			}
		}
	}
	return e
}
