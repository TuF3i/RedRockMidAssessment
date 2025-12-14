---
title: 学生选课系统
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 学生选课系统

Base URLs:

* <a href="http://localhost:8080">测试环境: http://localhost:8080</a>

* <a href="http://dev-cn.your-api-server.com">开发环境: http://dev-cn.your-api-server.com</a>

# Authentication

# 公共接口

## GET 刷新AccessToken

GET /v1/api/public/refresh

> ***AccessToken*** 的刷新接口，请将***RefreshToken*** 写入`Authorization: Bearer <token>` 中

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMTdjYjgzN2ItOWZiMy00MjlmLTliYTAtMGRlZWYzYjYzZThlIiwidWlkIjoiMTIzNDU2Nzg5MCIsInJvbGUiOiJzdHVkZW50IiwidHlwZSI6ImFjY2VzcyIsImlzcyI6IlN0dUNsYXNzTVMuQXBpTm9kZSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJhdWQiOlsic3R1ZGVudCIsImFkbWluIl0sImV4cCI6MTc2NTQ1ODQwNywibmJmIjoxNzY1NDU3NTAyLCJpYXQiOjE3NjU0NTc1MDd9.FPitrgXBEwnLe4Z8poPXBSRrXZngI4EwB2uR1vDSyCQ",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiNGZjYzg1OTItZjM4ZS00YzYzLTgzOTMtNzNlY2Q0NDA0MTljIiwidWlkIjoiMTIzNDU2Nzg5MCIsInJvbGUiOiJzdHVkZW50IiwidHlwZSI6InJlZnJlc2giLCJpc3MiOiJTdHVDbGFzc01TLkFwaU5vZGUiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbInN0dWRlbnQiLCJhZG1pbiJdLCJleHAiOjE3NjYwNjIzMDcsIm5iZiI6MTc2NTQ1NzUwMiwiaWF0IjoxNzY1NDU3NTA3fQ.-77dT8zahiHhX-SeTIQ4-BNaWA2zLWQ3BGGcOdo8b64"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» info|string|true|none||none|
|» data|object|true|none||none|
|»» access_token|string|true|none||none|
|»» refresh_token|string|true|none||none|

## POST 登录

POST /v1/api/public/login

> ***Admin*** 和 ***Student*** 的统一登录接口，接受`stu_id`和`password`两个参数，返回`access_token`和`refresh_token`

> Body 请求参数

