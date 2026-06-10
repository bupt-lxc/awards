<template>
  <a-layout class="admin-layout">
    <!-- 侧边栏 -->
    <a-layout-sider 
      v-if="userInfo" 
      v-model:collapsed="collapsed" 
      :trigger="null" 
      collapsible 
      class="sidebar"
      width="240"
    >
      <!-- Logo区域 -->
      <div class="logo-container">
        <div class="logo">
          <span class="logo-icon">🏆</span>
          <span v-if="!collapsed" class="logo-text">获奖管理系统</span>
        </div>
      </div>

      <!-- 导航菜单 -->
      <a-menu 
        v-model:selectedKeys="selectedKeys" 
        mode="inline" 
        class="sidebar-menu"
        theme="dark"
        @click="handleMenuClick"
      >
        <a-menu-item-group key="main" title="主要功能">
          <a-menu-item key="upload">
            <template #icon>
              <upload-outlined />
            </template>
            <span>上传记录</span>
          </a-menu-item>
          <a-menu-item key="statistics">
            <template #icon>
              <bar-chart-outlined />
            </template>
            <span>统计分析</span>
          </a-menu-item>
        </a-menu-item-group>
        
        <a-divider class="menu-divider" />
        
        <a-sub-menu key="management" title="数据管理">
          <template #icon>
            <database-outlined />
          </template>
          <a-menu-item key="manage">
            <template #icon>
              <table-outlined />
            </template>
            <span>记录管理</span>
          </a-menu-item>
          <a-menu-item key="export">
            <template #icon>
              <download-outlined />
            </template>
            <span>数据导出</span>
          </a-menu-item>
        </a-sub-menu>
        
        <a-divider class="menu-divider" />
        
        <a-menu-item-group key="system" title="系统设置">
          <a-menu-item key="monitor">
            <template #icon>
              <monitor-outlined />
            </template>
            <span>系统监控</span>
          </a-menu-item><a-menu-item key="settings">
            <template #icon>
              <setting-outlined />
            </template>
            <span>系统设置</span>
          </a-menu-item>
        </a-menu-item-group>
      </a-menu>

      <!-- 用户信息区域 -->
      <div class="user-section">
        <a-dropdown :trigger="['click']" placement="topLeft">
          <div class="user-info-card">
            <a-avatar class="user-avatar" :size="collapsed ? 32 : 40">
              <template #icon>
                <user-outlined />
              </template>
            </a-avatar>
            <div v-if="!collapsed" class="user-details">
              <div class="username">{{ userInfo.username }}</div>
              <a-tag color="blue" size="small" class="user-role">
                {{ userInfo.role === 'admin' ? '管理员' : '用户' }}
              </a-tag>
            </div>
          </div>
          <template #overlay>
            <a-menu @click="handleUserMenuClick">
              <a-menu-item key="profile">
                <user-outlined />
                <span>个人资料</span>
              </a-menu-item>
              <a-menu-item key="change-password">
                <lock-outlined />
                <span>修改密码</span>
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item key="logout">
                <logout-outlined />
                <span>退出登录</span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </a-layout-sider>

    <!-- 主内容区域 -->
    <a-layout :class="['main-layout', { 'collapsed': collapsed, 'no-sidebar': !userInfo }]">
      <!-- 顶部栏 -->
      <a-layout-header class="header">
        <div class="header-left">
          <a-button 
            v-if="userInfo" 
            type="text" 
            class="trigger-btn" 
            @click="collapsed = !collapsed"
          >
            <menu-unfold-outlined v-if="collapsed" />
            <menu-fold-outlined v-else />
          </a-button>
          <h1 v-if="!userInfo" class="welcome-title">
            <img src="/25551751939512_.pic_hd.jpg" alt="logo" class="logo-img" />
            中国国际青少年合唱比赛
          </h1>
          <a-breadcrumb v-else class="breadcrumb">
            <a-breadcrumb-item>
              <a href="#">首页</a>
            </a-breadcrumb-item>
            <a-breadcrumb-item>{{ getCurrentPageName() }}</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right">
          <a-button v-if="!userInfo" type="primary" @click="showLogin" class="login-btn">
            登录
          </a-button>
        </div>
      </a-layout-header>

      <!-- 内容区域 -->
      <a-layout-content class="content">
        <!-- 未登录时显示欢迎页面 -->
        <div v-if="!userInfo" class="welcome-container">
          <div v-if="!showLoginForm" class="welcome-content">
            <div class="welcome-icon">🏆</div>
            <h2 class="welcome-heading">欢迎使用获奖信息管理系统</h2>
            <p class="welcome-description">高效管理和查询获奖记录，让荣誉更有价值</p>
            <div class="feature-list">
              <div class="feature-item">
                <span class="feature-icon">📤</span>
                <span>便捷上传获奖记录</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">📊</span>
                <span>实时统计分析</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">📋</span>
                <span>智能记录管理</span>
              </div>
            </div>
          </div>
          <LoginForm v-else @login-success="onLoginSuccess" class="login-form" />
        </div>

        <!-- 已登录时显示功能页面 -->
        <div v-else class="page-container">
          <div v-if="currentPage === 'upload'" class="page-content">
            <InputForm />
          </div>
          <div v-else-if="currentPage === 'statistics'" class="page-content">
            <StatisticsView />
          </div>
          <div v-else-if="currentPage === 'manage'" class="page-content">
            <RecordManagement />
          </div>
          <div v-else-if="currentPage === 'change-password'" class="page-content">
            <ChangePassword />
          </div>
          <div v-else-if="currentPage === 'monitor'" class="page-content">
            <MonitorView />
          </div>
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import {
  UploadOutlined,
  BarChartOutlined,
  DatabaseOutlined,
  TableOutlined,
  DownloadOutlined,
  SettingOutlined,
  UserOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LockOutlined,
  MonitorOutlined
} from '@ant-design/icons-vue';
import InputForm from './components/InputForm.vue';
import LoginForm from './components/LoginForm.vue';
import StatisticsView from './components/StatisticsView.vue';
import RecordManagement from './components/RecordManagement.vue';
import ChangePassword from './components/ChangePassword.vue';
import MonitorView from './components/MonitorView.vue';

