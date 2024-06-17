# Cosys - Common
The package to be imported in cosys projects.

## Documentation

## type Cosys

```go

```

Cosys is the main app, where the database, module service, server, models, services, configs can be accessed from.

When the Cosys app is run, the Cosys struct undergoes the following procedure.

1. Import - The database, module service and server are registered using the blank import method (init function).
2. Initialization - The Cosys struct is initialized, the config yaml files are parsed and loaded into the struct.
3. Registration - The modules are registered. The models and services from each module are loaded into the Cosys struct. The 'OnRegister' function of each module is run.
4. Listening - The server from the Cosys struct is started.
5. Destruction - Clean-up processes before the app is shutdown. The 'OnDestroy' function of each module is run.

### **func Database**

```go
func (c Cosys) Database() Database
```

Returns the Database struct.

### func ModuleService

```go
func (c Cosys) ModuleService() ModuleService
```

Returns the ModuleService struct

### func Server

```go
func (c Cosys) Server() Server
```

Returns the Server struct

### func Register

```go
func (c Cosys) Register(modules map[string]*Module) (*Cosys, error)
```

Register the modules. The models and services of each module are loaded into the Cosys struct and the ‘OnRegister’ functions of each module are run.

### func Start

```go
func (c Cosys) Start() error
```

Starts the server.

### func Destroy

```go
func (c Cosys) Destroy() (*Cosys, error)
```

Runs clean-up processes. The ‘OnDestroy’ functions of each module are run.

## type Module

```go
type Module struct {
	Routes      []*Route
	Controllers map[string]*Controller
	Middlewares map[string]Middleware
	Policies    map[string]Policy

	Models   map[string]Model
	Services map[string]Service
	
	OnRegister func(Cosys) error
	OnDestroy  func(Cosys) error
}
```

The entrypoint into modules.

Routes, Controllers, Middlewares, Policies are used by the server to create API endpoints.

Models and Services from the model can be accessed by any other module through the Cosys struct.

OnRegister and OnDestroy are called when the app is started and closed respectively.

## type Database

```go
type Database interface {
	FindOne(uid string, params DBParams) (Entity, error)
	FindMany(uid string, params DBParams) ([]Entity, error)
	Create(uid string, data Entity, params DBParams) (Entity, error)
	CreateMany(uid string, data []Entity, params DBParams) ([]Entity, error)
	Update(uid string, data Entity, params DBParams) (Entity, error)
	UpdateMany(uid string, data Entity, params DBParams) ([]Entity, error)
	Delete(uid string, params DBParams) (Entity, error)
	DeleteMany(uid string, params DBParams) ([]Entity, error)
}
```

Database provides an API to interact with the database layer directly.

Database will be registered using blank imports (call RegisterDatabase in init function) in a mandatory Database module.

## type ModuleService

```go
type ModuleService interface {
	FindOne(uid string, id int, params MSParams) (Entity, error)
	FindMany(uid string, params MSParams) ([]Entity, error)
	Create(uid string, data Entity, params MSParams) (Entity, error)
	Update(uid string, data Entity, id int, params MSParams) (Entity, error)
	Delete(uid string, id int, params MSParams) (Entity, error)
}
```

ModuleService will be registered using blank imports (call RegisterModuleService in init function) in a mandatory ModuleService module.

## type Server

```go
type Server interface {
	Start(port string) error
}
```

Server will be registered using blank imports (call RegisterServer in init function) in a mandatory Server module.

## type Controller

```go
type Controller map[string]Action
```

Controller is a grouping of similar actions.

## type Action

```go
type Action func(Cosys, context.Context) http.HandlerFunc
```

Action handles incoming server requests.

Action receives a Cosys struct and Context struct containing the http.Request struct, the http.ResponseWriter struct, and a state variable of any type.

## type Middleware

```go
type Middleware func(Cosys, context.Context) func(http.HandlerFunc) http.HandlerFunc
```

Middleware runs processes before and after the controller handles incoming requests.

Middleware receives a Cosys struct and Context struct containing the http.Request struct, the http.ResponseWriter struct, and a state variable of any type.

## type Policy

```go
type Policy func(Cosys, context.Context) bool
```

Policy checks specific conditions before the controller handles incoming requests.

Policy receives a Cosys struct and Context struct containing the http.Request struct, the http.ResponseWriter struct, and a state variable of any type.

## type Route

```go
type Route struct {
	Method      string
	Regex       *regexp.Regexp
	Action      string
	Middlewares []string
	Policies    []string
}
```

Route specifies which controller and action should handle an incoming request based on its path and method, and can be configured with policies and middleware.

Request paths are specified using regex, and actions are specified using their uid `controller_name.action_name`.

## type Lifecycle

```go
type Lifecycle map[string]LifecycleFunc

type LifecycleFunc func(params DBParams, result any, state any) (afterState any, err error)
```

Lifecycle is the collection of lifecycle functions of a content type. LifecycleFunc is a function called before or after certain actions are performed on a content type. LifecycleFunc takes in the DBParams passed to the Database, the result of the Database call (only available for “after” actions, set to `nil` for “before” actions), and a state variable that is initialised as `nil`  and will be passed from the “before” to the “after” action.

## type Model

```go
type Model interface {
	// Returns an entity of the corresponding content type with zero values.
	New_() Entity
	
	// Returns an array containing all the attributes of the content type.
	All_() []Attribute
	
	// Returns the id attribute of the content type.
	Id_() *IntAttribute
	
	// Returns the schema name of the content type.
	Name_() string
	
	// Returns the schema of the content type.
	Schema_() *ModelSchema
	
	// Returns the lifecycle of the content type.
	Lifecycle_() Lifecycle
}
```

Model represents a content type.

## type Entity

```go
type Entity interface {
}
```

Entity represents an instance of a content type.

## type Attribute

```go
type Attribute interface {
	// Returns the schema name of the attribute.
	Name() string
	
	// Returns the name of the field in the corresponding entity struct.
	FieldName() string

	// Returns the "ascending" order-by condition for this attribute. 
	// To be used in DBParams or MSParams.
	Asc() *Order
	
	// Returns the "descending" order-by condition for this attribute. 
	// To be used in DBParams or MSParams.
	Desc() *Order

	// Returns the "is null" where condition for this attribute. 
	// To be used in DBParams or MSParams.
	Null() Condition
	
	// Returns the "is not null" where condition for this attribute. 
	// To be used in DBParams or MSParams.
	NotNull() Condition
}
```

Attribute represents an attribute of a content-type.