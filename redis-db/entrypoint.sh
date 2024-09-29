#!/bin/sh

# Create the directory for the Redis configuration file
mkdir -p /usr/local/etc/redis/

# Create a custom redis.conf with environment variables
cat <<EOF > /usr/local/etc/redis/redis.conf
user default off
user ${REDIS_USERNAME} on >${REDIS_PASSWORD} ~* +@all
dir /data
requirepass ${REDIS_PASSWORD}
EOF

# Start Redis with the generated configuration
exec redis-server /usr/local/etc/redis/redis.conf
