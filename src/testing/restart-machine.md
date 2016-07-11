# 用于测试Mdv显示效果
## 问题

- TOC 无法生成
- JSON没有高亮
- sh没有高亮

﻿# 机柜重启与集群恢复步骤

标签（空格分隔）： 集群 重启 还原

---


[TOC]


---

## 关机
进入下面主机，执行

> shutdown -h now

### VMs
#### 管理节点
```sh
192.168.16.20
192.168.16.21
```
#### 用户集群Test
```sh
192.168.16.22
192.168.16.23
192.168.16.24
```
#### 用户集群clusterh
```sh
192.168.16.30
192.168.16.31
192.168.16.34
```
### 物理机
```sh
192.168.16.12
192.168.16.13
```
### 树莓派
```sh
192.168.16.88
```

---

## 开机
### 启动物理机
按电源按钮，等待约5min
### 验证物理机网络
如果物理机无法与局域网其他主机连通，尝试
> iptables -t nat -A POSTROUTING -o ens1f1 -j MASQUERADE
> systemctl restart network

### 启动VMs
#### 物理机12
```sh
ssh root@192.168.16.12
virsh start dcos-1
virsh start dcos-2
virsh start dcos-3
```
#### 物理机13
```sh
ssh root@192.168.16.13
virsh start dcos1
virsh start dcos2
virsh start user-vm1
virsh start user-vm2
virsh start user-vm4
```

### 还原Snapshots
使用**&&**拼接命令，顺序还原master节点与slave节点，能够较少恢复后错误
#### 还原管理节点
```sh
ssh root@192.168.16.13
virsh snapshot-revert dcos1 1467844460 && virsh snapshot-revert dcos2 1467844465
```
#### 还原用户集群Test
```sh
virsh snapshot-revert dcos-1 1467841692 && virsh snapshot-revert dcos-2 1467841692 && virsh snapshot-revert dcos-3 1467841692
```
#### 还原用户集群clusterh
```sh
virsh snapshot-revert user-vm1 1467886058 && virsh snapshot-revert user-vm2 1467886050 && virsh snapshot-revert user-vm4 1467886050
```
### 启动树莓派
接通电池电源，接通MicroUSB口电源

## 验证
### 验证Mesos任务
> 192.168.16.22:5050
> 192.168.16.30:5050

**重点查看Spark Streaming相关Task是否为RUNNING**
如果相关Task不存在，进入集群管理Master节点 --> 容器client --> 执行下面命令：
```sh
dcos spark run --submit-args="--driver-cores 0.5 --driver-memory 2048M --executor-memory 2048M --class IoTStreaming http://192.168.16.20:8081/artifacts/spark/main_streaming_iot.jar --zookeeper master.mesos:2181 --consumer_group iot-streaming --topics iot --max_cores 2 --batch_interval 10 --db_name iot"
```

### 验证Grafana数据
#### 登录
admin/123456
#### 载入IoT.json
左上角Dashboard -> Import
选择 IoT.json，如果找不到文件，见附录
#### 验证数据
查看Disk占用是否显示，查看PM2.5有没有数据

## 其他
### 切换树莓派连接集群
#### 登录
ssh pi@192.168.16.88
密码：raspberry
sudo -s切换至root
#### 修改文件
- /boot/interfaces
- /boot/linkeriot
- /boot/.logagent
- /root/change-ip-port-docker.sh
- **执行脚本change-ip-port-docker.sh**
- reboot

如果树莓派中docker无法启动，删除目录/var/lib/docker/network，重启

### 如需同步时区
> **Ubuntu/Debian:**
> dpkg-reconfigure tzdata
> **CentOS/RedHat:**
> ln -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
> 或者
> timedatectl set-timezone Asia/Shanghai

