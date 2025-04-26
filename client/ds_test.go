package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerDataStore(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "badger-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试创建数据存储
	ds, err := NewBadgerDataStore(tempDir, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ds)
	defer ds.Close()

	// 测试数据存储的 Put 方法
	key := []byte("test-key")
	value := []byte("test-value")
	err = ds.Put(key, value)
	assert.NoError(t, err)

	// 测试数据存储的 Get 方法
	retrievedValue, err := ds.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// 测试数据存储的 ListPrefix 方法
	prefixKey := []byte("test")
	keys, err := ds.ListPrefix(prefixKey)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0])

	// 测试数据存储的 Delete 方法
	err = ds.Delete(key)
	assert.NoError(t, err)

	// 验证删除是否成功
	_, err = ds.Get(key)
	assert.Error(t, err) // 应该返回错误，因为键已被删除
}

func TestBadgerDataStoreWithEncryption(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "badger-test-encrypted")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建加密密钥 (32字节)
	encKey := make([]byte, 32)
	for i := range encKey {
		encKey[i] = byte(i)
	}

	// 测试创建加密数据存储
	ds, err := NewBadgerDataStore(tempDir, encKey)
	assert.NoError(t, err)
	assert.NotNil(t, ds)

	// 测试基本操作
	key := []byte("encrypted-key")
	value := []byte("encrypted-value")

	err = ds.Put(key, value)
	assert.NoError(t, err)

	retrievedValue, err := ds.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	err = ds.Close()
	assert.NoError(t, err)

	// 尝试用错误的密钥打开数据库
	wrongKey := make([]byte, 32)
	_, err = NewBadgerDataStore(tempDir, wrongKey)
	assert.Error(t, err) // 应该返回错误，因为密钥不匹配
}
