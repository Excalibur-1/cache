# cache
缓存公共服务封装

## 快速使用
1. 手动传入配置信息

```go
import (
    "context"
    "fmt"
    "github.com/Excalibur-1/cache"
)

func example() {
	ctx := context.Background()
	eng := cache.Engine(ctx, Config{
		Provider: cache.Redis,
		Host:     "localhost",
		Port:     "6379",
	})
	eng.Set(ctx, "test", "test test")
	value := eng.Get(ctx, "test")
	fmt.Println(value)
}

```

2.1 读取配置中心配置(mock配置) 
```go
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/cache"
	"github.com/Excalibur-1/configuration"
)

func example() {
	mockEngine := configuration.MockEngine(configuration.StoreConfig{Exp: map[string]string{
		"/myconf/base/cache/1000": "{\"provider\":\"redis\",\"host\":\"127.0.0.1\",\"port\":\"6379\",\"password\":\"\"}",
	}})
	s, err := mockEngine.String("base", "cache", "", "1000")
	if err != nil {
		panic(err)
	}
	if s == "" {
		panic("配置信息不存在")
	}
	var conf cache.Config
	if err = json.Unmarshal([]byte(s), &conf); err != nil {
		panic(err)
	}
	ctx := context.Background()
	eng := cache.Engine(ctx, conf)
	eng.Set(ctx, "test", "test test")
	value := eng.Get(ctx, "test")
	fmt.Println(value)
}

```

2.2 读取配置中心配置(zookeeper配置,需要事先添加 configuration.uaf 文件，配置方法见 [configuration](https://github.com/Excalibur-1/configuration) 项目)

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/cache"
	"github.com/Excalibur-1/configuration"
)

func example() {
	zkEngine := configuration.ZkEngine(configuration.NewStoreConfig())
	s, err := zkEngine.String("base", "cache", "", "1000")
	if err != nil {
		panic(err)
	}
	if s == "" {
		panic("配置信息不存在")
	}
	var conf cache.Config
	if err = json.Unmarshal([]byte(s), &conf); err != nil {
		panic(err)
	}
	ctx := context.Background()
	eng := cache.Engine(ctx, conf)
	eng.Set(ctx, "test", "test test")
	value := eng.Get(ctx, "test")
	fmt.Println(value)
}

```

## 2021-03-29 更新日志
* 添加redis缓存实现
* 添加默认缓存实现

## 2021-04-06 更新日志
* 优化代码逻辑
* 添加 readme 使用文档说明
* 添加接入配置中心单元测试