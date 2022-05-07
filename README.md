# portscan
全端口扫描器

结合Go并发特点 实现的全端口扫描器

## 使用说明
* 支持命令行传递IP地址
```
go run .\portScan.go -domain="IP地址"
```

* 不同系统的运行包已打
> windows amd64架构
> ```go
> ./port_windows_amd64.exe -domain="IP地址"```
