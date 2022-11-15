package interceptorTool

import (
	"context"
	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc_tools/common"
	"strconv"
)

var TokenKey = "hwdhy-0426-0125"

type AuthInterceptor struct {
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

// Unary 拦截器
func (interceptor *AuthInterceptor) Unary(enforcer *casbin.Enforcer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 从context中获取请求头信息
		md, _ := metadata.FromIncomingContext(ctx)
		var token string
		var userId uint64
		var role string

		// http接口调用，解析Cookie
		if len(md["grpcgateway-cookie"]) > 0 {
			token = md["grpcgateway-cookie"][0]
			// 解析用户ID、角色ID
			userId, role = common.GetUserID(token, TokenKey)
			if role == "" {
				role = "tourists"
			}
		} else {
			// rpc调用 解析角色信息
			if len(md["role"]) > 0 {
				role = md["role"][0]
			}
			if len(md["user-id"]) > 0 {
				userId, _ = strconv.ParseUint(md["user-id"][0], 10, 64)
			}
		}

		// 将角色id写入context中
		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "userId", userId)
		// 接口权限校验
		res, err := enforcer.Enforce(role, info.FullMethod, info.Server)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "permission verification failure")
		}
		if res {
			return handler(ctx, req)
		} else {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}
	}
}
