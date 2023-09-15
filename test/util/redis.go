package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/util"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func RedisClient(cfg *rest.Config, podKey types.NamespacedName) (*redis.Server, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, 6379)
	if err != nil {
		return nil, nil, err
	}

	rs, err := redis.NewServer(fmt.Sprintf("redis://localhost:%d", localPort), nil)
	if err != nil {
		return nil, nil, err
	}

	return rs, stopCh, nil
}

func SentinelClient(cfg *rest.Config, podKey types.NamespacedName) (*redis.Server, chan struct{}, error) {
	localPort, stopCh, err := PortForward(cfg, podKey, saasv1alpha1.SentinelPort)
	if err != nil {
		return nil, nil, err
	}

	ss, err := redis.NewServer(fmt.Sprintf("redis://localhost:%d", localPort), nil)
	if err != nil {
		return nil, nil, err
	}

	return ss, stopCh, nil
}

func LoadRedisDataset(ctx context.Context, srv *redis.Server, file string) error {

	// read csv file
	csvfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	data, err := CSVToMap(csvfile)
	if err != nil {
		return err
	}

	for key, value := range data {
		err := srv.RedisSet(ctx, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func CSVToMap(reader io.Reader) (map[string]string, error) {
	r := csv.NewReader(reader)
	m := map[string]string{}
	var headers []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if headers == nil {
			headers = record
		} else {
			m = util.MergeMaps(m, RecordToKeyValues(headers, record))
		}
	}
	return m, nil
}

func RecordToKeyValues(headers []string, record []string) map[string]string {
	m := map[string]string{}

	//skip first header, which is the row key
	for i := range headers[1:] {
		key := url.QueryEscape(strings.Join([]string{record[0], headers[i]}, "."))
		m[key] = record[i]
	}
	return m
}
