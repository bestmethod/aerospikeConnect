# AerospikeConnect

The library is a wrapper for Aerospike golang library. It is a yaml-support system.

With a single, or embedded yaml-struct, the Aerospike config can be set and parsed.

The library will parse the yaml config, sanitise it, connect to Aerospike and precreate client policies.

## Usage

1. Grab an example config.yml from here

2. Embed the aerospike config struct in your config struct:

```go
type fullConfig struct {
	AerospikeConfig *aerospikeConnect.AerospikeConfig `yaml:"aerospike"`
}
```

3. Parse the config:

```go
config := new(fullConfig)
err := ParseConfig("/path/to/config.yml", config)
if err != nil {
    log.Fatal(err)
}
```

4. Connect to Aerospike:

```go
aero, err := aerospikeConnect.Connect(fullConfig.AerospikeConfig)
if err != nil {
    log.Fatal(err)
}
```

The `aero` object contains the client under `aero.Client` and all policies under `aero.Policies`.

## Example configuration file

See [config.yml](config.yml) for a full reference.

The only mandatory configuration is the `Host` and `Port` parts.

## Function Reference:

```go
func ParseConfig(filePath string, ac interface{}) error
func Connect(ac *AerospikeConfig) (*Aerospike, error)
```

## Aerospike struct:

```go
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
```

## AerospikeConfig struct, used to parse the yaml file:

```go
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
```