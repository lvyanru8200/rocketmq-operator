package configs

import (
	"os"
	"strings"

	errors2 "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

var globalConfig Config

const (
	DEBUG = "DEBUG"

	// broker env
	JVM_XMS = "Xms"
	JVM_XMX = "Xmx"
	// nameserver env
	NsEnvPrefix = "NS_"
	NS_JVM_XMS  = NsEnvPrefix + JVM_XMS
	NS_JVM_XMX  = NsEnvPrefix + JVM_XMX

	// exporter
	SECRET_KEY = "SECRET_KEY"
	ACCESS_KEY = "ACCESS_KEY"

	Empty = "EMPTY" // crd有bug，空的env会被填上值，使用empty占位符
)

func GetGlobalConfig() Config {
	return globalConfig
}

func init() {
	globalConfig = configFromEnv()
}

type Config struct {
	MOCK_RANDOM_PORT string

	SERVICE_ACCOUNT string
	CLUSTER_ROLE    string

	IMAGE_ROCKETMQ     string
	IMAGE_EXPORTER     string
	STORAGE_CLASS_NAME string

	BROKER_CONFIG_MAP string
	ACL_CONFIG_MAP    string

	INSTANCE_ENV string
	InstanceEnv  []corev1.EnvVar
}

func configFromEnv() Config {
	c := Config{
		MOCK_RANDOM_PORT: getEnv("MOCK_RANDOM_PORT", ""),

		IMAGE_ROCKETMQ:     getEnv("IMAGE_ROCKETMQ", "harbor.dsp.local/middleware/rocketmq:4.6.1"),
		IMAGE_EXPORTER:     getEnv("IMAGE_EXPORTER", "harbor.dsp.local/middleware/rocketmq-exporter:0.0.1"),
		STORAGE_CLASS_NAME: getEnv("STORAGE_CLASS_NAME", "managed-nfs-storage"),

		SERVICE_ACCOUNT:   getEnv("SERVICE_ACCOUNT", "rocketmq-operator-instance"),
		CLUSTER_ROLE:      getEnv("CLUSTER_ROLE", "rocketmq-operator-instance"),
		BROKER_CONFIG_MAP: getEnv("BROKER_CONFIG_MAP", "rocketmq-default-broker-config"),
		ACL_CONFIG_MAP:    getEnv("ACL_CONFIG_MAP", "rocketmq-default-plain-acl"),
		INSTANCE_ENV:      getEnv("INSTANCE_ENV", ""),
	}

	kvs, err := parseKV(c.INSTANCE_ENV)
	if err != nil {
		panic(err)
	}
	c.InstanceEnv = func(kvs [][]string) []corev1.EnvVar {
		var r []corev1.EnvVar
		for _, kv := range kvs {
			r = append(r, corev1.EnvVar{
				Name:  kv[0],
				Value: kv[1],
			})
		}
		return r
	}(kvs)

	return c
}

func getEnv(name, defaultVal string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultVal
	}
	return val
}

func parseKV(s string) ([][]string, error) {
	kvarr := strings.Split(s, ";")
	var r [][]string

	for _, kvstr := range kvarr {
		if kvstr == "" {
			continue
		}

		kv := strings.SplitN(kvstr, "=", 2)
		if len(kv) < 2 || kv[1] == "" {
			return nil, errors2.Errorf("invalid kv: %s", kvstr)
		}
		r = append(r, kv)
	}

	return r, nil
}

func SetEnvIfUnset(env []corev1.EnvVar, key, def string) []corev1.EnvVar {
	val, ok := LookupEnv(env, key)
	if ok {
		if val == "" {
			return SetEnv(env, key, Empty)
		}
		return env
	}

	if def == "" {
		def = Empty
	}

	return SetEnv(env, key, def)
}

func LookupEnv(env []corev1.EnvVar, key string) (string, bool) {
	for _, e := range env {
		if e.Name == key {
			return e.Value, true
		}
	}
	return "", false
}

func SetEnv(env []corev1.EnvVar, key, val string) []corev1.EnvVar {
	var found bool
	for i := range env {
		if env[i].Name == key {
			found = true
			env[i].Value = val
			break
		}
	}

	if !found {
		env = append(env, corev1.EnvVar{
			Name:  key,
			Value: val,
		})
	}

	return env
}

func MergeEnv(src, dst []corev1.EnvVar) []corev1.EnvVar {
	//  使用set会去重
	set := make(map[string]string)
	for _, e := range src {
		set[e.Name] = e.Value
	}

	for _, e := range dst {
		_, ok := set[e.Name]
		if !ok {
			set[e.Name] = e.Value
		}
	}

	var r []corev1.EnvVar
	for k, v := range set {
		r = append(r, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return r
}
