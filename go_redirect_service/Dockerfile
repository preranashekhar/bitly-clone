FROM golang:latest
EXPOSE 80
RUN mkdir /app
ADD ./src /app/src
WORKDIR /app
ENV GOPATH /app
RUN cd /app ; go install server
RUN rm -r /app/src
CMD ["/app/bin/server"]