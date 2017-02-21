FROM debian:jessie
COPY ./dist/maestrod/ /opt/maestrod/
COPY ./dist/plugin.d/ /opt/maestrod/plugin.d/
ENTRYPOINT ["/opt/maestrod/maestrod"]
