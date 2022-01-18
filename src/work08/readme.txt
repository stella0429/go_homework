【week08作业】

1.使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。
2.写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。


解答：

1)题一：redis benchmark工具测试结果，查看同目录下《redis-benchmark性能测试结果.xlsx》。
a.命令：redis-benchmark  -t get,set  -n 100000
b.配置说明：
系统：mac
处理器：2.9 GHz Intel Core i5
内存：16 GB 1867 MHz DDR3

2)题二：代码见同级目录/project/cmd/internalstoragetest/下的main.go文件。
a.实现方法：设定key长度相等情况，分别设置value大小为1w,2w,5w,10w,20w,30w,50w数据大小，计算set前后内存差异（即占用内存大小），每种分别跑10组，分别计算每个key的占用内存空间。
b.分析结果：详细分析结果打印日志，查看同目录下《calculate.log》
c.分析总结：
    1w: 平均每个 key 的占用内存空间 10393
    2w：平均每个 key 的占用内存空间 28662
    5w: 平均每个 key 的占用内存空间 67267
    10w: 平均每个 key 的占用内存空间 167619
    20w: 平均每个 key 的占用内存空间 368272
    30w：平均每个 key 的占用内存空间 872080
    50w：平均每个 key 的占用内存空间 1572496 
d.大致总结：key大小相同的情况下，value值越大，占用内存也越大。
