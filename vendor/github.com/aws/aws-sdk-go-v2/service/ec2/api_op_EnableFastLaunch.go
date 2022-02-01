// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// When you enable faster launching for a Windows AMI, images are pre-provisioned,
// using snapshots to launch instances up to 65% faster. To create the optimized
// Windows image, Amazon EC2 launches an instance and runs through Sysprep steps,
// rebooting as required. Then it creates a set of reserved snapshots that are used
// for subsequent launches. The reserved snapshots are automatically replenished as
// they are used, depending on your settings for launch frequency.
func (c *Client) EnableFastLaunch(ctx context.Context, params *EnableFastLaunchInput, optFns ...func(*Options)) (*EnableFastLaunchOutput, error) {
	if params == nil {
		params = &EnableFastLaunchInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "EnableFastLaunch", params, optFns, c.addOperationEnableFastLaunchMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*EnableFastLaunchOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type EnableFastLaunchInput struct {

	// The ID of the image for which you’re enabling faster launching.
	//
	// This member is required.
	ImageId *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// The launch template to use when launching Windows instances from pre-provisioned
	// snapshots. Launch template parameters can include either the name or ID of the
	// launch template, but not both.
	LaunchTemplate *types.FastLaunchLaunchTemplateSpecificationRequest

	// The maximum number of parallel instances to launch for creating resources.
	MaxParallelLaunches *int32

	// The type of resource to use for pre-provisioning the Windows AMI for faster
	// launching. Supported values include: snapshot, which is the default value.
	ResourceType *string

	// Configuration settings for creating and managing the snapshots that are used for
	// pre-provisioning the Windows AMI for faster launching. The associated
	// ResourceType must be snapshot.
	SnapshotConfiguration *types.FastLaunchSnapshotConfigurationRequest

	noSmithyDocumentSerde
}

type EnableFastLaunchOutput struct {

	// The image ID that identifies the Windows AMI for which faster launching was
	// enabled.
	ImageId *string

	// The launch template that is used when launching Windows instances from
	// pre-provisioned snapshots.
	LaunchTemplate *types.FastLaunchLaunchTemplateSpecificationResponse

	// The maximum number of parallel instances to launch for creating resources.
	MaxParallelLaunches *int32

	// The owner ID for the Windows AMI for which faster launching was enabled.
	OwnerId *string

	// The type of resource that was defined for pre-provisioning the Windows AMI for
	// faster launching.
	ResourceType types.FastLaunchResourceType

	// The configuration settings that were defined for creating and managing the
	// pre-provisioned snapshots for faster launching of the Windows AMI. This property
	// is returned when the associated resourceType is snapshot.
	SnapshotConfiguration *types.FastLaunchSnapshotConfigurationResponse

	// The current state of faster launching for the specified Windows AMI.
	State types.FastLaunchStateCode

	// The reason that the state changed for faster launching for the Windows AMI.
	StateTransitionReason *string

	// The time that the state changed for faster launching for the Windows AMI.
	StateTransitionTime *time.Time

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationEnableFastLaunchMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpEnableFastLaunch{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpEnableFastLaunch{}, middleware.After)
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
	if err = addOpEnableFastLaunchValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opEnableFastLaunch(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opEnableFastLaunch(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "EnableFastLaunch",
	}
}