const selectedKeys = ref(['upload']);
const currentPage = ref('upload');
const showLoginForm = ref(false);
const userInfo = ref(null);
const collapsed = ref(false);

// 页面加载时检查登录状态
onMounted(() => {
  const token = localStorage.getItem('access_token');
  const savedUserInfo = localStorage.getItem('user_info');

  if (token && savedUserInfo) {
    try {
      userInfo.value = JSON.parse(savedUserInfo);
    } catch (error) {
      console.error('解析用户信息失败:', error);
      localStorage.removeItem('access_token');
      localStorage.removeItem('user_info');
    }
  }
});

// 处理菜单点击
const handleMenuClick = ({ key }) => {
  // 只有登录用户才能切换页面
  if (userInfo.value) {
    // 处理实际存在的页面
    if (['upload', 'statistics', 'manage', 'monitor'].includes(key)) {
      currentPage.value = key;
      selectedKeys.value = [key];
      showLoginForm.value = false;
    } else if (key === 'export') {
      // 数据导出功能
      message.info('数据导出功能开发中...');
    } else if (key === 'settings') {
      // 系统设置功能
      message.info('系统设置功能开发中...');
    }
  }
};

// 显示登录表单
const showLogin = () => {
  showLoginForm.value = true;
  currentPage.value = 'login';
};

// 登录成功回调
const onLoginSuccess = (user) => {
  userInfo.value = user;
  showLoginForm.value = false;
  currentPage.value = 'upload';
  selectedKeys.value = ['upload'];
};
// 处理用户菜单点击
const handleUserMenuClick = ({ key }) => {
  if (key === 'logout') {
    logout();
  } else if (key === 'change-password') {
    currentPage.value = 'change-password';
    selectedKeys.value = [];
  } else if (key === 'profile') {
    message.info('个人资料功能开发中...');
  }
};

// 退出登录
const logout = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('user_info');
  userInfo.value = null;
  message.success('已退出登录');
  currentPage.value = 'upload';
  showLoginForm.value = false;
};

// 获取当前页面名称
const getCurrentPageName = () => {
  const pageNames = {
    upload: '上传获奖记录',
    statistics: '统计信息',
    manage: '记录管理',
    monitor: '系统监控',
    login: '登录'
  };
  return pageNames[currentPage.value] || '登录';
};
</script>
<style scoped>
/* 主布局 */
.admin-layout {
  min-height: 100vh;
  background: #f8f9fa;
}

