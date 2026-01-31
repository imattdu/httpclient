# httpclient

`httpclient` 是一个 **基于 Middleware 的 Go HTTP Client**，用于构建**可组合、可测试**的 HTTP 调用逻辑。  
当前实现重点放在 **清晰分层、可扩展性、以及客户端侧治理能力**。

> 本 README 仅描述 **当前代码中已经存在或已经实现的能力**，不包含规划或假设功能。

---

## 目录结构

```text
httpclient/
├── client/        # Client 实现，负责组合 Handler 和 Middleware
├── transport/     # Transport 抽象（基于 http.Client）
├── middleware/    # Middleware（Retry / Log 等）
├── request/       # 请求抽象
├── response/      # 响应抽象
├── codec/         # 编解码接口（JSON / Raw 等）
├── err/           # 统一错误定义
```

---

## 核心设计

### 分层职责

- **Client**
	- 负责将 Transport 与 Middleware 组合成最终的 Handler
- **Transport**
	- 只负责发送 `*http.Request`
	- 不包含重试、熔断等策略逻辑
- **Middleware**
	- 基于 `Handler -> Handler` 的函数式组合
	- 承载重试、熔断、Mock 等横切逻辑
- **Request / Response**
	- 对 HTTP 请求 / 响应的抽象封装
- **Codec**
	- 负责请求体编码、响应体解码

---

## 快速开始

```go
client := client.NewDefault()

req := request.New(
    "GET",
    "https://httpbin.org/get",
)

resp, err := client.Do(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(resp.RawBody))
```

---

## Request

### 创建请求

```go
req := request.New(
    "POST",
    "https://api.example.com/items",
    request.WithHeader("Authorization", "Bearer xxx"),
    request.WithBody(
        map[string]any{"name": "demo"},
        codec.NewJSON(),
    ),
)
```

### Query 参数

```go
req := request.New(
    "GET",
    "https://api.example.com/search",
    request.WithQuery("q", "golang"),
)
```

### 流式请求

```go
file, _ := os.Open("large.json")

req := request.New(
    "POST",
    url,
    request.WithBody(file, codec.NewRaw()),
)
```

当请求体为流式数据时，编码过程将直接写入 `io.Writer`，不会在内存中缓存完整 body。

---

## Response

### 读取原始响应体

```go
fmt.Println(string(resp.RawBody))
```

### 解码响应体

```go
var out Result
if err := resp.Decode(&out); err != nil {
    log.Fatal(err)
}
```

### 流式响应

```go
req := request.New(
    "GET",
    url,
    request.WithStream(),
)

resp, err := client.Do(ctx, req)
defer resp.Stream.Close()

io.Copy(os.Stdout, resp.Stream)
```

当 `Stream=true` 时：

- Client 不会提前读取 `http.Response.Body`
- 由调用方自行负责读取和关闭

---

## Codec

### Codec 接口

```go
type Codec interface {
    ContentType() string
    Encode(w io.Writer, v any) error
    Decode(r io.Reader, v any) error
}
```

### JSON Codec

```go
codec.NewJSON()
```

### Raw Codec（透传 io.Reader / []byte）

```go
codec.NewRaw()
```

---

## Transport

### Transport 接口

```go
type Transport interface {
    Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
```

### HTTPTransport

```go
transport := transport.NewHTTPTransport(http.DefaultClient)
client := client.New(transport)
```

---

## Middleware

### Retry

```go
client := client.New(
    transport,
    middleware.Retry(middleware.RetryConfig{
        MaxRetries: 3,
        Interval:   100 * time.Millisecond,
    }),
)
```

Retry 仅针对可重试错误生效，流式请求不会参与重试。

---

### 





---

### Mock（测试使用）

```go
client := client.New(
    nil,
    middleware.Mock(func(
        ctx context.Context,
        req *request.Request,
    ) (*response.Response, error) {
        return &response.Response{
            StatusCode: 200,
            RawBody:    []byte(`{"ok":true}`),
        }, nil
    }),
)
```

Mock Middleware 可以在不引入网络的情况下测试 Client / Middleware 行为。



---

## 错误处理

所有错误统一为：

```go
*httpclient.Error
```

通过 `ErrorKind` 区分类型：

- network
- timeout
- http
- encode
- decode
- config

支持 `errors.Is` / `errors.As`。

---

## 设计原则

- Middleware 承载策略
- Transport 只负责 IO
- 显式优于隐式
- 组合优于继承
- 易测试优先

---

## License

MIT
