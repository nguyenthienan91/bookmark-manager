# Bookmark Manager - Learning Notes

File này dùng để ghi lại quá trình học theo từng lecture. Mục tiêu là sau này mở lại vẫn nhớ được bài đã làm gì, vì sao làm như vậy, và luồng code chạy qua những file nào.

## Lecture 01 - API Server, Clean Architecture, Health Check Endpoint

### Mục Tiêu Bài Tập

Khởi tạo một API server có endpoint:

```http
GET /health-check
```

Response cần trả về:

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "18c15eaf-f401-45bf-a827-265b9c6c9333"
}
```

Yêu cầu chính:

- Status code phải là `200`.
- `message` luôn là `"OK"`.
- `service_name` được cấu hình bằng environment variable.
- `instance_id` được cấu hình bằng environment variable.
- Nếu không có `instance_id` trong environment variable thì app tự generate UUID.
- Project follow kiến trúc clean architecture.
- Có unit test và integration test.
- App config được đọc từ environment variable.

### Project Hiện Tại Đang Là Gì?

Project là một Go backend API dùng Gin framework.

Hiện tại project có 2 endpoint chính:

```http
GET /generate-password
GET /health-check
```

Endpoint cũ:

- `/generate-password`: generate password ngẫu nhiên.

Endpoint mới trong Lecture 01:

- `/health-check`: kiểm tra service còn sống và trả về thông tin service.

### Clean Architecture Trong Bài Này

Luồng code được tách theo trách nhiệm:

```text
cmd/api/main.go
  -> internal/api/config.go
  -> internal/api/api.go
  -> internal/handler/health.go
  -> internal/service/health.go
  -> internal/model/health.go
```

Ý nghĩa từng layer:

| Layer | Nhiệm vụ |
| --- | --- |
| `cmd/api` | Điểm khởi động app |
| `internal/api` | Load config, tạo Gin engine, đăng ký route |
| `internal/handler` | Nhận HTTP request và trả HTTP response |
| `internal/service` | Chứa business logic |
| `internal/model` | Định nghĩa struct dữ liệu trả về |
| `internal/integration_test` | Test endpoint thông qua HTTP engine |

Điểm cần nhớ:

- Handler không nên tự xử lý business logic phức tạp.
- Service không nên biết HTTP là gì.
- Model giúp response có cấu trúc rõ ràng.
- API layer là nơi nối route với handler.

### Các File Đã Thêm Hoặc Chỉnh

#### 1. `internal/api/config.go`

File này dùng để đọc config từ environment variable.

Config hiện tại gồm:

```go
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-service" envconfig:"SERVICE_NAME"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
}
```

Ý nghĩa:

- `AppPort`: port để chạy server, default là `8080`.
- `ServiceName`: tên service, default là `bookmark-service`.
- `InstanceID`: ID của instance đang chạy.

Do code dùng:

```go
envconfig.Process("api", cfg)
```

nên env chính sẽ có prefix `API_`:

```text
API_APP_PORT
API_SERVICE_NAME
API_INSTANCE_ID
```

Nếu `API_INSTANCE_ID` rỗng, code tự generate UUID:

```go
if cfg.InstanceID == "" {
	cfg.InstanceID = uuid.NewString()
}
```

Package đã thêm:

```go
github.com/google/uuid
```

Lệnh cài:

```powershell
go get github.com/google/uuid
```

#### 2. `internal/model/health.go`

File này định nghĩa response body cho endpoint `/health-check`.

```go
type HealthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}
```

Điểm cần nhớ:

- Field Go thường viết kiểu PascalCase: `ServiceName`.
- JSON trả về thường dùng snake_case: `service_name`.
- Struct tag `json:"service_name"` giúp map field Go sang key JSON.

#### 3. `internal/service/health.go`

File này chứa logic tạo health check response.

Interface:

```go
type HealthCheck interface {
	Check() model.HealthCheckResponse
}
```

Struct implementation:

```go
type healthCheckService struct {
	serviceName string
	instanceID  string
}
```

Constructor:

```go
func NewHealthCheck(serviceName string, instanceID string) HealthCheck
```

Method xử lý logic:

```go
func (h *healthCheckService) Check() model.HealthCheckResponse
```

Điểm cần nhớ:

- Interface giúp code dễ test và dễ thay implementation sau này.
- Constructor nhận `serviceName` và `instanceID` từ config.
- Service chỉ trả data, không gọi `gin.Context`.

#### 4. `internal/handler/health.go`

File này nhận request HTTP và trả response.

Handler có dependency là service:

```go
type healthCheckHandler struct {
	healthCheckService service.HealthCheck
}
```

Method xử lý endpoint:

```go
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	response := h.healthCheckService.Check()

	c.JSON(http.StatusOK, response)
}
```

Điểm cần nhớ:

- `c *gin.Context` là object Gin dùng để đọc request và trả response.
- `http.StatusOK` tương đương status code `200`.
- `c.JSON(status, data)` trả JSON response cho client.

#### 5. `internal/api/api.go`

File này đăng ký route vào Gin engine.

Trong `initRoutes()`, đã thêm:

```go
healthCheckService := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)
healthCheckHandler := handler.NewHealthCheck(healthCheckService)

