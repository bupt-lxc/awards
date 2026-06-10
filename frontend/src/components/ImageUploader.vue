<template>
  <div class="image-uploader">
    <a-upload
      :before-upload="handleBeforeUpload"
      :show-upload-list="false"
      accept="image/*"
      list-type="picture-card"
      class="uploader"
    >
      <div v-if="modelValue">
        <img :src="modelValue" alt="preview" class="preview-image" />
      </div>
      <div v-else>
        <plus-outlined />
        <div class="ant-upload-text">{{ label }}</div>
      </div>
    </a-upload>
    <div class="actions" v-if="modelValue">
      <a-button size="small" type="link" @click="previewImage">预览</a-button>
      <a-button size="small" type="link" danger @click="removeImage">删除</a-button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';
import { message } from 'ant-design-vue';
import { PlusOutlined } from '@ant-design/icons-vue';

const props = defineProps({
  modelValue: String, // The base64 string
  label: {
    type: String,
    default: '上传图片',
  },
  maxSizeMB: {
    type: Number,
    default: 5,
  },
});

const emit = defineEmits(['update:modelValue']);

const convertImageToBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
};

const handleBeforeUpload = async (file) => {
  if (!file.type.startsWith('image/')) {
    message.error(`请选择图片文件 (${file.name} 不是图片格式)`);
    return false;
  }
  if (file.size > props.maxSizeMB * 1024 * 1024) {
    message.error(`图片大小不能超过${props.maxSizeMB}MB (${file.name} 过大)`);
    return false;
  }
  try {
    const base64 = await convertImageToBase64(file);
    emit('update:modelValue', base64);
    message.success(`成功选择图片: ${file.name}`);
  } catch (error) {
    message.error(`图片 ${file.name} 转换失败`);
    console.error('Base64 conversion failed:', error);
  }
  return false; // Prevent automatic upload
};

const removeImage = () => {
  emit('update:modelValue', '');
};

const previewImage = () => {
  if (!props.modelValue) return;
  const newWindow = window.open();
  newWindow.document.write(`<body style="margin:0; display:flex; justify-content:center; align-items:center; background-color:#333;"><img src="${props.modelValue}" style="max-width:100%; max-height:100%;" /></body>`);
};
</script>

<style scoped>
.image-uploader {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.uploader > :deep(.ant-upload.ant-upload-select-picture-card) {
  width: 128px;
  height: 128px;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.actions {
  margin-top: 8px;
}
</style> 