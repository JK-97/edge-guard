# jxcore

## 使用

go >= 1.13

```shell
git submodule update --init --recursive
```

### bootstrap

* 在第一次使用前需要执行bootstrap生成设备信息(/edge/init)
    * local:本地模式：不会进行心跳,不需要输入ticket
    * wireguard，vpn模式
    * openvpn，vpn模式
    * 心跳断联会重新获取 vpn配置
* 测试ticket: jiangxing123
* 换模式需要删除之前生成的init文件，再执行bootstrap。
* bootstrap命令会有装机过程.

```shell
jxcore bootstrap -m {local,wireguard,openvpn} -t {ticket}
```

### serve

* 默认读取同目录下的settings.yaml文件，启动组件
* setting.yaml文件，只需要将不需要启动的组件注释掉，或改为false
* jxcore会自动拉起异常退出，或被杀死的组件

```shell
jxcore serve {--port} {--config}
```

## 接口

* /edgenode/exec/restore 
* /edgenode/exec/clean
* /edgenode/exec/update
* /edgenode/exec/migrate
* /edgenode/exec/reload

* /edgenode/version
* /edgenode/changelog
* /edgenode/componentstate

## 约定
### 目录结构

* 每个目录结构与yaml配置文件一一对应。
* component目录下会有bin文件夹与version描述文件
```
 edge
 ├── cri
 ├── devicemanagement
 │   ├── aiserving
 │   ├── camera
 │   │   ├── bin
 │   │   └── version
 │   └── rs485
 │       └── version
 ├── mnt
 ├── monitor
 │   └── telegraf
 │       ├── bin
 │       │   ├── telegraf
 │       │   └── telegraf.cfg
 │       └── version
 ├── synctools
 │   ├── config
 │   ├── db
 │   ├── fs
 │   ├── mq
 │   ├── tsdb
 │   └── vpn
 ├── synctools.zip
 ├── tools
 │   ├── mcutools
 │   │   ├── mcuserver
 │   │   ├── powermanagement
 │   │   └── watchdog
 │   ├── nettools
 │   │   └── ifpllugd
 │   └── nodetools
 │       ├── cleaner
 │       ├── filelistener
 │       │   ├── bin
 │       │   │   ├── filelistener
 │       │   │   └── filelistener.cfg
 │       │   └── version
 │       └── usblistener
 │           ├── bin
 │           │   ├── usblistener
 │           │   └── usblistener.cfg
 │           └── version
 └── version
```

### version 描述文件
version文件更改时会自动触发 changelog的更替

```yaml
name: xxx
version: xxx
```

### bin运行
统一运行命令
```shell script
xxxx -c xxxx.cfg
```


### component 更新

* 更新的文件以zip压缩包的形式
* 压缩宝与目录结构要一致
* 会对旧版本进行覆盖
```
/bin/
version
```
### componentstate
返回当前的运行状态
```
{
    "data": {
        "filelistener": "STARTING",
        "telegraf": "RUNNING"
    },
    "desc": "success"
}
```
### reload
重新载入所有conponent的配置文件


### clean

* 删除所有的docker 容器与镜像
* 删除所有的sdk包
* 删除所有的指定目录
