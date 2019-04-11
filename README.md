# x
自用golang工具库

## xconfig
TODO

## xswagger
安装start_swagger
``` bash
# go get github.com/webee/x/xswagger/start_swagger
```
自定义cmd start swagger
``` golang
package main

import (
    "github.com/webee/x/xswagger/cmd"

    //change to your path
    "path/to/your/docs"
)

func main() {
	docs.SwaggerInfo.Host = cmd.GetConfig().APP.Host
	cmd.Start()
}
```