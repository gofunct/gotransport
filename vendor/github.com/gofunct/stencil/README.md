# stencil
--
    import "github.com/gofunct/stencil"


## Usage

```go
const (
	// CaptureStdout is a bitmask to capture STDOUT
	CaptureStdout = 1
	// CaptureStderr is a bitmask to capture STDERR
	CaptureStderr = 2
	// CaptureBoth captures STDOUT and STDERR
	CaptureBoth = CaptureStdout + CaptureStderr
)
```

```go
var Env string
```
Env is the default environment to use for all commands. That is, the effective
environment for all commands is the merged set of (parent environment, Env, func
specified environment). Whitespace or newline separate key value pairs. $VAR
interpolation is allowed.

Env = "GOOS=linux GOARCH=amd64" Env = `

    GOOS=linux
    GOPATH=./vendor:$GOPATH

`

```go
var InheritParentEnv bool
```
InheritParentEnv whether to inherit parent's environment

```go
var PathListSeparator = "::"
```
PathListSeparator is a cross-platform path list separator. On Windows,
PathListSeparator is replacd by ";". On others, PathListSeparator is replaced by
":"

```go
var Processes = make(map[string]*os.Process)
```
Processes are the processes spawned by Start()

```go
var Verbose = false
```
Verbose indicates whether to log verbosely

```go
var Version = "v0.0.1"
```
Version is the current version

```go
var (

	// WaitMs is the default time (1500 ms) to debounce task events in watch mode.
	Wait time.Duration
)
```

#### func  Bash

```go
func Bash(script string, options ...map[string]interface{}) (string, error)
```
Bash executes a bash script (string).

#### func  BashOutput

```go
func BashOutput(script string, options ...map[string]interface{}) (string, error)
```
BashOutput executes a bash script and returns the output

#### func  EffectiveEnv

```go
func EffectiveEnv(funcEnv []string) []string
```
EffectiveEnv is the effective environment for an exec function.

#### func  GetWatchDelay

```go
func GetWatchDelay() time.Duration
```
GetWatchDelay gets the watch delay

#### func  Getenv

```go
func Getenv(key string) string
```
Getenv environment variable from a string array.

#### func  GoThrottle

```go
func GoThrottle(throttle int, fns ...func() error) error
```
GoThrottle starts to run the given list of fns concurrently, at most n fns at a
time.

#### func  Halt

```go
func Halt(v interface{})
```
Halt is a soft panic and stops a task.

#### func  Inside

```go
func Inside(dir string, lambda func()) error
```
Inside temporarily changes the working directory and restores it when lambda
finishes.

#### func  ParseStringEnv

```go
func ParseStringEnv(s string) []string
```
ParseStringEnv parse the package Env string and converts it into an environment
slice.

#### func  Prompt

```go
func Prompt(prompt string) string
```
Prompt prompts user for input with default value.

#### func  PromptPassword

```go
func PromptPassword(prompt string) string
```
PromptPassword prompts user for password input.

#### func  Run

```go
func Run(commandstr string, options ...map[string]interface{}) (string, error)
```
Run runs a command.

#### func  RunOutput

```go
func RunOutput(commandstr string, options ...map[string]interface{}) (string, error)
```
RunOutput runs a command and returns output.

#### func  SetEnviron

```go
func SetEnviron(envstr string, inheritParent bool)
```
SetEnviron sets the environment for child processes. Note that SetEnviron(Env,
InheritParentEnv) is called once automatically.

#### func  SetWatchDelay

```go
func SetWatchDelay(delay time.Duration)
```
SetWatchDelay sets the time duration between watches.

#### func  Start

```go
func Start(commandstr string, options ...map[string]interface{}) error
```
Start starts an async command. If executable has suffix ".go" then it will be
"go install"ed then executed. Use this for watching a server task.

If Start is called with the same command it kills the previous process.

The working directory is optional.

#### func  Stencil

```go
func Stencil(tasksFunc func(*Project))
```
Stencil runs a project of tasks.

#### func  Usage

```go
func Usage(tasks string)
```
Usage prints a usage screen with task descriptions.

#### type Asset

```go
type Asset struct {
	bytes.Buffer
	Info *glob.FileAsset
	// WritePath is the write destination of the asset.
	WritePath string
	Pipeline  *Pipeline
}
```

Asset is any file which can be loaded and processed by a filter.

#### func (*Asset) ChangeExt

```go
func (asset *Asset) ChangeExt(newExt string)
```
ChangeExt changes the extension of asset.WritePath. ChangExt is used by filters
which transpile source. For example, a filter for Markdown would use
ChangeExt(".html") to write the asset as an HTML file.

#### func (*Asset) Dump

```go
func (asset *Asset) Dump() string
```
Dump returns a console friendly representation of asset. Note, String() returns
the string value of Buffer.

