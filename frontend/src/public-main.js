import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import PublicApp from './PublicApp.vue'

const app = createApp(PublicApp)
app.use(Antd)
app.mount('#app')