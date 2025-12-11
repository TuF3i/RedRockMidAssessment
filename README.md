# 学生选课系统后端项目-规划

> **红岩网校后端中期考核项目**

## 1. 业务模块划分

- 学生管理模块 - 学生
- 课程管理模块 - 学生
- 学生管理模块 - 管理员
  
- 课程管理模块 - 管理员

## 2. 依赖

| 软件包                | 安装方式                                     |       备注        |
| :-------------------- | -------------------------------------------- | :---------------: |
| **go-hertz**          | `go get github.com/cloudwego/hertz`          |    **API引擎**    |
| **go-gorm**           | `go get gorm.io/gorm`                        |  **数据库操作**   |
| **gorm-driver-mysql** | `gorm.io/driver/mysql`                       |   **mysql驱动**   |
| **go-redis**          | `go get github.com/go-redis/redis/v8`        |   **redis驱动**   |
| **kafka**             | `go get github.com/IBM/sarama`               |   **消息队列**    |
| **go-viper**          | `go get github.com/spf13/viper`              | **配置文件管理**  |
| **go-zap**            | `go get -u go.uber.org/zap`                  |   **日志支持**    |
| **hertz-zap**         | `go get github.com/hertz-contrib/logger/zap` | **Hertz框架日志** |
| **go-lumberjack**     | `github.com/natefinch/lumberjack`            |   **日志轮转**    |
| **go-color**          | `github.com/fatih/color`                     |  **好看的输出**   |
| **go-jwt**            | `go get -u github.com/golang-jwt/jwt/v5`     |  **jwt认证支持**  |
| **go-i18n**           | `github.com/nicksnyder/go-i18n/v2/i18n`      |  **国际化支持**   |

## 3. 接口设计

### 3.0 公共接口 (Auth-free)

#### 3.0.1 注册

```
POST /v1/api/public/register
```

- 请求方法：`POST`

- 请求体JSON:

  ```json
  {
      "name":"张三",
      "stu_id":"1699888",
      "stu_class":"07021580",
      "password":"Cy******",
      "sex": 1,
      "grade": "大一",
      "age": 18
  }
  ```

  - 字段含义: 

    | 字段      | 数据类型 | 备注                                       |
    | --------- | -------- | ------------------------------------------ |
    | name      | String   | 姓名                                       |
    | stu_class | String   | 班级                                       |
    | stu_id    | String   | 学号                                       |
    | password  | String   | 密码                                       |
    | sex       | uint     | 性别 (1表示男，2表示女，0表示沃尔玛购物袋) |
    | grade     | String   | 年级                                       |
    | age       | uint     | 年龄                                       |

- 响应体data字段为`null`

#### 3.0.2 登录

```
POST /v1/api/public/login
```

- 请求方法：`POST`

- 请求体JSON：

  ```json
  {
      "stu_id":"1699888",
      "password":"Cy******"
  }
  ```

  - 字段含义：

    | 字段     | 数据类型 | 备注 |
    | -------- | -------- | ---- |
    | stu_id   | String   | 学号 |
    | password | String   | 密码 |

- 响应data字段结构：

  ```json
  {
      "access_token": "****",
      "refresh_token": "****"
  }
  ```

  - 字段含义：

    | 字段          | 数据类型 | 备注                                     |
    | ------------- | -------- | ---------------------------------------- |
    | access_token  | String   | 认证Token,用于操作鉴权                   |
    | refresh_token | String   | 刷新Token,用于access_token过期后进行刷新 |

#### 3.0.3 刷新AccessToken

```
GET /v1/api/public/refresh
```

- 请求方法: `GET`

- 请求体data为`null`

- 响应data字段结构：

  ```json
  {
      "access_token": "****",
      "refresh_token": "****"
  }
  ```

  - 字段含义：

    | 字段          | 数据类型 | 备注                                     |
    | ------------- | -------- | ---------------------------------------- |
    | access_token  | String   | 认证Token,用于操作鉴权                   |
    | refresh_token | String   | 刷新Token,用于access_token过期后进行刷新 |


