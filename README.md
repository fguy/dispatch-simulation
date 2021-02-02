Dispatch Simulation
======================
# Author
Taehoon Kang (flowerguy@gmail.com)

# How to operate the system
You need Go: https://golang.org/doc/install

Please place this project directory `dispatch-simulation` under `$GOPATH/src/github.com/fguy`

## Running the simulation (No build required)
```
go run .
```
## Running tests
```
go test -v -cover ./...
```
## Build (Drops a binary in the working directory)
```
go build
```

# How design decisions are made
I chose Go programming language because of the real-time requirements. Go has elegant, lightweight channels to meet the requirements.
I tried to use vanilla Go, but Uber FX for dependency injection. It simplified dependancy graph, thus overall application complexity got lower.

# How to adjust the configurations
You can change the dispatch strategy, arrival time range and logging configs by editing `configs/base.yaml`.
You can define other enviroments by creating `production.yaml`, `qa.yaml`, `alpha.yaml`, `local.yaml`, `integration.yaml` files. They overrides configurations from `base.yaml`. The enviroment can be set with `APP_ENV` var. e.g. `export APP_ENV=production` to use `production.yaml`