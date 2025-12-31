Gin使用 go-playground/validator/v10 进行验证

项目地址：https://github.com/go-playground/validato

文档地址：https://pkg.go.dev/github.com/go-playground/validator/v10

# 常用参数

Gin 的 `binding:"..."` 实际上传给的就是 `validator/v10` 的 `validate` 规则

## 一、基础必会（90% 接口都会用）

### 1. `required`

必填

```
Name string `binding:"required"`
```

- 字段必须存在
- 不能是零值（`""`, `0`, `nil`）

---

### 2. `omitempty`

可选字段（非常重要）

```
Age int `binding:"omitempty,min=18"`
```

- 没传：不校验
- 传了：必须满足后续规则

这是做 PATCH / 可选参数的核心 tag

---

### 3. `min` / `max`

最小 / 最大值

```
Age int `binding:"min=1,max=120"`
Name string `binding:"min=2,max=20"`
```

- 数字：数值大小
- 字符串 / slice：长度

---

### 4. `len`

固定长度

```
Code string `binding:"len=6"`
```

---

### 5. `oneof`

枚举值（极其常用）

```
Sex string `binding:"required,oneof=man woman"`
Status string `binding:"oneof=on off pending"`
```

---

## 字符串 / 格式校验（接口常用）

### 6. `email`

```
Email string `binding:"required,email"`
```

---

### 7. `url`

```
Avatar string `binding:"omitempty,url"`
```

---

### 8. `uuid`（v4）

```
UserID string `binding:"required,uuid4"`
```

---

### 9. `alphanum` / `alphanumunicode`

```
Username string `binding:"required,alphanum"`
```

---

### 10. `startswith` / `endswith`

```
Path string `binding:"startswith=/api"`
```

### 11 `contains` / `excludes`

---

## 三、数值 & 类型相关（很实用）

### 11. `gt` / `gte` / `lt` / `lte`

```
Age int `binding:"gte=18,lte=65"`
```

---

### 12. `eq` / `ne`

```
Role string `binding:"ne=admin"`
```

---

### 13. `numeric`

```
Phone string `binding:"numeric,len=11"`
```

---

## 四、切片 / 数组 / Map 校验

### 14. `dive`

校验每个元素（非常重要）

```
IDs []int `binding:"required,dive,gt=0"`
```

含义：

- slice 必须存在
- 每个元素都 > 0

---

### 15. `min` / `max` + slice

```
Tags []string `binding:"min=1,max=5,dive,min=2,max=10"`
```

---

## 16. `unique`

```
IDs []int `binding:"unique"`
```

---

## 字段间关系校验（进阶但常用）

### 17. `eqfield` / `nefield`

```
Password string `binding:"required"`
Confirm  string `binding:"required,eqfield=Password"`
```

---

### 18. `required_with` / `required_without`

```
Start string `binding:"required_with=End"`
End   string `binding:"required_with=Start"`
```

---

### 19. `required_if`

```
Reason string `binding:"required_if=Status off"`
```

---

# 时间 & 日期（项目常见）

### 20. `datetime`

```
Date string `binding:"datetime=2006-01-02"`
```

---

### 21. `time.Duration`

自动解析

```
Timeout time.Duration `binding:"gt=0"`
```

---

## Struct & 嵌套校验

### 22. `required` + struct

```
type Profile struct {
Age int `binding:"gte=18"`
}

type User struct {
Profile Profile `binding:"required"`
}
```

---

### 23. `omitempty` + 嵌套 struct

```
Profile *Profile `binding:"omitempty"`
```