### 3.1 学生管理模块 - 学生 (Auth-required)

#### 3.1.1 学生信息修改 

> [!tip]
>
> 由于学生修改自己的信息这个功能需要先获取学生的信息，然后渲染到页面上供学生修改，最后提交给后端，所以将 **学生信息接口**`/v1/api/stu-manager/stu-info` 与 **修改学生信息接口**`/v1/api/stu-manager/stu-update` 统一写到**学生信息修改**这个子业务模块中

> [!Note]
>
> 学生信息的修改分为三步组成：
>
> 第一步：从`/v1/api/stu-manager/stu-info`获取学生的所有信息并由前端渲染到页面上
>
> 第二步：前端检测差异并生成修改过的字段和修改的值组成的列表
>
> 第三部：用`/v1/api/stu-manager/stu-update`进行更新

- 3.1.1 - 1 获取学生所有信息

  ```
  GET /v1/api/stu-manager/stu-info
  ```

  - 请求方法：`GET`

  - 请求体为`null`

  - 响应体data字段结构:

    ```json
    {
        "name":"张三",
        "stu_id":"1699888",
        "password":"Cy******",
        "sex": 1,
        "grade": "大一",
        "age": 18
    }
    ```

    - 字段含义：

      | 字段     | 数据类型 | 备注                                       |
      | -------- | -------- | ------------------------------------------ |
      | name     | String   | 姓名                                       |
      | stu_id   | String   | 学号                                       |
      | password | String   | 密码                                       |
      | sex      | uint     | 性别 (1表示男，2表示女，0表示沃尔玛购物袋) |
      | grade    | String   | 年级                                       |
      | age      | uint     | 年龄                                       |

- 3.1.1 -2 更新学生指定字段的信息

  ```
  PATCH /v1/api/stu-manager/stu-update
  ```

  - 请求方法：`PATCH`

  - 请求体JSON:

    ```json
    {
      "update_columns": [
        {"field": "key_1", "value": "value_1"},
        {"field": "key_2", "value": "value_2"}
      ]
    }
    ```

    - 字段含义：

      | 字段           | 数据类型 | 备注           |
      | -------------- | -------- | -------------- |
      | update_columns | list     | 更新的字段列表 |
      | field          | String   | 字段名称       |
      | value          | String   | 值             |

  - 响应体data字段为`null`

#### 3.1.2 学生注销 

> [!Note]
>
> 登出后将两个token全部写入redis的黑名单

```
GET /v1/api/stu-manager/stu-logout
```

- 请求方法：`GET`
- 请求体为`null`
- 响应体data字段为`null`

### 3.2 课程管理模块 - 学生 (Auth-required)

#### 3.2.1 查看可选课程

```
GET /v1/api/class-manager/get-selectable-classes
```

- 请求方法: `GET`

- 请求体为`null`

- 响应体data字段结构：

  ```json
  {
      selectable_classes: [
          {
              "class_name":"高等数学",
              "class_id":"xxxxxxx",
              "class_location":"2106",
              "class_time":"1-6,7",
              "class_teacher":"God",
              "class_capcity":80,
              "class_selsetion":68
          },
          {
              "class_name":"线性代数",
              "class_id":"xxxxxxx",
              "class_location":"2301",
              "class_time":"4-4,5",
              "class_teacher":"Godess",
              "class_capcity":60,
              "class_selsetion":60
          }
      ]
  }
  ```

  - 字段含义:

    | 字段               | 数据类型 | 备注                             |
    | ------------------ | -------- | -------------------------------- |
    | selectable_classes | list     | 可选的课的列表                   |
    | class_name         | String   | 课程名称                         |
    | class_id           | String   | 课程ID                           |
    | class_location     | String   | 上课地点                         |
    | class_time         | String   | 上课时间(星期-第几节课,第几节课) |
    | class_teacher      | String   | 上课老师                         |
    | class_capcity      | uint     | **选课容量**                     |
    | class_selsetion    | uint     | **选课人数**                     |

