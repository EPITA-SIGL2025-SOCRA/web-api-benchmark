#! /bin/sh
PORT=$1
if [ -z $PORT ]; then
    echo "MISSING PORT ARG"
    exit 1
fi
curl -XPOST -H "Content-Type: application/json" -d '{"position": {"lat": 48.8115336, "lng": 2.3681119}}' "http://127.0.0.1:$PORT/v1/tractors?radius=340"
