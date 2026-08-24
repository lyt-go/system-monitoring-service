# Bug Reproduction

## Bug

批量删除样本不具备整体语义：提交两条现有样本并夹带 missing，服务竟然报告成功，现有样本也全部消失了；同一 ID 在列表里重复出现时同样没有提示。文件先别改，帮我查清

## Trigger

1. 创建两个真实样本并构造含缺失 ID 的集合
2. 再构造重复 ID 集合并在每次调用后读取原样本

## Observed Error

断言含缺失 ID 的删除返回 ErrNotFound