#### 3.2.2 添加选课

```
POST /v1/api/class-manager/subscribe-class
```

- 请求方法：`POST`

- 请求体JSON:

  ```json
  {
      "class_id":"XXXXXXX"
  }
  ```

  - 字段含义：

    | 字段     | 数据类型 | 备注   |
    | -------- | -------- | ------ |
    | class_id | String   | 课程ID |

- 响应体data字段为`null`

#### 3.2.3 删除选课

```
DEL /v1/api/class-manager/del-class
```

- 请求方法：`DELETE`

- 请求体JSON:

  ```json
  {
      "class_id":"XXXXXXX"
  }
  ```

  - 字段含义：

    | 字段     | 数据类型 | 备注   |
    | -------- | -------- | ------ |
    | class_id | String   | 课程ID |

- 响应体data字段为`null`

#### 3.2.4 查看已选课程

```
GET /v1/api/class-manager/get-subscribed-classes
```

- 请求方法：`GET`

- 请求体为`null`

- 响应体data字段JSON:

  ```json
  {
      selected_classes: [
          {
              "class_name":"高等数学",
              "class_id":"xxxxxxx",
              "class_location":"2106",
              "class_time":"1-6,7",
              "class_teacher":"God",
          },
          {
              "class_name":"线性代数",
              "class_id":"xxxxxxx",
              "class_location":"2301",
              "class_time":"4-4,5",
              "class_teacher":"Godess",
          }
      ]
  }
  ```

  - 字段含义：

    | 字段               | 数据类型 | 备注                             |
    | ------------------ | -------- | -------------------------------- |
    | selectable_classes | list     | 可选的课的列表                   |
    | class_name         | String   | 课程名称                         |
    | class_id           | String   | 课程ID                           |
    | class_location     | String   | 上课地点                         |
    | class_time         | String   | 上课时间(星期-第几节课,第几节课) |
    | class_teacher      | String   | 上课老师                         |

### 3.3 学生管理模块 - Admin (Auth-required)

#### 3.3.1 查看学生列表

```
GET /v1/api/admin/stu-manager/get-stu-list?page=1&resNum=15
```

- 请求方法：`GET`

- 请求体为`null`

- 响应体data字段JSON:

  ```json
  {	
      total: 30,
      page: 1,
      page_size: 10,
      students_list: [
          {"stu_id":"1899778", "stu_name":"xxx", "stu_class":"07021580", "grade":"大一"},
          {"stu_id":"1799778", "stu_name":"xxx", "stu_class":"06021580", "grade":"大二"}
      ]
  }
  ```
  
  - 字段含义：
  
    | 字段          | 数据类型 | 备注     |
    | ------------- | -------- | -------- |
    | students_list | list     | 学生列表 |
    | stu_id        | String   | 学生ID   |
    | stu_name      | String   | 姓名     |
    | stu_class     | String   | 班级     |
    | grade         | String   | 年级     |

#### 3.3.2 修改学生信息

> [!Tip]
>
> 由于学生修改自己的信息这个功能需要先获取学生的信息，然后渲染到页面上供学生修改，最后提交给后端，所以将 **学生信息接口**`/v1/api/stu-manager/stu-info` 与 **修改学生信息接口**`/v1/api/stu-manager/stu-update` 统一写到**学生信息修改**这个子业务模块中

> [!note]
>
> 学生信息的修改分为三步组成：
>
> 第一步：从`/v1/api/admin/stu-manager/get-stu-list`获取学生列表并渲染到页面上
>
> 第二步：选中的学生从`/v1/api/admin/stu-manager/get-stu-info`获取学生的所有信息并由前端渲染到页面上
>
> 第三步：前端检测差异并生成修改过的字段和修改的值组成的列表
>
> 第四部：用`/v1/api/admin/stu-manager/update-stu-info`进行更新

