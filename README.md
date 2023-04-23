<p align='center'>
   Feishu-OpenAI ×（Embeddings + Qdrant）
<br>
<br>
    🚀 Feishu Vector Knowledge Management 🚀
</p>

## 👻 机器人功能

除了 [Feishu-OpenAI](https://github.com/ConnectAI-E/Feishu-OpenAI) 支持内功能外，
增加了知识库功能，其中具体有：

💬 知识库问答：使用 /faq 或 知识库 作为查询指令

🗣 知识库 CSV 导入：支持导入 CSV 格式向量数据

🎭 知识库 CSV 创建：支持重建 CSV 格式向量文件，降低 token 成本 🚧

📝 知识库管理：支持导入 URL 网页、文件数据 🚧

🔒 知识库管理：支持查询数据库记录并增删改记录 🚧

🍊 缓存问题向量：降低 token 成本，减少重复查询 🚧

## 🌟 项目特点

- 🥒 基于 Embeddings + Qdrant 查询上下文


## 项目部署

### 导入数据

```sh
git clone https://github.com/ConnectAI-E/Feishu-Vector-Knowledge-Management
cd Feishu-Vector-Knowledge-Management
go run ./cmd - prepare csv -f <csvfile>
```

#### CSV 文件表头

```csv
id,url,title,text,title_vector,content_vector,vector_id
```


[样例数据下载](https://cdn.openai.com/API/examples/data/vector_database_wikipedia_articles_embedded.zip)

### Qdrant 测试

在线swagger文档：

https://ui.qdrant.tech/#/

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
