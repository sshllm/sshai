package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"

	"sshai/pkg/config"
)

// Server SSH服务器结构体
type Server struct {
	config *ssh.ServerConfig
}

// NewServer 创建新的SSH服务器
func NewServer() (*Server, error) {
	cfg := config.Get()

	// 生成主机密钥
	hostKey, err := generateHostKey()
	if err != nil {
		return nil, fmt.Errorf("生成主机密钥失败: %v", err)
	}

	// SSH服务器配置
	sshConfig := &ssh.ServerConfig{
		// 设置自定义SSH Banner
		ServerVersion: "SSH-2.0-SSHAI.TOP",
	}

	// 根据配置决定认证方式
	if cfg.Auth.Password == "" {
		// 无密码认证 - 接受所有连接
		sshConfig.NoClientAuth = true
		log.Printf("SSH服务器配置：无密码认证模式")
	} else {
		// 密码认证模式
		sshConfig.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			log.Printf("密码认证尝试: user=%s", conn.User())
			if string(password) == cfg.Auth.Password {
				log.Printf("用户 %s 认证成功", conn.User())
				return nil, nil
			}
			log.Printf("用户 %s 认证失败", conn.User())
			return nil, fmt.Errorf("密码错误")
		}
		log.Printf("SSH服务器配置：密码认证模式")
	}

	// 公钥认证（可选，暂时禁用）
	sshConfig.PublicKeyCallback = func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		// 暂时不支持公钥认证
		return nil, fmt.Errorf("不支持公钥认证")
	}

	sshConfig.AddHostKey(hostKey)

	return &Server{
		config: sshConfig,
	}, nil
}

// Start 启动SSH服务器
func (s *Server) Start() error {
	cfg := config.Get()

	// 监听端口
	listener, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("监听端口失败 %s: %v", cfg.Server.Port, err)
	}
	defer listener.Close()

	log.Printf("SSH AI Server listening on port %s", cfg.Server.Port)
	log.Printf("Connect with: ssh localhost -p %s", cfg.Server.Port)

	// 接受连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 处理每个连接
		go s.handleConnection(conn)
	}
}

// handleConnection 处理SSH连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// SSH握手
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		log.Printf("SSH handshake failed: %v", err)
		return
	}
	defer sshConn.Close()

	// 获取用户名
	username := sshConn.User()
	log.Printf("New SSH connection from %s, user: %s", sshConn.RemoteAddr(), username)

	// 处理全局请求
	go ssh.DiscardRequests(reqs)

	// 处理通道
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		// 处理会话
		go HandleSession(channel, requests, username)
	}
}

// generateHostKey 生成或加载RSA主机密钥
func generateHostKey() (ssh.Signer, error) {
	cfg := config.Get()
	keyFile := cfg.Security.HostKeyFile

	// 尝试加载现有密钥
	if keyData, err := os.ReadFile(keyFile); err == nil {
		privateKey, err := ssh.ParsePrivateKey(keyData)
		if err == nil {
			log.Printf("已加载现有主机密钥")
			return privateKey, nil
		}
		log.Printf("无法解析现有密钥文件，将生成新密钥: %v", err)
	}

	// 生成新的RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// 将私钥转换为SSH签名器
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, err
	}

	// 保存私钥到文件
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})

	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		log.Printf("警告：无法保存主机密钥: %v", err)
	} else {
		log.Printf("已生成并保存新的主机密钥")
	}

	return signer, nil
}
