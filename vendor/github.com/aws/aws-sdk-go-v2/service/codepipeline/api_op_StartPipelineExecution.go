// Code generated by smithy-go-codegen DO NOT EDIT.

package codepipeline

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Starts the specified pipeline. Specifically, it begins processing the latest
// commit to the source location specified as part of the pipeline.
func (c *Client) StartPipelineExecution(ctx context.Context, params *StartPipelineExecutionInput, optFns ...func(*Options)) (*StartPipelineExecutionOutput, error) {
	if params == nil {
		params = &StartPipelineExecutionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "StartPipelineExecution", params, optFns, c.addOperationStartPipelineExecutionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*StartPipelineExecutionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Represents the input of a StartPipelineExecution action.
type StartPipelineExecutionInput struct {

	// The name of the pipeline to start.
	//
	// This member is required.
	Name *string

	// The system-generated unique ID used to identify a unique execution request.
	ClientRequestToken *string

	// A list that allows you to specify, or override, the source revision for a
	// pipeline execution that's being started. A source revision is the version with
	// all the changes to your application code, or source artifact, for the pipeline
	// execution.
	SourceRevisions []types.SourceRevisionOverride

	// A list that overrides pipeline variables for a pipeline execution that's being
	// started. Variable names must match [A-Za-z0-9@\-_]+ , and the values can be
	// anything except an empty string.
	Variables []types.PipelineVariable

	noSmithyDocumentSerde
}

// Represents the output of a StartPipelineExecution action.
type StartPipelineExecutionOutput struct {

	// The unique system-generated ID of the pipeline execution that was started.
	PipelineExecutionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationStartPipelineExecutionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpStartPipelineExecution{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpStartPipelineExecution{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "StartPipelineExecution"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIdempotencyToken_opStartPipelineExecutionMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpStartPipelineExecutionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opStartPipelineExecution(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
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
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

type idempotencyToken_initializeOpStartPipelineExecution struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpStartPipelineExecution) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpStartPipelineExecution) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*StartPipelineExecutionInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *StartPipelineExecutionInput ")
	}

	if input.ClientRequestToken == nil {
		t, err := m.tokenProvider.GetIdempotencyToken()
		if err != nil {
			return out, metadata, err
		}
		input.ClientRequestToken = &t
	}
	return next.HandleInitialize(ctx, in)
}
func addIdempotencyToken_opStartPipelineExecutionMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpStartPipelineExecution{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opStartPipelineExecution(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "StartPipelineExecution",
	}
}