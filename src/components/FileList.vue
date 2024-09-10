<template>
  <div class="file-list-container">
    <h2>文件列表</h2>
    <ul v-if="files.length">
      <li v-for="file in files" :key="file.id">
        {{ file.name }}
        <button @click="downloadFile(file.id)">下载</button>
        <button @click="deleteFile(file.id)">删除</button>
      </li>
    </ul>
    <p v-else>暂无文件</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'FileList',
  data() {
    return {
      files: []
    }
  },
  mounted() {
    this.fetchFiles();
  },
  methods: {
    async fetchFiles() {
      try {
        const response = await axios.get('/api/files');
        this.files = response.data;
      } catch (error) {
        console.error('获取文件列表失败', error);
      }
    },
    async downloadFile(fileID) {
      try {
        const response = await axios.get(`/api/files/download/${fileID}`, {
          responseType: 'blob'
        });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'file'); // 你可能需要从服务器获取实际的文件名
        document.body.appendChild(link);
        link.click();
      } catch (error) {
        console.error('下载文件失败', error);
      }
    },
    async deleteFile(fileID) {
      try {
        await axios.delete(`/api/files/${fileID}`);
        this.fetchFiles(); // 重新获取文件列表
      } catch (error) {
        console.error('删除文件失败', error);
      }
    }
  }
}
</script>

<style scoped>
.file-list-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}
/* 添加更多样式... */
</style>