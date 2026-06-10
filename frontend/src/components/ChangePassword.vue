<template>
  <div class="change-password-container">
    <a-card title="修改密码" class="change-password-card">
      <a-form
        ref="formRef"
        :model="form"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item label="当前密码" name="oldPassword">
          <a-input-password
            v-model:value="form.oldPassword"
            placeholder="请输入当前密码"
            size="large"
          >
            <template #prefix>
              <lock-outlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item label="新密码" name="newPassword">
          <a-input-password
            v-model:value="form.newPassword"
            placeholder="请输入新密码"
            size="large"
          >
            <template #prefix>
              <lock-outlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item label="确认新密码" name="confirmPassword">
          <a-input-password
            v-model:value="form.confirmPassword"
            placeholder="请再次输入新密码"
            size="large"
          >
            <template #prefix>
              <lock-outlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button
              type="primary"
              html-type="submit"
              :loading="loading"
              size="large"
            >
              <template #icon>
                <check-outlined />
              </template>
              确认修改
            </a-button>
            <a-button
              @click="resetForm"
              size="large"
            >
              <template #icon>
                <redo-outlined />
              </template>
              重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { LockOutlined, CheckOutlined, RedoOutlined } from '@ant-design/icons-vue'
import { authAPI } from '../api/index.js'

const formRef = ref()
const loading = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const rules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value) => {
        if (value && value !== form.newPassword) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      },
      trigger: 'blur'
    }
  ]
}

const handleSubmit = async () => {
  try {
    loading.value = true
    
    const data = {
      old_password: form.oldPassword,
      new_password: form.newPassword
    }
    
    console.log('发送修改密码请求:', data)
    const response = await authAPI.changePassword(data)
    console.log('修改密码响应:', response)
    
    message.success('密码修改成功！')
    resetForm()
    
  } catch (error) {
    console.error('修改密码失败:', error)
    console.error('错误详情:', {
      status: error.response?.status,
      data: error.response?.data,
      message: error.message
    })
    const errorMessage = error.response?.data?.error || '修改密码失败，请重试'
    message.error(errorMessage)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formRef.value?.resetFields()
  form.oldPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
}
</script>

<style scoped>
.change-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
  padding: 20px;
}

.change-password-card {
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.change-password-card :deep(.ant-card-head) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px 8px 0 0;
}

.change-password-card :deep(.ant-card-head-title) {
  color: white;
  font-size: 18px;
  font-weight: 600;
  text-align: center;
}

.change-password-card :deep(.ant-card-body) {
  padding: 32px;
}

:deep(.ant-form-item-label > label) {
  font-weight: 500;
  color: #333;
}

:deep(.ant-input-affix-wrapper) {
  border-radius: 6px;
}

:deep(.ant-btn) {
  border-radius: 6px;
  font-weight: 500;
}

:deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

:deep(.ant-btn-primary:hover) {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}
</style>