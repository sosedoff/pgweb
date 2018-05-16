package instances

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/stretchr/testify/assert"
)

type mockInstance struct {
	name    string
	host    string
	port    int64
	status  string
	cluster string
}

type mockCluster struct {
	name   string
	host   string
	reader string
	port   int64
	status string
}

type mockRDS struct {
	instances []mockInstance
	clusters  []mockCluster
}

func Test_FilterUnavailableInstances(t *testing.T) {
	m := &mockRDS{
		instances: []mockInstance{{
			name:   "ready",
			host:   "foo",
			port:   5432,
			status: "available",
		}, {
			name:   "notready",
			host:   "bar",
			port:   5432,
			status: "creating",
		}},
	}
	result, err := getEndpoints(m)

	assert.Equal(t, nil, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, "ready")
	assert.NotContains(t, result, "notready")
}

func Test_FilterClusterMembers(t *testing.T) {
	m := &mockRDS{
		instances: []mockInstance{{
			name:   "standalone",
			host:   "foo",
			port:   5432,
			status: "available",
		}, {
			name:    "cluster-member",
			host:    "bar",
			port:    5432,
			status:  "available",
			cluster: "cluster",
		}},
		clusters: []mockCluster{{
			name:   "cluster",
			host:   "writer",
			reader: "reader",
			port:   5432,
			status: "available",
		}},
	}
	result, err := getEndpoints(m)

	assert.Equal(t, nil, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "standalone")
	assert.Contains(t, result, "cluster")
	assert.NotContains(t, result, "cluster-member")
}

func Test_ReturnWriterForCluster(t *testing.T) {
	command.Opts.ReadOnly = false

	m := &mockRDS{
		clusters: []mockCluster{{
			name:   "cluster",
			host:   "writer",
			reader: "reader",
			port:   5432,
			status: "available",
		}},
	}
	result, err := getEndpoints(m)

	assert.Equal(t, nil, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, "cluster")
	assert.Equal(t, result["cluster"].Host, "writer")
}

func Test_ReturnReaderForClusterIfInReadonlyMode(t *testing.T) {
	command.Opts.ReadOnly = true

	m := &mockRDS{
		clusters: []mockCluster{{
			name:   "cluster",
			host:   "writer",
			reader: "reader",
			port:   5432,
			status: "available",
		}},
	}
	result, err := getEndpoints(m)

	assert.Equal(t, nil, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, "cluster")
	assert.Equal(t, result["cluster"].Host, "reader")
}

func (m *mockRDS) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	instances := []*rds.DBInstance{}

	for _, mock := range m.instances {
		var cluster *string

		if mock.cluster != "" {
			cluster = aws.String(mock.cluster)
		}

		instances = append(instances, &rds.DBInstance{
			DBInstanceStatus:     aws.String(mock.status),
			DBInstanceIdentifier: aws.String(mock.name),
			DBClusterIdentifier:  cluster,
			Endpoint: &rds.Endpoint{
				Address: aws.String(mock.host),
				Port:    aws.Int64(mock.port),
			},
		})
	}

	return &rds.DescribeDBInstancesOutput{
		DBInstances: instances,
	}, nil
}

func (m *mockRDS) DescribeDBClusters(input *rds.DescribeDBClustersInput) (*rds.DescribeDBClustersOutput, error) {
	clusters := []*rds.DBCluster{}

	for _, mock := range m.clusters {
		clusters = append(clusters, &rds.DBCluster{
			Status:              aws.String(mock.status),
			DBClusterIdentifier: aws.String(mock.name),
			Endpoint:            aws.String(mock.host),
			ReaderEndpoint:      aws.String(mock.reader),
			Port:                aws.Int64(mock.port),
		})
	}

	return &rds.DescribeDBClustersOutput{
		DBClusters: clusters,
	}, nil
}
