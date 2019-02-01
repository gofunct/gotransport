package cloud

import (
	"os"
	"path/filepath"
)

type Env int

const (
	Local  Env = 0
	Gcloud Env = 1
	Aws    Env = 2
)

func (e Env) String() string {
	envs := [...]string{
		"local",
		"gcloud",
		"aws",
	}
	return envs[e]
}

type API interface {
	Env() Env
	Bucket() string
	DbName() string
	DbUser() string
	DbPassword() string
	DbHost() string
	DbRegion() string
	RunVarValue() string
	RunVarTimeout() string
	RunVarName() string
}

type ContextFunc func(*Context)

func NewApi(ctx ...ContextFunc) API {
	c := &Context{}
	for _, f := range ctx {
		f(c)
	}
	return c
}

type Context struct {
	Environment       Env
	Storage           string
	DatabaseName      string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseHost      string
	DatabaseRegion    string
	RuntimeVarValue   string
	RuntimeVarTimeout string
	RuntimeVarName    string
}

func (c *Context) Env() Env {
	if c.Environment.String() == "" {

	}
	return c.Environment
}

func (c *Context) Bucket() string {
	return c.Storage
}

func (c *Context) DbName() string {
	return c.DatabaseName
}

func (c *Context) DbUser() string {
	return c.DatabaseUser
}

func (c *Context) DbPassword() string {
	return c.DatabasePassword
}

func (c *Context) DbHost() string {
	return c.DatabaseHost
}

func (c *Context) DbRegion() string {
	return c.DatabaseRegion
}

func (c *Context) RunVarValue() string {
	return c.RuntimeVarValue
}

func (c *Context) RunVarTimeout() string {
	return c.RuntimeVarTimeout
}

func (c *Context) RunVarName() string {
	return c.RuntimeVarName
}

func InitializeApi() API {
	return NewApi(
		func(c *Context) {
			c.Environment = SelectEnv()
			if c.Storage == "" {
				c.Storage = Ask(
					"Please provide a storage bucket name",
					"default-storage",
					true,
				)
			}
			if c.DatabaseName == "" {
				c.DatabaseName = Ask(
					"Please provide a database name",
					filepath.Base(os.Getenv("PWD")),
					true,
				)
			}
			if c.DatabaseUser == "" {
				c.DatabaseUser = Ask(
					"Please provide a database username",
					os.Getenv("USER"),
					true,
				)
			}
			if c.DatabasePassword == "" {
				c.DatabasePassword = Ask(
					"Please provide a database password",
					os.Getenv("USER")+"-password",
					true,
				)
			}
			if c.DatabaseHost == "" {
				c.DatabaseHost = Ask(
					"Please provide a database host",
					"0.0.0.0:3306",
					true,
				)
			}
			if c.DatabaseRegion == "" {
				c.DatabaseRegion = Ask(
					"Please provide a database region",
					"us-central1-a",
					false,
				)
			}

			if c.RuntimeVarTimeout == "" {
				c.RuntimeVarTimeout = "15s"
			}

			if c.RuntimeVarName == "" {
				c.RuntimeVarName = Ask(
					"Please provide a runtime variable name",
					"runtime-variable",
					true,
				)
			}
			if c.RuntimeVarValue == "" {
				c.RuntimeVarValue = Ask(
					"Please provide a runtime variable value",
					"Default Runtime Variable",
					true,
				)
			}

		})
}
