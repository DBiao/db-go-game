package detcd

import (
	"context"
	"db-go-game/pkg/common/dlog"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"strconv"
	"time"
)

type RegEtcd struct {
	cli          *clientv3.Client
	ctx          context.Context
	cancel       context.CancelFunc
	endpoints    []string
	serviceValue string
	serviceKey   string
	ttl          int
	host         string
	port         int
	schema       string
	serviceName  string
	reChan       chan bool
	connected    bool
}

var rEtcd *RegEtcd

// "%s:///%s/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema string, endpoints []string, myHost string, myPort int, serviceName string, ttl int) error {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return RegisterEtcd(schema, endpoints, myHost, myPort, serviceName, ttl)
}

func GetTarget(schema, myHost string, myPort int, serviceName string) string {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return serviceName
}

func RegisterEtcd(schema string, endpoints []string, myHost string, myPort int, serviceName string, ttl int) (err error) {
	serviceValue := net.JoinHostPort(myHost, strconv.Itoa(myPort))
	serviceKey := GetPrefix(schema, serviceName) + serviceValue
	rEtcd = &RegEtcd{
		endpoints:    endpoints,
		serviceValue: serviceValue,
		serviceKey:   serviceKey,
		ttl:          ttl,
		host:         myHost,
		port:         myPort,
		schema:       schema,
		serviceName:  serviceName,
		reChan:       make(chan bool),
	}
	rEtcd.ReRegister()
	rEtcd.Register()
	return
}

func (r *RegEtcd) Register() (err error) {
	defer func() {
		if err != nil {
			r.reChan <- true
		}
	}()

	var (
		cli    *clientv3.Client
		ctx    context.Context
		cancel context.CancelFunc
		resp   *clientv3.LeaseGrantResponse
		kresp  <-chan *clientv3.LeaseKeepAliveResponse
	)

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: 5 * time.Second})
	if err != nil {
		dlog.Error(err.Error())
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	resp, err = cli.Grant(ctx, int64(r.ttl))
	if err != nil {
		dlog.Error(err.Error())
		return
	}

	if _, err = cli.Put(ctx, r.serviceKey, r.serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		dlog.Error(err.Error())
		return
	}
	kresp, err = cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		dlog.Error(err.Error())
		return
	}
	rEtcd.cli = cli
	rEtcd.ctx = ctx
	rEtcd.cancel = cancel

	go func() {
		for {
			select {
			case _, ok := <-kresp:
				r.connected = ok
				if ok {
					//xlog.Info("续约成功")
				} else {
					dlog.Error("租约失效")
					r.UnRegisterEtcd()
					r.reChan <- true
					return
				}
			}
		}
	}()
	return
}

func (r *RegEtcd) ReRegister() {
	go func() {
		var (
			ok bool
		)
		select {
		case _, ok = <-r.reChan:
			if ok == false {
				return
			}
			time.Sleep(2 * time.Second)
			r.Register()
		}
	}()
}

func (r *RegEtcd) UnRegisterEtcd() {
	r.cancel()
	r.cli.Delete(rEtcd.ctx, rEtcd.serviceKey)
}
