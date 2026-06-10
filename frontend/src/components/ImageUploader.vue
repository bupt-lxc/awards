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
import { awardAPI } from '../api/index.js';

const props = defineProps({
  modelValue: String, // 证书图片的静态文件路径
  label: {
    type: String,
    default: '上传图片',
  },
  maxSizeMB: {
    type: Number,
    default: 10,
  },
});

const emit = defineEmits(['update:modelValue']);

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
    const response = await awardAPI.uploadCertificate(file, 'record');
    emit('update:modelValue', response.data.file_path);
    message.success(`成功上传图片: ${file.name}`);
  } catch (error) {
    message.error(`图片 ${file.name} 上传失败`);
    console.error('Image upload failed:', error);
  }
  return false;
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