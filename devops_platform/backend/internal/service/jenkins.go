package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
)

type JenkinsService struct {
	client *gojenkins.Jenkins
	mu     sync.RWMutex
}

var jenkinsService *JenkinsService

func GetJenkinsService() *JenkinsService {
	if jenkinsService == nil {
		jenkinsService = &JenkinsService{}
	}
	return jenkinsService
}

func (s *JenkinsService) Connect(url, username, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	jenkins := gojenkins.CreateJenkins(nil, url, username, password)
	_, err := jenkins.Init(ctx)
	if err != nil {
		return fmt.Errorf("连接Jenkins失败: %w", err)
	}

	s.client = jenkins
	return nil
}

func (s *JenkinsService) GetNodes(ctx context.Context) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return nil, fmt.Errorf("Jenkins未连接")
	}

	nodes, err := s.client.GetAllNodes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		offline, _ := node.IsOffline(ctx)
		result = append(result, map[string]interface{}{
			"name":    node.GetName(),
			"offline": offline,
		})
	}

	return result, nil
}

func (s *JenkinsService) GetNodeInfo(ctx context.Context, nodeName string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return nil, fmt.Errorf("Jenkins未连接")
	}

	node, err := s.client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	offline, _ := node.IsOffline(ctx)
	info := map[string]interface{}{
		"name":    node.GetName(),
		"offline": offline,
	}

	return info, nil
}

func (s *JenkinsService) ToggleNode(ctx context.Context, nodeName string, offline bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return fmt.Errorf("Jenkins未连接")
	}

	node, err := s.client.GetNode(ctx, nodeName)
	if err != nil {
		return err
	}

	if offline {
		_, err = node.Disable(ctx, "通过DevOps平台禁用")
	} else {
		_, err = node.Enable(ctx)
	}

	return err
}
