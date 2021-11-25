FROM golang:1.17-alpine as helper
WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -ldflags="-s -w" -trimpath .

FROM gcr.io/distroless/static:latest-amd64

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="bdwyertech/eks-lifecycle-labeler" \
      org.opencontainers.image.description="Labels EKS Nodes with the Instance Lifecycle" \
      org.opencontainers.image.authors="Brian Dwyer <bdwyertech@github.com>" \
      org.opencontainers.image.url="https://hub.docker.com/r/bdwyertech/eks-lifecycle-labeler" \
      org.opencontainers.image.source="https://github.com/bdwyertech/eks-lifecycle-labeler.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE \
      org.label-schema.name="bdwyertech/eks-lifecycle-labeler" \
      org.label-schema.description="Labels EKS Nodes with the Instance Lifecycle" \
      org.label-schema.url="https://hub.docker.com/r/bdwyertech/eks-lifecycle-labeler" \
      org.label-schema.vcs-url="https://github.com/bdwyertech/eks-lifecycle-labeler.git" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.build-date=$BUILD_DATE \
      org.tooling.user=tfkit \
      org.tooling.uid=1000 \
      org.tooling.gid=1000

COPY --from=helper /go/src/eks-lifecycle-labeler /.
CMD ["/eks-lifecycle-labeler"]