- 3.3.2 - 2 查看学生信息

  ```
  GET /v1/api/admin/stu-manager/get-stu-info/:stuID
  ```

  - 请求方法：`GET`

  - 请求体为`null`

  - 响应体data字段JSON:

    ```json
    {
        "name":"张三",
        "stu_id":"1699888",
        "stu_class":"07021580",
        "password":"Cy******",
        "sex": 1,
        "grade": "大一",
        "age": 18
    }
    ```

  - 字段含义:

    | 字段      | 数据类型 | 备注                                       |
    | --------- | -------- | ------------------------------------------ |
    | name      | String   | 姓名                                       |
    | stu_class | String   | 班级                                       |
    | stu_id    | String   | 学号                                       |
    | password  | String   | 密码                                       |
    | sex       | uint     | 性别 (1表示男，2表示女，0表示沃尔玛购物袋) |
    | grade     | String   | 年级                                       |
    | age       | uint     | 年龄                                       |

- 3.3.2 - 3 更新学生指定字段的信息

  ```
  PATCH /v1/api/admin/stu-manager/update-stu-info
  ```

  - 请求方法：`PATCH`

  - 请求体JSON:

    ```json
    {
        stu_id: "xxxxx",
        update_columns: [
            {"key_1":"value_1"},
            {"key_2":"value_2"},
        ]
    }
    ```

    - 字段含义：

      | 字段           | 数据类型 | 备注           |
      | -------------- | -------- | -------------- |
      | update_columns | list     | 更新的字段列表 |
      | stu_id         | String   | 学生ID         |
      | key            | String   | 字段名称       |
      | value          | String   | 值             |

  - 响应体data字段为`null`

#### 3.3.3 学生创建

```
POST /v1/api/admin/stu-manager/create-stu
```

- 请求方法：`POST`

- 请求体JSON:

  ```json
  {
      "name":"张三",
      "stu_id":"1699888",
      "stu_class":"07021580",
      "password":"Cy******",
      "sex": 1,
      "grade": "大一",
      "age": 18
  }
  ```

  - 字段含义：

    | 字段      | 数据类型 | 备注                                       |
    | --------- | -------- | ------------------------------------------ |
    | name      | String   | 姓名                                       |
    | stu_id    | String   | 学生ID                                     |
    | stu_class | String   | 班级                                       |
    | password  | String   | 密码                                       |
    | sex       | uint     | 性别 (1表示男，2表示女，0表示沃尔玛购物袋) |
    | grade     | String   | 年级                                       |
    | age       | uint     | 年龄                                       |

- 响应体data字段为`null`

#### 3.3.4 学生删除

```
DEL /v1/api/admin/stu-manager/del-stu
```

- 请求方法：`DELETE`

- 请求体JSON:

  ```json
  {
      "stu_id": "xxxxxx"
  }
  ```

  - 字段含义：

    | 字段   | 数据类型 | 备注   |
    | ------ | -------- | ------ |
    | stu_id | String   | 学生ID |

- 响应体data字段为`null`

### 3.4 课程管理模块 - Admin (Auth-required)

#### 3.4.1 查看选课情况

```
GET /v1/api/admin/classes-manager/get-class-status
```

- 请求方法: `GET`

- 请求体为`null`

- 响应体data字段结构：

  ```json
  {
      selectable_classes: [
          {
              "class_name":"高等数学",
              "class_id":"xxxxxxx",
              "class_location":"2106",
              "class_time":"1-6,7",
              "class_teacher":"God",
              "class_capcity":80,
              "class_selsetion":68
          },
          {
              "class_name":"线性代数",
              "class_id":"xxxxxxx",
              "class_location":"2301",
              "class_time":"4-4,5",
              "class_teacher":"Godess",
              "class_capcity":60,
              "class_selsetion":60
          }
      ]
  }
  ```

  - 字段含义:

    | 字段               | 数据类型 | 备注                             |
    | ------------------ | -------- | -------------------------------- |
    | selectable_classes | list     | 可选的课的列表                   |
    | class_name         | String   | 课程名称                         |
    | class_id           | String   | 课程ID                           |
    | class_location     | String   | 上课地点                         |
    | class_time         | String   | 上课时间(星期-第几节课,第几节课) |
    | class_teacher      | String   | 上课老师                         |
    | class_capcity      | uint     | **选课容量**                     |
    | class_selsetion    | uint     | **选课人数**                     |

