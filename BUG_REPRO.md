# Bug Reproduction

## Bug

采集器改名有个奇怪现象：把 node-a 改成已经存在的 node-b 时接口返回冲突，但随后查询原采集器，名称和主机却已经变成了这次失败

## Trigger

1. 先创建名称不同的两个采集器
2. 再用第二个名称更新第一个采集器并在冲突后重新查询

## Observed Error

断言失败更新必须返回 store.ErrConflict
