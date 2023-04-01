// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Describes the progress of the AMI store tasks. You can describe the store tasks
// for specified AMIs. If you don't specify the AMIs, you get a paginated list of
// store tasks from the last 31 days. For each AMI task, the response indicates if
// the task is InProgress, Completed, or Failed. For tasks InProgress, the response
// shows the estimated progress as a percentage. Tasks are listed in reverse
// chronological order. Currently, only tasks from the past 31 days can be viewed.
// To use this API, you must have the required permissions. For more information,
// see Permissions for storing and restoring AMIs using Amazon S3
// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-store-restore.html#ami-s3-permissions)
// in the Amazon EC2 User Guide. For more information, see Store and restore an AMI
// using Amazon S3
// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-store-restore.html) in
// the Amazon EC2 User Guide.
func (c *Client) DescribeStoreImageTasks(ctx context.Context, params *DescribeStoreImageTasksInput, optFns ...func(*Options)) (*DescribeStoreImageTasksOutput, error) {
	if params == nil {
		params = &DescribeStoreImageTasksInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeStoreImageTasks", params, optFns, c.addOperationDescribeStoreImageTasksMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeStoreImageTasksOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeStoreImageTasksInput struct {

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// The filters.
	//
	// * task-state - Returns tasks in a certain state (InProgress |
	// Completed | Failed)
	//
	// * bucket - Returns task information for tasks that targeted
	// a specific bucket. For the filter value, specify the bucket name.
	Filters []types.Filter

	// The AMI IDs for which to show progress. Up to 20 AMI IDs can be included in a
	// request.
	ImageIds []string

	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see Pagination
	// (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination).
	// You cannot specify this parameter and the ImageIDs parameter in the same call.
	MaxResults *int32

	// The token returned from a previous paginated request. Pagination continues from
	// the end of the items returned by the previous request.
	NextToken *string

	noSmithyDocumentSerde
}

type DescribeStoreImageTasksOutput struct {

	// The token to include in another request to get the next page of items. This
	// value is null when there are no more items to return.
	NextToken *string

	// The information about the AMI store tasks.
	StoreImageTaskResults []types.StoreImageTaskResult

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeStoreImageTasksMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeStoreImageTasks{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeStoreImageTasks{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeStoreImageTasks(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

// DescribeStoreImageTasksAPIClient is a client that implements the
// DescribeStoreImageTasks operation.
type DescribeStoreImageTasksAPIClient interface {
	DescribeStoreImageTasks(context.Context, *DescribeStoreImageTasksInput, ...func(*Options)) (*DescribeStoreImageTasksOutput, error)
}

var _ DescribeStoreImageTasksAPIClient = (*Client)(nil)

// DescribeStoreImageTasksPaginatorOptions is the paginator options for
// DescribeStoreImageTasks
type DescribeStoreImageTasksPaginatorOptions struct {
	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see Pagination
	// (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination).
	// You cannot specify this parameter and the ImageIDs parameter in the same call.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeStoreImageTasksPaginator is a paginator for DescribeStoreImageTasks
type DescribeStoreImageTasksPaginator struct {
	options   DescribeStoreImageTasksPaginatorOptions
	client    DescribeStoreImageTasksAPIClient
	params    *DescribeStoreImageTasksInput
	nextToken *string
	firstPage bool
}

// NewDescribeStoreImageTasksPaginator returns a new
// DescribeStoreImageTasksPaginator
func NewDescribeStoreImageTasksPaginator(client DescribeStoreImageTasksAPIClient, params *DescribeStoreImageTasksInput, optFns ...func(*DescribeStoreImageTasksPaginatorOptions)) *DescribeStoreImageTasksPaginator {
	if params == nil {
		params = &DescribeStoreImageTasksInput{}
	}

	options := DescribeStoreImageTasksPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeStoreImageTasksPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeStoreImageTasksPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeStoreImageTasks page.
func (p *DescribeStoreImageTasksPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeStoreImageTasksOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.DescribeStoreImageTasks(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opDescribeStoreImageTasks(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "DescribeStoreImageTasks",
	}
}
