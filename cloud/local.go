//+build wireinject

package cloud

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/requestlog"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/filevar"
	"gocloud.dev/server"
	"time"
)

func LOCAL(ctx context.Context, runvar string) (*Runtime, func(), error) {
	wire.Build(
		InitializeApi,
		trace.AlwaysSample,
		AppHealthChecks,
		wire.InterfaceValue(new(requestlog.Logger), requestlog.Logger(nil)),
		wire.InterfaceValue(new(trace.Exporter), trace.Exporter(nil)),
		server.Set,
		NewRuntime,
		DialLocalSQL,
		LocalBucket,
		LocalRuntimeVar,
	)
	return nil, nil, nil
}

// localBucket is a Wire provider function that returns a directory-based bucket
// based on the command-line a.
func LocalBucket(a API) (*blob.Bucket, error) {
	return fileblob.OpenBucket(a.Bucket(), nil)
}

// dialLocalSQL is a Wire provider function that connects to a MySQL database
// (usually on localhost).
func DialLocalSQL(a API) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 a.DbHost(),
		DBName:               a.DbName(),
		User:                 a.DbUser(),
		Passwd:               a.DbPassword(),
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

func LocalRuntimeVar(a API) (*runtimevar.Variable, func(), error) {
	dur, err := time.ParseDuration(a.RunVarTimeout())
	if err != nil {
		dur = 15 * time.Second
	}
	v, err := filevar.New(a.RunVarValue(), runtimevar.StringDecoder, &filevar.Options{
		WaitDuration: dur,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