## 附录
### IoT.json
```json
{
  "id": 1,
  "title": "IoT",
  "originalTitle": "IoT",
  "tags": [
    "linker",
    "iot",
    "pm",
    "usage",
    "humid",
    "temperature"
  ],
  "style": "dark",
  "timezone": "browser",
  "editable": true,
  "hideControls": true,
  "sharedCrosshair": false,
  "rows": [
    {
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "aliasColors": {
            "pm.pm25": "#AEA2E0"
          },
          "bars": false,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)",
            "thresholdLine": false
          },
          "hideTimeOverride": false,
          "id": 1,
          "isNew": true,
          "legend": {
            "avg": false,
            "current": false,
            "hideEmpty": false,
            "hideZero": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": false,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "null as zero",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 6,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "select pm25 from pm where time > '2016-06-20'",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "PM2.5",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "value_type": "cumulative"
          },
          "transparent": true,
          "type": "graph",
          "xaxis": {
            "show": true
          },
          "yaxes": [
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ]
        },
        {
          "aliasColors": {
            "hddusage.avail": "#B7DBAB",
            "hddusage.used": "#F29191"
          },
          "cacheTimeout": null,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fontSize": "80%",
          "format": "kbytes",
          "id": 4,
          "interval": null,
          "isNew": true,
          "legend": {
            "show": false,
            "values": false
          },
          "legendType": "Right side",
          "links": [],
          "maxDataPoints": 3,
          "nullPointMode": "connected",
          "pieType": "pie",
          "span": 6,
          "strokeWidth": "0",
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "SELECT avail,used FROM hddusage",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "title": "Disk Usage",
          "transparent": true,
          "type": "grafana-piechart-panel",
          "valueName": "current"
        }
      ],
      "showTitle": false,
      "title": "Row"
    },
    {
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "aliasColors": {
            "dht.humid": "#F9934E"
          },
          "bars": false,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "id": 3,
          "isNew": true,
          "legend": {
            "avg": false,
            "current": false,
            "hideEmpty": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": false,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "null as zero",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 6,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "SELECT humid FROM dht",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "Humidity",
          "tooltip": {
            "msResolution": false,
            "shared": true,
            "value_type": "cumulative"
          },
          "transparent": true,
          "type": "graph",
          "xaxis": {
            "show": true
          },
          "yaxes": [
            {
              "format": "humidity",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": false
            }
          ]
        },
        {
          "aliasColors": {
            "dht.temp": "#E24D42"
          },
          "bars": false,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "id": 2,
          "isNew": true,
          "legend": {
            "alignAsTable": false,
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": false,
            "sideWidth": null,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "null as zero",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 6,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "SELECT temp FROM dht",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "Temperature",
          "tooltip": {
            "msResolution": false,
            "shared": true,
            "value_type": "cumulative"
          },
          "transparent": true,
          "type": "graph",
          "xaxis": {
            "show": true
          },
          "yaxes": [
            {
              "format": "celsius",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ]
        }
      ],
      "title": "New row"
    },
    {
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "aliasColors": {
            "pm.pm10": "#9AC48A",
            "pm.pm25": "#AEA2E0"
          },
          "bars": false,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)",
            "thresholdLine": false
          },
          "id": 5,
          "isNew": true,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": false,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "null as zero",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 6,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "select pm10 from pm where time > '2016-06-20'",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "PM10",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "value_type": "cumulative"
          },
          "transparent": true,
          "type": "graph",
          "xaxis": {
            "show": true
          },
          "yaxes": [
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ]
        },
        {
          "aliasColors": {
            "pm.pm1": "#70DBED",
            "pm.pm25": "#AEA2E0"
          },
          "bars": false,
          "datasource": "IoT",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)",
            "thresholdLine": false
          },
          "id": 6,
          "isNew": true,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": false,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "null as zero",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 6,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "dsType": "influxdb",
              "groupBy": [
                {
                  "params": [
                    "$interval"
                  ],
                  "type": "time"
                },
                {
                  "params": [
                    "null"
                  ],
                  "type": "fill"
                }
              ],
              "policy": "default",
              "query": "select pm1 from pm where time > '2016-06-20'",
              "rawQuery": true,
              "refId": "A",
              "resultFormat": "time_series",
              "select": [
                [
                  {
                    "params": [
                      "value"
                    ],
                    "type": "field"
                  },
                  {
                    "params": [],
                    "type": "mean"
                  }
                ]
              ],
              "tags": []
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "PM1",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "value_type": "cumulative"
          },
          "transparent": true,
          "type": "graph",
          "xaxis": {
            "show": true
          },
          "yaxes": [
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ]
        }
      ],
      "title": "New row"
    }
  ],
  "time": {
    "from": "now-3h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "1s",
      "2s",
      "3s",
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "templating": {
    "list": []
  },
  "annotations": {
    "list": []
  },
  "refresh": "10s",
  "schemaVersion": 12,
  "version": 13,
  "links": [
    {
      "icon": "external link",
      "tags": [],
      "type": "dashboards"
    }
  ]
}
```




