package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	configPath := "test_config.json"
	defer os.Remove(configPath)
	// 测试创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// 测试关闭客户端
	err = client.Close()
	assert.NoError(t, err)

	cfg := client.GetConfig()
	os.RemoveAll(cfg.DataStorePath)
}

func TestClientAddrs(t *testing.T) {
	// 创建临时配置文件
	configPath := "test_config.json"
	defer os.Remove(configPath)

	// 创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)
	defer client.Close()

	// 测试获取地址
	addrs := client.Addrs()
	assert.NotEmpty(t, addrs)
}

func TestClientPeers(t *testing.T) {
	// 创建临时配置文件
	configPath := "test_config.json"
	defer os.Remove(configPath)

	// 创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)
	defer client.Close()

	// 测试获取对等点
	peers := client.Peers()
	// 初始可能没有连接的对等点，所以这里只验证返回类型
	assert.IsType(t, []peer.ID{}, peers)
}

func TestClientID(t *testing.T) {
	// 创建临时配置文件
	configPath := "test_config.json"
	defer os.Remove(configPath)

	// 创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)
	defer client.Close()

	// 测试获取客户端ID
	id := client.ID()
	assert.NotEmpty(t, id, "客户端ID不应为空")
}

func TestClientAdvertiseAndFindPeers(t *testing.T) {
	// 创建临时配置文件
	configPath := "test_config.json"
	defer os.Remove(configPath)

	// 创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)
	defer client.Close()

	// 测试广播服务
	serviceName := "test-service"
	client.Advertise(serviceName)

	// 等待广播完成
	time.Sleep(time.Second)

	// 测试查找对等点
	peerChan, err := client.FindPeers(serviceName)
	assert.NoError(t, err)
	assert.NotNil(t, peerChan)

	// 尝试从通道读取，但不阻塞测试
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case peer := <-peerChan:
		// 如果找到对等点，验证其有效性
		assert.NotEmpty(t, peer.ID)
	case <-ctx.Done():
		// 超时是可接受的，因为在测试环境中可能找不到对等点
	}
}

func TestClientGetConfig(t *testing.T) {
	// 创建临时配置文件
	configPath := "test_config.json"
	defer os.Remove(configPath)

	// 创建客户端
	client, err := NewClient(configPath)
	assert.NoError(t, err)

	// 测试获取配置
	config := client.GetConfig()
	assert.NotNil(t, config, "配置不应为空")

	// 清理数据存储目录
	client.Close()
	err = os.RemoveAll(config.DataStorePath)
	assert.NoError(t, err)
}
