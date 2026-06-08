# TCP Demo
## 通信协议
- 需要注意MSS分段大小
  - MSS= MTU - ip首部 - tcp首部，MTU根据链路接口层的不同而不同
  - TCP传输的时候需要双方协商MSS值
  - 应用层传输数据大于MSS时，需要将数据分段传输
  - 报文格式(segment)：TCP首部 + TCP数据[不能超过MSS]
- 链路层的MTU大小
