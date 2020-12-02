package main

import (
	"database/sql"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
)

func IntService() (int, error) {
	i, err := IntDao()
	if err != nil {
		return 0, err
	}
	return i, nil
}

var ErrUserNotFound = errors.New("User Not Found")

func IntDao() (int, error) {
	_, err := DBQuery()
	if err != nil {
		return 0, xerrors.Wrapf(ErrUserNotFound, "with db error: %s", err)
	}
	return 1, nil
}
func DBQuery() (int, error) {
	return 0, sql.ErrNoRows
}

func main() {
	i, err := IntService()
	if err != nil {
		// 上层代码无法获取sql.ErrNoRows，仅可在日志中查看
		fmt.Printf("%+v\n\n %+v\n\n %+v", xerrors.Cause(err), errors.Unwrap(err), err)
	}
	fmt.Println(i)
}
