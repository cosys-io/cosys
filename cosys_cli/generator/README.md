# cosys_cli - generator
This package contains tools for code generation.

## Documentation

### Generator
```go
type Generator []Action
```
Generator is a sequence of actions for code generation.

```go
func (g Generator) Generate() error
```
Generate performs the sequence of actions.

```go
func NewGenerator(actions ...Action) Generator
```
NewGenerator returns a new Generator from a sequence of actions.

### Action
```go
type Action interface {
	Act() error
}
```
Action performs a file system change for code generation.

```go
var DeleteIfExists genOption
```
DeleteIfExists is a configuration for actions, specifying that the given file should be deleted if exists.

```go
var SkipIfExists genOption
```
SkipIfExists is a configuration for actions, specifying that the action should be skipped if the given file exists.

```go
var GenHeadOnly genOption
```
GenHeadOnly is a configuration for actions, specifying that only the head of the path should be created.

### New Directory Action
```go
type NewDirAction struct {
    // contains filtered or unexported fields
}
```
NewDirAction is an action that creates a new directory.

```go
func (a NewDirAction) Act() error
```
Act creates a new directory.

```go
func NewDir(path string, options ...genOption) *NewDirAction
```
NewDir returns a new directory action that creates a new directory at the given path.

### New File Action
```go
type NewFileAction struct {
    // contains filtered or unexported fields
}
```
NewFileAction is an action that create a new file.

```go
func (a NewFileAction) Act() error
```
Act creates a new file.

```go
func NewFile(path string, tmplString string, ctx any, options ...genOption) *NewFileAction
```
NewFile returns a new file action that creates a file at the given path using the given template and context.

### Modify File Action
```go
type ModifyFileAction struct {
    // contains filtered or unexported fields
}
```
ModifyFileAction is an action that modifies a file.

```go
func (a ModifyFileAction) Act() error
```
Act modifies a file.

```go
func ModifyFile(path, patternString, tmplString string, ctx any) *ModifyFileAction
```
ModifyFile returns a modify file action that modifies the file at the given path, replacing all text matching the given regex pattern with the given template and context.