ARG BUILDER=cloudfoundry/cnb:bionic


FROM concourse/registry-image-resource as registry-image-resource

FROM golang:stretch as pack-lifecycle-resource-builder
COPY ./ /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go get -d ./...
RUN go build -o /assets/out ./cmd/out
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM $BUILDER as resource
COPY --from=registry-image-resource /opt/resource/ /opt/resource
COPY --from=pack-lifecycle-resource-builder /assets/ /opt/resource
WORKDIR /opt/resource

FROM resource AS tests
ARG DOCKER_USERNAME
ARG DOCKER_PASSWORD
COPY --from=pack-lifecycle-resource-builder /tests /tests
RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
 	done

FROM resource
