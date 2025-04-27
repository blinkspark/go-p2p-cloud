package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"
)

// 测试 DefaultConfig 函数
func TestDefaultConfig(t *testing.T) {
	cfg, err := DefaultConfig()
	assert.NoError(t, err, "DefaultConfig 不应返回错误")
	assert.NotNil(t, cfg, "返回的配置不应为 nil")
	assert.NotEmpty(t, cfg.PrivateKey, "私钥不应为空")
	assert.Equal(t, 12233, cfg.Port, "默认端口应为 12233")

	// 验证私钥是否有效
	_, err = crypto.UnmarshalPrivateKey(cfg.PrivateKey)
	assert.NoError(t, err, "私钥应能被成功解析")
}

// 测试 LoadConfig 函数 - 配置文件不存在的情况
func TestLoadConfig_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 确保文件不存在
	_, err := os.Stat(configPath)
	assert.True(t, os.IsNotExist(err), "配置文件初始时不应存在")

	// 加载配置（应创建默认配置）
	cfg, err := LoadConfig(configPath)
	assert.NoError(t, err, "LoadConfig 在创建新文件时不应返回错误")
	assert.NotNil(t, cfg, "返回的配置不应为 nil")

	// 验证返回的是默认配置
	defaultCfg, _ := DefaultConfig()
	// 比较端口，私钥是随机生成的不比较
	assert.Equal(t, defaultCfg.Port, cfg.Port, "加载的配置端口应与默认配置相同")
	assert.NotEmpty(t, cfg.PrivateKey, "加载的配置私钥不应为空")


	// 验证文件是否已创建并包含内容
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "配置文件应已被创建")

	fileBytes, err := os.ReadFile(configPath)
	assert.NoError(t, err, "读取创建的配置文件时不应出错")
	assert.NotEmpty(t, fileBytes, "创建的配置文件不应为空")

	// 验证文件内容是否是有效的 JSON 并且可以解析回 HostConfig
	var loadedFromFile HostConfig
	err = json.Unmarshal(fileBytes, &loadedFromFile)
	assert.NoError(t, err, "配置文件内容应为有效的 JSON")
	assert.Equal(t, cfg.Port, loadedFromFile.Port, "文件中的端口应与返回的配置匹配")
	assert.Equal(t, cfg.PrivateKey, loadedFromFile.PrivateKey, "文件中的私钥应与返回的配置匹配")
}

// 测试 LoadConfig 函数 - 配置文件已存在的情况
func TestLoadConfig_ExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建一个自定义的配置文件内容
	privKey, _, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	privKeyBytes, _ := crypto.MarshalPrivateKey(privKey)
	customCfg := &HostConfig{
		PrivateKey: privKeyBytes,
		Port:       54321,
	}
	configBytes, _ := json.MarshalIndent(customCfg, "", "  ")

	// 写入自定义配置文件
	err := os.WriteFile(configPath, configBytes, 0644)
	assert.NoError(t, err, "写入自定义配置文件时不应出错")

	// 加载配置
	loadedCfg, err := LoadConfig(configPath)
	assert.NoError(t, err, "LoadConfig 加载现有文件时不应返回错误")
	assert.NotNil(t, loadedCfg, "返回的配置不应为 nil")

	// 验证加载的配置是否与文件内容匹配
	assert.Equal(t, customCfg.Port, loadedCfg.Port, "加载的端口应与文件中的端口匹配")
	assert.Equal(t, customCfg.PrivateKey, loadedCfg.PrivateKey, "加载的私钥应与文件中的私钥匹配")
}

// 测试 LoadConfig 函数 - 配置文件格式无效的情况
func TestLoadConfig_InvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 写入无效的 JSON 内容
	invalidJSON := []byte("{invalid json")
	err := os.WriteFile(configPath, invalidJSON, 0644)
	assert.NoError(t, err, "写入无效配置文件时不应出错")

	// 加载配置
	_, err = LoadConfig(configPath)
	assert.Error(t, err, "LoadConfig 加载无效文件时应返回错误")
	assert.ErrorContains(t, err, "invalid character", "错误信息应指示 JSON 无效")
}


// 测试 BuildOptionFromConfig 函数
func TestBuildOptionFromConfig(t *testing.T) {
	// 使用默认配置进行测试
	cfg, err := DefaultConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	options, err := BuildOptionFromConfig(cfg)
	assert.NoError(t, err, "BuildOptionFromConfig 不应返回错误")
	assert.NotEmpty(t, options, "返回的 libp2p 选项列表不应为空")

	// 可以添加更具体的检查，例如检查是否包含 Identity 和 ListenAddrStrings 选项
	// 但这会使测试更依赖于 libp2p 的内部实现细节，所以保持简单可能更好

	// 测试无效私钥的情况
	invalidCfg := &HostConfig{
		PrivateKey: []byte("invalid key"),
		Port:       12345,
	}
	_, err = BuildOptionFromConfig(invalidCfg)
	assert.Error(t, err, "BuildOptionFromConfig 使用无效私钥时应返回错误")
}