FROM gcr.io/distroless/static-debian11
WORKDIR /
ADD build .
EXPOSE 8080

ENTRYPOINT ["/rinha-app"]