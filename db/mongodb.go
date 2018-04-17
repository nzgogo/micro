package db

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"net"
	"strings"
	"time"
)

const (
	DefaultConnTimeout = 60 * time.Second
)

type MgoDB interface {
	Connect() error
	Close()
	Session() *mgo.Session
	DB(string) *MicroDB
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

	if d.dialInfo.Timeout == 0 {
		d.dialInfo.Timeout = DefaultConnTimeout
	}

	var err error
	d.conn, err = mgo.DialWithInfo(d.dialInfo)
	return err
}

func (d *mgodb) Close() {
	d.conn.Close()
}

func (d *mgodb) Session() *mgo.Session {
	return d.conn
}

func (d *mgodb) DB(name string) *MicroDB {
	return &MicroDB{d.conn.DB(name)}
}

func NewMongoDB(url string, opts ...Option) MgoDB {
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

type MicroDB struct {
	*mgo.Database
}

func (d *MicroDB) C(name string) *MicroCollect {
	return &MicroCollect{d.Database.C(name)}
}

type MicroCollect struct {
	*mgo.Collection
}

// Count returns the total number of documents in the collection.
func (m *MicroCollect) Count() (n int, err error) {
	return m.Find(nil).Count()
}

// Find prepares a query using the provided document. A additional condition
// is added to the query -> { delete_at: { $exists: false } }.
// The document may be a map or a struct value capable of being marshalled with bson.
// The map may be a generic one using interface{} for its key and/or values, such as
// bson.M, or it may be a properly typed map.  Providing nil as the document
// is equivalent to providing an empty document such as bson.M{}.
//
// Further details of the query may be tweaked using the resulting Query value,
// and then executed to retrieve results using methods such as One, For,
// Iter, or Tail.
//
// In case the resulting document includes a field named $err or errmsg, which
// are standard ways for MongoDB to return query errors, the returned err will
// be set to a *QueryError value including the Err message and the Code.  In
// those cases, the result argument is still unmarshalled into with the
// received document so that any other custom values may be obtained if
// desired.
func (m *MicroCollect) Find(query interface{}) *mgo.Query {
	if s, ok := query.(bson.M); ok {
		return m.Collection.Find(bson.M{"$and": []bson.M{
			s,
			{"delete_at": bson.M{"$exists": false}},
		}})
	} else {
		bytes, _ := bson.Marshal(query)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		return m.Collection.Find(bson.M{"$and": []bson.M{
			origin,
			{"delete_at": bson.M{"$exists": false}},
		}})
	}
	return nil
}

// FindId is a convenience helper equivalent to:
//
//     query := MicroCollect.Find(bson.M{"_id": id,"delete_at": bson.M{"$exists":false}},)
//
// See the Find method for more details.
func (m *MicroCollect) FindId(id interface{}) *mgo.Query {
	//return  m.Collection.Find(bson.D{{Name: "_id", Value: id}})
	return m.Collection.Find(bson.M{"$and": []bson.M{
		{"_id": id},
		{"delete_at": bson.M{"$exists": false}},
	}})
}

// See details in m.Collection.Find()
func (m *MicroCollect) FindWithTrash(query interface{}) *mgo.Query {
	return m.Collection.Find(query)
}

// See details in m.Collection.FindId()
func (m *MicroCollect) FindIdWithTrash(id interface{}) *mgo.Query {
	return m.Collection.FindId(id)
}

// Remove finds a single document matching the provided selector document
// and performs a soft delete to the matched document (add a pair of
// key/value "delete_at":time.Now()).
//
// If the session is in safe mode (see SetSafe) a ErrNotFound error is
// returned if a document isn't found, or a value of type *LastError
// when some other error is detected.
func (m *MicroCollect) Remove(selector interface{}) error {
	update := bson.M{"$set": bson.M{"delete_at": time.Now()}}
	return m.Collection.Update(selector, update)
}

// RemoveId is a convenience helper equivalent to:
//
//     err := MicroCollect.Remove(bson.M{"_id": id})
//
// See the Remove method for more details.
func (m *MicroCollect) RemoveId(id interface{}) error {
	update := bson.M{"$set": bson.M{"delete_at": time.Now()}}
	return m.Collection.UpdateId(bson.D{{Name: "_id", Value: id}}, update)
}

// RemoveAll finds all documents matching the provided selector document
// and performs soft delete to the matched documents.
//
// In case the session is in safe mode (see the SetSafe method) and an
// error happens when attempting the change, the returned error will be
// of type *LastError.
func (m *MicroCollect) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	update := bson.M{"$set": bson.M{"delete_at": time.Now()}}
	return m.Collection.UpdateAll(selector, update)
}

// See details in m.Collection.Remove()
func (m *MicroCollect) ForceRemove(selector interface{}) error {
	return m.Collection.Remove(selector)
}

// See details in m.Collection.RemoveId()
func (m *MicroCollect) ForceRemoveId(id interface{}) error {
	return m.Collection.RemoveId(id)
}

// See details in m.Collection.RemoveAll()
func (m *MicroCollect) ForceRemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	return m.Collection.RemoveAll(selector)
}

// Update finds a single document matching the provided selector document
// that is not marked as deleted (without field deleted_at) and modifies
// it according to the update document.

