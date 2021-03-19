package aerospikeConnect

import (
	"github.com/aerospike/aerospike-client-go"
)

type Aerospike struct {
	Client   *aerospike.Client
	Policies struct {
		Read  *aerospike.BasePolicy
		Write *aerospike.WritePolicy
		Scan  *aerospike.ScanPolicy
		Query *aerospike.QueryPolicy
		Info  *aerospike.InfoPolicy
	}
}

type AerospikeConfig struct {
	Host                        *string `yaml:"host"`
	Port                        *int    `yaml:"port"`
	ConnectionQueueSize         *int    `yaml:"connectionQueueSize"`
	LimitConnectionsToQueueSize *bool   `yaml:"limitConnectionsToQueueSize"`
	ClusterName                 *string `yaml:"clusterName"`
	MaxConnectAttempts          *int    `yaml:"maxConnectAttempts"`
	MaxConnectRetryTimeMs       *int    `yaml:"maxConnectRetryTimeMs"`
	ConnectRetrySleepMs         *int    `yaml:"connectRetrySleepMs"`
	Policies                    *struct {
		Base *struct {
			TimeoutMs *struct {
				Connect *int `yaml:"connect"`
				Idle    *int `yaml:"idle"`
				Login   *int `yaml:"login"`
			} `yaml:"timeoutMs"`
		} `yaml:"base"`
		Read  *AerospikeConfigTransactionPolicy `yaml:"read"`
		Write *AerospikeConfigTransactionPolicy `yaml:"write"`
		Scan  *AerospikeConfigTransactionPolicy `yaml:"scan"`
		Query *AerospikeConfigTransactionPolicy `yaml:"query"`
		Info  *struct {
			TimeoutMs *int `yaml:"timeoutMs"`
		} `yaml:"info"`
	} `yaml:"policies"`
	Security *struct {
		Username         *string `yaml:"username"`
		Password         *string `yaml:"password"`
		AuthModeExternal *bool   `yaml:"authModeExternal"`
	} `yaml:"security"`
	TLS *struct {
		CaFile             *string `yaml:"caFile"`
		CertFile           *string `yaml:"certFile"`
		KeyFile            *string `yaml:"keyFile"`
		ServerName         *string `yaml:"serverName"`
		InsecureSkipVerify *bool   `yaml:"InsecureSkipVerify"`
	} `yaml:"tls"`
}

type AerospikeConfigTransactionPolicy struct {
	TimeoutMs *struct {
		Socket              *int `yaml:"socket"`
		Total               *int `yaml:"total"`
		SleepBetweenRetries *int `yaml:"sleepBetweenRetries"`
	} `yaml:"timeoutMs"`
	MaxRetries *int `yaml:"maxRetries"`
}
