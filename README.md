# System Monitor

纯 Go 标准库实现的系统监控后端服务，零第三方依赖。

## 运行方式

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/metrics | 创建指标 |
| GET | /api/metrics | 指标列表（支持 type/status/keyword 筛选 + 分页） |
| GET | /api/metrics/{id} | 获取指标 |
| PUT | /api/metrics/{id} | 更新指标 |
| DELETE | /api/metrics/{id} | 删除指标 |
| PATCH | /api/metrics/{id}/status | 指标状态流转（active/inactive） |
| POST | /api/collectors | 创建采集器 |
| GET | /api/collectors | 采集器列表（支持 status/host/keyword 筛选 + 分页） |
| GET | /api/collectors/{id} | 获取采集器 |
| PUT | /api/collectors/{id} | 更新采集器 |
| DELETE | /api/collectors/{id} | 删除采集器 |
| PATCH | /api/collectors/{id}/status | 采集器状态流转（active/stopped） |
| POST | /api/thresholds | 创建阈值规则（需关联存在的指标） |
| GET | /api/thresholds | 阈值列表（支持 metric_id/operator/status 筛选 + 分页） |
| GET | /api/thresholds/{id} | 获取阈值 |
| PUT | /api/thresholds/{id} | 更新阈值 |
| DELETE | /api/thresholds/{id} | 删除阈值 |
| PATCH | /api/thresholds/{id}/status | 阈值状态流转（active/disabled） |
| POST | /api/samples | 创建样本（需关联存在的指标） |
| GET | /api/samples | 样本列表（支持 metric_id/host 筛选 + 分页） |
| GET | /api/samples/{id} | 获取样本 |
| DELETE | /api/samples/{id} | 删除样本 |
| POST | /api/samples/batch-delete | 批量删除样本 |
| POST | /api/alerts | 创建告警（需关联存在的指标和阈值） |
| GET | /api/alerts | 告警列表（支持 metric_id/level/status/keyword 筛选 + 分页） |
| GET | /api/alerts/{id} | 获取告警 |
| PUT | /api/alerts/{id} | 更新告警 |
| DELETE | /api/alerts/{id} | 删除告警 |
| POST | /api/alerts/batch-status | 批量更新告警状态 |
| GET | /api/stats/overview | 综合概览（实体总数 + 未关闭告警数） |
| GET | /api/stats/samples-by-metric | 按指标统计样本数 |
| GET | /api/stats/samples-by-host | 按主机统计样本数 |
| GET | /api/stats/top-alert-metrics | TOP N 告警最多指标（默认 N=5） |

## 响应格式

统一返回：
```json
{"code": 0, "message": "ok", "data": ...}
```

错误码映射：
- 400 ValidationError
- 404 ErrNotFound
- 409 ErrConflict
- 500 其他内部错误
