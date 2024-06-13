# Cosys - Common
The package to be imported in cosys projects.

## Documentation

### Types

**type Action**

Action handles incoming server requests. 

Implemented as a function that takes in a Cosys and context.Context struct and returns a http.HandlerFunc. 

**type Controller**

Controller is a grouping of similar actions. 

Implemented as a string map of actions.

**type Cosys**

Cosys is the main app, where the database, module service, server, models, services, configs can be accessed from.

When the Cosys app is run, the Cosys struct undergoes the following procedure.
1. Initialization - The Cosys struct is initialized, the config yaml files are parsed and loaded into the struct.
2. Registration - The modules are registered. The database, module service, server and the models and services from each module are loaded into the Cosys struct. The 'OnRegister' function of each module is run.  
3. Listening - The server from the Cosys struct is started.
4. Destruction - Clean-up processes before the app is shutdown. The 'OnDestroy' function of each module is run.
   
**type Database**

Database provides an API to interact with the database layer directly.

**type Lifecycle**

Lifecycle is the collection of lifecycle functions of a content type.

Implemented as a string map of LifecycleFunc.

**type LifecycleFunc**

LifecycleFunc is a function called before or after certain actions are performed on a content type.

Implemented as a function that takes in an Event struct.

**type Middleware**

Middleware runs processes before and after the controller handles incoming requests.

Implemented as a function that takes in a Cosys and context.Context struct and returns a middleware function.

**type Model**

Model represents a content type.

**type Policy**

Policy checks specific conditions before the controller handles incoming requests.

Implemented as a function that takes in a Cosys and context.Context struct and returns a bool value.

**type Route**

Route specifies which controller and action should handle an incoming request based on its path and method.

Routes can be configured with policies and middleware.

Request paths are specified using regex, and actions are specified using their uid `controller_name.action_name`.

**type RouteOption**

RouteOption allow routes to be configured with policies and middleware.

### Functions

**func NewRoute**

Creates a new Route. Takes in the route method, route path as a regex string, action uid, and any RouteOptions.

**func UseMiddlewares**

Creates a RouteOption that configures a route to use certain middleware. Takes in any number of middleware uids.

**func UsePolicies**

Creates a RouteOption that configures a route to use certain policies. Takes in any number of policy uids.