#### 3.4.2 修改学生的选课 

> [!Note]
>
> 修改学生选课分三步：
>
> 1. 通过`/v1/api/admin/stu-manager/get-stu-list`查询学生列表 🆗
> 2. 通过`/v1/api/admin/classes-manager/get-stu-classes`查询单个学生选课情况 🆗
> 3. 通过`/v1/api/admin/classes-manager/update-stu-classes`修改选课 🆗

- 3.4.3 - 1 获取学生列表
  - **此次操作方法同 `3.3.1` **

- 3.4.3 - 2 查看学生选课

  ```
  GET /v1/api/admin/classes-manager/get-stu-classes/:stuID
  ```

  - 请求方法：`GET`

  - 请求体为`null`

  - 响应体data字段JSON:

    ```json
    {
        selected_classes: [
            {
                "class_name":"高等数学",
                "class_id":"xxxxxxx",
                "class_location":"2106",
                "class_time":"1-6,7",
                "class_teacher":"God",
                "class_capcity":80,
                "class_selsetion":68
            },
            {
                "class_name":"线性代数",
                "class_id":"xxxxxxx",
                "class_location":"2301",
                "class_time":"4-4,5",
                "class_teacher":"Godess",
                "class_capcity":60,
                "class_selsetion":60
            }
        ]
    }
    ```

    - 字段含义：

      | 字段               | 数据类型 | 备注                             |
      | ------------------ | -------- | -------------------------------- |
      | selectable_classes | list     | 可选的课的列表                   |
      | class_name         | String   | 课程名称                         |
      | class_id           | String   | 课程ID                           |
      | class_location     | String   | 上课地点                         |
      | class_time         | String   | 上课时间(星期-第几节课,第几节课) |
      | class_teacher      | String   | 上课老师                         |
      | class_capcity      | uint     | **选课容量**                     |
      | class_selsetion    | uint     | **选课人数**                     |

- 3.4.3 - 3 添加学生选课

  ```
  PATCH /v1/api/admin/classes-manager/update-stu-classes
  ```

  - 请求方法：`PATCH`

  - 请求体JSON:

    ```json
    {
        "stu_id": "xxxxx",
        "update_class_id":"xxxxx",
    }
    ```
    
    - 字段含义：
    
      | 字段            | 数据类型 | 备注   |
      | --------------- | -------- | ------ |
      | stu_id          | String   | 学生ID |
      | update_class_id | string   | 课程ID |
    
  - 响应体data字段为`null`
  
- 3.4.3 - 4 删除学生选课

  ```
  DELETE /v1/api/admin/classes-manager/update-stu-classes
  ```

  - 请求方法：`DELETE`

  - 请求体JSON:

    ```json
    {
        "stu_id": "xxxxx",
        "update_class_id":"xxxxx",
    }
    ```
  
    - 字段含义：
  
      | 字段            | 数据类型 | 备注   |
      | --------------- | -------- | ------ |
      | stu_id          | String   | 学生ID |
      | update_class_id | string   | 课程ID |
  
  - 响应体data字段为`null`

#### 3.4.3 修改课程信息

> [!note]
>
> 修改课程信息分三步：
>
> 1. 通过`/v1/api/admin/classes-manager/get-class-status`查看可选课程
> 2. 更新课程信息

> [!Warning]
>
> 1. ⚠`/v1/api/classes-manager/edit-class-info`仅可用修改课程信息，**不可以修改课程容量**
> 2. ⚠`/v1/api/classes-manager/edit-class-stock`才**可以修改课程容量**

