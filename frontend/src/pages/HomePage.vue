<template>
  <div class="home-page">
    <header class="header">
      <h1>Novel AI</h1>
      <p>Nền tảng Light Novel Tương Tác</p>
      <nav class="nav">
        <router-link to="/login" class="btn">Đăng nhập</router-link>
        <router-link to="/register" class="btn btn-primary">Đăng ký</router-link>
      </nav>
    </header>

    <main class="main">
      <section class="hero">
        <h2>Khám phá thế giới truyện tương tác</h2>
        <p>Đọc, lựa chọn và thay đổi cốt truyện theo cách của bạn</p>
        <router-link to="/register" class="btn btn-large btn-primary">Bắt đầu ngay</router-link>
      </section>

      <section class="features">
        <div class="feature">
          <h3>📚 Đọc tương tác</h3>
          <p>Trải nghiệm câu chuyện với các lựa chọn ảnh hưởng đến kết cục</p>
        </div>
        <div class="feature">
          <h3>✍️ Sáng tạo dễ dàng</h3>
          <p>Công cụ trực quan để xây dựng cốt truyện rẽ nhánh</p>
        </div>
        <div class="feature">
          <h3>🎨 Thiết kế trực quan</h3>
          <p>Sơ đồ node-based giúp quản lý cấu trúc truyện phức tạp</p>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: rgba(0, 0, 0, 0.2);
}

.header h1 {
  font-size: 1.5rem;
}

.header p {
  opacity: 0.9;
}

.nav {
  display: flex;
  gap: 1rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  text-decoration: none;
  color: white;
  border: 1px solid white;
  transition: all 0.3s;
}

.btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-primary {
  background: white;
  color: #667eea;
}

.btn-primary:hover {
  background: #f0f0f0;
}

.btn-large {
  padding: 1rem 2rem;
  font-size: 1.1rem;
}

.main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 4rem 2rem;
}

.hero {
  text-align: center;
  padding: 4rem 0;
}

.hero h2 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.25rem;
  opacity: 0.9;
  margin-bottom: 2rem;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 4rem;
}

.feature {
  background: rgba(255, 255, 255, 0.1);
  padding: 2rem;
  border-radius: 1rem;
  text-align: center;
}

.feature h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.feature p {
  opacity: 0.9;
  line-height: 1.6;
}
</style>
