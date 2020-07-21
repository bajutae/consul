package autoconf

import (
	"context"
	"net"

	"github.com/hashicorp/consul/agent/config"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/go-hclog"
)

// DirectRPC is the interface that needs to be satisifed for AutoConfig to be able to perform
// direct RPCs against individual servers. This should not use
type DirectRPC interface {
	RPC(dc string, node string, addr net.Addr, method string, args interface{}, reply interface{}) error
}

type CertMonitor interface {
	Update(*structs.SignedResponse) error
	Start(context.Context) (<-chan struct{}, error)
	Stop() bool
}

type Config struct {
	Logger      hclog.Logger
	DirectRPC   DirectRPC
	BuilderOpts config.BuilderOpts
	Waiter      *lib.RetryWaiter
	Overrides   []config.Source
	CertMonitor CertMonitor
}

// WithLogger will cause the created AutoConfig type to use the provided logger
func (c *Config) WithLogger(logger hclog.Logger) *Config {
	c.Logger = logger
	return c
}

// WithConnectionPool will cause the created AutoConfig type to use the provided connection pool
func (c *Config) WithDirectRPC(directRPC DirectRPC) *Config {
	c.DirectRPC = directRPC
	return c
}

// WithBuilderOpts will cause the created AutoConfig type to use the provided CLI builderOpts
func (c *Config) WithBuilderOpts(builderOpts config.BuilderOpts) *Config {
	c.BuilderOpts = builderOpts
	return c
}

// WithRetryWaiter will cause the created AutoConfig type to use the provided retry waiter
func (c *Config) WithRetryWaiter(waiter *lib.RetryWaiter) *Config {
	c.Waiter = waiter
	return c
}

// WithOverrides is used to provide a config source to append to the tail sources
// during config building. It is really only useful for testing to tune non-user
// configurable tunables to make various tests converge more quickly than they
// could otherwise.
func (c *Config) WithOverrides(overrides ...config.Source) *Config {
	c.Overrides = overrides
	return c
}

// WithCertMonitor is used to provide a certificate monitor to the auto-config.
// This monitor is responsible for renewing the agents TLS certificate and keeping
// the connect CA roots up to date.
func (c *Config) WithCertMonitor(certMonitor CertMonitor) *Config {
	c.CertMonitor = certMonitor
	return c
}
