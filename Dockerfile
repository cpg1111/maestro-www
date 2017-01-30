FROM golang:1.7
COPY . /opt/maestro-www/
WORKDIR /opt/maestro-www
ENTRYPOINT ["make"]
