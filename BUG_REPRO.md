# Bug Reproduction

## Bug

指标更新失败后数据还会变化：已有 cpu 和 memory 两个指标，把 cpu 政名为 memory 会收到冲突错误，可再次读取 cpu 那条记录时，名称、单位和描述都成了失败

## Trigger

1. 建立 cpu 与 memory 两个不同名称指标
2. 用 memory 更新 cpu 并在冲突后按 ID 和旧名称查询

## Observed Error

断言重名更新返回 store.ErrConflict
