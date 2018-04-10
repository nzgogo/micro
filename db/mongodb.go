package db

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	"net"
	"strings"
)

type MgoDB interface {
	Connect() error
	Close()
	Session() *mgo.Session
	DB()
}

type mgodb struct {
	conn     *mgo.Session
	opts     Options
	dialInfo *mgo.DialInfo
}

func (d *mgodb) Connect() error {
	var tlsConfig *tls.Config
	if d.opts.TLS != nil {
		tlsConfig = d.opts.TLS
	} else {
		tlsConfig = &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
	}

	if d.opts.sslMgo {
		d.dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial(d.opts.Protocol, addr.String(), tlsConfig)
			return conn, err
		}
	}
	var err error
	d.conn, err = mgo.DialWithInfo(d.dialInfo)
	return err
}

func (d *mgodb) Close() {
	d.conn.Close()
}

func (d *mgodb) DB() *mgo.Session {
	return d.conn
}

func NewMongoDB(url string, opts ...Option) *mgodb {
	options := Options{
		Protocol: DefaultProtocol,
		sslMgo:   strings.Contains(url, "ssl=true"),
	}
	url = strings.Replace(url, "ssl=true", "", -1)
	dialOp, err := mgo.ParseURL(url)
	if err != nil {
		panic("Failed to parse URI: " + err.Error())
	}

	for _, o := range opts {
		o(&options)
	}

	return &mgodb{
		opts:     options,
		dialInfo: dialOp,
	}
}

// Collection method override. Soft delete
type MgoCURD interface {

}

type MgoCollect struct {
	mgo.Collection
}

//func (m *MgoCollect) FindId(){
//}

func (m *MgoCollect) Find(query interface{}) {

	m.Collection.Find(query)

}