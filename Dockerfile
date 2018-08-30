FROM alpine:latest
EXPOSE 3611
RUN apk add --update --no-cache ca-certificates
RUN apk add --update --no-cache tzdata
ENV TZ America/New_York
ADD bin/* /bin/
ADD script/* /script/
ADD docs/* /docs/
CMD ["/bin/dapi", "serve"]
