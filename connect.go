package aerospikeConnect

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aerospike/aerospike-client-go"
)

func sanityCheck(ac *AerospikeConfig) error {
	if ac.Host == nil {
		return errors.New("aerospike seed host not specified")
	}
	if ac.Port == nil {
		return errors.New("aerospike seed port not specified")
	}
	return nil
}

func setConnectPolicy(ac *AerospikeConfig) (*aerospike.ClientPolicy, error) {
	connectPolicy := aerospike.NewClientPolicy()
	connectPolicy.FailIfNotConnected = true
	if ac.Policies != nil && ac.Policies.Base != nil && ac.Policies.Base.TimeoutMs != nil {
		if ac.Policies.Base.TimeoutMs.Idle != nil {
			connectPolicy.IdleTimeout = time.Millisecond * time.Duration(*ac.Policies.Base.TimeoutMs.Idle)
		}
		if ac.Policies.Base.TimeoutMs.Connect != nil {
			connectPolicy.Timeout = time.Millisecond * time.Duration(*ac.Policies.Base.TimeoutMs.Connect)
		}
		if ac.Security != nil && (ac.Security.Username != nil || ac.Security.Password != nil) && ac.Policies.Base.TimeoutMs.Login != nil {
			connectPolicy.LoginTimeout = time.Millisecond * time.Duration(*ac.Policies.Base.TimeoutMs.Login)
		}
	}
	if ac.Security != nil && (ac.Security.Username != nil || ac.Security.Password != nil) {
		if ac.Security.Username != nil {
			connectPolicy.User = *ac.Security.Username
		}
		if ac.Security.Password != nil {
			connectPolicy.Password = *ac.Security.Password
		}
		if ac.Security.AuthModeExternal != nil {
			if *ac.Security.AuthModeExternal {
				connectPolicy.AuthMode = aerospike.AuthModeExternal
			} else {
				connectPolicy.AuthMode = aerospike.AuthModeInternal
			}
		}
	}
	if ac.ConnectionQueueSize != nil {
		connectPolicy.ConnectionQueueSize = *ac.ConnectionQueueSize
	}
	if ac.LimitConnectionsToQueueSize != nil {
		connectPolicy.LimitConnectionsToQueueSize = *ac.LimitConnectionsToQueueSize
	}
	if ac.ClusterName != nil {
		connectPolicy.ClusterName = *ac.ClusterName
	}
	if ac.TLS != nil && (ac.TLS.CaFile != nil || ac.TLS.CertFile != nil || ac.TLS.KeyFile != nil || ac.TLS.ServerName != nil) {
		nTLS := new(tls.Config)
		if ac.TLS.InsecureSkipVerify != nil {
			nTLS.InsecureSkipVerify = *ac.TLS.InsecureSkipVerify
		}
		if ac.TLS.ServerName != nil {
			nTLS.ServerName = *ac.TLS.ServerName
		}
		if ac.TLS.CaFile != nil {
			caCert, err := ioutil.ReadFile(*ac.TLS.CaFile)
			if err != nil {
				return nil, fmt.Errorf("tls: loadca: %s", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			nTLS.RootCAs = caCertPool
		}
		if ac.TLS.CertFile != nil || ac.TLS.KeyFile != nil {
			certFile := ""
			keyFile := ""
			if ac.TLS.CertFile != nil {
				certFile = *ac.TLS.CertFile
			}
			if ac.TLS.KeyFile != nil {
				keyFile = *ac.TLS.KeyFile
			}
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("tls: loadkeys: %s", err)
			}
			nTLS.Certificates = []tls.Certificate{cert}
			nTLS.BuildNameToCertificate()
		}
		connectPolicy.TlsConfig = nTLS
	}
	return connectPolicy, nil
}

func connect(ac *AerospikeConfig, connectPolicy *aerospike.ClientPolicy) (*Aerospike, error) {
	var err error
	maxConnectAttempts := 1
	if ac.MaxConnectAttempts != nil && *ac.MaxConnectAttempts > 0 {
		maxConnectAttempts = *ac.MaxConnectAttempts
	}

	maxConnectRetryTimeMs := 30000
	if ac.Policies != nil && ac.Policies.Base != nil && ac.Policies.Base.TimeoutMs != nil && (ac.Policies.Base.TimeoutMs.Connect != nil || ac.Policies.Base.TimeoutMs.Login != nil) {
		maxConnectRetryTimeMs = 0
		if ac.Policies.Base.TimeoutMs.Connect != nil {
			maxConnectRetryTimeMs = *ac.Policies.Base.TimeoutMs.Connect
		}
		if ac.Policies.Base.TimeoutMs.Login != nil {
			maxConnectRetryTimeMs += *ac.Policies.Base.TimeoutMs.Login
		}
	}
	if ac.MaxConnectRetryTimeMs != nil {
		maxConnectRetryTimeMs = *ac.MaxConnectRetryTimeMs
	}

	connectRetrySleepMs := 0
	if ac.ConnectRetrySleepMs != nil {
		connectRetrySleepMs = *ac.ConnectRetrySleepMs
	}

	aero := new(Aerospike)
	attempt := 1
	attemptStart := time.Now()
	for {
		aero.Client, err = aerospike.NewClientWithPolicy(connectPolicy, *ac.Host, *ac.Port)
		if err == nil {
			break
		}
		attempt++
		if attempt > maxConnectAttempts || (maxConnectRetryTimeMs != 0 && time.Since(attemptStart).Nanoseconds()*1000 > int64(maxConnectRetryTimeMs)) {
			break
		}
		if connectRetrySleepMs > 0 {
			time.Sleep(time.Millisecond * time.Duration(connectRetrySleepMs))
		}
	}
	if err != nil {
		return nil, err
	}
	return aero, nil
}

func makePolicies(ac *AerospikeConfig, aero *Aerospike) error {
	aero.Policies.Read = aerospike.NewPolicy()
	aero.Policies.Write = aerospike.NewWritePolicy(0, 0)
	aero.Policies.Query = aerospike.NewQueryPolicy()
	aero.Policies.Scan = aerospike.NewScanPolicy()
	aero.Policies.Info = aerospike.NewInfoPolicy()
	if ac.Policies != nil {
		if ac.Policies.Read != nil {
			if ac.Policies.Read.MaxRetries != nil {
				aero.Policies.Read.MaxRetries = *ac.Policies.Read.MaxRetries
			}
			if ac.Policies.Read.TimeoutMs != nil {
				if ac.Policies.Read.TimeoutMs.SleepBetweenRetries != nil {
					aero.Policies.Read.SleepBetweenRetries = time.Millisecond * time.Duration(*ac.Policies.Read.TimeoutMs.SleepBetweenRetries)
				}
				if ac.Policies.Read.TimeoutMs.Socket != nil {
					aero.Policies.Read.SocketTimeout = time.Millisecond * time.Duration(*ac.Policies.Read.TimeoutMs.Socket)
				}
				if ac.Policies.Read.TimeoutMs.Total != nil {
					aero.Policies.Read.TotalTimeout = time.Millisecond * time.Duration(*ac.Policies.Read.TimeoutMs.Total)
				}
			}
		}
		if ac.Policies.Write != nil {
			if ac.Policies.Write.MaxRetries != nil {
				aero.Policies.Write.MaxRetries = *ac.Policies.Write.MaxRetries
			}
			if ac.Policies.Write.TimeoutMs != nil {
				if ac.Policies.Write.TimeoutMs.SleepBetweenRetries != nil {
					aero.Policies.Write.SleepBetweenRetries = time.Millisecond * time.Duration(*ac.Policies.Write.TimeoutMs.SleepBetweenRetries)
				}
				if ac.Policies.Write.TimeoutMs.Socket != nil {
					aero.Policies.Write.SocketTimeout = time.Millisecond * time.Duration(*ac.Policies.Write.TimeoutMs.Socket)
				}
				if ac.Policies.Write.TimeoutMs.Total != nil {
					aero.Policies.Write.TotalTimeout = time.Millisecond * time.Duration(*ac.Policies.Write.TimeoutMs.Total)
				}
			}
		}
		if ac.Policies.Query != nil {
			if ac.Policies.Query.MaxRetries != nil {
				aero.Policies.Query.MaxRetries = *ac.Policies.Query.MaxRetries
			}
			if ac.Policies.Query.TimeoutMs != nil {
				if ac.Policies.Query.TimeoutMs.SleepBetweenRetries != nil {
					aero.Policies.Query.SleepBetweenRetries = time.Millisecond * time.Duration(*ac.Policies.Query.TimeoutMs.SleepBetweenRetries)
				}
				if ac.Policies.Query.TimeoutMs.Socket != nil {
					aero.Policies.Query.SocketTimeout = time.Millisecond * time.Duration(*ac.Policies.Query.TimeoutMs.Socket)
				}
				if ac.Policies.Query.TimeoutMs.Total != nil {
					aero.Policies.Query.TotalTimeout = time.Millisecond * time.Duration(*ac.Policies.Query.TimeoutMs.Total)
				}
			}
		}
		if ac.Policies.Scan != nil {
			if ac.Policies.Scan.MaxRetries != nil {
				aero.Policies.Scan.MaxRetries = *ac.Policies.Scan.MaxRetries
			}
			if ac.Policies.Scan.TimeoutMs != nil {
				if ac.Policies.Scan.TimeoutMs.SleepBetweenRetries != nil {
					aero.Policies.Scan.SleepBetweenRetries = time.Millisecond * time.Duration(*ac.Policies.Scan.TimeoutMs.SleepBetweenRetries)
				}
				if ac.Policies.Scan.TimeoutMs.Socket != nil {
					aero.Policies.Scan.SocketTimeout = time.Millisecond * time.Duration(*ac.Policies.Scan.TimeoutMs.Socket)
				}
				if ac.Policies.Scan.TimeoutMs.Total != nil {
					aero.Policies.Scan.TotalTimeout = time.Millisecond * time.Duration(*ac.Policies.Scan.TimeoutMs.Total)
				}
			}
		}
		if ac.Policies.Info != nil {
			if ac.Policies.Info.TimeoutMs != nil {
				aero.Policies.Info.Timeout = time.Duration(*ac.Policies.Info.TimeoutMs) * time.Millisecond
			}
		}
	}
	return nil
}

func Connect(ac *AerospikeConfig) (*Aerospike, error) {
	if err := sanityCheck(ac); err != nil {
		return nil, fmt.Errorf("config check: %s", err)
	}
	connectPolicy, err := setConnectPolicy(ac)
	if err != nil {
		return nil, fmt.Errorf("create connect policy: %s", err)
	}

	aero, err := connect(ac, connectPolicy)
	if err != nil {
		return nil, fmt.Errorf("connect: %s", err)
	}

	err = makePolicies(ac, aero)
	if err != nil {
		return aero, fmt.Errorf("makePolicies: %s", err)
	}

	return aero, nil
}
