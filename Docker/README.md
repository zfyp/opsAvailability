# Docker跨主机通讯

```
graph LR
A[Container] --- B[eth0@if14]
A --- C[eth1@if16]
B --- D[Overlay Network: mynet]
C --- E[docker_gwbridge: 172.18.0.1/16]
```