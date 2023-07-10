# Gobaboon
Baboon aims to kickstart the development of your go web and api projects without dependency on this project. 

You can start your own project by cloning or forking this repo or use baboonctl to kickstart your project. Additionally, with the baboonctl you can easily add user and token based authentication, models, persistent session store and caching.

## GOTO
* [baboonctl](#baboonctl)
* [Getting started](#getting-started)
* [Authentication](#authentication)
* [Sessions](#sessions)
* [Middleware](#middleware)
* [Routes and handlers](#routes-and-handlers)

<hr>

### Working on updates
> Working on this project to include the things I've learned after building and deploying a REST API made on top of gobaboon.

- [x] Update project structure
- [ ] After creating a new project there should be no dependency on this repo:
    - [x] Move baboonapp template to this repo (generated project)
    - [x] Update the generated project structure
    - [x] Include baboon logic in the result (merge baboon.go and app.go as result)
    - [ ] Expand application struct with pkgs functionality
    - [x] Check/Fix/Test baboonctl make commands
    - [ ] Check/Fix/Test baboonctl migrate commands
    - [ ] Check/Fix/Test baboonctl rpc commands
    - [ ] Document new result of creating a new project
- [x] Add easy containerization support
- [ ] Check/fix rpc issue #21
- [ ] Revisit config loading and default values
- [ ] Revisit logger pkg (include new slog from 1.21?)

<hr>

### Used libraries
| subject         | library                                                         |
|-----------------|-----------------------------------------------------------------|
| baboonctl       | [spf13/cobra](https://github.com/spf13/cobra)                   |
| logging         | [rs/zerolog](https://github.com/rs/zerolog)                     |
| log rotating    | [natefinch/lumberjack](https://github.com/natefinch/lumberjack) |
| Routing         | [go-chi/chi](https://github.com/go-chi/chi)                     |
| CSRF protection | [justinas/nosurdf](https://github.com/justinas/nosurf)          |
| Databases       | [upper/db](https://github.com/upper/db)                         |
| Sessions        | [alexedwards/scs](https://github.com/alexedwards/scs)           |
   
### Supported databases
| Type            | Comments                                   |
|-----------------|--------------------------------------------|
| Postgresql      |                                            |
| Mysql - MariaDB |                                            |
| SQLite3         |                                            |
| SQL Server      | _not support with baboonctl make commands_ | 

### Supported caching
| Type     | Comments                                   |
|----------|--------------------------------------------|
| Redis    |                                            |
| BadgerDB |                                            |

## baboonctl
Refer to the full documentation: [docs](https://github.com/martijnkorbee/gobaboon/tree/master/tools/baboonctl)

## Getting started

#### Clone the repo.
```bash
git clone https://github.com/martijnkorbee/gobaboon
```

#### CD to the repo directory.
```bash
cd gobaboon
```

#### Build or install baboonctl
Build will compile baboonctl in ./bin.  
After compiling you can cp or mv the binary to your preferred bin location.
```bash
make baboonctl_build 
```

Install will build the file to your $GOPATH.
```bash
make baboonctl_install
```

#### Create a new project
Create a new project in the current directory.
```bash
bobo make new --name [application name]
```

## Authentication
If you want you can add the pre-made authentication to your project.  
This will create up and down migrations for the auth tables and add user and token models in the models directory.  
Should be called from the application's root directory.
> Note: after running the command - _add the models in the models struct and activate the middleware on your desired routes._

```bash
bobo make auth --db-type [value]
```

## Sessions
If you want you can activate the persistent session store to your project.  
Creates up and down migrations for the server's persistent sessions.
> Note: _Set the persistent session store type in the .config.properties file._

```bash
bobo make session --db-type [value]
```

## Middleware
Add your middleware to http/middleware.
* CSRF protection is automatically enforced by the baboon server, exempts all API routes.
* CORS is automatically enabled by the baboon server on API routes.

## Routes and Handlers
Add your routes and handlers to their respective folders in the http directory.