- 3.4.3 - 1 修改课程信息

  ```
  PATCH /v1/api/classes-manager/edit-class-info
  ```

  - 请求方法：`PATCH`


  - 请求体JSON:

    ```json
    {
      "class_id": "xxxxx",
      "update_columns": [
        {"field": "key_1", "value": "value_1"},
        {"field": "key_2", "value": "value_2"}
      ]
    }
    ```

    - 字段含义：

      | 字段           | 数据类型 | 备注           |
      | -------------- | -------- | -------------- |
      | update_columns | list     | 更新的字段列表 |
      | field          | String   | 字段名称       |
      | value          | String   | 值             |


  - 响应体data字段为`null`

- 3.4.3 - 2 修改课程容量

  ```
  PATCH /v1/api/classes-manager/edit-class-stock
  ```

  - 请求方法：`PATCH`

  - 请求体JSON: 

    ```json
    {
        "class_id": "xxxxx",
        "stock": 10,
    }
    ```

    - 字段含义：

      | 字段     | 数据类型 | 备注     |
      | -------- | -------- | -------- |
      | class_id | String   | 课程ID   |
      | stock    | uint     | 课程容量 |

  - 响应体data字段为`null`

#### 3.4.4 添加课程

```
POST /v1/api/admin/classes-manager/add-course
```

- 请求方法 `POST`

- 请求体JSON: 

  ```json
  {
    "class_name": "高等数学",
    "class_id": "xxxxxxx",
    "class_location": "2106",
    "class_time": "1-6,7",
    "class_teacher": "God",
    "class_capcity": 80
  }
  ```

  - 字段含义

    | 字段           | 数据类型 | 备注     |
    | -------------- | -------- | -------- |
    | class_name     | String   | 课程名称 |
    | class_id       | String   | 课程ID   |
    | class_location | String   | 上课地点 |
    | class_time     | String   | 上课时间 |
    | class_teacher  | String   | 课程教师 |
    | class_capcity  | uint     | 课程容量 |

- 响应体data字段为`null`

####  3.4.5 删除课程

```
DELETE /v1/api/admin/classes-manager/delete-course
```

- 请求方法 `DELETE`

- 请求体JSON:

  ```json
  {
      "class_id":"XXXXXXX"
  }
  ```

  - 字段含义：

    | 字段     | 数据类型 | 备注   |
    | -------- | -------- | ------ |
    | class_id | String   | 课程ID |

- 响应体data字段为`null`

#### 3.4.6 开始选课

```
GET /v1/api/admin/classes-manager/start-course-select-event
```

- 请求方法`GET`
- 请求体data为`null`
- 响应体data字段为`null`

#### 3.4.6 结束选课

```
GET /v1/api/admin/classes-manager/stop-course-select-event
```

- 请求方法`GET`
- 请求体data为`null`
- 响应体data字段为`null`

## 4. 错误传递规范

### 4.1 错误传递和日志规范

1. 每次调用产生的错误由：

   ```
   [dao] --> [service] --> [handle]
   ```

2. 每次调用产生的traceID传递方向：

   ```
   [handle] --> [service] --> [dao]
   ```

3. traceID包装在`context.Context()`中逐层向下传递

4. 每次返回的日志都用`response.Response{}`进行包装逐层上抛

5. 在`return response.Response`的地方不进行日志记录，在出现`err != nil`时才进行日志记录

6. 出现**非业务报错**时用`func ServerInternalError(err error) Response生成封装`

7. 错误封装：

   ```go
   type Response struct {
   	Status uint   `json:"status"`
   	Info   string `json:"info"`
   }
   ```

8. 请求层错误封装：

   ```go
   type FinalResponse struct {
   	Status uint        `json:"status"`
   	Info   string      `json:"info"`
   	Data   interface{} `json:"data"`
   }
   ```

## 5. 数据校验规范

###  5.1 用户名校验规范

- **相关带代码块：**

