package sqlite

import (
	"context"
	"database/sql"
	"io/ioutil"

	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func BackupSQLite3(ctx context.Context) (string, error) {
	// 备份到临时文件
	tmpFile, err := ioutil.TempFile(`./`, `taoblog-*`)
	if err != nil {
		return ``, err
	}
	tmpFile.Close()

	// 目的数据库
	dstDB, err := sql.Open(`sqlite3`, tmpFile.Name())
	if err != nil {
		return ``, err
	}
	defer dstDB.Close()

	dstConn, err := dstDB.Conn(ctx)
	if err != nil {
		return ``, err
	}
	defer dstConn.Close()

	if err := dstConn.Raw(func(dstDC interface{}) error {
		rawDstConn := dstDC.(*sqlite3.SQLiteConn)

		// 源数据库
		srcConn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		defer srcConn.Close()

		if err := srcConn.Raw(func(srcDC interface{}) error {
			rawSrcConn := srcDC.(*sqlite3.SQLiteConn)

			// 备份函数调用
			backup, err := rawDstConn.Backup(`main`, rawSrcConn, `main`)
			if err != nil {
				return err
			}

			// errors can be safely ignored.
			_, _ = backup.Step(-1)

			if err := backup.Close(); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ``, err
	}

	zap.L().Info(`backuped to file`, zap.String(`path`, tmpFile.Name()))

	return tmpFile.Name(), nil
}
