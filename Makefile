.PHONY: clean
clean:
	bazel clean

.PHONY: build
build:
	bazel build //...

.PHONY: test
test:
	bazel test //...

.PHONY: deploy-buildkite-oidc-stack
deploy-buildkite-oidc-stack:
	sam deploy \
		--template-file ./sam/app/oidc-buildkite.yaml --stack-name ${OIDC_STACK_NAME} \
		--resolve-s3

.PHONY: deploy-cache-bucket-stack
deploy-cache-bucket-stack:
	$(eval BUILDKITE_OIDC := $(shell aws cloudformation describe-stacks --stack-name ${OIDC_STACK_NAME} --query 'Stacks[].Outputs[?OutputKey==`BuildkiteOidc`].OutputValue' --output text))
	sam deploy \
		--template-file ./sam/app/bazel-cache-bucket.cfn.yaml \
		--stack-name ${ORGANIZATION_SLUG}-${PIPELINE_SLUG}-bazel-cache-bucket \
    --capabilities CAPABILITY_IAM \
		--resolve-s3 \
		--parameter-overrides \
			BuildkiteOidc=${BUILDKITE_OIDC} \
			OrginizationSlug=${ORGANIZATION_SLUG} \
			PipelineSlug=${PIPELINE_SLUG}
