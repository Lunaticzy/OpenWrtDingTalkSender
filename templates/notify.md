# 系统运行状态
```
平均负载: {{ .loadAvg }}
内存占用: {{ .memoryUsedPercent }} %
CPU占用：{{ .cpuUsedPercent }} %
运行时间：{{ .uptime }}
```

---

# 设备温度
```
CPU：{{ .cpuTemp}} ℃
```

---

# WAN 口信息
```
接口ip: {{ .wanIP }}

```

