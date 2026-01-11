# API Reference

## Service

Service 是 Service Operator 的核心资源，用于定义和管理一个完整的服务部署。

### ServiceSpec

ServiceSpec 定义了服务的期望状态。

| 字段 | 类型 | 必需 | 默认值 | 描述 |
|------|------|------|--------|------|
| `image` | string | 是 | - | 容器镜像名称和标签 |
| `replicas` | *int32 | 否 | 1 | Pod 副本数量 |
| `port` | int32 | 否 | 8080 | 容器和服务监听的端口 |
| `serviceType` | string | 否 | "ClusterIP" | Kubernetes Service 类型 (ClusterIP, NodePort, LoadBalancer) |
| `configData` | map[string]string | 否 | - | 配置数据，将创建为 ConfigMap |
| `env` | []EnvVar | 否 | - | 环境变量列表 |
| `resources` | *ResourceRequirements | 否 | - | 资源请求和限制 |
| `ingress` | *IngressSpec | 否 | - | Ingress 配置 |

### EnvVar

环境变量定义。

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `name` | string | 是 | 环境变量名称 |
| `value` | string | 是 | 环境变量值 |

### ResourceRequirements

资源要求定义。

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `requests` | map[string]string | 否 | 资源请求 (cpu, memory) |
| `limits` | map[string]string | 否 | 资源限制 (cpu, memory) |

#### 资源格式示例

```yaml
resources:
  requests:
    cpu: "100m"      # 100 millicores
    memory: "128Mi"   # 128 MiB
  limits:
    cpu: "500m"      # 500 millicores
    memory: "512Mi"   # 512 MiB
```

### IngressSpec

Ingress 配置定义。

| 字段 | 类型 | 必需 | 默认值 | 描述 |
|------|------|------|--------|------|
| `enabled` | bool | 是 | - | 是否启用 Ingress |
| `host` | string | 否 | - | 主机名 |
| `path` | string | 否 | "/" | 路径匹配规则 |
| `annotations` | map[string]string | 否 | - | Ingress 注解 |
| `tls` | *TLSSpec | 否 | - | TLS 配置 |

### TLSSpec

TLS 配置定义。

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `enabled` | bool | 是 | 是否启用 TLS |
| `secretName` | string | 否 | TLS 证书 Secret 名称 |

### ServiceStatus

ServiceStatus 定义了服务的观察状态。

| 字段 | 类型 | 描述 |
|------|------|------|
| `phase` | string | 当前阶段 (Ready, Pending) |
| `readyReplicas` | int32 | 就绪的副本数量 |
| `url` | string | 外部访问 URL (如果启用了 Ingress) |
| `conditions` | []metav1.Condition | 状态条件列表 |

## 完整示例

### 基础 Web 服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: web-app
  namespace: default
spec:
  image: nginx:1.21
  replicas: 3
  port: 80
  serviceType: ClusterIP
```

### 带配置的应用服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: api-service
  namespace: production
spec:
  image: myapp/api:v1.0.0
  replicas: 5
  port: 8080
  serviceType: ClusterIP
  configData:
    application.yml: |
      server:
        port: 8080
      database:
        url: jdbc:postgresql://db:5432/myapp
    logback.xml: |
      <configuration>
        <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
          <encoder>
            <pattern>%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
          </encoder>
        </appender>
        <root level="INFO">
          <appender-ref ref="STDOUT" />
        </root>
      </configuration>
  env:
    - name: SPRING_PROFILES_ACTIVE
      value: "production"
    - name: JVM_OPTS
      value: "-Xmx2g -Xms1g"
  resources:
    requests:
      cpu: "500m"
      memory: "1Gi"
    limits:
      cpu: "2000m"
      memory: "4Gi"
```

### 带 Ingress 的服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: web-service
  namespace: default
spec:
  image: myapp/web:latest
  replicas: 2
  port: 3000
  serviceType: ClusterIP
  env:
    - name: NODE_ENV
      value: "production"
  resources:
    requests:
      cpu: "200m"
      memory: "256Mi"
    limits:
      cpu: "1000m"
      memory: "1Gi"
  ingress:
    enabled: true
    host: myapp.example.com
    path: /
    annotations:
      kubernetes.io/ingress.class: nginx
      nginx.ingress.kubernetes.io/rewrite-target: /
      cert-manager.io/cluster-issuer: letsencrypt-prod
    tls:
      enabled: true
      secretName: myapp-tls
```

## 状态查询

### 查看服务列表

```bash
kubectl get services.apps.example.com
```

输出示例：
```
NAME          IMAGE              REPLICAS   READY   PHASE   AGE
web-service   myapp/web:latest   2          2       Ready   5m
api-service   myapp/api:v1.0.0   5          5       Ready   10m
```

### 查看详细状态

```bash
kubectl describe services.apps.example.com web-service
```

### 查看服务 URL

如果启用了 Ingress，可以从状态中获取 URL：

```bash
kubectl get services.apps.example.com web-service -o jsonpath='{.status.url}'
```

## 常用操作

### 扩缩容

```bash
kubectl patch services.apps.example.com web-service -p '{"spec":{"replicas":5}}'
```

### 更新镜像

```bash
kubectl patch services.apps.example.com web-service -p '{"spec":{"image":"myapp/web:v2.0.0"}}'
```

### 更新配置

```bash
kubectl patch services.apps.example.com web-service -p '{"spec":{"configData":{"app.conf":"new config content"}}}'
```

### 启用/禁用 Ingress

```bash
# 启用 Ingress
kubectl patch services.apps.example.com web-service -p '{"spec":{"ingress":{"enabled":true,"host":"myapp.example.com"}}}'

# 禁用 Ingress
kubectl patch services.apps.example.com web-service -p '{"spec":{"ingress":{"enabled":false}}}'
```