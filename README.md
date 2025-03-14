# NapCatShellUpdater

###### 仅考虑适配最新版
###### 开发初衷只是为了更加方便的更新Windows平台下的NapCat.Shell所以只考虑Windows平台

---

* 自动更新你的NapCat.Shell
* 更新后自动快速登录你的NapCat(支持从logs文件获取)

```
将程序放在你的NapCat.Shell工作目录双击打开

或使用命令./NapCatShellUpdater
-proxy=http://url:port(遵循proxy url规范)
-version=v0.0.0(指定更新版本, 留空自动获得最新版)
-debug=true(debug模式日志)
-login=true(更新完成后等待进程启动自动登录NapCat)
-skipcheck=false(跳过版本检查, 比如仅登录NapCat)
-ncpanel=http://127.0.0.1:6099(NapCat默认url, 留空自动从logs获取)
-nctoken=napcat(NapCat默认token, 留空自动从logs获取)
-sleep=30s(等待NapCat登录超时的时间, 在此之后NapCatShellUpdater自动尝试登录)
```

---

 - [NapNeko/NapCatQQ](https://github.com/NapNeko/NapCatQQ)