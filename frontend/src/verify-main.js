import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import VerifyApp from './VerifyApp.vue'

const app = createApp(VerifyApp)
app.use(Antd)
app.mount('#app')