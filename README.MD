<div align="center">
  <h1 align="center">Gobaboon</h1>
</div>

<div>
  <p align="left">Baboon aims to kickstart the development of go web and api projects by laying out a project structure and doing the scaffolding for tasks that could be seen as repetitive in various projects.</p>
</div>

<div>
  <p>At the moment of writing I am working on this project to include the things I've learned after building and deploying a REST API made on top of Baboon.</p>

- [x] Update project structure
- [ ] After creating a new project there should be no dependency on this repo:
    - [x] Move baboonapp template to this repo (generated project)
    - [x] Update the generated project structure
    - [ ] Include baboon logic in the result (merge baboon.go and app.go as result)
    - [ ] Check/Fix/Test bobo make commands
    - [ ] Release new project structure and result of creating a new project
- [ ] Check/fix rpc issue #21
- [ ] Revisit config loading and default values
- [ ] Revisit logger pkg (include new slog from 1.21?)

</div>

<div>
  <p>The Baboon project evolved from https://github.com/martijnkorbee/goracoon which is created following a golang course. Compared to Racoon, the project structure and inner workings had a major overhaul, inspired by: https://github.com/golang-standards/project-layout.</p>
  <p>🟠 Gopher by: https://github.com/MariaLetta/free-gophers-pack<br></p>
</div>

<div>
  <h2>Tech-stack</h2>
  <p>
  🐍 CLI-tools with: https://github.com/spf13/cobra<br>
  🕸️ Routing with: https://github.com/go-chi/chi<br>
  📰 Logging with: https://github.com/rs/zerolog and https://github.com/natefinch/lumberjack<br>
  🏠 Databases with: https://github.com/upper/db<br>
  💡 Sessions with: https://github.com/alexedwards/scs<br>
  🛑 CSRF protection with: https://github.com/justinas/nosurf<br>
  </p>
  
  #### Supported databases
  * Postgresql
  * Mysql/MariaDB
  * SQLite3
  * SQL Server -> (not supported by bobo-cli)

  #### Supported caching
  * Redis
  * BadgerDB
</div>

## GOTO
* [Bobo-cli](#bobo-cli)
* [Getting started](#getting-started)
* [Project structure](#project-structure)
* [Authentication](#authentication) 
* [Sessions](#sessions)
* [Middleware](#middleware)
* [Routes and handlers](#routes-and-handlers)

## Bobo-cli
Refer to the full documentation: [bobo-cli docs](https://github.com/martijnkorbee/gobaboon/tree/master/cmd/cli/bobo)

## Getting started
Clone the repo.
```
git clone https://github.com/martijnkorbee/gobaboon
```
CD to the repo directory.
```
cd gobaboon
```
Build will compile the cli in ./cmd/cli/bobo/bin.<br>
After you can cp or mv the bin to your required bin directory.
```
make build_cli 
```
Install will build the file to ${HOME}/bin (the directory should exist and exported to your $PATH.
```
make install_cli
```
Create a new skeleton app in the current directory.
```
bobo make new -n <appname>
```

## Project structure
This is the project structure of a fresh baboon kickstarted project.
```
├── cmd
│   └── web
│       ├── app
│       │   ├── app.go
│       │   └── init-app.go
│       ├── bin
│       │   └── testapp
│       └── main.go
├── database
│   └── models
│       └── models.go
├── go.mod
├── go.sum
├── http
│   ├── handlers
│   │   ├── handlers-api.go
│   │   └── handlers.go
│   ├── middleware
│   │   └── middleware.go
│   └── routes
│       ├── routes-api.go
│       └── routes.go
├── LICENSE.md
├── Makefile
├── public
│   └── static
│       ├── html
│       │   └── maintenance.html
│       └── images
│           └── gobaboon.jpg
├── README.md
└── templates
    ├── mail
    │   ├── mail.html.tmpl
    │   └── mail.plain.tmpl
    └── views
        ├── home.jet
        ├── home.page.tmpl
        └── layouts
            └── base.jet
```

## Authentication
Make sure your database is available (setup the .env), and in your project root directory execute the bobo command.
This will create migration files, runs the migrations and adds the user and token models and middleware to your project.<br>
```
bobo make auth
```
_Note: you still have to add the models in your models type and activate the middleware on your desired routes._

## Sessions
Make sure your cache/database is available (setup the .env). For cache store no extra steps are required.<br>
For SQL DB; in your project root directory execute the bobo command. This will create migration files and runs the migrations.
```
bobo make session
```
_Note: Set persistent session store by configuring the correct session type in the .env file._

## Middleware
Add your middleware to http/middleware.
* CSRF protection is automatically enforced by the baboon server, exempts all API routes.
* CORS is automatically enabled by the baboon server on API routes.

## Routes and Handlers
Add your routes and handlers to their respective folders in the http directory.
