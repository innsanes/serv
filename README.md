# serv

## Prepare
通过`go get`命令下载依赖包
```shell
go get -u github.com/innsanes/serv
```

## Usage

1. 创建结构体符合serv的interface
```go
type TestServ struct {
    // 通过匿名字段继承serv.Service的方法
    // 如果需要在某个阶段执行逻辑, 可以重写方法
    serv.Service
}
```

2. 重写serv的方法
```go
// BeforeServe 在Serve之前执行
func (s *TestServ) BeforeServe() error {
	fmt.Println("BeforeServe")
	return nil
}

// AfterServe 在Serve之后执行
func (s *TestServ) AfterServe() error {
	fmt.Println("AfterServe")
	return nil
}

// Serve 执行服务启动逻辑
func (s *TestServ) Serve() error {
	fmt.Println("Serve")
	return nil
}

// BeforeStop 在Stop之前执行
func (s *TestServ) BeforeStop() error {
	fmt.Println("BeforeStop")
	return nil
}
```

3. 启动服务
```go
func main() {
    serv.Serve(&TestServ{})
}
```