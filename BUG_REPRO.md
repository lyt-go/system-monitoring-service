# Bug Reproduction

## Bug

告警编辑失败后内容被偷偷改了：把一条 open 告警的消息改为 replacement、级别改成 unknown，接口按预期返回级别不合法；可随后查询和列表筛选都能看到 replacement 和 unknown。

## Trigger

1. 创建一条 warn 级别的 open 告警
2. 用 replacement 与 unknown 更新并在失败后查询和列举

## Observed Error

断言 unknown 级别返回 ValidationError
