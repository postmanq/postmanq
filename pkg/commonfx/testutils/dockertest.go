package testutils

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"net"
	"strings"
	"time"
)

var (
	pool *Pool
)

const (
	subnet      = "174.28.0.0/16"
	ipRange     = "174.28.5.0/24"
	gateway     = "174.28.5.254"
	networkName = "postmanq-network"
)

func GetPool() *Pool {
	if pool == nil {
		p, err := dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not construct pool: %s", err)
		}

		err = p.Client.Ping()
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}

		containers, err := p.Client.ListContainers(docker.ListContainersOptions{
			All: true,
		})
		if err != nil {
			log.Fatalf("Could not get containers: %s", err)
		}

		for _, container := range containers {
			if strings.Contains(container.Names[0], "postmanq") {
				err := p.Client.RemoveContainer(docker.RemoveContainerOptions{
					ID:            container.ID,
					RemoveVolumes: true,
					Force:         true,
				})
				if err != nil {
					log.Fatalf("Could not remove container: %s", err)
				}
			}
		}

		networks, err := p.Client.ListNetworks()
		if err != nil {
			log.Fatalf("Could not get list networks: %s", err)
		}

		for _, n := range networks {
			if strings.Contains(n.Name, networkName) {
				err = p.Client.RemoveNetwork(n.ID)
				if err != nil {
					log.Fatalf("Could not remove network: %s", err)
				}
			}
		}

		network, err := p.CreateNetwork(
			networkName,
			func(config *docker.CreateNetworkOptions) {
				config.IPAM = &docker.IPAMOptions{
					Config: []docker.IPAMConfig{
						{
							Subnet:  subnet,
							IPRange: ipRange,
							Gateway: gateway,
						},
					},
				}
			},
		)
		if err != nil {
			log.Fatalf("Could not create network: %s", err)
		}

		pool = &Pool{
			pool:    p,
			network: network,
		}
	}
	return pool
}

type Options struct {
	Name  string
	Image string
	Tag   string
	Env   map[string]interface{}
}

type Option func(opt *dockertest.RunOptions)

type Env map[string]interface{}

func WithName(name string) Option {
	return func(opt *dockertest.RunOptions) {
		opt.Name = name
	}
}

func WithImage(image string) Option {
	return func(opt *dockertest.RunOptions) {
		opt.Repository = image
	}
}

func WithTag(tag string) Option {
	return func(opt *dockertest.RunOptions) {
		opt.Tag = tag
	}
}

func WithEnv(env Env) Option {
	return func(opt *dockertest.RunOptions) {
		opt.Env = make([]string, 0)
		for k, v := range env {
			opt.Env = append(opt.Env, fmt.Sprintf("%s=%v", k, v))
		}
	}
}

func WithMount(mount string) Option {
	return func(opt *dockertest.RunOptions) {
		opt.Mounts = append(opt.Mounts, mount)
	}
}

type Resource struct {
	res *dockertest.Resource
}

type Pool struct {
	pool    *dockertest.Pool
	network *dockertest.Network
}

func (r *Resource) GetHost(portId string) string {
	hostAndPort := r.res.GetHostPort(portId)
	host, _, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		log.Fatalf("Could not get host: %s", err)
	}

	if host == "localhost" {
		return "127.0.0.1"
	}

	return host
}

func (r *Resource) GetPort(portId string) string {
	hostAndPort := r.res.GetHostPort(portId)
	_, port, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		log.Fatalf("Could not get port: %s", err)
	}

	return port
}

func (p *Pool) Run(opts ...Option) *Resource {
	runOptions := &dockertest.RunOptions{
		User:     "root",
		Networks: []*dockertest.Network{p.network},
		Mounts:   make([]string, 0),
	}
	for _, opt := range opts {
		opt(runOptions)
	}

	res, err := p.pool.RunWithOptions(runOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	return &Resource{res: res}
}

func (p *Pool) Check(f func() error) {
	p.pool.MaxWait = 60 * time.Second
	err := p.pool.Retry(f)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func (p *Pool) Purge(resource *Resource) {
	err := p.pool.Purge(resource.res)
	if err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (p *Pool) GetContainerByName(name string) (*Resource, bool) {
	res, ok := p.pool.ContainerByName(name)
	if ok {
		return &Resource{res: res}, true
	}

	return nil, false
}
