使用 redis benchmark 工具| 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

写入一定量的 kv 数据| 根据数据大小 1w-50w 自己评估| 结合写入前后的 info memory 信息 | 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

测试环境 windows10 专业工作站版 i7 8550u 16Gddr4-2400
Redis 5.0.14 for windows 默认配置
命令及结果

.\redis-benchmark.exe -q -d 10 -r 1000000 -n 500000 -t get,set -p 6379 SET: 55456.96 requests per second GET: 60790.27 requests per second

.\redis-benchmark.exe -q -d 20 -r 1000000 -n 500000 -t get,set -p 6379 SET: 48063.06 requests per second GET: 49251.38 requests per second

.\redis-benchmark.exe -q -d 50 -r 1000000 -n 500000 -t get,set -p 6379 SET: 49726.50 requests per second GET: 48742.44 requests per second

.\redis-benchmark.exe -q -d 100 -r 1000000 -n 500000 -t get,set -p 6379 SET: 50266.41 requests per second GET: 47851.47 requests per second

.\redis-benchmark.exe -q -d 200 -r 1000000 -n 500000 -t get,set -p 6379 SET: 43580.58 requests per second GET: 48524.84 requests per second

.\redis-benchmark.exe -q -d 1024 -r 1000000 -n 500000 -t get,set -p 6379 SET: 50005.00 requests per second GET: 58370.30 requests per second

.\redis-benchmark.exe -q -d 5120 -r 1000000 -n 500000 -t get,set -p 6379 SET: 56554.69 requests per second GET: 46330.62 requests per second

10字节，随机写入39w，占用内存34M,平均94字节

20字节，随机写入39w，占用内存37M,平均99字节

50字节，随机写入39w，占用内存52M,平均140字节

100字节，随机写入39w，占用内存70M,平均188字节

200字节，随机写入39w，占用内存115M,平均310字节

1k字节，随机写入39w，占用内存509M,平均1368字节

5k字节，随机写入39w，占用内存2.28G,平均6277字节
