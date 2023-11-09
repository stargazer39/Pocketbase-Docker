FROM alpine:latest
RUN mkdir /app
RUN adduser -D myuser
RUN chown myuser:myuser /app
USER myuser
WORKDIR /app
RUN wget "https://github.com/pocketbase/pocketbase/releases/download/v0.19.3/pocketbase_0.19.3_linux_amd64.zip"
RUN unzip pocketbase* -d /app
EXPOSE 8090
CMD ["./pocketbase", "serve", "--http", ":8090"]