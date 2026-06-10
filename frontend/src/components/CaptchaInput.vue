<template>
  <div class="captcha-container">
    <a-form-item 
      :label="label" 
      name="captcha" 
      :rules="[{ required: true, message: '请输入验证码' }]"
    >
      <div class="captcha-input-group">
        <a-input 
          v-model:value="captchaValue" 
          placeholder="请输入验证码" 
          size="large"
          @change="handleChange"
          class="captcha-input"
        />
        <div class="captcha-right">
          <div class="captcha-image" @click="refreshCaptcha">
            <canvas ref="captchaCanvas" width="120" height="40"></canvas>
          </div>
          <a-button 
            type="link" 
            size="small" 
            @click="refreshCaptcha" 
            class="refresh-btn"
          >
            看不清？换一张
          </a-button>
        </div>
      </div>
    </a-form-item>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';

const props = defineProps({
  label: {
    type: String,
    default: '验证码'
  }
});

const emit = defineEmits(['update:value', 'refresh']);

const captchaCanvas = ref(null);
const captchaValue = ref('');
const captchaText = ref('');

// 生成随机验证码
const generateCaptcha = () => {
  const canvas = captchaCanvas.value;
  const ctx = canvas.getContext('2d');
  
  // 清空画布
  ctx.fillStyle = '#f0f2f5';
  ctx.fillRect(0, 0, canvas.width, canvas.height);
  
  // 生成4位随机数字
  const code = Math.random().toString().slice(2, 6);
  captchaText.value = code;
  
  // 绘制文字
  ctx.font = 'bold 24px Arial';
  ctx.fillStyle = '#1890ff';
  
  // 随机偏移每个数字
  for (let i = 0; i < code.length; i++) {
    const x = 20 + i * 25 + Math.random() * 5;
    const y = 28 + Math.random() * 5;
    const rotate = (Math.random() - 0.5) * 0.3;
    
    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(rotate);
    ctx.fillText(code[i], 0, 0);
    ctx.restore();
  }
  
  // 绘制干扰线
  ctx.strokeStyle = '#1890ff';
  ctx.lineWidth = 1;
  ctx.globalAlpha = 0.5;
  
  for (let i = 0; i < 3; i++) {
    ctx.beginPath();
    ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height);
    ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height);
    ctx.stroke();
  }
  
  // 绘制干扰点
  ctx.fillStyle = '#1890ff';
  ctx.globalAlpha = 0.5;
  
  for (let i = 0; i < 30; i++) {
    const x = Math.random() * canvas.width;
    const y = Math.random() * canvas.height;
    ctx.beginPath();
    ctx.arc(x, y, 1, 0, Math.PI * 2);
    ctx.fill();
  }
  
  ctx.globalAlpha = 1;
};

// 刷新验证码
const refreshCaptcha = () => {
  captchaValue.value = '';
  emit('update:value', '');
  generateCaptcha();
};

// 处理输入变化
const handleChange = (e) => {
  emit('update:value', e.target.value);
};

// 监听值变化
watch(captchaValue, (newVal) => {
  emit('update:value', newVal);
});

onMounted(() => {
  generateCaptcha();
});

defineExpose({
  refreshCaptcha,
  getCaptchaText: () => captchaText.value
});
</script>

<style scoped>
.captcha-container {
  width: 100%;
}

.captcha-input-group {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.captcha-input {
  flex: 1;
  min-width: 120px;
}

.captcha-right {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.captcha-image {
  cursor: pointer;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  background: #fff;
  padding: 2px;
}

.captcha-image:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.captcha-image canvas {
  display: block;
  border-radius: 2px;
}

.refresh-btn {
  padding: 0;
  height: 20px;
  line-height: 20px;
  font-size: 12px;
  color: #1890ff;
}

.refresh-btn:hover {
  color: #40a9ff;
}
</style> 