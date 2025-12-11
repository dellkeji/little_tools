import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Jenkins API
export const jenkinsAPI = {
  connect: (data: { url: string; username: string; password: string }) =>
    api.post('/jenkins/connect', data),
  getNodes: () => api.get('/jenkins/nodes'),
  getNodeInfo: (name: string) => api.get(`/jenkins/nodes/${name}`),
  toggleNode: (name: string, offline: boolean) =>
    api.post(`/jenkins/nodes/${name}/toggle`, { offline }),
};

// MySQL API
export const mysqlAPI = {
  connect: (data: { host: string; port: string; username: string; password: string; database: string }) =>
    api.post('/mysql/connect', data),
  executeSQL: (query: string) => api.post('/mysql/execute', { query }),
  validateSQL: (query: string) => api.post('/mysql/validate', { query }),
  getDatabases: () => api.get('/mysql/databases'),
  getTables: (database: string) => api.get('/mysql/tables', { params: { database } }),
};

// Redis API
export const redisAPI = {
  connect: (data: { host: string; port: string; password?: string; db: number }) =>
    api.post('/redis/connect', data),
  getKeys: (pattern?: string) => api.get('/redis/keys', { params: { pattern } }),
  getValue: (key: string) => api.get(`/redis/key/${key}`),
  setValue: (data: { key: string; value: string; ttl?: number }) =>
    api.post('/redis/key', data),
  deleteKey: (key: string) => api.delete(`/redis/key/${key}`),
  getInfo: () => api.get('/redis/info'),
};

export default api;
