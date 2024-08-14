# cosys_cli - commands
This package contains the commands for the cosys cli tool.

## Commands

### New
`cosys new project_name -M module_name [flags]`

Create a new cosys project.

### Build
`cosys build [flags]`

Build Golang binaries and Content Management UI deployment

### Start
`cosys start [flags]`

Deploy the Golang server and the Content Management UI server

### Run
`cosys run command_name [arguments] -- [flags]`

Run commands registered from modules.

### Get Configuration
`cosys config get config_name`

Get configurations for the cli tool.

### Set Configuration
`cosys config set config_name config_value`

Set configurations for the cli tool.