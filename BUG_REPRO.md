# Bug Reproduction

## Bug

阈值规则更新有副作用：把一条 gt 规则改成不支持的操作符 bad，同时带上新阈值，调用会报参数错误；但紧接着查询这条规则，操作符和数值已经变成 bad 与 99。当前项目先别动代码，

## Trigger

1. 创建一条可用的 gt 阈值规则
2. 提交 bad 操作符和新数值后分别按 ID 与列表读取

## Observed Error

断言非法操作符触发 ValidationError
