package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

// 导出 gzip 常量
const (
	// HuffmanOnly 仅使用哈夫曼编码
	HuffmanOnly = gzip.HuffmanOnly
	// BestSpeed 最快压缩速度
	BestSpeed = gzip.BestSpeed
	// BestCompression 最佳压缩率
	BestCompression = gzip.BestCompression
	// DefaultCompression 默认压缩率
	DefaultCompression = gzip.DefaultCompression
)

// Compressor 压缩器接口
 type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

// GzipCompressor Gzip压缩器
 type GzipCompressor struct {
	level int
}

// NewGzipCompressor 创建Gzip压缩器
func NewGzipCompressor(level int) *GzipCompressor {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	return &GzipCompressor{
		level: level,
	}
}

// Compress 压缩数据
func (c *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	// 创建带有指定压缩级别的 gzip.Writer
	writer, err := gzip.NewWriterLevel(&buf, c.level)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decompress 解压缩数据
func (c *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

// CompressObject 压缩对象
func CompressObject(obj interface{}, compressor Compressor) ([]byte, error) {
	// 序列化对象
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// 压缩数据
	return compressor.Compress(data)
}

// DecompressObject 解压缩对象
func DecompressObject(data []byte, obj interface{}, compressor Compressor) error {
	// 解压缩数据
	decompressed, err := compressor.Decompress(data)
	if err != nil {
		return err
	}

	// 反序列化对象
	return json.Unmarshal(decompressed, obj)
}

// CalculateCompressionRatio 计算压缩比率
func CalculateCompressionRatio(original []byte, compressed []byte) float64 {
	if len(original) == 0 {
		return 0
	}

	return float64(len(compressed)) / float64(len(original))
}

// CalculateSpaceSavings 计算空间节省百分比
func CalculateSpaceSavings(original []byte, compressed []byte) float64 {
	if len(original) == 0 {
		return 0
	}

	savings := float64(len(original)-len(compressed)) / float64(len(original)) * 100
	return savings
}
