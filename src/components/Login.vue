<template>
  <div class="login-container">
    <h2>登录</h2>
    <form @submit.prevent="handleLogin">
      <div>
        <label for="username">用户名:</label>
        <input type="text" id="username" v-model="username" required>
      </div>
      <div>
        <label for="password">密码:</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <button type="submit">登录</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const response = await axios.post('/login', {
          username: this.username,
          password: this.password
        });
        // 处理登录成功
        console.log('登录成功', response.data);
        // 保存token
        localStorage.setItem('token', response.data.token);
        // 设置axios默认headers
        axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`;
        // 跳转到文件列表页面
        this.$router.push('/files');
      } catch (error) {
        console.error('登录失败', error);
        // 这里可以添加错误提示
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  max-width: 300px;
  margin: 0 auto;
  padding: 20px;
}
/* 添加更多样式... */
</style>