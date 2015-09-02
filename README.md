###本地文件变更监控

---

基于 `fsnotify` 的监控功能, 实现文件或文件夹新增, 修改, 删除, 重命名的变更监控. 相关变更内容, 提供电子邮件发送到目的邮箱上.

获取依赖
```bash
 $ go get github.com/go-fsnotify/fsnotify
```

