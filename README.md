# ORM Benchmark

### ORMs

All package run in no-cache mode, forked & refreshed based https://github.com/yusaer/orm-benchmark

### Run

```go
# setup env
docker-compose up

# all
go run main.go -multi=20 -orm=all

# portion
go run main.go -multi=20 -orm=gorm
```
