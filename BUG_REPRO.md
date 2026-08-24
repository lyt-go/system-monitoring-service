# Bug Reproduction

## Bug

删除仍被告警引用的阈值时，接口虽然成功，已有告警却失去了有效的阈值关联。

## Trigger

1. 先准备依赖实体或边界输入
2. 再执行目标请求并读取后续状态

## Observed Error

断言目标异常被拒绝