// If the session is in safe mode (see SetSafe) a ErrNotFound error is
// returned if a document isn't found, or a value of type *LastError
// when some other error is detected.
func (m *MicroCollect) Update(selector interface{}, update interface{}) error {
	var newSelector interface{}
	if s, ok := selector.(bson.M); ok {
		newSelector = bson.M{"$and": []bson.M{s, {"delete_at": bson.M{"$exists": false}}}}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newSelector = bson.M{"$and": []bson.M{origin, {"delete_at": bson.M{"$exists": false}}}}
	}
	return m.Collection.Update(newSelector, update)
}

// Update finds a single document matching the provided selector document
// that is not marked as deleted (without field deleted_at) and partially
// modifies it according to the update document.

// If the session is in safe mode (see SetSafe) a ErrNotFound error is
// returned if a document isn't found, or a value of type *LastError
// when some other error is detected.
func (m *MicroCollect) UpdateParts(selector interface{}, update interface{}) error {
	var newSelector interface{}
	if s, ok := selector.(bson.M); ok {
		newSelector = bson.M{"$and": []bson.M{s, {"delete_at": bson.M{"$exists": false}}}}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newSelector = bson.M{"$and": []bson.M{origin, {"delete_at": bson.M{"$exists": false}}}}
	}
	var newUpdate interface{}
	if s, ok := selector.(bson.M); ok {
		newUpdate = bson.M{"set": s}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newUpdate = bson.M{"set": origin}
	}
	return m.Collection.Update(newSelector, newUpdate)
}

// See more details in m.Collection.Update
func (m *MicroCollect) UpdateWithTrash(selector interface{}, update interface{}) error {
	return m.Collection.Update(selector, update)
}

// IncrementUpdate finds a single document matching the provided selector document
// and performs a soft delete, then inserts the update document. Do not
// use Update Operators here since it's actually an insertion operation.
//
// If the session is in safe mode (see SetSafe) a ErrNotFound error is
// returned if a document isn't found, or a value of type *LastError
// when some other error is detected.
func (m *MicroCollect) IncrementUpdate(selector interface{}, update interface{}) error {
	var newSelector interface{}
	if s, ok := selector.(bson.M); ok {
		newSelector = bson.M{"$and": []bson.M{s, {"delete_at": bson.M{"$exists": false}}}}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newSelector = bson.M{"$and": []bson.M{origin, {"delete_at": bson.M{"$exists": false}}}}
	}
	if err := m.Remove(newSelector); err != nil {
		return err
	}

	var newDoc interface{}
	if s, ok := update.(bson.M); ok {
		s["_id"] = bson.NewObjectId()
		newDoc = s
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		origin["_id"] = bson.NewObjectId()
		newDoc = origin
	}

	if err := m.Insert(newDoc); err != nil {
		return err
	}
	return nil
}

// UpdateId is a convenience helper equivalent to:
//
//     err := MicroCollect.Update(bson.M{"_id": id}, update)
//
// See the Update method for more details.
func (m *MicroCollect) IncrementUpdateId(id interface{}, update interface{}) error {
	return m.IncrementUpdate(bson.D{{Name: "_id", Value: id}}, update)
}

// UpdateAll finds all documents matching the provided selector document
// and performs soft delete to them, then inserts the update document. Do
// not use Update Operators here since it's actually an insertion operation.
//
// If the session is in safe mode (see SetSafe) details of the executed
// operation are returned in info or an error of type *LastError when
// some problem is detected. It is not an error for the update to not be
// applied on any documents because the selector doesn't match.
func (m *MicroCollect) IncrementUpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	var newSelector interface{}
	if s, ok := selector.(bson.M); ok {
		newSelector = bson.M{"$and": []bson.M{s, {"delete_at": bson.M{"$exists": false}}}}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newSelector = bson.M{"$and": []bson.M{origin, {"delete_at": bson.M{"$exists": false}}}}
	}

	info, err = m.RemoveAll(newSelector)
	if err != nil {
		return
	}

	var newDoc interface{}
	if s, ok := update.(bson.M); ok {
		newDoc = s
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newDoc = origin
	}

	for i := 0; i < info.Updated; i++ {
		if err = m.Insert(newDoc); err != nil {
			return
		}
	}

	return
}

// IncreUpsert finds a single document matching the provided selector document
// and performs a soft delete, then inserts the update document.  If no
// document matching the selector is found, the update document is inserted
// in the collection.
//
// If the session is in safe mode (see SetSafe) details  of the executed
// operation are returned in info, or an error of type *LastError when
// some problem is detected.
func (m *MicroCollect) IncreUpsert(selector interface{}, update interface{}) (err error) {
	var newSelector interface{}
	if s, ok := selector.(bson.M); ok {
		newSelector = bson.M{"$and": []bson.M{s, {"delete_at": bson.M{"$exists": false}}}}
	} else {
		bytes, _ := bson.Marshal(update)
		origin := bson.M{}
		bson.Unmarshal(bytes, origin)
		newSelector = bson.M{"$and": []bson.M{origin, {"delete_at": bson.M{"$exists": false}}}}
	}

	err = m.Remove(newSelector)
	if err != nil && err != mgo.ErrNotFound {
		return
	}
	err = m.Insert(update)
	return
}

// IncreUpsertId is a convenience helper equivalent to:
//
//     info, err := MicroCollect.Upsert(bson.M{"_id": id}, update)
//
// See the Upsert method for more details.
func (m *MicroCollect) IncreUpsertId(id interface{}, update interface{}) (err error) {
	return m.IncreUpsert(bson.M{"_id": id}, update)
}
