//+build wireinject

package cloud

import (
	"context"
	"go.opencensus.io/trace"
	"time"

	"github.com/google/wire"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/gcpcloud"
	"gocloud.dev/mysql/cloudmysql"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/runtimeconfigurator"
	pb "google.golang.org/genproto/googleapis/cloud/runtimeconfig/v1beta1"
)

func GCP(ctx context.Context, runvar string) (*Runtime, func(), error) {
	wire.Build(
		InitializeApi,
		trace.AlwaysSample,
		AppHealthChecks,
		gcpcloud.GCP,
		// Provides:
		// server.Server,
		// *certs.RemoteCertSource,
		// http.RoundTripper,
		// gcp.HTTPClient,
		// pb.RuntimeConfigManagerClient,
		// func(),
		// requeslog.Logger,
		// trace.Exporter,
		// monitoredresource.Interface
		// error,
		cloudmysql.Open,
		//Provides:
		// *sql.Db
		NewRuntime,
		// Provides:
		// *Runtime
		GcpBucket,
		// Provides:
		// *blob.Bucket,
		// error
		GcpRunVar,
		//Provides:
		// *cloudmysql.Params
		GcpSqlParams,
		// Provides:
		// *runtimevar.Variable,
		// func(),
		// error
	)
	return nil, nil, nil
}

func GcpBucket(ctx context.Context, a API, client *gcp.HTTPClient) (*blob.Bucket, error) {
	return gcsblob.OpenBucket(ctx, client, a.Bucket(), nil)
}

func GcpSqlParams(id gcp.ProjectID, a API) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    a.DbRegion(),
		Instance:  a.DbHost(),
		Database:  a.DbName(),
		User:      a.DbUser(),
		Password:  a.DbPassword(),
	}
}

func GcpRunVar(ctx context.Context, client pb.RuntimeConfigManagerClient, project gcp.ProjectID, a API) (*runtimevar.Variable, func(), error) {
	name := runtimeconfigurator.ResourceName{
		ProjectID: string(project),
		Config:    a.RunVarName(),
		Variable:  a.RunVarValue(),
	}
	dur, err := time.ParseDuration(a.RunVarTimeout())
	if err != nil {
		dur = 15 * time.Second
	}
	v, err := runtimeconfigurator.NewVariable(client, name, runtimevar.StringDecoder, &runtimeconfigurator.Options{
		WaitDuration: dur,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
