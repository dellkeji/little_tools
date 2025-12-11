import React, { useState } from 'react';
import { Card, Form, Input, Button, message, Space, Table } from 'antd';
import { mysqlAPI } from '../services/api';

const { TextArea } = Input;

const MySQLPage: React.FC = () => {
  const [form] = Form.useForm();
  const [sqlForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);
  const [result, setResult] = useState<any>(null);

  const handleConnect = async (values: any) => {
    setLoading(true);
    try {
      await mysqlAPI.connect(values);
      message.success('连接成功');
      setConnected(true);
    } catch (error: any) {
      message.error(error.response?.data?.error || '连接失败');
    } finally {
      setLoading(false);
    }
  };

  const handleValidateSQL = async () => {
    const query = sqlForm.getFieldValue('query');
    if (!query) {
      message.warning('请输入SQL语句');
      return;
    }

    try {
      const response = await mysqlAPI.validateSQL(query);
      if (response.data.valid) {
        message.success('SQL语句合法');
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '验证失败');
    }
  };

  const handleExecuteSQL = async (values: any) => {
    setLoading(true);
    try {
      const response = await mysqlAPI.executeSQL(values.query);
      setResult(response.data.result);
      message.success('执行成功');
    } catch (error: any) {
      message.error(error.response?.data?.error || '执行失败');
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  const renderResult = () => {
    if (!result) return null;

    if (Array.isArray(result)) {
      if (result.length === 0) {
        return <div>查询结果为空</div>;
      }

      const columns = Object.keys(result[0]).map(key => ({
        title: key,
        dataIndex: key,
        key,
      }));

      return <Table columns={columns} dataSource={result} rowKey={(_, index) => index!.toString()} />;
    }

    return <pre>{JSON.stringify(result, null, 2)}</pre>;
  };

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card title="MySQL连接配置">
        <Form form={form} onFinish={handleConnect} layout="vertical">
          <Form.Item label="主机" name="host" rules={[{ required: true }]}>
            <Input placeholder="localhost" />
          </Form.Item>
          <Form.Item label="端口" name="port" rules={[{ required: true }]} initialValue="3306">
            <Input />
          </Form.Item>
          <Form.Item label="用户名" name="username" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item label="密码" name="password">
            <Input.Password />
          </Form.Item>
          <Form.Item label="数据库" name="database" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              连接
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {connected && (
        <Card title="SQL执行">
          <Form form={sqlForm} onFinish={handleExecuteSQL} layout="vertical">
            <Form.Item label="SQL语句" name="query" rules={[{ required: true }]}>
              <TextArea rows={6} placeholder="输入SQL语句..." />
            </Form.Item>
            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" loading={loading}>
                  执行
                </Button>
                <Button onClick={handleValidateSQL}>
                  验证SQL
                </Button>
              </Space>
            </Form.Item>
          </Form>

          {result && (
            <Card title="执行结果" style={{ marginTop: 16 }}>
              {renderResult()}
            </Card>
          )}
        </Card>
      )}
    </Space>
  );
};

export default MySQLPage;
