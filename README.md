# Containers

## Basic Usage

```go
import (
	"context"
	
	"github.com/creawill-com/containers"
	"github.com/creawill-com/containers/driver"
)

func Test_FeeMethod(t *testing.T) {
    ctx := context.Background()
    
    c := containers.NewGorm(ctx, driver.NewPostgres())
    defer func () { c.Terminate(ctx) }()
	
    c.DB().AutoMigrate(
        ...
    ),
	
	...
}
```