```go
var (
    	//长度控制，最短2个，最长15个
		miniLength = 2 
		maxLength  = 15
    	// 0. 系统保留字（bob是测试账号）
		reserved   = map[string]struct{}{"admin": {}, "root": {}, "user": {}, "api": {}, "bob": {}}  
    	// 1. 合法字符：字母数字中文下划线；禁止首尾下划线 *AI生成的正则表达式*
		userRe     = regexp.MustCompile(`^[a-zA-Z0-9\p{Han}]([a-zA-Z0-9_\p{Han}]*[a-zA-Z0-9\p{Han}])?$`) 
	)
```

#### 5.1.1 长度控制

> [!Note]
>
> 最长不超过`15`个字符
>
> 最短不少于`2`个字符

#### 5.1.2 系统保留字

> [!Note]
>
> 系统中的：
>
> 1. `admin`  管理员
> 2. `root`
> 3. `user`
> 4. `api`
> 5. `bob` 测试账号
>
> 都为不可用的非法用户名

#### 5.1.3 合法字符

> [!note]
>
> 允许：
>
> 1. **字母**
> 2.  **数字** 
> 3. **中文**
> 4.  **下划线**
>
> 禁止：
>
> 1. **首字符**为**下划线**
> 2. **尾字符**为**下划线**

### 5.2 学生ID校验规范

- **相关带代码块：**

```go
func VerifyUserID(stuID string) bool {
	// 10位用户ID
	if len(stuID) == 10 {
		return true
	}
	return false
}
```

#### 5.2.1 长度规范

> [!note]
>
> `10` 位长的 **学生ID**

### 5.3 学生班级校验规范

- **相关代码块：**

```go
func VerifyStudentClass(stuClass string) bool {
	// 班级字符串在3~15个字符间
	length := utf8.RuneCountInString(stuClass)
	if !(length < 3 || length > 10) {
		return true
	}
	return false
}
```

#### 5.3.1 长度规范

> [!note]
>
> 最长不超过`10`个字符
>
> 最短不少于`2`个字符

### 5.4 密码校验规范

- 相关代码块

```go
var (
		miniLength       = 6 //长度控制，最短6个，最长20个
		maxLength        = 20
		miniNumberNum    = 1                                // 最少含有一个数字
		miniLowerCharNum = 1                                // 至少一个小写字母
		miniUpperCharNum = 1                                // 至少一个大写字母
		miniSpecialsNum  = 1                                // 至少一个特殊字符
		specials         = "!@#$%^&*()_+-=[]{}|;':\",./<>?" // 特殊字符集
	)
```

#### 5.4.1 长度规范

> [!note]
>
> 最长不超过`20`个字符
>
> 最短不少于`6`个字符

#### 5.4.2 特殊字符要求

> [!note]
>
> 1. 最少含有一个数字
> 2. 至少一个小写字母
> 3. 至少一个大写字母
> 4. 至少一个特殊字符

> [!Important]
>
> **特殊字符集：**`!@#$%^&*()_+-=[]{}|;':\",./<>?`

### 5.5 性别校验规范

- **相关代码块：**

```go	
func VerifySexSetting(sex uint) bool {
	if sex > 2 {
		return false
	}
	return true
}
```

- **标识码：**

| 标识码 | 性别         |
| ------ | ------------ |
| 1      | 男           |
| 2      | 女           |
| 0      | 沃尔玛购物袋 |

#### 5.5.1 取值规范

> [!note]
>
> **取值在：** 0~2 间的整数

### 5.6 年级校验规范

- **相关代码块：**

```go
func VerifyGrade(grade uint) bool {
	if grade > 0 && grade < 5 {
		return true
	}
	return false
}
```

- **标识码：**

| 标识码 | 年级 |
| ------ | ---- |
| 1      | 大一 |
| 2      | 大二 |
| 3      | 大三 |
| 4      | 大四 |

#### 5.6.1 取值规范

> [!note]
>
> **取值在：** 1~4 间的整数

### 5.7 年龄校验规范

- **相关代码块：**

