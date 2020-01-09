ARG BUILDER=cloudfoundry/cnb:bionic

FROM golang:stretch AS lifecycle-build-task-builder
COPY ./ /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go get -d ./...
RUN go build -o /assets/lifecycle-build-task ./cmd/task
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM $BUILDER AS resource
USER root
COPY --from=lifecycle-build-task-builder /assets/ /usr/bin
WORKDIR /opt/resource

FROM resource AS tests
COPY --from=lifecycle-build-task-builder /tests /tests
ARG DOCKER_USERNAME
ARG DOCKER_PASSWORD
WORKDIR /tests
RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
 	done

FROM resource
