# Repository

Repository is a Go open source library that uses generics to implement a repository layer
using [GORM V2](https://gorm.io). It provides a simple and consistent interface for accessing data from a database,
making it easy to write reusable and maintainable code.

# Getting started

To get started with this package, you first need to add it to your dependencies:

```
go get github.com/marcoshuck/repository
```

Once Repository is installed, you can create a new repository instance and start interacting with your database right
away,
without worrying about SQL.

```go
type User struct {
    gorm.Model
    FirstName string
    LastName string
}
    
func main() {
    const dsn = "..."
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalln("Failed to open MySQL connection:", err)
    }
    
    userRepository := NewRepositorySQL[User](db)
    
    ctx := context.Background()
    user := User{
        FirstName: "Marcos",
        LastName: "Huck",
    }
    
    result, err := userRepository.Create(ctx, user)
    if err != nil {
        log.Fatalln("Failed to create user:", err)
    }
    log.Println("User created:", result)
}
```