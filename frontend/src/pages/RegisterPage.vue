<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Đăng ký</h1>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>Tên người dùng</label>
          <input
            v-model="username"
            type="text"
            required
            placeholder="Nhập tên người dùng"
          />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="Nhập email của bạn"
          />
        </div>
        <div class="form-group">
          <label>Mật khẩu</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="Nhập mật khẩu"
          />
        </div>
        <div class="form-group">
          <label>Vai trò</label>
          <select v-model="role" required>
            <option value="reader">Người đọc</option>
            <option value="creator">Tác giả</option>
          </select>
        </div>
        <div v-if="authStore.error" class="error-message">
          {{ authStore.error }}
        </div>

        <button type="submit" class="btn btn-primary btn-block" :disabled="authStore.isLoading">
          {{ authStore.isLoading ? 'Đang đăng ký...' : 'Đăng ký' }}
        </button>
      </form>
      <p class="auth-link">
        Đã có tài khoản? <router-link to="/login">Đăng nhập</router-link>
      </p>
      <router-link to="/" class="back-link">← Quay lại trang chủ</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const email = ref('')
const password = ref('')
const role = ref<'reader' | 'creator'>('reader')

const handleRegister = async () => {
  try {
    await authStore.register(username.value, email.value, password.value, role.value)
    void router.push('/')
  } catch (err) {
    // Error is handled by store
    console.error('Register failed:', err)
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
}

.auth-card {
  background: white;
  padding: 2.5rem;
  border-radius: 1rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.auth-card h1 {
  text-align: center;
  color: #333;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
}

.btn-block {
  width: 100%;
}

.auth-link {
  text-align: center;
  margin-top: 1.5rem;
  color: #666;
}

.auth-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.auth-link a:hover {
  text-decoration: underline;
}

.back-link {
  display: block;
  text-align: center;
  margin-top: 1rem;
  color: #999;
  text-decoration: none;
  font-size: 0.9rem;
}

.back-link:hover {
  color: #667eea;
}

.error-message {
  background: #fee2e2;
  color: #dc2626;
  padding: 0.75rem;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
</style>
