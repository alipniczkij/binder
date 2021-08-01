package bbolt

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"time"
)

type Mapper struct {
	db *bbolt.DB
}

func New(path string) *Mapper {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Can't open bbolt storage with path: %s", path)
	}
	return &Mapper{
		db: db,
	}
}

func (m *Mapper) Get(key, bucket string) ([]string, bool) {
	k := m.bytes(key)
	bucketName := m.bytes(bucket)
	var foundValue []byte

	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("can't find bucket %s", bucket)
		}
		foundValue = b.Get(k)
		if foundValue == nil {
			return fmt.Errorf("no value for this key %s", key)
		}
		return nil
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, false
	}
	res, err := m.decodeValue(foundValue)
	if err != nil {
		log.Printf("Can't decode founded value %s", string(foundValue))
		return nil, false
	}
	return res, true
}

func (m *Mapper) Store(key, value, bucket string) error {
	bucketName := m.bytes(bucket)
	err := m.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucket(bucketName)
		if err != nil {
			if errors.Is(err, bbolt.ErrBucketExists) {
				b = tx.Bucket(bucketName)
			} else {
				return fmt.Errorf("get bucket fail: %s. key: %s. Value: %s", err.Error(), key, value)
			}
		}
		existed, found := m.Get(key, bucket)
		var data []byte
		if already := contains(existed, value); already {
			return nil
		}
		if found {
			updated := append(existed, value)
			data, err = m.encodeValue(updated)
		} else {
			data, err = m.encodeValue([]string{value})
		}
		if err != nil {
			return fmt.Errorf("error while encoding %s: %s", value, err.Error())
		}
		err = b.Put(m.bytes(key), data)
		if err != nil {
			return fmt.Errorf("put fail: %s. key: %s. Value: %s", err.Error(), key, value)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Mapper) Delete(key, value, bucket string) error {
	k := m.bytes(key)
	bucketName := m.bytes(bucket)
	err := m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}
		existed, found := m.Get(key, bucket)
		if !found {
			return fmt.Errorf("not found")
		}
		updated := m.deleteValue(existed, value)
		data, err := m.encodeValue(updated)
		if err != nil {
			return fmt.Errorf("error while encoding %s: %s", value, err.Error())
		}
		err = b.Put(k, data)
		if err != nil {
			return fmt.Errorf("put fail: %s. key: %s. Value: %s", err.Error(), key, value)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Mapper) Close() error {
	return m.db.Close()
}

func (m *Mapper) deleteValue(s []string, value string) []string {
	updated := s
	for i := range s {
		if s[i] == value {
			updated = remove(s, i)
		}
	}
	return updated
}

func (m *Mapper) bytes(s string) []byte {
	return []byte(s)
}

func (m *Mapper) encodeValue(v []string) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func (m *Mapper) decodeValue(v []byte) ([]string, error) {
	buf := &bytes.Buffer{}
	res := make([]string, 0)
	buf.Write(v)
	err := gob.NewDecoder(buf).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func contains(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
