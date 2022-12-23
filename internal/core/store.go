package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

type BoltBucket string

func (bucket BoltBucket) String() string {
	return string(bucket)
}

var Buckets []BoltBucket

func NewBucket(name string) BoltBucket {
	b := BoltBucket(name)
	Buckets = append(Buckets, b)
	return b
}

func GetDB() *bolt.DB {
	return db
}

func InitStore() {
	var err error
	db, err = bolt.Open("sillyGirl.cache", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func (bucket BoltBucket) Set(key interface{}, value interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		k := fmt.Sprint(key)
		if _, ok := value.([]byte); !ok {
			v := fmt.Sprint(value)
			if v == "" {
				if err := b.Delete([]byte(k)); err != nil {
					return err
				}
			} else {
				if err := b.Put([]byte(k), []byte(v)); err != nil {
					return err
				}
			}
		} else {
			if len(value.([]byte)) == 0 {
				if err := b.Delete([]byte(k)); err != nil {
					return err
				}

			} else {
				if err := b.Put([]byte(k), value.([]byte)); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (bucket BoltBucket) Push2Array(key, value string) {
	_ = bucket.Set(key, strings.Join(append(strings.Split(bucket.GetString(key), ","), value), ","))
}

func (bucket BoltBucket) GetArray(key string) []string {
	return strings.Split(bucket.GetString(key), ",")
}
func (bucket BoltBucket) Get(kv ...interface{}) string {
	return bucket.GetString(kv)
}
func (bucket BoltBucket) GetString(kv ...interface{}) string {
	var key, value string
	for i := range kv {
		if i == 0 {
			key = fmt.Sprint(kv[0])
		} else {
			value = fmt.Sprint(kv[1])
		}
	}
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		if v := string(b.Get([]byte(key))); v != "" {
			value = v
		}
		return nil
	})
	return value
}

func (bucket BoltBucket) GetBytes(key string) []byte {
	var value []byte
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		if v := b.Get([]byte(key)); v != nil {
			value = v
		}
		return nil
	})
	return value
}

func (bucket BoltBucket) GetInt(key interface{}, vs ...int) int {
	var value int
	if len(vs) != 0 {
		value = vs[0]
	}
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := Int(string(b.Get([]byte(fmt.Sprint(key)))))
		if v != 0 {
			value = v
		}
		return nil
	})
	return value
}

func (bucket BoltBucket) GetBool(key interface{}, vs ...bool) bool {
	var value bool
	if len(vs) != 0 {
		value = vs[0]
	}
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := string(b.Get([]byte(fmt.Sprint(key))))
		if v == "true" {
			value = true
		} else if v == "false" {
			value = false
		}
		return nil
	})
	return value
}

func (bucket BoltBucket) Foreach(f func(k, v []byte) error) {
	var bs [][][]byte
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				bs = append(bs, [][]byte{k, v})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	for i := range bs {
		_ = f(bs[i][0], bs[i][1])
	}
}

var Int = func(s interface{}) int {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return i
}

func (bucket BoltBucket) Create(i interface{}) error {
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	sequence := s.FieldByName("Sequence")
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		if _, ok := id.Interface().(int); ok {
			key := id.Int()
			sq, err := b.NextSequence()
			if err != nil {
				return err
			}
			if key == 0 {
				key = int64(sq)
				id.SetInt(key)
			}
			if sequence != reflect.ValueOf(nil) {
				sequence.SetInt(int64(sq))
			}
			buf, err := json.Marshal(i)
			if err != nil {
				return err
			}
			return b.Put(itob(uint64(key)), buf)
		} else {
			key := id.String()
			sq, err := b.NextSequence()
			if err != nil {
				return err
			}
			if key == "" {
				key = fmt.Sprint(sq)
				id.SetString(key)
			}
			if sequence != reflect.ValueOf(nil) {
				sequence.SetInt(int64(sq))
			}
			buf, err := json.Marshal(i)
			if err != nil {
				return err
			}
			return b.Put([]byte(key), buf)
		}
	})
}

func itob(i uint64) []byte {
	return []byte(fmt.Sprint(i))
}

func (bucket BoltBucket) First(i interface{}) error {
	var err error
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	if v, ok := id.Interface().(int); ok {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(fmt.Sprint(v)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			return json.Unmarshal(data, i)
		})
	} else {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(id.Interface().(string)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			return json.Unmarshal(data, i)
		})
	}
	return err
}