```json
{
  "stu_id": "1234567890",
  "password": "P@ssw0rd=Ping12345"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|body|body|object| 是 ||none|
|» stu_id|body|string| 是 | 学生ID|10位学生ID|
|» password|body|string| 是 | 密码|登录密码|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiYzQ3OWE0ZTItMWM4OS00MDg3LWFkZGEtNGQwMjQ3NjUzOGY3IiwidWlkIjoiMTIzNDU2Nzg5MCIsInJvbGUiOiJzdHVkZW50IiwidHlwZSI6ImFjY2VzcyIsImlzcyI6IlN0dUNsYXNzTVMuQXBpTm9kZSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJhdWQiOlsic3R1ZGVudCIsImFkbWluIl0sImV4cCI6MTc2NTM4MzI2NCwibmJmIjoxNzY1MzgyMzU5LCJpYXQiOjE3NjUzODIzNjR9.oS98u8MTJ-KP5EwYnSvfvrSaJ39Cbe1GPFduvlb7xzU",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMjNiNDUxNTItYzA4Ni00MzE2LWJmYWYtYmMxYjVhNWJjZTI2IiwidWlkIjoiMTIzNDU2Nzg5MCIsInJvbGUiOiJzdHVkZW50IiwidHlwZSI6InJlZnJlc2giLCJpc3MiOiJTdHVDbGFzc01TLkFwaU5vZGUiLCJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjpbInN0dWRlbnQiLCJhZG1pbiJdLCJleHAiOjE3NjU5ODcxNjQsIm5iZiI6MTc2NTM4MjM1OSwiaWF0IjoxNzY1MzgyMzY0fQ.WfZj48oNlt7aR4wIQLX-61JLiCI0JjnMDJDSnD1X8eI"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|标识响应状态的状态码|
|» info|string|true|none|响应信息|响应的信息|
|» data|object|true|none|数据段|含有响应数据|
|»» access_token|string|true|none|验证Token|用于身份验证|
|»» refresh_token|string|true|none|刷新Token|用于刷新AccessToken|

## POST 注册

POST /v1/api/public/register

> 学生注册接口

> Body 请求参数

```json
{
  "name": "zjq",
  "stu_id": "9876543210",
  "stu_class": "计科8班",
  "password": "P@ssw0rd=Ping12345",
  "sex": 1,
  "grade": 1,
  "age": 19
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|body|body|object| 是 ||none|
|» name|body|string| 是 | 姓名|学生姓名|
|» stu_id|body|string| 是 | 学生ID|学生的ID|
|» stu_class|body|string| 是 | 班级|学生的班级|
|» password|body|string| 是 | 密码|登录密码|
|» sex|body|integer| 是 | 性别|性别|
|» grade|body|integer| 是 | 年级|就读年级|
|» age|body|integer| 是 | 年龄|年龄|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "info": "string",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|null|true|none|数据段|none|

# 学生接口

## GET 获取学生所有信息

GET /v1/api/stu-manager/stu-info

> 获取学生的所有信息

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "create_at": "2025-12-10T14:15:07.447+08:00",
    "name": "admin",
    "stu_id": "1234567890",
    "stu_class": "adminGroup",
    "password": "71483d953f4843cf69231307e77e8665",
    "sex": 0,
    "grade": 4,
    "age": 19,
    "courses": null
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||状态码|
|» info|string|true|none||响应信息|
|» data|object|true|none||数据段|
|»» create_at|string|true|none||注册时间|
|»» name|string|true|none||姓名|
|»» stu_id|string|true|none||学生ID|
|»» stu_class|string|true|none||班级|
|»» password|string|true|none||密码|
|»» sex|integer|true|none||性别|
|»» grade|integer|true|none||年级|
|»» age|integer|true|none||年龄|
|»» courses|null|true|none||选择的课程|

## GET 学生注销（退出登录）

GET /v1/api/stu-manager/stu-logout

> 学生退出登录

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|null|true|none|数据段|none|

## PATCH 学生修改学生信息

PATCH /v1/api/stu-manager/stu-update

> 更新学生信息

> Body 请求参数

```json
{
  "update_columns": [
    {
      "field": "student_class",
      "value": "计科7班"
    },
    {
      "field": "grade",
      "value": 2
    }
  ]
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» update_columns|body|[object]| 是 | 更新的字段列表|none|
|»» field|body|string| 是 | 字段|none|
|»» value|body|string| 是 | 值|none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|null|true|none|数据段|none|

# 学生选课接口

## GET 查看可选课程

GET /v1/api/class-manager/get-selectable-classes

> 学生查看可选课程

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "info": "string",
  "data": {
    "selectable_classes": [
      {
        "ID": 0,
        "create_at": "string",
        "class_name": "string",
        "class_id": "string",
        "class_location": "string",
        "class_time": "string",
        "class_teacher": "string",
        "class_capacity": 0,
        "class_selection": 0,
        "students": null
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» info|string|true|none||none|
|» data|object|true|none||none|
|»» selectable_classes|[object]|true|none||none|
|»»» ID|integer|false|none||none|
|»»» create_at|string|false|none||none|
|»»» class_name|string|false|none||none|
|»»» class_id|string|false|none||none|
|»»» class_location|string|false|none||none|
|»»» class_time|string|false|none||none|
|»»» class_teacher|string|false|none||none|
|»»» class_capacity|integer|false|none||none|
|»»» class_selection|integer|false|none||none|
|»»» students|null|false|none||none|

## GET 查看已选课程

GET /v1/api/class-manager/get-subscribed-classes

> 学生查看已选课程

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 添加选课

POST /v1/api/class-manager/subscribe-class

> 学生添加选课

> Body 请求参数

```json
{
  "class_id": "CS101"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_id|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "info": "string",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» info|string|true|none||none|
|» data|null|true|none||none|

## DELETE 退课

DELETE /v1/api/class-manager/del-class

> 学生退课

> Body 请求参数

```json
{
  "class_id": "CS101"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_id|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 学生管理接口

## GET 获取学生列表

GET /v1/api/admin/stu-manager/get-stu-list

> 获取所有学生的列表

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|page|query|string| 否 ||页码|
|resNum|query|string| 否 ||每页条数|
|Authorization|header|string| 是 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "total": 2,
    "page": 1,
    "page_size": 15,
    "students_list": [
      {
        "stu_id": "1234567890",
        "student_name": "admin",
        "student_class": "计科7班",
        "grade": "2"
      },
      {
        "stu_id": "9876543210",
        "student_name": "zjq",
        "student_class": "计科8班",
        "grade": "1"
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|object|true|none|数据段|none|
|»» total|integer|true|none|条目总数|none|
|»» page|integer|true|none|当前页码|none|
|»» page_size|integer|true|none|页面数据条数|none|
|»» students_list|[object]|true|none|学生列表|none|
|»»» stu_id|string|true|none|学生ID|none|
|»»» student_name|string|true|none|学生名字|none|
|»»» student_class|string|true|none|班级|none|
|»»» grade|string|true|none|年级|none|

## POST 创建学生

POST /v1/api/admin/stu-manager/create-stu

> 创建学生接口

> Body 请求参数

```json
{
  "name": "dzjq",
  "student_id": "9876543210",
  "student_class": "计科2班",
  "password": "P@ssw0rd=Ping12345",
  "sex": 1,
  "grade": 1,
  "age": 19
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» name|body|string| 是 | 姓名|none|
|» student_id|body|string| 是 | 学生ID|none|
|» student_class|body|string| 是 | 班级|none|
|» password|body|string| 是 | 密码|none|
|» sex|body|integer| 是 | 性别|none|
|» grade|body|integer| 是 | 年级|none|
|» age|body|integer| 是 | 年龄|none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PATCH 更新学生信息

PATCH /v1/api/admin/stu-manager/update-stu-info

> 管理员更新学生信息

> Body 请求参数

```json
{
  "stu_id": "9876543210",
  "update_columns": [
    {
      "field": "student_class",
      "value": "计科9班"
    }
  ]
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» stu_id|body|string| 是 | 学生ID|none|
|» update_columns|body|[object]| 是 | 更新字段列表|none|
|»» field|body|string| 否 | 字段|none|
|»» value|body|string| 否 | 值|none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|null|true|none|数据段|none|

## DELETE 删除学生

DELETE /v1/api/admin/stu-manager/del-stu

> 管理员删除学生接口

> Body 请求参数

```json
{
  "stu_id": "1122334455"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 课程管理接口

## GET 查看选课情况

GET /v1/api/admin/classes-manager/get-class-status

> 管理员查看选课的情况

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "selectable_classes": [
      {
        "ID": 1,
        "create_at": "2025-12-12T10:42:10.954+08:00",
        "class_name": "测试课程",
        "class_id": "CS101",
        "class_location": "2106教室",
        "class_time": "周一6，7节",
        "class_teacher": "God老师",
        "class_capacity": 11,
        "class_selection": 0,
        "students": null
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|object|true|none|数据段|none|
|»» selectable_classes|[object]|true|none|可选择的课程|none|
|»»» ID|integer|false|none|主键ID|none|
|»»» create_at|string|false|none|创建时间|none|
|»»» class_name|string|false|none|课程名称|none|
|»»» class_id|string|false|none|课程ID|none|
|»»» class_location|string|false|none|上课地点|none|
|»»» class_time|string|false|none|上课时间|none|
|»»» class_teacher|string|false|none|授课教师|none|
|»»» class_capacity|integer|false|none|课程容量|none|
|»»» class_selection|integer|false|none|已选课人数|none|
|»»» students|null|false|none|已选课学生|none|

## GET 查看单个学生的选课

GET /v1/api/admin/classes-manager/get-stu-classes/{stuID}

> 管理员查看单个学生的选课

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|stuID|path|string| 是 ||none|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 开始选课

GET /v1/api/admin/classes-manager/start-course-select-event

> 管理员开始选课

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 停止选课

GET /v1/api/admin/classes-manager/stop-course-select-event

> 管理员停止选课

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 添加课程

POST /v1/api/admin/classes-manager/add-course

> 管理员添加课程的接口

> Body 请求参数

```json
{
  "class_name": "测试课程",
  "class_id": "CS101",
  "class_location": "2106",
  "class_time": "周一6，7节",
  "class_teacher": "God",
  "class_capcity": 2
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_name|body|string| 是 ||none|
|» class_id|body|string| 是 ||none|
|» class_location|body|string| 是 ||none|
|» class_time|body|string| 是 ||none|
|» class_teacher|body|string| 是 ||none|
|» class_capcity|body|integer| 是 ||none|

> 返回示例

> 200 Response

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none|状态码|none|
|» info|string|true|none|响应信息|none|
|» data|null|true|none|数据段|none|

## PATCH 修改课程信息

PATCH /v1/api/admin/classes-manager/edit-class-info

> 管理员修改课程信息

> Body 请求参数

```json
{
  "class_id": "CS101",
  "update_columns": [
    {
      "field": "class_teacher",
      "value": "BigGod老师"
    },
    {
      "field": "class_location",
      "value": "9201教室"
    }
  ]
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_id|body|string| 是 ||none|
|» update_columns|body|[object]| 是 ||none|
|»» field|body|string| 是 ||none|
|»» value|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PATCH 添加学生选课

PATCH /v1/api/admin/classes-manager/update-stu-classes

> 管理员添加学生选课

> Body 请求参数

```json
{
  "stu_id": "9876543210",
  "update_class_id": "CS101"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» stu_id|body|string| 是 ||none|
|» update_class_id|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## DELETE 删除学生选课

DELETE /v1/api/admin/classes-manager/update-stu-classes

> 管理员删除学生选课

> Body 请求参数

```json
{
  "stu_id": "9876543210",
  "update_class_id": "CS101"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» stu_id|body|string| 是 ||none|
|» update_class_id|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PATCH 修改课程容量

PATCH /v1/api/admin/classes-manager/edit-class-stock

> 管理员修改课程容量

> Body 请求参数

```json
{
  "class_id": "CS101",
  "stock": 20
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_id|body|string| 是 ||none|
|» stock|body|integer| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## DELETE 删除课程

DELETE /v1/api/admin/classes-manager/delete-course

> 管理员删除课程接口

> Body 请求参数

```json
{
  "class_id": "CS101"
}
```

### 请求参数

|名称|位置|类型|必选|中文名|说明|
|---|---|---|---|---|---|
|Authorization|header|string| 否 ||none|
|body|body|object| 是 ||none|
|» class_id|body|string| 是 ||none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 数据模型

