import axios from 'axios';

// 获取API基础URL
const getBaseUrl = () => {
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  const port = window.location.port;
  return `${protocol}//${hostname}${port ? ':' + port : ''}`;
};

// 创建axios实例
const api = axios.create({
  baseURL: `${getBaseUrl()}/api/v1`,
  timeout: 120000, // 增加超时时间到120秒，避免大文件上传误超时
  headers: {
    'Content-Type': 'application/json',
    'Cache-Control': 'no-cache', // 禁用浏览器缓存
    'Pragma': 'no-cache'
  },
  // 添加重试配置
  retry: 3, // 重试次数
  retryDelay: 1000, // 重试间隔（毫秒）
});

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data;
  },
  async error => {
    const config = error.config;

    // 如果没有设置重试配置，就返回错误
    if (!config || !config.retry) {
      // 专门处理401：不再自动刷新页面
      if (error.response?.status === 401) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('user_info');
        if (typeof window !== 'undefined' && window.$message) {
          window.$message.error('登录已过期，请重新登录');
        }
      }
      return Promise.reject(error);
    }

    // 设置用于重试计数的变量
    config.__retryCount = config.__retryCount || 0;

    // 检查是否已经达到最大重试次数
    if (config.__retryCount >= config.retry) {
      // 如果是认证错误，清除token（不刷新页面）
      if (error.response?.status === 401) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('user_info');
        if (typeof window !== 'undefined' && window.$message) {
          window.$message.error('登录已过期，请重新登录');
        }
      }
      return Promise.reject(error);
    }

    // 增加重试计数
    config.__retryCount += 1;

    // 创建延迟执行的 promise
    const backoff = new Promise(resolve => {
      setTimeout(() => {
        resolve();
      }, config.retryDelay || 1000);
    });

    // 等待延迟时间过后重试请求
    await backoff;
    
    // 添加时间戳参数，避免浏览器缓存
    const timestamp = new Date().getTime();
    if (config.params) {
      config.params._t = timestamp;
    } else {
      config.params = { _t: timestamp };
    }

    // 重新发起请求
    return api(config);
  }
);

// API接口
export const authAPI = {
  // 管理员登录
  login: (data) => api.post('/auth/login', data),
  // 管理员登出
  logout: () => api.post('/auth/logout'),
  // 获取当前用户信息
  getCurrentUser: () => api.get('/auth/me'),
  // 修改密码
  changePassword: (data) => api.put('/auth/change-password', data)
};

export const awardAPI = {
  // 查询个人获奖记录（公开接口）
  queryPersonAwards: (data) => api.post('/awards/query/person', data),
  // 查询团体获奖记录（公开接口）
  queryGroupAwards: (data) => api.post('/awards/query/group', data),
  // 上传个人获奖记录（需要认证）
  uploadPersonAward: (data) => api.post('/awards/upload/person', data),
  // 上传团体获奖记录（需要认证）
  uploadGroupAward: (data) => api.post('/awards/upload/group', data),
  // 批量上传个人获奖记录（需要认证）
  batchUploadPersonAwards: (data) => api.post('/awards/batch-upload/person', data),
  // 批量上传团体获奖记录（需要认证）
  batchUploadGroupAwards: (data) => api.post('/awards/batch-upload/group', data),
  // 修改个人获奖记录（需要认证）
  updatePersonAward: (data) => api.put('/awards/update/person', data),
  // 修改团体获奖记录（需要认证）
  updateGroupAward: (data) => api.put('/awards/update/group', data),
  // 删除获奖记录（需要认证）
  deleteAward: (id, award_type) => {
    const url = award_type ? `/awards/${id}?award_type=${award_type}` : `/awards/${id}`;
    return api.delete(url);
  },
  // 获取单个获奖记录详情（公开接口）
  getAwardDetail: (id) => api.get(`/awards/${id}`),
  // 获取所有获奖记录（管理员接口）
  getAllAwards: (timestamp) => api.get('/awards/manage/all', {
    params: {
      _t: timestamp || new Date().getTime()
    }
  }),
  // 获取所有获奖记录（轻量级版本，管理员接口）
  getAllAwardsLight: (timestamp) => api.get('/awards/manage/all-light', {
    params: {
      _t: timestamp || new Date().getTime()
    }
  }),
  // 获取团体成员及其证书信息（需要认证）
  getGroupMembersWithCertificates: (awardId) => api.get(`/awards/group/${awardId}/members`),
  // 更新成员证书（需要认证）
  updateMemberCertificate: (data) => api.put('/awards/group/member/certificate', data),
  // 批量上传证书图片（新增）
  uploadCertificate: (file, idNumber) => {
    const formData = new FormData();
    formData.append('certificate', file); // 这里必须用'certificate'，与后端一致
    formData.append('id_number', idNumber);
    return api.post('/awards/upload/certificate', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
  }
};

// 监控API
export const monitorAPI = {
  // 获取监控数据
  getMonitorData: () => api.get('/monitor/data'),
  // 获取数据库统计
  getDatabaseStats: () => api.get('/monitor/database'),
  // 获取系统健康状态
  getSystemHealth: () => api.get('/monitor/health'),
  // 获取在线用户
  getOnlineUsers: () => api.get('/monitor/users'),
  // 获取请求历史
  getRequestHistory: (limit = 50) => api.get('/monitor/requests', {
    params: { limit }
  }),
  // 重置统计
  resetStats: () => api.post('/monitor/reset'),
  // 踢出用户
  kickUser: (userId) => api.delete(`/monitor/users/${userId}`)
};

// 导出独立的查询函数，方便组件直接使用
export const queryPersonAwards = awardAPI.queryPersonAwards;
export const queryGroupAwards = awardAPI.queryGroupAwards;

export default api;