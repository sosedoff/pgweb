package instances

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/patrickmn/go-cache"

	"github.com/sosedoff/pgweb/pkg/command"
)

type Instance struct {
	Host string `json:"host"` // Server hostname
	Port int64  `json:"port"` // Server port
}

type RDS interface {
	DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(input *rds.DescribeDBClustersInput) (*rds.DescribeDBClustersOutput, error)
}

var c = cache.New(1*time.Minute, 2*time.Minute)

func GetAll() (map[string]Instance, error) {
	if !command.Opts.AwsInstanceDiscovery {
		return nil, nil
	}

	cached, found := c.Get("instances")
	if found {
		return cached.(map[string]Instance), nil
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	region := command.Opts.AwsRegion
	if region == "" {
		m := ec2metadata.New(sess)
		region, err = m.Region()
		if err != nil {
			return nil, err
		}
	}

	r := rds.New(sess, &aws.Config{Region: aws.String(region)})

	result, err := getEndpoints(r)
	if err != nil {
		return nil, err
	}

	c.Add("instances", result, cache.DefaultExpiration)

	return result, nil
}

func getEndpoints(r RDS) (map[string]Instance, error) {
	result := map[string]Instance{}

	resp, err := r.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		Filters: []*rds.Filter{{
			Name:   aws.String("engine"),
			Values: []*string{aws.String("postgres")},
		}},
	})
	if err != nil {
		return nil, err
	}

	for _, instance := range resp.DBInstances {
		if instance.DBClusterIdentifier != nil {
			// skip instances that are part of a cluster, clusters get added later
			continue
		}
		if aws.StringValue(instance.DBInstanceStatus) != "available" {
			// skip instances that aren't up
			continue
		}
		result[aws.StringValue(instance.DBInstanceIdentifier)] = Instance{
			Host: aws.StringValue(instance.Endpoint.Address),
			Port: aws.Int64Value(instance.Endpoint.Port),
		}
	}

	resp2, err := r.DescribeDBClusters(&rds.DescribeDBClustersInput{
		Filters: []*rds.Filter{{
			Name:   aws.String("engine"),
			Values: []*string{aws.String("aurora-postgresql")},
		}},
	})
	if err != nil {
		return nil, err
	}

	for _, cluster := range resp2.DBClusters {
		if aws.StringValue(cluster.Status) != "available" {
			// skip clusters that aren't up
			continue
		}

		var host string

		if command.Opts.ReadOnly {
			host = aws.StringValue(cluster.ReaderEndpoint)
		} else {
			host = aws.StringValue(cluster.Endpoint)
		}

		result[aws.StringValue(cluster.DBClusterIdentifier)] = Instance{
			Host: host,
			Port: aws.Int64Value(cluster.Port),
		}
	}

	return result, nil
}
