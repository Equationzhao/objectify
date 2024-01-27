# V0

实现了上传/下载

**文件不能同名**

存储为 `object-storage/v0/filename`

日志为 

- storage-v0-YYYYMMDDHHmmss.log 
- storage-v0-info-YYYYMMDDHHmmss.log 
- storage-v0-error-YYYYMMDDHHmmss.log

## API

PUT v0/object/filename
```bash
curl PUT -T file localhost:8080/v0/object/file
```

GET v0/object/filename
```bash
curl GET localhost:8080/v0/object/file
```