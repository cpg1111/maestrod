FROM debian:jessie
COPY ./dist/maestrod/ /opt/maestrod/
ENTRYPOINT ["/opt/maestrod/maestrod"]
