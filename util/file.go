/*
@Time : 2019-01-25 16:43
@Author : seefan
@File : file
@Software: microgo
*/
package util

import "os"

func FileIsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}
