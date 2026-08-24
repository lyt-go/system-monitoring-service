# Bug Reproduction

## Bug

批量确认告警时出现半成功：

## Trigger

1. 创建 open 与 resolved 两种状态告警
2. 分别提交包含缺失 ID 的批次和 resolved 到 open 的批次后重新查询

## Observed Error

断言包含缺失 ID 的批次必须返回 ErrNotFound
