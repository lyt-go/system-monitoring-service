# Bug Reproduction

## Bug

采集器正在 active 状态时直接删除会成功，之后监控端还会继续把它当成可运行采集器。

## Trigger

1. 先准备依赖实体或边界输入
2. 再执行目标请求并读取后续状态

## Observed Error

断言目标异常被拒绝
