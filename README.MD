<div align="center">
  <img src="https://github.com/martijnkorbee/gobaboon/blob/master/gopher.png" width="150">
  <h1 align="center">Gobaboon</h1>
</div>

<div>
  <p>
  💥 Gopher by: https://github.com/MariaLetta/free-gophers-pack<br>
  🐍 CLI-tools with: https://github.com/spf13/cobra<br>
  🕸️ Routing with: https://github.com/go-chi/chi<br>
  📰 Logging with: https://github.com/rs/zerolog and https://github.com/natefinch/lumberjack<br>
  🏠 Databases with: https://github.com/upper/db<br>
  💡 Sessions with: https://github.com/alexedwards/scs<br>
  🛑 CSRF protection with: https://github.com/justinas/nosurf<br>
  </p>
</div>

## GOTO
* [Bobo-cli](#bobo-cli)
* [Getting started](#getting-started)
* [Authentication](#authentication) 
* [Sessions](#sessions)
* [Middleware](#middleware)
* [Routes and handlers](#routes-and-handlers)

## Database support
Supported databases:
* Postgresql
* Mysql/MariaDB
* SQLite3
* SQL Server -> (not supported by bobo-cli)

## Caching
Supported caching:
* Redis
* BadgerDB

## Bobo-cli
Bobo is a quality of life improvement for building web apps and services.<br>
Refer to the full documentation: [docs](https://github.com/martijnkorbee/gobaboon/tree/master/cmd/cli/bobo)

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
