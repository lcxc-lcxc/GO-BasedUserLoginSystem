/**
 @author: 15973
 @date: 2022/07/16
 @note:
**/
package setting

import "time"

type DBSetting struct {
	DBType       string
	Username     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxOpenConns int
	MaxIdleConns int
}

type HttpServerSetting struct {
	Host         string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type CacheSetting struct {
	Host    string
	DBIndex int
}

type RpcServerSetting struct {
	Host         string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type RpcClientSetting struct {
	RPCHost string
	ConnNum int
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	return s.vp.UnmarshalKey(k, v)
}
