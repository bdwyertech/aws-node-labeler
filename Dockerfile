FROM golang:1.17-alpine as helper
WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -ldflags="-s -w" -trimpath .

FROM gcr.io/distroless/static:latest-amd64

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="bdwyertech/aws-node-labeler" \
      org.opencontainers.image.description="Labels EKS Nodes with the Instance Attributes" \
      org.opencontainers.image.authors="Brian Dwyer <bdwyertech@github.com>" \
      org.opencontainers.image.url="https://hub.docker.com/r/bdwyertech/eks-lifecycle-labeler" \
      org.opencontainers.image.source="https://github.com/bdwyertech/eks-lifecycle-labeler.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE \
      org.label-schema.name="bdwyertech/aws-node-labeler" \
      org.label-schema.description="Labels EKS Nodes with Instance Attributes" \
      org.label-schema.url="https://hub.docker.com/r/bdwyertech/aws-node-labeler" \
      org.label-schema.vcs-url="https://github.com/bdwyertech/aws-node-labeler.git" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.build-date=$BUILD_DATE

COPY --from=helper /go/src/aws-node-labeler /.
CMD ["/aws-node-labeler"]
