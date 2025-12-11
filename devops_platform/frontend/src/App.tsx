import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import { DatabaseOutlined, CloudServerOutlined, ApiOutlined } from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import JenkinsPage from './pages/Jenkins';
import MySQLPage from './pages/MySQL';
import RedisPage from './pages/Redis';
import './App.css';

const { Header, Content, Sider } = Layout;

const App: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    {
      key: '/jenkins',
      icon: <ApiOutlined />,
      label: 'Jenkins管理',
    },
    {
      key: '/mysql',
      icon: <DatabaseOutlined />,
      label: 'MySQL管理',
    },
    {
      key: '/redis',
      icon: <CloudServerOutlined />,
      label: 'Redis管理',
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <div style={{ color: 'white', fontSize: '20px', fontWeight: 'bold' }}>
          DevOps管理平台
        </div>
      </Header>
      <Layout>
        <Sider width={200} style={{ background: '#fff' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content
            style={{
              padding: 24,
              margin: 0,
              minHeight: 280,
              background: '#fff',
            }}
          >
            <Routes>
              <Route path="/" element={<Navigate to="/jenkins" replace />} />
              <Route path="/jenkins" element={<JenkinsPage />} />
              <Route path="/mysql" element={<MySQLPage />} />
              <Route path="/redis" element={<RedisPage />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default App;
