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

type RedisConfig struct {
	name  string
	value string
}

type RedisEventMessage struct {
	event  string
	target RedisInstanceDetails
	master RedisInstanceDetails
	config RedisConfig
}

func NewRedisEventMessage(msg *goredis.Message) (RedisEventMessage, error) {
	rem := RedisEventMessage{
		event:  msg.Channel,
		target: RedisInstanceDetails{},
		master: RedisInstanceDetails{},
		config: RedisConfig{},
	}

	if rem.event == "" {
		return RedisEventMessage{}, fmt.Errorf("emtpy event name")
	}

	if err := rem.parsePayload(strings.Split(msg.Payload, " ")); err != nil {
		return RedisEventMessage{}, err
	}

	return rem, nil

}

func (rem *RedisEventMessage) parsePayload(payload []string) error {
	switch rem.event {
	case "+tilt", "-tilt":
		return rem.parseEmptyPayload(payload)
	case "+switch-master":
		return rem.parseSwitchPayload(payload)
	case "+monitor", "+set", "+new-epoch", "+vote-for-leader":
		return rem.parseConfigurationPayload(payload)
	default:
		return rem.parseInstanceDetailsPayload(payload)
	}
}

func (rem *RedisEventMessage) parseEmptyPayload(payload []string) error {
	if len(payload) > 0 && payload[0] != "" {
		return fmt.Errorf("payload should be empty")
	}
	return nil
}

func (rem *RedisEventMessage) parseConfigurationPayload(payload []string) error {

	// <master name> <oldip> <oldport> <newip> <newport>
	// https://redis.io/docs/manual/sentinel/

	switch len(payload) {
	case 6:
		rem.config = RedisConfig{
			name:  payload[4],
			value: payload[5],
		}
		return rem.parseInstanceDetailsPayload(payload)
	case 1:
		rem.config = RedisConfig{
			value: payload[0],
		}
	case 2:
		rem.config = RedisConfig{
			name:  payload[0],
			value: payload[1],
		}
	default:
		return fmt.Errorf("invalid payload for this configuration message: %s", payload)
	}

	return nil

}

func (rem *RedisEventMessage) parseSwitchPayload(payload []string) error {

	// <master name> <oldip> <oldport> <newip> <newport>
	// https://redis.io/docs/manual/sentinel/

	if len(payload) != 5 {
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

	if len(payload) < 4 {
		return fmt.Errorf("invalid payload for instance details event: %s", payload)
	}

	rem.target = RedisInstanceDetails{
		role: payload[0],
		name: payload[1],
		ip:   payload[2],
		port: payload[3],
	}

	if rem.target.role != "master" && len(payload) < 8 {
		return fmt.Errorf("invalid payload for non-master instance event: %s", payload)
	}

	if len(payload) == 8 && payload[4] == "@" {
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
