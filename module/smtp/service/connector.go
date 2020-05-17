package service

import "sync"

type Connector interface {
	Connect(ScannerResult) (ClientPool, error)
}

type connector struct {
	mtx   *sync.Mutex
	pools map[string]*clientPool
}

func (c *connector) Connect(result ScannerResult) (ClientPool, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	pool, ok := c.pools[result.GetHostname()]
	if !ok {
		pool = &clientPool{}
		c.pools[result.GetHostname()] = pool
	}

	return pool, nil
}
