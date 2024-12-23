
.PHONY: bazel-remote-up
bazel-remote-up:
	mkdir -p "${HOME}/.cache/bazel/remote"
	USER_ID="$(shell id -u)" GROUP_ID="$(shell id -g)" docker compose -f bazel-remote.yml up -d

.PHONY: bazel-remote-down
bazel-remote-down:
	USER_ID="$(shell id -u)" GROUP_ID="$(shell id -g)" docker compose -f bazel-remote.yml down

.PHONY: build
build: bazel-remote-up
	bazel build //...

.PHONY: test
test: bazel-remote-up
	bazel test //...
