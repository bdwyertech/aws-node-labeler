// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Adds the specified inbound (ingress) rules to a security group. An inbound rule
// permits instances to receive traffic from the specified IPv4 or IPv6 CIDR
// address range, or from the instances that are associated with the specified
// destination security groups. When specifying an inbound rule for your security
// group in a VPC, the IpPermissions must include a source for the traffic. You
// specify a protocol for each rule (for example, TCP). For TCP and UDP, you must
// also specify the destination port or port range. For ICMP/ICMPv6, you must also
// specify the ICMP/ICMPv6 type and code. You can use -1 to mean all types or all
// codes. Rule changes are propagated to instances within the security group as
// quickly as possible. However, a small delay might occur. For more information
// about VPC security group quotas, see Amazon VPC quotas
// (https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html). We
// are retiring EC2-Classic. We recommend that you migrate from EC2-Classic to a
// VPC. For more information, see Migrate from EC2-Classic to a VPC
// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/vpc-migrate.html) in the
// Amazon Elastic Compute Cloud User Guide.
func (c *Client) AuthorizeSecurityGroupIngress(ctx context.Context, params *AuthorizeSecurityGroupIngressInput, optFns ...func(*Options)) (*AuthorizeSecurityGroupIngressOutput, error) {
	if params == nil {
		params = &AuthorizeSecurityGroupIngressInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "AuthorizeSecurityGroupIngress", params, optFns, c.addOperationAuthorizeSecurityGroupIngressMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*AuthorizeSecurityGroupIngressOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type AuthorizeSecurityGroupIngressInput struct {

	// The IPv4 address range, in CIDR format. You can't specify this parameter when
	// specifying a source security group. To specify an IPv6 address range, use a set
	// of IP permissions. Alternatively, use a set of IP permissions to specify
	// multiple rules and a description for the rule.
	CidrIp *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// If the protocol is TCP or UDP, this is the start of the port range. If the
	// protocol is ICMP, this is the type number. A value of -1 indicates all ICMP
	// types. If you specify all ICMP types, you must specify all ICMP codes.
	// Alternatively, use a set of IP permissions to specify multiple rules and a
	// description for the rule.
	FromPort *int32

	// The ID of the security group. You must specify either the security group ID or
	// the security group name in the request. For security groups in a nondefault VPC,
	// you must specify the security group ID.
	GroupId *string

	// [EC2-Classic, default VPC] The name of the security group. You must specify
	// either the security group ID or the security group name in the request. For
	// security groups in a nondefault VPC, you must specify the security group ID.
	GroupName *string

	// The sets of IP permissions.
	IpPermissions []types.IpPermission

	// The IP protocol name (tcp, udp, icmp) or number (see Protocol Numbers
	// (http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)). To
	// specify icmpv6, use a set of IP permissions. [VPC only] Use -1 to specify all
	// protocols. If you specify -1 or a protocol other than tcp, udp, or icmp, traffic
	// on all ports is allowed, regardless of any ports you specify. Alternatively, use
	// a set of IP permissions to specify multiple rules and a description for the
	// rule.
	IpProtocol *string

	// [EC2-Classic, default VPC] The name of the source security group. You can't
	// specify this parameter in combination with the following parameters: the CIDR IP
	// address range, the start of the port range, the IP protocol, and the end of the
	// port range. Creates rules that grant full ICMP, UDP, and TCP access. To create a
	// rule with a specific IP protocol and port range, use a set of IP permissions
	// instead. For EC2-VPC, the source security group must be in the same VPC.
	SourceSecurityGroupName *string

	// [nondefault VPC] The Amazon Web Services account ID for the source security
	// group, if the source security group is in a different account. You can't specify
	// this parameter in combination with the following parameters: the CIDR IP address
	// range, the IP protocol, the start of the port range, and the end of the port
	// range. Creates rules that grant full ICMP, UDP, and TCP access. To create a rule
	// with a specific IP protocol and port range, use a set of IP permissions instead.
	SourceSecurityGroupOwnerId *string

	// [VPC Only] The tags applied to the security group rule.
	TagSpecifications []types.TagSpecification

	// If the protocol is TCP or UDP, this is the end of the port range. If the
	// protocol is ICMP, this is the code. A value of -1 indicates all ICMP codes. If
	// you specify all ICMP types, you must specify all ICMP codes. Alternatively, use
	// a set of IP permissions to specify multiple rules and a description for the
	// rule.
	ToPort *int32

	noSmithyDocumentSerde
}

type AuthorizeSecurityGroupIngressOutput struct {

	// Returns true if the request succeeds; otherwise, returns an error.
	Return *bool

	// Information about the inbound (ingress) security group rules that were added.
	SecurityGroupRules []types.SecurityGroupRule

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationAuthorizeSecurityGroupIngressMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpAuthorizeSecurityGroupIngress{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpAuthorizeSecurityGroupIngress{}, middleware.After)
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opAuthorizeSecurityGroupIngress(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opAuthorizeSecurityGroupIngress(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "AuthorizeSecurityGroupIngress",
	}
}
