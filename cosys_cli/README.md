# Cosys CLI
CLI tools for managing a cosys project.

## Installation
Run the following to install the Cosys CLI binary into `$GOPATH/bin`. 

`go install github.com/cosys-io/cosys/cosys_cli@latest`

Run the following to add `$GOPATH/bin` to `$PATH`.

`export $PATH=$GOPATH/bin:$PATH`



## Commands
### New
`cosys_cli new project_name -P profile_name [flags]`

Create a new cosys project

### Build
`cosys_cli build [flags]`

Build Golang binaries and Content Management UI deployment

### Start
`cosys_cli start [flags]`

Deploy the Golang server and the Content Management UI server

### Generate

Generate files and boilerplate code for new modules, collections...

**Collection Type**

```
cosys_cli generate collectiontype collection_name [attributes] [flags]

Attributes:
    attribute_name:data_type[:options]

Options:
    required                 
    max=max_value
    min=min_value
    maxlength=max_length
    minlength=min_length
    private
    notconfigurable
    default=default_value
    notnullable
    unsigned
    unique

Flags:
    -D, --description string   description of the new content type
    -N, --display string       display name of the new content type (mandatory)
    -P, --plural string        plural name of the new content type (mandatory)
    -S, --singular string      singular name of the new content type (mandatory)
```

Generate files and boilerplate code for a new collection type

### Install/Uninstall

```yaml
cosys_cli install [module_names]

cosys_cli uninstall [module_names]
```

Install and uninstall modules