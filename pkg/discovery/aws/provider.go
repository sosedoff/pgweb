package aws

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/go-multierror"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/discovery"
)

var (
	errAccessKeyMissing  = errors.New("AWS access key is not set")
	errSecretKeyMissing  = errors.New("AWS secret key is not set")
	errRegionMissing     = errors.New("AWS region is not set")
	errInstanceNotFound  = errors.New("Database instance not found")
	errInvalidResourceID = errors.New("Invalid resource ID")
)

// Provider represents AWS RDS instance discovery provider
type Provider struct {
	service *rds.RDS
}

// New returns a new AWS provider instance
func New(opts command.Options) (*Provider, error) {
	// Try to read configuration options from env vars first
	if opts.AWSAccessKey == "" {
		opts.AWSAccessKey = os.Getenv("AWS_ACCESS_KEY")
	}
	if opts.AWSSecretKey == "" {
		opts.AWSSecretKey = os.Getenv("AWS_SECRET_KEY")
	}
	if opts.AWSRegion == "" {
		opts.AWSRegion = os.Getenv("AWS_REGION")
	}

	// Load configuration options from the AWS CLI profile
	readProfile(&opts)

	// Validate options
	if opts.AWSAccessKey == "" {
		return nil, errAccessKeyMissing
	}
	if opts.AWSSecretKey == "" {
		return nil, errSecretKeyMissing
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
	var (
		mainerr error
		resLock sync.Mutex
		wg      sync.WaitGroup
	)

	resources := []discovery.Resource{}

	wg.Add(2)

	// Fetch postgres instances
	go func() {
		defer wg.Done()

		resp, err := p.service.DescribeDBInstances(&rds.DescribeDBInstancesInput{
			Filters: []*rds.Filter{{
				Name:   aws.String("engine"),
				Values: []*string{aws.String("postgres")},
			}},
		})
		if err != nil {
			mainerr = multierror.Append(mainerr, err)
			return
		}

		resLock.Lock()
		defer resLock.Unlock()

		for _, instance := range resp.DBInstances {
			// Skip instances that are part of a cluster
			if instance.DBClusterIdentifier != nil {
				continue
			}

			// Skip instances that aren't up
			if aws.StringValue(instance.DBInstanceStatus) != "available" {
				continue
			}

			resources = append(resources, discovery.Resource{
				ID:   fmt.Sprintf("instance/%s", *instance.DBInstanceIdentifier),
				Name: fmt.Sprintf("[instance] %s", *instance.DBInstanceIdentifier),
				Meta: map[string]interface{}{
					"instance_type":     *instance.DBInstanceClass,
					"availability_zone": *instance.AvailabilityZone,
				},
			})
		}
	}()

	// Fetch aurora clusters
	go func() {
		defer wg.Done()

		resp, err := p.service.DescribeDBClusters(&rds.DescribeDBClustersInput{
			Filters: []*rds.Filter{{
				Name:   aws.String("engine"),
				Values: []*string{aws.String("aurora-postgresql")},
			}},
		})

		if err != nil {
			mainerr = multierror.Append(mainerr, err)
			return
		}

		resLock.Lock()
		resLock.Unlock()

		for _, cluster := range resp.DBClusters {
			// Skip clusters that aren't up
			if aws.StringValue(cluster.Status) != "available" {
				continue
			}

			resources = append(resources, discovery.Resource{
				ID:   fmt.Sprintf("cluster/%s", *cluster.DBClusterIdentifier),
				Name: fmt.Sprintf("[cluster] %s", *cluster.DBClusterIdentifier),
			})
		}
	}()

	wg.Wait()
	return resources, mainerr
}

// Get returns RDS instance details
func (p Provider) Get(id string) (*bookmarks.Bookmark, error) {
	// Require resource ID formatted as "resouceType/resourceID".
	// Examples: cluster/mycluster, instance/myinstance.
	// This way we dont have to pass ARNs.
	items := strings.SplitN(id, "/", 2)
	if len(items) < 2 {
		return nil, errInvalidResourceID
	}

	switch items[0] {
	case "instance":
		resp, err := p.service.DescribeDBInstances(&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(items[1]),
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

	case "cluster":
		resp, err := p.service.DescribeDBClusters(&rds.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(items[1]),
		})
		if err != nil {
			return nil, err
		}
		if len(resp.DBClusters) == 0 {
			return nil, errInstanceNotFound
		}

		cluster := resp.DBClusters[0]

		var host string
		if command.Opts.ReadOnly {
			host = aws.StringValue(cluster.ReaderEndpoint)
		} else {
			host = aws.StringValue(cluster.Endpoint)
		}

		return &bookmarks.Bookmark{
			User: "postgres",
			Host: host,
			Port: int(aws.Int64Value(cluster.Port)),
			Ssl:  "verify-full",
		}, nil
	}

	return nil, errInvalidResourceID
}
