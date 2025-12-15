# 🌐学生选课系统🌐

> ## 红岩网校后端中期考核

## 1. 基本信息

> [!Note]
>
> ##### 作者：TuF3i
>
> ##### 名称：学生选课系统
>
> ##### 前端 + 后端

## 2. 运行截图

![feature](/img/feature.jpg)

## 3. 相关文档

| 文档                                                         | 备注           |
| ------------------------------------------------------------ | -------------- |
| [Plan.md](docs\Plan.md)                                      | 项目规划文档   |
| [接口文档_Apifox.md](docs\接口文档_Apifox.md)                | Api接口文档    |
| [RedRockMidAccessment.apifox.json](docs\RedRockMidAccessment.apifox.json) | Apifox接口JSON |

## 4. 部署

### 4.1 克隆仓库

```shell
git clone https://github.com/TuF3i/RedRockMidAssessment.git
```

### 4.2 进入目录

```shell
cd RedRockMidAssessment
```

### 4.3 首次安装执行

```shell
bash start.sh
```

### 4.4 停止服务

```shell
docker-compose down
```

> [!Important]
>
> #### 首次启动后二次启动
>
> ```shell
> docker-compose up -d
> ```

## 5. 不在本地构建镜像

> [!Tip]
>
> 修改docker-compose

```yaml
services:
  # 消费者节点
  consumer:
    image: ghcr.io/tuf3i/consumer-node:latest
    container_name: redrock-consumer
    restart: unless-stopped
    volumes:
      - /opt/data/consumer/config/config.yaml:/root/data/config/config.yaml
      - /opt/data/consumer/logs:/root/data/logs
    environment:
      - TZ=Asia/Shanghai
    depends_on:
      - mysql
      - kafka
      - zookeeper
    networks:
      - redrock-network

  # 路由节点
  gateway:
    image: ghcr.io/tuf3i/gateway-node:latest
    container_name: redrock-gateway
    ports:
      - "8080:8080"
    restart: unless-stopped
    volumes:
      - /opt/data/gateway/config/config.yaml:/root/data/config/config.yaml
      - /opt/data/gateway/logs:/root/data/logs
    environment:
      - TZ=Asia/Shanghai
    depends_on:
      - consumer
      - synchronizer
      - redis
      - mysql
    networks:
      - redrock-network

  # 同步节点
  synchronizer:
    image: ghcr.io/tuf3i/synchronizer-node:latest
    container_name: redrock-synchronizer
    restart: unless-stopped
    volumes:
      - /opt/data/synchronizer/config/config.yaml:/root/data/config/config.yaml
      - /opt/data/synchronizer/logs:/root/data/logs
    environment:
      - TZ=Asia/Shanghai
    depends_on:
      - kafka
      - zookeeper
      - redis
    networks:
      - redrock-network

  # Redis服务
  redis:
    image: redis:7-alpine
    container_name: redrock-redis
    ports:
      - "6379:6379"
    restart: unless-stopped
    command: redis-server --appendonly yes
    volumes:
      - /opt/data/redis-data:/data
    networks:
      - redrock-network

  # MySQL服务
  mysql:
    image: mysql:8.0
    container_name: redrock-mysql
    ports:
      - "3306:3306"
    restart: unless-stopped
    environment:
      MYSQL_ALLOW_EMPTY_PASSWORD: yes
      MYSQL_DATABASE: stuClass
    volumes:
      - /opt/data/mysql-data:/var/lib/mysql
    networks:
      - redrock-network

  # 可选：添加Kafka服务
  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    container_name: redrock-zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - redrock-network

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    container_name: redrock-kafka
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - redrock-network

networks:
  redrock-network:
    driver: bridge
```

> [!Tip]
>
> 手动复制`data`到`/opt`下

```shell
cp -r ./data /opt/data
```

> [!Tip]
>
> 启动容器

```shell
docker-compose up -d
```

## 6. 系统入口

### 选课系统入口：`IP_ADDR:8080`
