## baboonctl make auth

Make table migrations and models for authentication

### Synopsis

Creates up and down migrations for the auth tables, and adds user and token models in models directory.
Should be called from the root directory of your application.

SUPPORTED DATABASES: [postgres, mysql, mariadb, sqlite]


```
baboonctl make auth [flags]
```

### Options

```
  -t, --db-type string   specify your database type
  -h, --help             help for auth
```

### SEE ALSO

* [baboonctl make](baboonctl_make.md)	 - Make all kinds of things

###### Auto generated by spf13/cobra on 10-Jul-2023
