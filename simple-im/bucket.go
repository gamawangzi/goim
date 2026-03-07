package main

import (
	"encoding/json"
	"hash/fnv"
	"sync"

	"github.com/gorilla/websocket"
)

// Bucket 单个Bucket结构
type Bucket struct {
	connections map[string]*websocket.Conn
	lock sync.RWMutex
}

// BucketManager 管理Bucket 从而扩大连接并发量 
type BucketManager struct {
	buckets []*Bucket 
	bucketCount int // bucket数组数量  
}

// NewBucket 新建Bucket实例 
func NewBucket() *Bucket {
	return &Bucket{
		connections: make(map[string]*websocket.Conn),
		lock: sync.RWMutex{},
	}
}

func NewBucketManager(bucketCount int )*BucketManager {
	bm := &BucketManager{
		buckets: make([]*Bucket, bucketCount),
		bucketCount: bucketCount,
	}
	// 初始化每个 Bucket
	for i := 0; i < bucketCount; i++ {
		bm.buckets[i] = NewBucket()
	}
	return bm
}

// Add 添加连接 
func (b *Bucket) Add(userID string,conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.connections[userID] = conn
}
// Remove 从Bucket中移除连接 
func (b *Bucket) Remove(userID string) {
	b.lock.Lock()
	defer b.lock.Unlock() 
	delete(b.connections,userID)
}
// Get 获取指定用户的连接 
func (b *Bucket) Get(userID string) (*websocket.Conn,bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	conn, ok := b.connections[userID]
	return conn, ok
}
// Count 获取Bucket的连接数 
func (b *Bucket) Count() int {
	b.lock.RLock() 
	defer b.lock.RUnlock() 
	return len(b.connections)
}


func (m *BucketManager) Add(userID string,conn *websocket.Conn) {
	bucket := m.getBucket(userID) 
	bucket.Add(userID,conn)
}

func (m *BucketManager) Remove(userID string) {
	bucket := m.getBucket(userID) 
	bucket.Remove(userID)
}

func (m *BucketManager) Get(userID string) *websocket.Conn{
	bucket := m.getBucket(userID) 
	conn, ok := bucket.Get(userID)
	if !ok {
		return nil 
	}
	return conn 
}

func (m *BucketManager) Push(userID string , msg Message) error{
	bucket := m.getBucket(userID)
	conn,ok := bucket.Get(userID)
	if !ok {
		return ErrUserNotOnline
	}
	data,err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage,data)
}

// 总连接数 
func (m *BucketManager) Count() int {
	num := 0 
	for _,val := range m.buckets {
		num += val.Count()
	}
	return num
}
func (m *BucketManager)getBucket(userID string) *Bucket{
	hash := fnv.New32a()
	hash.Write([]byte(userID))
	hashValue := hash.Sum32()
	return m.buckets[hashValue % uint32(m.bucketCount)]
}
