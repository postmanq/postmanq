package service

import (
	"net"
	"sync"
)

type ConnectorFactory interface {
	CreateConnector(string, []net.IP) Connector
}

func NewConnectorFactory() ConnectorFactory {
	return &connectorFactory{}
}

type connectorFactory struct{}

func (f *connectorFactory) CreateConnector(hostname string, ips []net.IP) Connector {
	c := &connector{
		connMtx:          new(sync.Mutex),
		pools:            make(map[string]*clientPool),
		processablePools: make(chan *clientPool),
		hostname:         hostname,
		ips:              ips,
	}

	return c
}

type Connector interface {
	Connect(ScannerResult) ClientPool
}

type connector struct {
	connMtx          *sync.Mutex
	pools            map[string]*clientPool
	processablePools chan *clientPool
	hostname         string
	ips              []net.IP
}

func (c *connector) processPools() {
	for pool := range c.processablePools {
		go c.processPool(pool)
	}
}

func (c *connector) processPool(pool *clientPool) {
	//for _, mx := range pool.mxs {
	//	for _, ip := range c.ips {
	//		tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ip.String(), "0"))
	//
	//		dialer := &net.Dialer{
	//			//Timeout:   common.App.Timeout().Connection,
	//			LocalAddr: tcpAddr,
	//		}
	//		hostname := net.JoinHostPort(mx.Host, "25")
	//
	//		connection, err := dialer.Dial("tcp", hostname)
	//
	//		//connection.SetDeadline(time.Now().Add(common.App.Timeout().Hello))
	//		client, err := smtp.NewClient(connection, mx.Host)
	//
	//		err = client.Hello(c.hostname)
	//
	//		useTLS, _ := client.Extension("STARTTLS")
	//
	//		client.StartTLS()
	//	}
	//}
}

func (c *connector) Connect(result ScannerResult) ClientPool {
	c.connMtx.Lock()
	defer c.connMtx.Unlock()
	pool, ok := c.pools[result.GetHostname()]
	if ok {
		if !pool.mxs.Equal(result.GetMxs()) {
			c.sendToReconnect(pool)
		}
	} else {
		pool = &clientPool{
			mxs: result.GetMxs(),
			wg:  new(sync.WaitGroup),
		}
		c.pools[result.GetHostname()] = pool
		c.sendToReconnect(pool)
	}

	defer pool.wg.Wait()
	return pool
}

func (c *connector) sendToReconnect(pool *clientPool) {
	pool.lock()
	c.processablePools <- pool
}
