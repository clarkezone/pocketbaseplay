FROM alpine
ADD pocketbase /pocketbase
cmd ./pocketbase serve --http 0.0.0.0:8080 --dir /data
EXPOSE 8080