e.app.GET("/health-check", healthCheckHandler.HealthCheck)
```

Luồng khi client gọi API:

```text
GET /health-check
  -> Gin route
  -> healthCheckHandler.HealthCheck
  -> healthCheckService.Check
  -> model.HealthCheckResponse
  -> JSON response
```

### Tests Đã Thêm

#### 1. Service Test

File:

```text
internal/service/health_test.go
```

Mục tiêu:

- Test `service.Check()` trả đúng `message`.
- Test `service.Check()` trả đúng `service_name`.
- Test `service.Check()` trả đúng `instance_id`.

Ý tưởng:

```text
Tạo service với data cố định
-> gọi Check()
-> assert response đúng
```

#### 2. Handler Test

File:

```text
internal/handler/health_test.go
```

Mục tiêu:

- Test handler trả status code `200`.
- Test JSON body đúng format.

Trong test có fake service:

```go
type fakeHealthCheckService struct{}
```

Fake service giúp test handler mà không phụ thuộc vào service thật.

Điểm cần nhớ:

- Đây là lợi ích của interface.
- Handler chỉ cần object nào có method `Check()` đúng interface là chạy được.

#### 3. Integration Test

File:

```text
internal/integration_test/health_ep_test.go
```

Mục tiêu:

- Test endpoint `/health-check` thông qua API engine thật.
- Không cần start server thật.
- Dùng `httptest.NewRequest` và `httptest.NewRecorder`.

Luồng test:

```text
api.NewEngine(...)
-> tạo request GET /health-check
-> recorder nhận response
-> apiEngine.ServeHTTP(rec, req)
-> assert status và JSON body
```

### Commands Quan Trọng

Cài UUID package:

```powershell
go get github.com/google/uuid
```

Chạy toàn bộ test:

```powershell
go test ./...
```

Chạy app:

```powershell
go run ./cmd/api
```

Chạy app với env variable trong PowerShell:

```powershell
$env:API_SERVICE_NAME = "bookmark_service"
$env:API_INSTANCE_ID = "18c15eaf-f401-45bf-a827-265b9c6c9333"
go run ./cmd/api
```

Gọi endpoint:

```powershell
curl http://localhost:8080/health-check
```

Response mong muốn:

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "18c15eaf-f401-45bf-a827-265b9c6c9333"
}
```

### Kiến Thức Go Cần Nhớ

#### Struct

Struct dùng để gom nhiều field vào một kiểu dữ liệu.

Ví dụ:

```go
type HealthCheckResponse struct {
	Message string
}
```

#### Struct Tag

Struct tag giúp điều khiển cách field được encode sang JSON hoặc đọc từ env.

Ví dụ:

```go
ServiceName string `json:"service_name"`
```

#### Interface

Interface định nghĩa behavior.

Ví dụ:

```go
type HealthCheck interface {
	Check() model.HealthCheckResponse
}
```

Object nào có method `Check()` đúng signature thì thỏa interface này.

#### Constructor Function

Trong Go thường dùng function `New...` để tạo object.

Ví dụ:

```go
func NewHealthCheck(serviceName string, instanceID string) HealthCheck
```

#### Method Receiver

Method receiver cho biết function thuộc về struct nào.

Ví dụ:

```go
func (h *healthCheckService) Check() model.HealthCheckResponse
```

Nghĩa là method `Check()` thuộc về `healthCheckService`.

#### Package

Mỗi folder thường là một package.

Ví dụ:

```text
internal/service
```

có:

```go
package service
```

### Những Điều Cần Chú Ý Sau Bài Này

- Tên service trong đề là `bookmark_service`, nhưng default hiện tại trong code là `bookmark-service`.
- Nếu muốn đúng y hệt đề khi không set env, có thể đổi default thành `bookmark_service`.
- Nên chạy `gofmt` để format code Go sau khi viết:

```powershell
gofmt -w internal/api/config.go internal/api/api.go internal/model/health.go internal/service/health.go internal/handler/health.go internal/service/health_test.go internal/handler/health_test.go internal/integration_test/health_ep_test.go
```

- Nên chạy `go test ./...` sau mỗi bài.
- Nếu test fail, đọc từ dòng đầu tiên có chữ `FAIL` hoặc error message cụ thể.

### Checklist Lecture 01

- [x] Tạo API server bằng Gin.
- [x] Thêm endpoint `GET /health-check`.
- [x] Trả status code `200`.
- [x] Trả response JSON gồm `message`, `service_name`, `instance_id`.
- [x] Đọc `service_name` từ env/config.
- [x] Đọc `instance_id` từ env/config.
- [x] Nếu thiếu `instance_id`, tự generate UUID.
- [x] Tách code theo hướng clean architecture.
- [x] Viết service test.
- [x] Viết handler test.
- [x] Viết integration test.

## Lecture 02 - Shorten URL Endpoint

### Ngày thực hiện

`11/09/2026`

### Mục tiêu bài tập

Tạo endpoint:

```http
POST /v1/links/shorten
```

