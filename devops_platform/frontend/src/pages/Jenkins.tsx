import React, { useState } from 'react';
import { Card, Form, Input, Button, Table, message, Space, Tag } from 'antd';
import { jenkinsAPI } from '../services/api';

interface JenkinsNode {
  name: string;
  offline: boolean;
}

const JenkinsPage: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);
  const [nodes, setNodes] = useState<JenkinsNode[]>([]);

  const handleConnect = async (values: any) => {
    setLoading(true);
    try {
      await jenkinsAPI.connect(values);
      message.success('连接成功');
      setConnected(true);
      loadNodes();
    } catch (error: any) {
      message.error(error.response?.data?.error || '连接失败');
    } finally {
      setLoading(false);
    }
  };

  const loadNodes = async () => {
    try {
      const response = await jenkinsAPI.getNodes();
      setNodes(response.data.nodes || []);
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取节点列表失败');
    }
  };

  const handleToggleNode = async (name: string, offline: boolean) => {
    try {
      await jenkinsAPI.toggleNode(name, !offline);
      message.success('操作成功');
      loadNodes();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: '节点名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '状态',
      dataIndex: 'offline',
      key: 'offline',
      render: (offline: boolean) => (
        <Tag color={offline ? 'red' : 'green'}>
          {offline ? '离线' : '在线'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: JenkinsNode) => (
        <Button
          size="small"
          onClick={() => handleToggleNode(record.name, record.offline)}
        >
          {record.offline ? '启用' : '禁用'}
        </Button>
      ),
    },
  ];

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card title="Jenkins连接配置">
        <Form form={form} onFinish={handleConnect} layout="vertical">
          <Form.Item label="Jenkins URL" name="url" rules={[{ required: true }]}>
            <Input placeholder="http://jenkins.example.com" />
          </Form.Item>
          <Form.Item label="用户名" name="username" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item label="密码" name="password" rules={[{ required: true }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              连接
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {connected && (
        <Card title="Jenkins节点列表" extra={<Button onClick={loadNodes}>刷新</Button>}>
          <Table columns={columns} dataSource={nodes} rowKey="name" />
        </Card>
      )}
    </Space>
  );
};

export default JenkinsPage;
