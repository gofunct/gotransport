//+build wireinject

package cloud

import (
	"context"
	"go.opencensus.io/trace"
	"time"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/google/wire"
	"gocloud.dev/aws/awscloud"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/mysql/rdsmysql"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/paramstore"
)

func AWS(ctx context.Context, runvar string) (*Runtime, func(), error) {
	wire.Build(
		InitializeApi,
		trace.AlwaysSample,
		AppHealthChecks,
		awscloud.AWS,
		rdsmysql.Open,
		NewRuntime,
		AwsBucket,    // *blob.Bucket
		AwsRunVar,    // runtimevar.Variable
		AwsSQLParams, // *rdsmysql.Params
	)
	return nil, nil, nil
}

func AwsBucket(ctx context.Context, cp awsclient.ConfigProvider, a API) (*blob.Bucket, error) {
	return s3blob.OpenBucket(ctx, cp, a.Bucket(), nil)
}

func AwsSQLParams(a API) *rdsmysql.Params {
	return &rdsmysql.Params{
		Endpoint: a.DbHost(),
		Database: a.DbName(),
		User:     a.DbUser(),
		Password: a.DbPassword(),
	}
}

func AwsRunVar(ctx context.Context, sess awsclient.ConfigProvider, a API) (*runtimevar.Variable, error) {
	dur, err := time.ParseDuration(a.RunVarTimeout())
	if err != nil {
		dur = 15 * time.Second
	}
	return paramstore.NewVariable(sess, a.RunVarValue(), runtimevar.StringDecoder, &paramstore.Options{
		WaitDuration: dur,
	})
}
