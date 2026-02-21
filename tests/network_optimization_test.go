package tests

import (
	"bytes"
	"io"
	"testing"

	"nogochain/network/compression"
	"nogochain/network/types"
)

// TestCompression 测试数据压缩功能
func TestCompression(t *testing.T) {
	// 创建测试数据 - 使用更长、更有可压缩性的数据
	var testData []byte
	for i := 0; i < 1000; i++ {
		testData = append(testData, []byte("This is a repeated test message for network compression. ")...)
	}

	// 压缩数据
	compressed, err := compression.Compress(testData)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 解压数据
	decompressed, err := compression.Decompress(compressed)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证数据是否一致
	if !bytes.Equal(testData, decompressed) {
		t.Fatalf("解压后的数据与原始数据不一致")
	}

	// 计算压缩率
	compressionRatio := float64(len(compressed)) / float64(len(testData))
	t.Logf("原始数据大小: %d 字节", len(testData))
	t.Logf("压缩后数据大小: %d 字节", len(compressed))
	t.Logf("压缩率: %.2f%%", compressionRatio*100)
	t.Logf("压缩效果: %.2f%%", (1-compressionRatio)*100)

	// 验证压缩率是否达到要求（≥30%的压缩效果，即压缩后的数据大小≤原始大小的70%）
	if compressionRatio > 0.7 {
		t.Fatalf("压缩率未达到要求，需要≥30%%的压缩效果")
	}
}

// TestMessageCreation 测试消息创建
func TestMessageCreation(t *testing.T) {
	// 创建测试数据
	testData := []byte("Test message payload")

	// 创建消息
	msg := &types.Msg{
		Code:    1,
		Size:    uint32(len(testData)),
		Payload: bytes.NewReader(testData),
	}

	// 验证消息字段
	if msg.Code != 1 {
		t.Fatalf("消息代码不正确")
	}

	if msg.Size != uint32(len(testData)) {
		t.Fatalf("消息大小不正确")
	}

	// 读取消息体
	payload, err := io.ReadAll(msg.Payload)
	if err != nil {
		t.Fatalf("读取消息体失败: %v", err)
	}

	if !bytes.Equal(testData, payload) {
		t.Fatalf("消息体与原始数据不一致")
	}
}
