package detcd

import (
	"context"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/conf"
	"fmt"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Resolver struct {
	cc                 resolver.ClientConn
	serviceName        string
	grpcClientConn     *grpc.ClientConn
	cli                *clientv3.Client
	schema             string
	endpoints          []string
	watchStartRevision int64
}

var (
	nameResolver        = make(map[string]*Resolver)
	rwNameResolverMutex sync.RWMutex
)

func NewResolver(opt *conf.GrpcDialOption) (r *Resolver, err error) {
	var (
		etcdCli *clientv3.Client
		opts    []grpc.DialOption
		conn    *grpc.ClientConn
	)

	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints: opt.Etcd.Endpoints,
	})
	if err != nil {
		dlog.Error(err.Error())
		return
	}

	r = new(Resolver)
	r.serviceName = opt.ServiceName
	r.cli = etcdCli
	r.schema = opt.Etcd.Schema
	r.endpoints = opt.Etcd.Endpoints
	resolver.Register(r)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if opt.Tracing.Enabled == true && opt.Tracing.Tracer != nil {
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opt.Tracing.Tracer)))
	}
	if opt.Cert.Enabled == true {
		var creds credentials.TransportCredentials
		creds, err = credentials.NewClientTLSFromFile(opt.Cert.CertFile, opt.Cert.ServerNameOverride)
		if err != nil {
			dlog.Error(err.Error())
		} else {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err = grpc.DialContext(ctx, GetPrefix(opt.Etcd.Schema, opt.ServiceName), opts...)
	if err == nil {
		r.grpcClientConn = conn
	}
	return
}

func (r1 *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r1 *Resolver) Close() {
}

func GetConn(opt *conf.GrpcDialOption) *grpc.ClientConn {
	var (
		r   *Resolver
		key = opt.Etcd.Schema + opt.ServiceName
		ok  bool
		err error
	)
	rwNameResolverMutex.RLock()
	r, ok = nameResolver[key]
	rwNameResolverMutex.RUnlock()
	if ok {
		return r.grpcClientConn
	}

	rwNameResolverMutex.Lock()
	r, ok = nameResolver[key]
	if ok {
		rwNameResolverMutex.Unlock()
		return r.grpcClientConn
	}

	r, err = NewResolver(opt)
	if err != nil {
		rwNameResolverMutex.Unlock()
		return nil
	}

	nameResolver[key] = r
	rwNameResolverMutex.Unlock()
	return r.grpcClientConn
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if r.cli == nil {
		return nil, fmt.Errorf("etcd clientv3 client failed, etcd:%s", target)
	}
	r.cc = cc
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	prefix := GetPrefix(r.schema, r.serviceName)
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err == nil {
		var addrList []resolver.Address
		for i := range resp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: string(resp.Kvs[i].Value)})
		}
		r.cc.UpdateState(resolver.State{Addresses: addrList})
		r.watchStartRevision = resp.Header.Revision + 1
		go r.watch(prefix, addrList)
	} else {
		return nil, fmt.Errorf("etcd get failed, prefix: %s", prefix)
	}
	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.schema
}

func exists(addrList []resolver.Address, addr string) bool {
	for _, v := range addrList {
		if v.Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func (r *Resolver) watch(prefix string, addrList []resolver.Address) {
	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrefix())
	for n := range rch {
		flag := 0
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				if !exists(addrList, string(ev.Kv.Value)) {
					flag = 1
					addrList = append(addrList, resolver.Address{Addr: string(ev.Kv.Value)})
				}
			case mvccpb.DELETE:
				i := strings.LastIndexAny(string(ev.Kv.Key), "/")
				if i < 0 {
					return
				}
				t := string(ev.Kv.Key)[i+1:]
				if s, ok := remove(addrList, t); ok {
					flag = 1
					addrList = s
				}
			}
		}

		if flag == 1 {
			r.cc.UpdateState(resolver.State{Addresses: addrList})
		}
	}
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}
func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}
