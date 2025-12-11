import React, { useState } from 'react';
import { Card, Form, Input, Button, message, Space, List, Modal, InputNumber } from 'antd';
import { redisAPI } from '../services/api';

const RedisPage: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);
  const [keys, setKeys] = useState<string[]>([]);
  const [selectedKey, setSelectedKey] = useState<any>(null);
  const [modalVisible, setModalVisible] = useState(false);

  const handleConnect = async (values: any) => {
    setLoading(true);
    try {
      await redisAPI.connect(values);
      message.success('连接成功');
      setConnected(true);
      loadKeys();
    } catch (error: any) {
      message.error(error.response?.data?.error || '连接失败');
    } finally {
      setLoading(false);
    }
  };

  const loadKeys = async (pattern?: string) => {
    try {
      const response = await redisAPI.getKeys(pattern);
      setKeys(response.data.keys || []);
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取键列表失败');
    }
  };

  const handleViewKey = async (key: string) => {
    try {
      const response = await redisAPI.getValue(key);
      setSelectedKey(response.data);
      setModalVisible(true);
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取值失败');
    }
  };

  const handleDeleteKey = async (key: string) => {
    try {
      await redisAPI.deleteKey(key);
      message.success('删除成功');
      loadKeys();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card title="Redis连接配置">
        <Form form={form} onFinish={handleConnect} layout="vertical">
          <Form.Item label="主机" name="host" rules={[{ required: true }]}>
            <Input placeholder="localhost" />
          </Form.Item>
          <Form.Item label="端口" name="port" rules={[{ required: true }]} initialValue="6379">
            <Input />
          </Form.Item>
          <Form.Item label="密码" name="password">
            <Input.Password />
          </Form.Item>
          <Form.Item label="数据库" name="db" initialValue={0}>
            <InputNumber min={0} max={15} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              连接
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {connected && (
        <Card title="Redis键列表" extra={<Button onClick={() => loadKeys()}>刷新</Button>}>
          <List
            dataSource={keys}
            renderItem={(key) => (
              <List.Item
                actions={[
                  <Button size="small" onClick={() => handleViewKey(key)}>查看</Button>,
                  <Button size="small" danger onClick={() => handleDeleteKey(key)}>删除</Button>,
                ]}
              >
                {key}
              </List.Item>
            )}
          />
        </Card>
      )}

      <Modal
        title="键详情"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        {selectedKey && (
          <Space direction="vertical" style={{ width: '100%' }}>
            <div><strong>键:</strong> {selectedKey.key}</div>
            <div><strong>类型:</strong> {selectedKey.type}</div>
            <div><strong>TTL:</strong> {selectedKey.ttl}秒</div>
            <div><strong>值:</strong></div>
            <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
              {JSON.stringify(selectedKey.value, null, 2)}
            </pre>
          </Space>
        )}
      </Modal>
    </Space>
  );
};

export default RedisPage;
