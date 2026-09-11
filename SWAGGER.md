# Swagger Documentation

## Tổng quan

Project này đã được tích hợp Swagger/OpenAPI documentation để mô tả các API endpoints.

## Các endpoints có Swagger documentation

1. **Generate Password** - `GET /generate-password`
   - Tạo một mật khẩu ngẫu nhiên an toàn với độ dài 12 ký tự
   - Response: `{ "password": "string" }`

2. **Health Check** - `GET /health-check`
   - Kiểm tra trạng thái service và trả về thông tin service
   - Response: `{ "message": "string", "service_name": "string", "instance_id": "string" }`

## Cách sử dụng

### 1. Truy cập Swagger UI

Sau khi start server, bạn có thể truy cập Swagger UI tại:

```
http://localhost:8080/swagger/index.html
```

### 2. Test API trực tiếp từ Swagger UI

- Mở trình duyệt và truy cập URL trên
- Click vào endpoint muốn test
- Click nút "Try it out"
- Click "Execute" để gửi request
- Xem response trả về

## Generate Swagger docs

Khi bạn thay đổi hoặc thêm endpoint mới với Swagger annotations, chạy lệnh sau để regenerate docs:

```bash
swag init -g cmd/api/main.go -o docs
```

## Swagger Annotations

Mỗi handler function có các annotations:

```go
// GeneratePassword godoc
// @Summary      Generate a random password
// @Description  Generate a secure random password with 12 characters
// @Tags         Password
// @Produce      json
// @Success      200  {object}  model.GeneratePasswordResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /generate-password [get]
```

## Files liên quan

- `docs/swagger.json` - Swagger specification (JSON format)
- `docs/swagger.yaml` - Swagger specification (YAML format)
- `docs/docs.go` - Generated Go code chứa swagger docs

## Dependencies

```
github.com/swaggo/swag/cmd/swag@latest
github.com/swaggo/gin-swagger
github.com/swaggo/files
```
