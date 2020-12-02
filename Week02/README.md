### 作业题目

>我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

### 我的回答

- dao层一般位于项目分层中的最底层，其在跟标准库或第三方库协作时，通常情况下应该将接收到的标准库或第三方库的错误值Wrap一下再向上抛。
- 因为标准库、第三方库返回的错误值一般是未经包装的根错误（root error），故doa层要Wrap记录堆栈、本层的错误信息或hit信息等。

针对于此问题中的`sql.ErrNoRows`错误值如何处理的问题，需要根据具体情况分析：

情况一：

当上层代码确保不会通过`errors.Is`或`errors.As`解开dao层抛出的经过包装的`sql.ErrNoRows`时，可以包装`sql.ErrNoRows`抛给上层并在最上层使用`%+v`记录日志，用于出错时排查问题。但这需要团队达成共识，并用代码审查等机制监督。

```go
func UserDaoGetUserByID(uid int) (user *User, err error) {
    var user *User
    rows, err := db.Query(sqlString, uid)
	if err != nil {
        // 包装向上抛，携带sql语句
        return nil, errors.Wrapf(err, "with sql statement %s", sqlString)
	}
    defer rows.Close()
    .....

    return user, nil
}
```
情况二:

如果dao层不认为`sql.ErrNoRows`是错误或者说将它与其他类型错误区别对待，不对`sql.ErrNoRows`Wrap也不向上抛，而是做降级处理，接收过旧项目就是这么写的。

```go
func UserDaoGetUserByID(uid int) (user *User, err error) {
    var user *User
    rows, err := db.Query(sqlString, uid)
    // 区别对待
    if err == sql.ErrNoRows {
        // 降级处理
        return user, nil
    }
	if err != nil {
        // 其他类型错误，包装向上抛
		return nil, errors.Wrapf(err, "User Not Found with sql statement %s", sqlString)
	}
    defer rows.Close()
    .....

    return user, nil
}
```

情况三：

为了让上层代码对底层所用数据库无感知、不对底层数据库错误进行硬编码而仅在排查错误时可以通过日志查看到所用数据库产生的错误信息，故用dao层自定义的根错误替换掉`sql.ErrNoRows`再包装再向上抛，包装时将数据库报错信息写入携带信息中。

```go
var ErrUserNotFound = errors.New("user dao: user not found")

func UserDaoGetUserByID(uid int) (user *User, err error) {
    var user *User
    rows, err := db.Query(sqlString, uid)
	if err != nil {
        // 对dao层自定义错误ErrUserNotFound包装向上抛，包装过程中将err信息写入附带信息里
        // 即使上层代码Unwrap得到也只是ErrUserNotFound
	    return nil, errors.Wrapf(ErrUserNotFound, "with sql %s got db error %s", sqlString, err)
	}
    defer rows.Close()
    .....

    return user, nil
}
```




