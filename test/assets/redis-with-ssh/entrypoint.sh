#!/bin/sh

set -x

# Start SSH service in the background
/usr/sbin/sshd -D &

# Redirect Redis logs to /dev/stdout
redis-server /redis/redis.conf