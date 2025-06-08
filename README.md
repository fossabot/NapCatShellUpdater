# NapCatShellUpdater

###### 仅考虑适配最新版NapCat
###### 开发初衷只是为了更加方便的更新Windows平台下的NapCat.Shell所以只考虑Windows平台

---

* 自动更新你的NapCat.Shell

```
将程序放在你的NapCat.Shell工作目录双击打开

或使用命令./NapCatShellUpdater
-path=xxx(指定NapCat.Shell工作目录, 默认为当前目录)
-version=v0.0.0(指定更新版本, 默认留空自动获得最新版)
-download-url=xxx(指定下载url, 默认为: https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip)
-proxy=http://url:port(遵循proxy url规范, 默认留空)
-version=v0.0.0(指定更新版本, 默认留空自动获得最新版)
-debug=true(debug模式日志, 默认启用)
-exclude=xxx(忽略文件, 多个用逗号分隔, 默认为: config,logs,quickLoginExample.bat,update.bat(强制加入更新的压缩包和程序本身)
```

### Login模块删除

新版NapCat的WebUI可以设置QuickLogin了, 所以废弃此功能

---

- [NapNeko/NapCatQQ](https://github.com/NapNeko/NapCatQQ)