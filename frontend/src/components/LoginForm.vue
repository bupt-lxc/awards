<template>
  <div class="login-container">
    <a-card title="管理员登录" style="width: 400px; margin: 0 auto;">
      <a-form :model="loginForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }" @finish="onLogin">
        <a-form-item label="用户名" name="username" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model:value="loginForm.username" placeholder="请输入用户名" />
        </a-form-item>

        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="loginForm.password" placeholder="请输入密码" />
        </a-form-item>

        <CaptchaInput
          ref="captchaRef"
          v-model:value="loginForm.captcha"
        />

        <a-form-item :wrapper-col="{ offset: 6, span: 18 }">
          <a-button type="primary" html-type="submit" :loading="loading" block>
            登录
          </a-button>
        </a-form-item>
      </a-form>
      

    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { authAPI } from '../api/index.js';
import CaptchaInput from './CaptchaInput.vue';

const emit = defineEmits(['login-success']);

const captchaRef = ref(null);
const loading = ref(false);
const loginForm = reactive({
  username: '',
  password: '',
  captcha: ''
});

const onLogin = async () => {
  try {
    // 验证码校验
    if (!loginForm.captcha) {
      message.warning('请输入验证码');
      return;
    }
    const captchaText = captchaRef.value.getCaptchaText();
    if (loginForm.captcha !== captchaText) {
      message.error('验证码错误，请重新输入');
      captchaRef.value.refreshCaptcha();
      loginForm.captcha = '';
      return;
    }

    loading.value = true;

    const response = await authAPI.login({
      username: loginForm.username,
      password: loginForm.password
    });

    // 保存token和用户信息到localStorage
    localStorage.setItem('access_token', response.access_token);
    localStorage.setItem('user_info', JSON.stringify(response.user_info));

    message.success('登录成功！');

    // 通知父组件登录成功
    emit('login-success', response.user_info);

  } catch (error) {
    console.error('登录失败:', error);

    if (error.response?.data?.detail) {
      message.error(`登录失败: ${error.response.data.detail}`);
    } else {
      message.error('登录失败，请检查用户名和密码');
    }
    // 刷新验证码
    captchaRef.value.refreshCaptcha();
    loginForm.captcha = '';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}


</style>