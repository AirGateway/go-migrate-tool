# go-migrate-tool

Go migration tool for mongodb

# Usage

To generate new migration use CLI command

```
go-migrate-tool --new=<value> [--path="<value>"]
```

Import client
```
import gmt "github.com/AirGateway/go-migrate-tool/client"
```

Add in your application code after mongodb connection initialization
```
migrator, err := gmt.New(<MongoConnection>, &gmt.Config{
    DatabaseName:    <MongoDB name>,
    MigrationPath:   "$GOPATH/src/path/to/app/migrations",
})
if err != nil {
    Log.Error("migration tool not loaded: ", err)
    return
}

count, err := migrator.Apply()
if err != nil {
    Log.Error(err)
}
if count > 0 {
    // do actions if new migration successfully applied
}
```

# Documentation

1. `--new` parameter specifies name for new migration file 
1. System automatically creates `migrations` folder in current working directory. You can specify custom path with `--path` parameter 
  
# Links

Read https://docs.mongodb.com/manual/reference/command/ for mongodb commands list

 