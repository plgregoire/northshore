// Copyright 2016 The NorthShore Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

func openDBBucket(bucket []byte) *bolt.DB {
	path := viper.GetString("BoltDBPath")
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return db
}

// Delete deletes key from boltdb Bucket
func Delete(bucket []byte, key []byte) error {
	db := openDBBucket(bucket)
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete(key)
	}); err != nil {
		log.Println("#DB,#Store,#Error", err)
		return err
	}
	return nil
}

// Load loads item from boltdb Bucket
func Load(bucket []byte, key []byte, v interface{}) error {
	db := openDBBucket(bucket)
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		buf := b.Get(key)
		if err := json.Unmarshal(buf, &v); err != nil {
			log.Println("#DB,#Load,#Error", err)
			return err
		}
		return nil
	})
	return err
}

// LoadBucket loads all items from boltdb Bucket
func LoadBucket(bucket []byte, buf *[]interface{}) error {
	db := openDBBucket(bucket)
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		b.ForEach(func(k, v []byte) error {
			var item map[string]interface{}
			if err := json.Unmarshal(v, &item); err != nil {
				log.Println("#DB,#LoadBucket,#Error", err)
				return err
			}
			*buf = append(*buf, item)
			return nil
		})
		return nil
	})
	return err
}

// Save stores item in boltdb Bucket as JSON
func Save(bucket []byte, key []byte, v interface{}) error {
	db := openDBBucket(bucket)
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		vEncoded, err := json.Marshal(v)
		if err != nil {
			log.Println("#DB,#Store,#Error", err)
			return err
		}
		if err := b.Put(key, vEncoded); err != nil {
			log.Println("#DB,#Store,#Error", err)
			return err
		}
		return nil
	})
	return err
}
