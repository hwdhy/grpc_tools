package etcd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"strconv"
)

// grpc调用将cookie放入请求体中
type tokenAuth struct {
	userId uint64
	role   string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"role":    t.role,
		"user-id": strconv.FormatUint(t.userId, 10),
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return false
}

// ClientConn 通过etcd获取grpc连接
func ClientConn(serviceName string, userId uint64, role string) *grpc.ClientConn {
	etcdResolverBuilder := NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)

	var conn *grpc.ClientConn
	var err error
	if role == "" || userId == 0 {
		conn, err = grpc.Dial("etcd:///"+serviceName,
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		conn, err = grpc.Dial("etcd:///"+serviceName,
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithPerRPCCredentials(&tokenAuth{
				role:   role,
				userId: userId,
			}))
	}
	if err != nil {
		logrus.Errorf("dial server err: %v", err)
		return nil
	}
	return conn
}
