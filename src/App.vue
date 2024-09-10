<script>
import axios from 'axios';

export default {
  name: 'App',
  created() {
    // 请求拦截器
    axios.interceptors.request.use(
      config => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
      },
      error => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    axios.interceptors.response.use(
      response => response,
      error => {
        if (error.response.status === 401) {
          // token失效，清除它并跳转到登录页面
          localStorage.removeItem('token');
          this.$router.push('/login');
        }
        return Promise.reject(error);
      }
    );
  }
}
</script>

<template>
  <router-view></router-view>
</template>