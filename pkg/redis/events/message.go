package events

import (
	"fmt"
	"strings"

	goredis "github.com/go-redis/redis/v8"
)

type RedisInstanceDetails struct {
	name string
	ip   string
	port string
	role string
}

type RedisEventMessage struct {
	event  string
	target RedisInstanceDetails
	master RedisInstanceDetails
}

func NewRedisEventMessage(msg *goredis.Message) (RedisEventMessage, error) {
	rem := RedisEventMessage{
		event:  msg.Channel,
		target: RedisInstanceDetails{},
		master: RedisInstanceDetails{},
	}
	err := rem.parsePayload(strings.Split(msg.Payload, " "))
	return rem, err
}

func (rem *RedisEventMessage) parsePayload(payload []string) error {
	switch rem.event {
	case "+ilt", "-tilt":
		return rem.parseTiltPayload(payload)
	case "+switch-master":
		return rem.parseSwitchPayload(payload)
	default:
		return rem.parseInstanceDetailsPayload(payload)
	}
}
func (rem *RedisEventMessage) parseTiltPayload(payload []string) error {

	if len(payload) > 0 {
		return fmt.Errorf("payload for tilt events should be empty")
	}

	rem.master = RedisInstanceDetails{}
	rem.target = RedisInstanceDetails{}

	return nil
}

func (rem *RedisEventMessage) parseSwitchPayload(payload []string) error {

	// <master name> <oldip> <oldport> <newip> <newport>
	// https://redis.io/docs/manual/sentinel/

	if len(payload) == 0 {
		return fmt.Errorf("empty payload for switch event")
	}

	if len(payload) < 5 {
		return fmt.Errorf("invalid payload for switch event: %s", payload)
	}

	rem.target = RedisInstanceDetails{
		role: "master",
		name: payload[0],
		ip:   payload[1],
		port: payload[2],
	}

	rem.master = RedisInstanceDetails{
		role: "master",
		name: payload[0],
		ip:   payload[3],
		port: payload[4],
	}

	return nil

}

func (rem *RedisEventMessage) parseInstanceDetailsPayload(payload []string) error {

	// <instance-type> <name> <ip> <port> @ <master-name> <master-ip> <master-port>
	// The part identifying the master (from the @ argument to the end) is optional
	// and is only specified if the instance is not a master itself.
	// https://redis.io/docs/manual/sentinel/

	if len(payload) == 0 {
		return fmt.Errorf("empty payload for intance event")
	}

	if len(payload) < 4 {
		return fmt.Errorf("invalid payload for instance event: %s", payload)
	}

	rem.target = RedisInstanceDetails{
		role: payload[0],
		name: payload[1],
		ip:   payload[2],
		port: payload[3],
	}

	if len(payload) == 8 {
		rem.master = RedisInstanceDetails{
			role: "master",
			name: payload[5],
			ip:   payload[6],
			port: payload[7],
		}
	} else {
		rem.master = RedisInstanceDetails{
			role: payload[0],
			name: payload[1],
			ip:   payload[2],
			port: payload[3],
		}
	}

	return nil

}
