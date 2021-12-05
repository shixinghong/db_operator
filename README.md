# Operator

手动实现一个Operator

## CRD的定义
[kubernetes官方链接](https://kubernetes.io/zh/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)

- 额外打印列（执行kubectl get < resource > 命令时候显示的列）
  kubectl 工具依赖服务器端的输出格式化。你的集群的 API 服务器决定 kubectl get 命令要显示的列有哪些。 
  为 CustomResourceDefinition 定制这些要打印的列。 下面的例子添加了 Spec、Replicas 和 Age 列：

- 优先级 
  每个列都包含一个 priority（优先级）字段。当前，优先级用来区分标准视图（Standard View）和宽视图（Wide View）（使用 -o wide 标志）中显示的列：
  - 优先级为 0 的列会在标准视图中显示。
  - 优先级大于 0 的列只会在宽视图中显示。