#### func (*Asset) Ext

```go
func (asset *Asset) Ext() string
```
Ext returns the extension of asset.WritePath.

#### func (*Asset) IsText

```go
func (asset *Asset) IsText() bool
```
IsText return true if it thinks this asset is text based, meaning it can be
manipulated with string functions.

#### func (*Asset) MimeType

```go
func (asset *Asset) MimeType() string
```
MimeType returns an educated guess of the content type of asset.

#### func (*Asset) Rewrite

```go
func (asset *Asset) Rewrite(bytes []byte)
```
Rewrite sets the buffer to bytes.

#### func (*Asset) RewriteString

```go
func (asset *Asset) RewriteString(s string)
```
RewriteString sets the buffer to a string value.

#### type Context

```go
type Context struct {
	// Task is the currently running task.
	Job *Job

	Pipeline *Pipeline

	// FileEvent is an event from the watcher with change details.
	FileEvent *watcher.FileEvent

	// Task command line arguments
	Args minimist.ArgMap

	Error error
}
```

Context is the data passed to a task.

#### func (*Context) AnyFile

```go
func (context *Context) AnyFile() []string
```
AnyFile returns either a non-DELETe FileEvent file or the WatchGlob patterns
which can be used by goa.Load()

#### func (*Context) Bash

```go
func (context *Context) Bash(cmd string, options ...map[string]interface{})
```
Bash runs a bash shell.

#### func (*Context) BashOutput

```go
func (context *Context) BashOutput(script string, options ...map[string]interface{}) string
```
BashOutput executes a bash script and returns the output

#### func (*Context) Check

```go
func (context *Context) Check(err error, msg string)
```
Check halts the task if err is not nil.

Do this

    Check(err, "Some error occured")

Instead of

    if err != nil {
    	Halt(err)
    }

#### func (*Context) Pipe

```go
func (context *Context) Pipe(filters ...interface{})
```
Bash runs a bash shell.

#### func (*Context) Run

```go
func (context *Context) Run(cmd string, options ...map[string]interface{})
```
Run runs a command

#### func (*Context) RunOutput

```go
func (context *Context) RunOutput(commandstr string, options ...map[string]interface{}) string
```
RunOutput runs a command and returns output.

#### func (*Context) Start

```go
func (context *Context) Start(cmd string, options ...map[string]interface{})
```
Start run aysnchronously.

#### type Dependency

```go
type Dependency interface {
	// contains filtered or unexported methods
}
```

Dependency marks an interface as a dependency.

#### type Handler

```go
type Handler interface {
	Handle(*Context)
}
```

Handler is the interface which all task handlers eventually implement.

#### type HandlerFunc

```go
type HandlerFunc func(*Context)
```

HandlerFunc is a Handler adapter.

#### func (HandlerFunc) Handle

```go
func (f HandlerFunc) Handle(ctx *Context)
```
Handle implements Handler.

#### type Job

```go
type Job struct {
	Name string

	Handler Handler

	// computed based on dependencies
	EffectiveWatchRegexps []*glob.RegexpInfo
	EffectiveWatchGlobs   []string

	// Complete indicates whether this task has already ran. This flag is
	// ignored in watch mode.
	Complete bool

	RunOnce bool

	SrcFiles   []*glob.FileAsset
	SrcGlobs   []string
	SrcRegexps []*glob.RegexpInfo

	DestFiles   []*glob.FileAsset
	DestGlobs   []string
	DestRegexps []*glob.RegexpInfo

	// used when a file event is received between debounce intervals, the file event
	// will queue itself and set this flag and force debounce to run it
	// when time has elapsed
	sync.Mutex
}
```

A Task is an operation performed on a user's project directory.

#### func  NewJob

```go
func NewJob(name string, argm minimist.ArgMap) *Job
```
NewTask creates a new Task.

#### func (*Job) DependencyNames

```go
func (j *Job) DependencyNames() []string
```
DependencyNames gets the flattened dependency names.

#### func (*Job) Deps

```go
func (j *Job) Deps(names ...interface{})
```
Deps are task dependencies and must specify how to run jobs in series or in
parallel.

#### func (*Job) Desc

```go
func (j *Job) Desc(desc string) *Job
```
Desc is alias for Description.

#### func (*Job) Description

```go
func (j *Job) Description(desc string) *Job
```
Description sets the description for the task.

#### func (*Job) Dest

```go
func (j *Job) Dest(globs ...string) *Job
```
Dest adds target globs which are used to calculated outdated files. The jobs is
not run unless ANY file Src are newer than ANY in DestN.

#### func (*Job) Run

```go
func (j *Job) Run() error
```
Run runs all the dependencies of this job and when they have completed, runs
this job.

#### func (*Job) RunWithEvent

