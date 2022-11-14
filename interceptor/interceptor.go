package interceptorTool

import (
	"context"
	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc_tools/common"
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
		token := ""
		if len(md["grpcgateway-cookie"]) > 0 {
			token = md["grpcgateway-cookie"][0]
		}
		// 解析用户ID、角色ID
		userId, role := common.GetUserID(token, TokenKey)
		if role == "" {
			role = "tourists"
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
