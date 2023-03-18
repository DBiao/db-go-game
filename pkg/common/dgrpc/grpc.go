package dgrpc

import (
	"db-go-game/pkg/common/detcd"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/common/dtracer"
	"db-go-game/pkg/conf"
	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io"
	"time"
)

func NewServer(c *conf.Grpc) (srv *grpc.Server, closer io.Closer) {
	var (
		keepParams grpc.ServerOption
		tracer     opentracing.Tracer
		creds      credentials.TransportCredentials
		opts       []grpc.ServerOption
		err        error
	)
	keepParams = grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(c.MaxConnectionIdle) * time.Millisecond,
		MaxConnectionAge:      time.Duration(c.MaxConnectionAge) * time.Millisecond,
		MaxConnectionAgeGrace: time.Duration(c.MaxConnectionAgeGrace) * time.Millisecond,
		Time:                  time.Duration(c.Time) * time.Millisecond,
		Timeout:               time.Duration(c.Timeout) * time.Millisecond,
	})
	opts = append(opts, keepParams)
	if c.Credential.Enabled == true {
		// TLS认证
		creds, err = credentials.NewServerTLSFromFile(c.Credential.CertFile, c.Credential.KeyFile)
		if err != nil {
			dlog.Error(err.Error())
		} else {
			opts = append(opts, grpc.Creds(creds))
		}
	}
	if c.Jaeger.Enabled == true {
		// 链路追踪
		tracer, closer, err = dtracer.NewTracer(c.Name, c.Jaeger)
		if err != nil {
			dlog.Error(err.Error())
		} else {
			opts = append(opts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
		}
	}
	if c.StreamsLimit > 0 {
		// 设置最大并发流
		opts = append(opts, grpc.MaxConcurrentStreams(c.StreamsLimit))
	}
	if c.MaxRecvMsgSize > 0 {
		// 设置服务器端可接收的最大请求体长度为4KB
		opts = append(opts, grpc.MaxRecvMsgSize(c.MaxRecvMsgSize))
	}
	srv = grpc.NewServer(opts...)
	return
}

func GetClientConn(opt *conf.GrpcDialOption) (clientConn *grpc.ClientConn) {
	clientConn = detcd.GetConn(opt)
	return
}

/*
type ServerParameters struct {
    // 当连接处于idle的时长超过 MaxConnectionIdle时，服务端就发送GOAWAY，关闭连接。默认值为无限大
    MaxConnectionIdle time.Duration
    // 一个连接只能使用 MaxConnectionAge 这么长的时间，服务端就会关闭这个连接。
    MaxConnectionAge time.Duration
    // 服务端优雅关闭连接时长
    MaxConnectionAgeGrace time.Duration
    // 这个时间是服务端用来ping 客户端的。默认值为2小时
    Time time.Duration
    // 默认值为20秒
    Timeout time.Duration
}

type ClientParameters struct {
    // 超过这个时长都没有活动的话，客户端就会ping服务端，这个值最小是10秒。
    Time time.Duration // The current default value is infinity.
    // 发出Ping后，客户端等待 这个时长，仍旧没有收到ping的回复消息。
    Timeout time.Duration // The current default value is 20 seconds.
    PermitWithoutStream bool // false by default.
}
*/