```go
func VerifyAge(age uint) bool {
	if age > 10 && age < 60 {
		return true
	}
	return false
}
```

#### 5.7.1 年龄取值规范

> [!Note]
>
> 最小不小于 `11` 岁
>
> 最大不大于 `60` 岁

### 5.8 课程ID校验规范

- **相关代码块：**

```go
func VerifyCourseID(classID string) bool {
	// ClassID中不可含中文
	for _, r := range classID {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}
	// classID在5~25位
	if len(classID) >= 5 && len(classID) <= 25 {
		return true
	}
	return false
}
```

#### 5.8.1 内容校验

> [!note]
>
> **不可以含中文**

#### 5.8.2 长度校验

> [!note]
>
> 最短不少于 `5` 个字符
>
> 最长不超过 `25` 个字符

### 5.9课程名称校验规范

- **相关代码块：**

```go	
func VerifyCourseName(courseName string) bool {
	allDigit := true
	for _, r := range []rune(courseName) {
		if !unicode.IsDigit(r) {
			allDigit = false
			break
		}
	}
	if allDigit {
		// 不能为纯数字
		return false
	}

	var ( // courseName在2~15个字符
		maxLength  = 15
		miniLength = 2
	)
	// 检测字符串长度
	if length := utf8.RuneCountInString(courseName); length <= miniLength || length >= maxLength {
		return false
	}

	return true
}
```

#### 5.9.1 内容规范

> [!note]
>
> **不能为纯数字**

#### 5.9.2 长度规范

> [!note]
>
> 最短不少于 `2` 个字符
>
> 最长不超过 `15` 个字符

### 5.10 上课地点校验规范

- **相关代码块：**

```go
func VerifyCourseLocation(courseLocation string) bool {
	var ( // courseName在4~10个字符
		maxLength  = 10
		miniLength = 4
	)
	// 检测字符串长度
	if length := utf8.RuneCountInString(courseLocation); length <= miniLength || length >= maxLength {
		return false
	}

	return true
}
```

#### 5.10.1 长度规范

> [!note]
>
> 最短不少于 `4` 个字符
>
> 最长不超过 `10` 个字符

### 5.11 上课时间校验规范

- **相关代码块：**

```go
func VerifyCourseTime(courseTime string) bool {
	var ( // courseTime在4~10个字符
		maxLength  = 10
		miniLength = 4
	)
	// 检测字符串长度
	if length := utf8.RuneCountInString(courseTime); length <= miniLength || length >= maxLength {
		return false
	}

	return true
}
```

#### 5.11.1 长度规范

> [!note]
>
> 最短不少于 `4` 个字符
>
> 最长不超过 `10` 个字符

### 5.12 授课教师姓名校验规范

- **相关代码块：**

```go	
func VerifyCourseTeacher(courseTeacher string) bool {
	allDigit := true
	for _, r := range []rune(courseTeacher) {
		if !unicode.IsDigit(r) {
			allDigit = false
			break
		}
	}
	if allDigit {
		// 不能为纯数字
		return false
	}

	var ( // courseTeacher在2~8个字符
		maxLength  = 10
		miniLength = 4
	)
	// 检测字符串长度
	if length := utf8.RuneCountInString(courseTeacher); length <= miniLength || length >= maxLength {
		return false
	}

	return true
}
```

#### 5.12.1 内容规范

> [!note]
>
> **不能为纯数字**

#### 5.12.2 长度规范

> [!note]
>
> 最短不少于 `4` 个字符
>
> 最长不超过 `10` 个字符

### 5.13 课程容量校验规范

- **相关代码块：**

```go
func VerifyCourseStock(courseStock uint) bool {
	// courseStock在10~500个字符
	var maxStock uint = 500
	var miniStock uint = 10

	return courseStock <= maxStock && courseStock >= miniStock
}
```

#### 5.13.1 长度规范

>  [!note]
>
> 最短不少于 `10` 个学生
>
> 最长不超过 `500` 个学生