/* 侧边栏样式 */
.sidebar {
  background: #001529;
  border-right: 1px solid #e8e9ea;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 100;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.sidebar :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Logo区域 */
.logo-container {
  padding: 20px 16px;
  border-bottom: 1px solid #1f1f1f;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 28px;
  line-height: 1;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

/* 侧边栏菜单 */
.sidebar-menu {
  flex: 1;
  border: none;
  background: transparent;
  padding: 16px 0;
}

.sidebar-menu :deep(.ant-menu-item-group-title) {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 8px 24px;
}

.sidebar-menu :deep(.ant-menu-item) {
  margin: 2px 8px;
  border-radius: 6px;
  transition: all 0.3s ease;
  height: 40px;
  line-height: 40px;
}

.sidebar-menu :deep(.ant-menu-item:hover) {
  background-color: #1890ff;
  color: #ffffff;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background-color: #1890ff !important;
  color: #ffffff !important;
}

.sidebar-menu :deep(.ant-menu-submenu-title) {
  margin: 2px 8px;
  border-radius: 6px;
  height: 40px;
  line-height: 40px;
}

.menu-divider {
  background-color: #1f1f1f;
  margin: 16px 0;
}

/* 用户信息区域 */
.user-section {
  padding: 16px;
  border-top: 1px solid #1f1f1f;
  margin-top: auto;
}

.user-info-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-info-card:hover {
  background-color: #1f1f1f;
}

.user-avatar {
  background: #1890ff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #ffffff;
  line-height: 1.2;
  margin-bottom: 4px;
}

.user-role {
  font-size: 12px;
  line-height: 1.2;
}

/* 主内容区域 */
.main-layout {
  margin-left: 240px;
  min-height: 100vh;
}

.main-layout.collapsed {
  margin-left: 80px;
}

.main-layout.no-sidebar {
  margin-left: 0;
}

/* 顶部栏 */
.header {
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.breadcrumb {
  margin: 0;
}

.breadcrumb :deep(.ant-breadcrumb-link) {
  color: #6b7280;
  text-decoration: none;
  transition: color 0.2s ease;
}

.breadcrumb :deep(.ant-breadcrumb-link:hover) {
  color: #3b82f6;
}

.breadcrumb :deep(.ant-breadcrumb-separator) {
  color: #d1d5db;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.trigger-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  color: #6b7280;
  font-size: 16px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.trigger-btn:hover {
  background: #f3f4f6;
  color: #1890ff;
  transform: scale(1.05);
}

.trigger-btn:active {
  transform: scale(0.95);
}

.welcome-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.page-title {
  font-size: 18px;
  font-weight: 500;
  color: #1a1a1a;
}

.login-btn {
  height: 36px;
  border-radius: 6px;
  font-weight: 500;
}

/* 内容区域 */
.content {
  padding: 24px;
  min-height: calc(100vh - 64px);
}

/* 欢迎页面 */
.welcome-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 112px);
}

.welcome-content {
  text-align: center;
  max-width: 600px;
  padding: 48px 24px;
}

.welcome-icon {
  font-size: 72px;
  margin-bottom: 24px;
}

.welcome-heading {
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 16px;
  line-height: 1.2;
}

.welcome-description {
  font-size: 18px;
  color: #6b7280;
  margin-bottom: 48px;
  line-height: 1.5;
}

.feature-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
  margin-top: 48px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e8e9ea;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.feature-item:hover {
  border-color: #3b82f6;
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 24px;
}

.feature-item span:last-child {
  font-size: 16px;
  font-weight: 500;
  color: #1a1a1a;
}

/* 登录表单容器 */
.login-form {
  max-width: 400px;
  width: 100%;
}

/* 页面内容 */
.page-container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-content {
  background: #ffffff;
  border: 1px solid #e8e9ea;
  border-radius: 12px;
  padding: 24px;
  min-height: 600px;
}

/* 动画效果 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-content {
  animation: slideIn 0.3s ease-out;
}

.sidebar-menu .ant-menu-item,
.sidebar-menu .ant-menu-submenu {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-layout {
    margin-left: 0 !important;
  }
  
  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }
  
  .sidebar.mobile-open {
    transform: translateX(0);
  }
  
  .header {
    padding: 0 16px;
  }
  
  .content {
    padding: 16px;
  }
  
  .welcome-heading {
    font-size: 24px;
  }
  
  .welcome-description {
    font-size: 16px;
  }
  
  .feature-list {
    grid-template-columns: 1fr;
  }
  
  .logo-text {
    display: none;
  }
  
  .user-details {
    display: none;
  }
}

/* 收起状态样式调整 */
:deep(.ant-layout-sider-collapsed) .user-section {
  padding: 16px 8px;
}

:deep(.ant-layout-sider-collapsed) .user-info {
  justify-content: center;
}

:deep(.ant-layout-sider-collapsed) .logout-btn {
  width: 40px;
  padding: 0;
}

.logo-img {
  height: 64px;
  margin-right: 16px;
  vertical-align: middle;
}
</style>