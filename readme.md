## AlienDB 
### 软工2022年春大作业
- 指导老师：鲍凌峰
- 小组成员：cyq，zyc，qxy，yxc，wcl
- based on etcd and sqlite

## 项目编译运行
1. 安装 etcd 并运行
2. 编译 master目录下的main文件, 运行
    ```cmd
    go build main.go
    ./main.exe
    ```
   监听client的端口号为 2224
   监听region的端口号为 2223
3. 分别编译 region 和 client 目录下的main文件，运行
   
   - 输入目标地址和端口号，即可通过client端进行试验
   - 需要至少运行3个region
   - 最好在同一个WiFi环境下进行操作