```go
func (j *Job) RunWithEvent(logName string, e *watcher.FileEvent) (err error)
```
RunWithEvent runs this task when triggered from a watch. *e* FileEvent contains
information about the file/directory which changed in watch mode.

#### func (*Job) Src

```go
func (j *Job) Src(globs ...string) *Job
```
Src adds a source globs to this task. The task is not run unless files are
outdated between Src and Dest globs.

#### func (*Job) Wait

```go
func (j *Job) Wait(duration time.Duration) *Job
```
Wait is minimum milliseconds before task can run again

#### type M

```go
type M map[string]interface{}
```

M is generic string to interface alias

#### type Message

```go
type Message struct {
	Event string
	Data  string
}
```

Message are sent on the Events channel

#### type P

```go
type P []interface{}
```

P is alias for Parallel

#### type Parallel

```go
type Parallel []interface{}
```

Parallel runs jobs in parallel

#### type Pipeline

```go
type Pipeline struct {
	Assets  []*Asset
	Filters []interface{}
}
```

Pipeline is a asset flow through which each asset is processed by one or more
filters. For text files this could be something as simple as adding a header or
minification. Some filters process assets in batches combining them, for example
concatenating JavaScript or CSS.

#### func  NewPipeline

```go
func NewPipeline() *Pipeline
```

#### func  Pipe

```go
func Pipe(filters ...interface{}) (*Pipeline, error)
```
Pipe creates a pipeline with filters and runs it.

#### func (*Pipeline) AddAsset

```go
func (pipeline *Pipeline) AddAsset(asset *Asset)
```
AddAsset adds an asset

#### func (*Pipeline) Pipe

```go
func (pipeline *Pipeline) Pipe(filters ...interface{}) *Pipeline
```
Pipe adds one or more filters to the pipeline. Pipe may be called more than
once.

Filters are simple function. Options are handle through closures. The supported
handlers are

1. Single asset handler. Use this to transorm each asset individually.

    AddHeader filter is an example.

      // signature
      func(*pipeline.Asset) error

2. Multi asset handler. Does not modify the number of elements. See

    Write filter is an example.

      //  signature
      func(assets []*pipeline.Asset) error

3. Pipeline handler. Use this to have unbridled control. Load filter

    is an example.

      // signature
      func(*Pipeline) error

#### func (*Pipeline) Run

```go
func (pipeline *Pipeline) Run()
```
Run runs assets through the pipeline.

#### func (*Pipeline) Truncate

```go
func (pipeline *Pipeline) Truncate()
```
Truncate removes all assets, resetting Assets to empty slice.

#### type Project

```go
type Project struct {
	sync.Mutex
	Jobs      map[string]*Job
	Namespace map[string]*Project
}
```

Project is a container for tasks.

#### func  NewProject

```go
func NewProject(tasksFunc func(*Project), exitFn func(code int), argm minimist.ArgMap) *Project
```
NewProject creates am empty project ready for tasks.

#### func (*Project) Define

```go
func (project *Project) Define(fn func(*Project))
```
Define defines tasks

#### func (*Project) Exit

```go
func (project *Project) Exit(code int)
```
Exit quits the project.

#### func (*Project) Func

```go
func (project *Project) Func(name string, dependencies Dependency, handler func(*Context)) *Job
```
Task adds a task to the project with dependencies and handler.

#### func (*Project) Func1

```go
func (project *Project) Func1(name string, handler func(*Context)) *Job
```
Do1 adds a simple task to the project.

#### func (*Project) FuncD

```go
func (project *Project) FuncD(name string, dependencies Dependency) *Job
```
TaskD adds a task which runs other dependencies with no handler.

#### func (*Project) Run

```go
func (project *Project) Run(name string) error
```
Run runs a task by name.

#### func (*Project) Use

```go
func (project *Project) Use(namespace string, tasksFunc func(*Project))
```
Use uses another project's task within a namespace.

#### func (*Project) Watch

```go
func (project *Project) Watch(names []string, isParent bool) bool
```
Watch watches the Files of a task and reruns the task on a watch event. Any
direct dependency is also watched. Returns true if watching.

TODO: 1. Only the parent task watches, but it gathers wath info from all
dependencies.

2. Anything without src files always run when a dependency is triggered by a
glob match.

    build [generate{*.go} compile] => go file changes =>  build, generate and compile

3. Tasks with src only run if it matches a src

    build [generate{*.go} css{*.scss} compile] => go file changes => build, generate and compile
    css does not need to run since no SCSS files ran

X depends on [A:txt, B] => txt changes A runs, X runs without deps X:txt on [A,
B] => txt changes A, B, X runs

#### type S

```go
type S []interface{}
```

S is alias for Series

#### type Series

```go
type Series []interface{}
```

Series are dependent jobs which must run in series.
