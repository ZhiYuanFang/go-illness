// Package ossclient 从 yaml oss 段与环境变量组装阿里云 OSS 客户端（上传直传用）。
package ossclient

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	envAccessKeyID     = "ILLNESS_OSS_ACCESS_KEY_ID"
	envAccessKeySecret = "ILLNESS_OSS_ACCESS_KEY_SECRET"
	envEndpoint        = "ILLNESS_OSS_ENDPOINT"
)

// Config 运行时 OSS 配置（密钥仅来自 env）。
type Config struct {
	Bucket          string
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	ObjectKeyPrefix string
	CdnBaseURL      string
}

// LoadConfig 读取 manifest 中 oss 段，并用 env 注入 AK/SK；ILLNESS_OSS_ENDPOINT 可覆盖 endpoint。
func LoadConfig(ctx context.Context) Config {
	cfg := Config{
		Bucket:          strings.TrimSpace(g.Cfg().MustGet(ctx, "oss.bucket").String()),
		Endpoint:        strings.TrimSpace(g.Cfg().MustGet(ctx, "oss.endpoint").String()),
		ObjectKeyPrefix: strings.TrimSpace(g.Cfg().MustGet(ctx, "oss.objectKeyPrefix").String()),
		CdnBaseURL:      strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(ctx, "oss.cdnBaseUrl").String()), "/"),
		AccessKeyID:     strings.TrimSpace(os.Getenv(envAccessKeyID)),
		AccessKeySecret: strings.TrimSpace(os.Getenv(envAccessKeySecret)),
	}
	if v := strings.TrimSpace(os.Getenv(envEndpoint)); v != "" {
		cfg.Endpoint = v
	}
	if cfg.ObjectKeyPrefix == "" {
		cfg.ObjectKeyPrefix = "illness/"
	}
	if !strings.HasSuffix(cfg.ObjectKeyPrefix, "/") {
		cfg.ObjectKeyPrefix += "/"
	}
	return cfg
}

// Validate 检查上传所需最小配置；缺密钥时返回可诊断错误。
func (c Config) Validate() error {
	if c.Bucket == "" {
		return fmt.Errorf("oss.bucket 未配置")
	}
	if c.Endpoint == "" {
		return fmt.Errorf("oss.endpoint 未配置（可用 %s 覆盖）", envEndpoint)
	}
	if c.AccessKeyID == "" || c.AccessKeySecret == "" {
		return fmt.Errorf("缺少 %s / %s，无法上传 OSS", envAccessKeyID, envAccessKeySecret)
	}
	if c.CdnBaseURL == "" {
		return fmt.Errorf("oss.cdnBaseUrl 未配置")
	}
	return nil
}

// ObjectKey 在前缀下拼 objectKey（自动去掉重复斜杠）。
func (c Config) ObjectKey(name string) string {
	name = strings.TrimPrefix(strings.TrimSpace(name), "/")
	return c.ObjectKeyPrefix + name
}

// CDNURL 将 objectKey 拼为 CDN 绝对 URL。
func (c Config) CDNURL(objectKey string) string {
	key := strings.TrimPrefix(objectKey, "/")
	return c.CdnBaseURL + "/" + key
}

// Client 持有已校验配置，供 Put / Exist 使用。
type Client struct {
	cfg Config
}

// New 校验配置并返回客户端封装。
func New(ctx context.Context) (*Client, error) {
	cfg := LoadConfig(ctx)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{cfg: cfg}, nil
}

// Config 返回当前配置副本。
func (c *Client) Config() Config {
	return c.cfg
}

func (c *Client) bucket() (*oss.Bucket, error) {
	client, err := oss.New(c.cfg.Endpoint, c.cfg.AccessKeyID, c.cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("OSS 客户端初始化失败: %w", err)
	}
	b, err := client.Bucket(c.cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("OSS Bucket 不可用: %w", err)
	}
	return b, nil
}

// Exists 判断 objectKey 是否已在桶中。
func (c *Client) Exists(objectKey string) (bool, error) {
	b, err := c.bucket()
	if err != nil {
		return false, err
	}
	return b.IsObjectExist(objectKey)
}

// PutBytes 上传字节到指定 objectKey，返回 CDN URL。
func (c *Client) PutBytes(objectKey string, data []byte, contentType string) (cdnURL string, err error) {
	b, err := c.bucket()
	if err != nil {
		return "", err
	}
	opts := []oss.Option{}
	if contentType != "" {
		opts = append(opts, oss.ContentType(contentType))
	}
	if err = b.PutObject(objectKey, bytes.NewReader(data), opts...); err != nil {
		return "", fmt.Errorf("OSS 上传失败: %w", err)
	}
	return c.cfg.CDNURL(objectKey), nil
}