Endpoint nhận một URL, sinh một mã rút gọn ngẫu nhiên gồm 7 ký tự, sau đó lưu mapping:

```text
short code -> original URL
```

Mapping được lưu trong Redis.

Request mẫu:

```json
{
  "exp": 604800,
  "url": "https://example.com/article"
}
```

Response thành công:

```json
{
  "code": "abc1234",
  "message": "Shorten URL generated successfully!"
}
```

### Những phần đã triển khai

- [x] Thêm dependency Redis bằng `go get github.com/redis/go-redis/v9`.
- [x] Tạo Redis config trong `pkg/redis/config.go`.
- [x] Tạo Redis client trong `pkg/redis/client.go`.
- [x] Tạo interface `Store` trong `pkg/redis/store.go`.
- [x] Cấu hình Redis bằng các biến môi trường:
  - `REDIS_ADDR`
  - `REDIS_PASSWORD`
  - `REDIS_DB`
- [x] Tạo interface và implementation cho repository trong `internal/repository`.
- [x] Dùng `SetNX` để chỉ lưu code nếu key chưa tồn tại, tránh ghi đè URL cũ.
- [x] Dùng prefix `shorturl:` cho key trong Redis.
- [x] Tạo request/response model trong `internal/model/shorten_link.go`.
- [x] Validate `url` bắt buộc và phải có format URL hợp lệ.
- [x] Reject `exp` âm với HTTP `400`.
- [x] Sinh short code dài đúng 7 ký tự.
- [x] Short code chỉ dùng chữ cái và chữ số.
- [x] Dùng `crypto/rand` để sinh code ngẫu nhiên.
- [x] Chuyển `exp` từ giây sang `time.Duration`.
- [x] Quy ước `exp = 0` là không đặt TTL.
- [x] Retry khi short code bị trùng, tối đa 10 lần.
- [x] Tạo handler cho endpoint `POST /v1/links/shorten`.
- [x] Đăng ký route trong group `/v1`.
- [x] Inject service vào API engine để endpoint dễ test.
- [x] Load Redis config, tạo client, ping Redis và đóng client trong `cmd/api/main.go`.
- [x] Regenerate mock bằng Mockery cho Redis store, repository, service và health-check service.
- [x] Tạo repository unit test.
- [x] Tạo service unit test.
- [x] Tạo handler unit test.
- [x] Tạo integration test cho endpoint.
- [x] Cập nhật Swagger documentation cho endpoint mới.

### Luồng xử lý

```text
HTTP POST request
    -> Gin router
    -> handler.ShortenLink
    -> validate request
    -> service.Shorten
    -> generate 7-character code
    -> repository.SaveLink
    -> Redis Store.SetNX
    -> return code and success message
```

### Phân chia trách nhiệm

#### `internal/model`

Chứa struct request và response của endpoint.

#### `internal/handler`

- Đọc JSON request.
- Validate dữ liệu đầu vào.
- Gọi service.
- Trả HTTP status code và JSON response.

Handler không tự sinh code và không gọi trực tiếp Redis.

#### `internal/service`

- Sinh short code.
- Chỉ sử dụng ký tự chữ và số.
- Chuyển expiration từ giây sang duration.
- Retry khi code bị collision.
- Gọi repository để lưu dữ liệu.

#### `internal/repository`

- Tạo Redis key từ prefix `shorturl:` và short code.
- Gọi `SetNX`.
- Trả về trạng thái lưu thành công hoặc key đã tồn tại.

#### `pkg/redis`

- Quản lý cấu hình và kết nối Redis.
- Cung cấp interface nhỏ `Store` để repository dễ mock.
- Thực hiện `Ping`, `Close` và `SetNX`.

### Các trường hợp đã test

- Request hợp lệ có expiration.
- Request hợp lệ không có expiration.
- Thiếu URL.
- JSON không hợp lệ.
- URL không đúng format.
- `exp` là số âm.
- Service trả lỗi.
- Redis trả lỗi.
- Redis key đã tồn tại.
- Redis key có đúng prefix.
- Expiration được truyền đúng.
- Short code có đúng 7 ký tự.
- Short code chỉ chứa chữ và số.
- Retry thành công sau một lần collision.
- Retry thất bại sau đủ 10 lần.

### Bài học rút ra

- Interface giúp tách business logic khỏi implementation cụ thể và giúp viết unit test bằng mock.
- `SetNX` phù hợp hơn `Set` khi short code phải unique, vì không ghi đè dữ liệu đang tồn tại.
- `exp` từ request là số giây, còn Redis client nhận `time.Duration`, nên phải chuyển đổi rõ ràng.
- Integration test có thể inject mock service để kiểm tra route mà không cần kết nối Redis thật.
- Dependency runtime như `go-redis` dùng `go get`; tool sinh mock như Mockery cài bằng `go install`.

## Lecture 03

Chưa học.

## Lecture 04

Chưa học.

## Lecture 05

Chưa học.

## Lecture 06

Chưa học.

## Lecture 07

Chưa học.

## Lecture 08

Chưa học.

## Lecture 09

Chưa học.

## Lecture 10

Chưa học.

## Lecture 11

Chưa học.
