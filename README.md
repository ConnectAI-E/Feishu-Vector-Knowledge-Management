

<p align='center'>
    <img src='https://user-images.githubusercontent.com/50035229/233825222-0add62d1-e12f-41ef-86d5-4bf3d0f96d84.png' alt='' width='800'/>
</p>

<p align='center'>
   Feishu-OpenAI × 私有知识库
<br>
<br>
    🚀 Feishu Vector Knowledge Management 🚀
</p>

## 👻 机器人功能

除了 [Feishu-OpenAI](https://github.com/ConnectAI-E/Feishu-OpenAI) 支持内功能外，
增加了知识库功能，其中具体有：

💬 知识库问答：使用 /faq 或 知识库 作为查询指令

🗣 知识库 CSV 导入：支持导入 CSV 格式向量数据

🎭 知识库 CSV 创建：支持重建 CSV 格式向量文件，降低 token 成本 

📝 知识库管理：支持导入 URL 网页、文件数据 🚧

🔒 知识库管理：支持查询数据库记录并增删改记录 🚧

🍊 缓存问题向量：降低 token 成本，减少重复查询

## 🌟 项目特点

- 🥒 基于 Embeddings + Qdrant 查询上下文


## 项目部署

### 项目初始化

```sh
git clone https://github.com/ConnectAI-E/Feishu-Vector-Knowledge-Management
cd Feishu-Vector-Knowledge-Management
```


### 导入数据
```sh
# 切割qa数据为csv文件 demo:raw.txt 
go run ./cmd - prepare split -f ./data/demo/raw.txt -o ./data/demo/raw.csv

# 将csv文件转换为向量数据(调用openai-embedding-api), raw.csv 必须包含title和content字段
go run ./cmd - prepare analyze -f ./data/demo/raw.csv -o ./data/demo/vector.csv

# 导入数据csv(向量)数据
go run ./cmd - prepare import -f ./data/demo/vector.csv
```

#### CSV 文件表头
```csv
id,url,title,content,title_vector,content_vector,vector_id
```
[样例数据下载](./data/demo/data.csv)

#### Qdrant 接口调试测试

在线swagger文档：https://ui.qdrant.tech/#/

#### 部署

<details>
    <summary>docker-compose 部署</summary>
<br>

编辑 docker-compose.yaml，通过 environment 配置相应环境变量（或者通过 volumes 挂载相应配置文件），然后运行下面的命令即可

```bash
# 构建镜像
docker compose build

# 启动服务
docker compose up -d

# 停止服务
docker compose down
```

事件回调地址: http://IP:9000/webhook/event
卡片回调地址: http://IP:9000/webhook/card


</details>

## 更多交流

更多结节请访问项目 [Feishu-OpenAI](https://github.com/ConnectAI-E/Feishu-OpenAI)

## 赞助感谢

友情感谢 'Find My Ai' 提供的部分经费赞助！

