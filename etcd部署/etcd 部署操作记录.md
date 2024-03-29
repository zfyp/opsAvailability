## etcd安装
### 主机配置(ALL)
172.29.29.3 etcd-29-3

172.29.29.9 etcd-29-9

172.29.29.11 etcd-29-11

将主机IP、Host信息写入/etc/hosts文件
```bash
echo "
172.29.29.3 etcd-29-3
172.29.29.9 etcd-29-9
172.29.29.11 etcd-29-11
" >> /etc/hosts
```

### 下载etcd二进制安装包(ALL)
```bash
#下载
wget https://github.com/etcd-io/etcd/releases/download/v3.4.31/etcd-v3.4.31-linux-amd64.tar.gz
#解压文件
tar zfx etcd-v3.4.31-linux-amd64.tar.gz

# 将etcd  etcdctl两个文件移动到/usr/bin/目录
cd etcd-v3.4.31-linux-amd64
mv etcd  etcdctl /usr/local/bin/

#查看安装是否正常
etcd --version
etcdctl version
```
### 创建etcd账号以及数据存放目录
```bash
# 创建etcd的运行账号
sudo groupadd etcd
sudo useradd -s /sbin/nologin -g etcd etcd

# 数据目录的路径: /data/etcd/data/
# 用于持久化存储的预写日志目录: /data/etcd/data/wal
mkdir -pv /data/etcd/data/wal

#设置目录权限
chown -R etcd:etcd /data/etcd/
```
## 配置文件
### 编辑/etc/systemd/system/etcd.service
<h4> 使用systemd来管理etcd服务</h4>

- 172.29.29.3 etcd-29-3
```shell
[Unit]
Description=etcd service
Documentation=https://github.com/etcd-io/etcd
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \
  --name etcd-29-3 \
  --data-dir /var/lib/etcd \
  --wal-dir /var/lib/etcd/wal \
  --snapshot-count 10000 \
  --heartbeat-interval 100 \
  --election-timeout 1000 \
  --quota-backend-bytes 0 \
  --listen-peer-urls http://172.29.29.3:2380 \
  --listen-client-urls http://172.29.29.3:2379,http://127.0.0.1:2379 \
  --max-snapshots 5 \
  --max-wals 5 \
  --initial-advertise-peer-urls http://172.29.29.3:2380 \
  --advertise-client-urls http://172.29.29.3:2379 \
  --initial-cluster etcd-29-3=http://172.29.29.3:2380,etcd-29-9=http://172.29.29.9:2380,etcd-29-11=http://172.29.29.11:2380 \
  --initial-cluster-token dockerNet \
  --initial-cluster-state new \
  --enable-pprof \
  --auto-compaction-mode periodic \
  --auto-compaction-retention "1" \
  --log-outputs stdout
StandardOutput=journal
Restart=on-failure
RestartSec=5
User=etcd
Group=etcd

[Install]
WantedBy=multi-user.target
```
- 172.29.29.9 etcd-29-9
```shell
[Unit]
Description=etcd service
Documentation=https://github.com/etcd-io/etcd
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \
  --name etcd-29-9 \
  --data-dir /var/lib/etcd \
  --wal-dir /var/lib/etcd/wal \
  --snapshot-count 10000 \
  --heartbeat-interval 100 \
  --election-timeout 1000 \
  --quota-backend-bytes 0 \
  --listen-peer-urls http://172.29.29.9:2380 \
  --listen-client-urls http://172.29.29.9:2379,http://127.0.0.1:2379 \
  --max-snapshots 5 \
  --max-wals 5 \
  --initial-advertise-peer-urls http://172.29.29.9:2380 \
  --advertise-client-urls http://172.29.29.9:2379 \
  --initial-cluster etcd-29-3=http://172.29.29.3:2380,etcd-29-9=http://172.29.29.9:2380,etcd-29-11=http://172.29.29.11:2380 \
  --initial-cluster-token dockerNet \
  --initial-cluster-state new \
  --enable-pprof \
  --auto-compaction-mode periodic \
  --auto-compaction-retention "1" \
  --log-outputs stdout  
StandardOutput=journal
Restart=on-failure
RestartSec=5
User=etcd
Group=etcd

[Install]
WantedBy=multi-user.target
```
- 172.29.29.11 etcd-29-11
```shell
[Unit]
Description=etcd service  
Documentation=https://github.com/etcd-io/etcd
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \
  --name etcd-29-11 \
  --data-dir /var/lib/etcd \
  --wal-dir /var/lib/etcd/wal \
  --snapshot-count 10000 \
  --heartbeat-interval 100 \
  --election-timeout 1000 \
  --quota-backend-bytes 0 \
  --listen-peer-urls http://172.29.29.11:2380 \
  --listen-client-urls http://172.29.29.11:2379,http://127.0.0.1:2379 \
  --max-snapshots 5 \
  --max-wals 5 \
  --initial-advertise-peer-urls http://172.29.29.11:2380 \
  --advertise-client-urls http://172.29.29.11:2379 \
  --initial-cluster etcd-29-3=http://172.29.29.3:2380,etcd-29-9=http://172.29.29.9:2380,etcd-29-11=http://172.29.29.11:2380 \
  --initial-cluster-token dockerNet \
  --initial-cluster-state new \
  --enable-pprof \
  --auto-compaction-mode periodic \
  --auto-compaction-retention "1" \
  --log-outputs stdout
StandardOutput=journal  
Restart=on-failure
RestartSec=5
User=etcd
Group=etcd

[Install]
WantedBy=multi-user.target
```

### 重载配置(ALL)
```bash
#重新加载systemd服务
sudo systemctl daemon-reload

#放入自启动
systemctl enable etcd

#启动
systemctl start etcd
```