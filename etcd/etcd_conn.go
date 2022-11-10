package etcd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// ClientConn 通过etcd获取grpc连接
func ClientConn(serviceName string) *grpc.ClientConn {
	etcdResolverBuilder := NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)

	conn, err := grpc.Dial("etcd:///"+serviceName,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("dial server err: %v", err)
		return nil
	}
	return conn
}
