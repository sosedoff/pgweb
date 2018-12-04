package aws

import (
	"errors"
	"os"

	"github.com/sosedoff/pgweb/pkg/bookmarks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/discovery"
)

var (
	errAccessKeyMissing = errors.New("AWS access key is missing")
	errSecretKeyMissing = errors.New("AWS secret key is missing")
	errRegionMissing    = errors.New("AWS region is missing")
	errInstanceNotFound = errors.New("Database instance not found")
)

// Provider represents AWS RDS instance discovery provider
type Provider struct {
	service *rds.RDS
}

// New returns a new AWS provider instance
func New(opts command.Options) (*Provider, error) {
	if opts.AWSAccessKey == "" {
		opts.AWSAccessKey = os.Getenv("AWS_ACCESS_KEY")
	}
	if opts.AWSAccessKey == "" {
		return nil, errAccessKeyMissing
	}

	if opts.AWSSecretKey == "" {
		opts.AWSSecretKey = os.Getenv("AWS_SECRET_KEY")
	}
	if opts.AWSSecretKey == "" {
		return nil, errSecretKeyMissing
	}

	if opts.AWSRegion == "" {
		opts.AWSRegion = os.Getenv("AWS_REGION")
	}
	if opts.AWSRegion == "" {
		return nil, errRegionMissing
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(opts.AWSAccessKey, opts.AWSSecretKey, ""),
		Region:      aws.String(opts.AWSRegion),
	})
	if err != nil {
		return nil, err
	}
	service := rds.New(sess)

	return &Provider{service}, nil
}

// ID returns the provider ID
func (p Provider) ID() string {
	return "aws"
}

// Name returns the provider name
func (p Provider) Name() string {
	return "Amazon Web Services"
}

// List returns list of all RDS instances
func (p Provider) List() ([]discovery.Resource, error) {
	resp, err := p.service.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		Filters: []*rds.Filter{{
			Name:   aws.String("engine"),
			Values: []*string{aws.String("postgres")},
		}},
	})
	if err != nil {
		return nil, err
	}

	resources := []discovery.Resource{}

	for _, instance := range resp.DBInstances {
		// Skip instances that are part of a cluster, clusters get added later
		if instance.DBClusterIdentifier != nil {
			continue
		}

		// Skip instances that aren't up
		if aws.StringValue(instance.DBInstanceStatus) != "available" {
			continue
		}

		resources = append(resources, discovery.Resource{
			ID:   *instance.DBInstanceArn,
			Name: *instance.DBInstanceIdentifier,
			Meta: map[string]interface{}{
				"instance_type":     *instance.DBInstanceClass,
				"availability_zone": *instance.AvailabilityZone,
			},
		})
	}

	return resources, nil
}

// Get returns RDS instance details
func (p Provider) Get(id string) (*bookmarks.Bookmark, error) {
	resp, err := p.service.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(id),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.DBInstances) == 0 {
		return nil, errInstanceNotFound
	}

	instance := resp.DBInstances[0]

	return &bookmarks.Bookmark{
		User: "postgres",
		Host: aws.StringValue(instance.Endpoint.Address),
		Port: int(aws.Int64Value(instance.Endpoint.Port)),
		Ssl:  "verify-full",
	}, nil
}
