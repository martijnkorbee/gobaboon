<div align="center">
  <h1 align="center">Gobaboon</h1>
</div>

<div>
  <p align="left">Baboon aims to kickstart the development of go web and api applications by laying out a project structure and doing the scafolding for tasks that could be seen as repetitive in various projects.</p>
  <p>At the moment of writing I am using Baboon for a (real-world) REST API.</p>
</div>

<div>
  <p>The Baboon project evolved from https://github.com/martijnkorbee/goracoon which is created following a golang course. Compared to Racoon, the project structure and inner workings had a major overhaul, inspired by: https://github.com/golang-standards/project-layout.</p>
  <p>π  Gopher by: https://github.com/MariaLetta/free-gophers-pack<br></p>
</div>

<div>
  <h2>Tech-stack</h2>
  <p>
  π CLI-tools with: https://github.com/spf13/cobra<br>
  πΈοΈ Routing with: https://github.com/go-chi/chi<br>
  π° Logging with: https://github.com/rs/zerolog and https://github.com/natefinch/lumberjack<br>
  π  Databases with: https://github.com/upper/db<br>
  π‘ Sessions with: https://github.com/alexedwards/scs<br>
  π CSRF protection with: https://github.com/justinas/nosurf<br>
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
βββ cmd
βΒ Β  βββ web
βΒ Β      βββ app
βΒ Β      βΒ Β  βββ app.go
βΒ Β      βΒ Β  βββ init-app.go
βΒ Β      βββ bin
βΒ Β      βΒ Β  βββ testapp
βΒ Β      βββ main.go
βββ database
βΒ Β  βββ models
βΒ Β      βββ models.go
βββ go.mod
βββ go.sum
βββ http
βΒ Β  βββ handlers
βΒ Β  βΒ Β  βββ handlers-api.go
βΒ Β  βΒ Β  βββ handlers.go
βΒ Β  βββ middleware
βΒ Β  βΒ Β  βββ middleware.go
βΒ Β  βββ routes
βΒ Β      βββ routes-api.go
βΒ Β      βββ routes.go
βββ LICENSE.md
βββ Makefile
βββ public
βΒ Β  βββ static
βΒ Β      βββ html
βΒ Β      βΒ Β  βββ maintenance.html
βΒ Β      βββ images
βΒ Β          βββ gobaboon.jpg
βββ README.md
βββ templates
    βββ mail
    βΒ Β  βββ mail.html.tmpl
    βΒ Β  βββ mail.plain.tmpl
    βββ views
        βββ home.jet
        βββ home.page.tmpl
        βββ layouts
            βββ base.jet
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
