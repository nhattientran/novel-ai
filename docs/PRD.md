# Tài Liệu Yêu Cầu Sản Phẩm (PRD)
**Tên dự án:** Nền tảng Light Novel Tương Tác
**Phiên bản:** 1.3 (AI-Ready Tech Stack & Kiến trúc GraphRAG)
**Ngày cập nhật:** 24/02/2026

---

## 1. Tổng Quan Sản Phẩm (Product Overview)
**Mục tiêu:** Xây dựng web application sáng tạo và trải nghiệm light novel có cốt truyện rẽ nhánh. Hệ thống quản lý điểm quan hệ, hồ sơ tính cách nhân vật và được thiết kế kiến trúc dữ liệu sẵn sàng cho AI sinh tạo (Generative AI) tích hợp sâu vào cốt truyện.
**Giải pháp:** Sử dụng giao diện bản đồ Node (Node-based Map) cho tác giả. Vận hành logic rẽ nhánh, điểm hảo cảm và phân loại tính cách thông qua cơ sở dữ liệu đồ thị Neo4j. Chuẩn bị sẵn nền tảng Vector Search để LLM (Mô hình ngôn ngữ lớn) có thể truy xuất ngữ cảnh chính xác.

---

## 2. Ngăn Xếp Công Nghệ (Tech Stack)



### 2.1. Frontend (Giao diện người dùng)
* **Core Framework:** Vue.js 3 (Composition API). Ưu tiên tốc độ phản hồi và hệ sinh thái phong phú.
* **State Management:** Pinia (Lưu trữ trạng thái rẽ nhánh, túi đồ, điểm hảo cảm của người đọc theo thời gian thực).
* **Story Mapping:** Vue Flow (`@vue-flow/core`). Xây dựng giao diện kéo-thả trực quan cho Creator Mode.
* **Styling:** TailwindCSS V4 (Khuyến nghị để build UI nhanh và đồng nhất).

### 2.2. Backend (Xử lý Logic & API)
* **Ngôn ngữ:** Golang.
* **Web Framework:** Fiber (Xử lý Concurrency xuất sắc, đáp ứng hàng ngàn lượt chọn rẽ nhánh cùng lúc).
* **Database Driver:** `neo4j-go-driver` (Kết nối chính thức và bảo mật với Neo4j).
* **AI Orchestration (Phase 2):** LangChainGo (`github.com/tmc/langchaingo`) hoặc tự viết Service gọi trực tiếp Google Gemini API.

### 2.3. Cơ Sở Dữ Liệu (Database)
* **Graph Database:** Neo4j (Phiên bản 5.x trở lên).
* **Tính năng bắt buộc:**
  * **Cypher Query Language:** Truy vấn mối quan hệ rẽ nhánh siêu tốc.
  * **Neo4j Vector Index:** (Chuẩn bị cho Phase 2) Lưu trữ Vector Embeddings của nội dung truyện và tính cách nhân vật để hỗ trợ tính năng tìm kiếm ngữ nghĩa (Semantic Search) cho AI.

---

## 3. Tính Năng Cốt Lõi (MVP Scope)

### 3.1. Creator Mode (Dành cho Tác giả)
| Tính năng | Mô tả chi tiết |
| :--- | :--- |
| **Story Management** | Tạo, sửa, xóa thông tin bộ truyện (Tên, Ảnh bìa, Tóm tắt). |
| **Character & Trait Setup** | Tạo nhân vật (NPCs), viết tiểu sử (Description) và gắn nhãn Tính cách (Traits). Việc viết mô tả chi tiết là tiền đề để làm Vector Embeddings cho AI sau này. |
| **Visual Story Map** | Giao diện kéo thả hiển thị toàn bộ cấu trúc rẽ nhánh (Nodes & Edges). |
| **Scene Editor** | Soạn thảo nội dung cho từng Node (Cảnh). Cài đặt cảnh Bắt đầu / Kết thúc. |
| **Choice & Effect Creator** | Tạo các liên kết rẽ nhánh (`LEADS_TO`). Thiết lập "Effect" (VD: `{"target": "char_A", "score": +5}`). |

### 3.2. Reader Mode (Dành cho Người đọc)
| Tính năng | Mô tả chi tiết |
| :--- | :--- |
| **Discovery by Traits** | Tìm kiếm truyện dựa trên tính cách nhân vật (VD: Lọc truyện có nhân vật "Tsundere"). |
| **Playthrough Management** | Tạo và quản lý các lượt chơi (Save files) độc lập. |
| **Interactive Reading** | Đọc nội dung Scene và chọn hành động tiếp theo. |
| **Affinity Tracking** | Tự động tính toán điểm quan hệ với các nhân vật dựa trên các quyết định đã chọn. |
| **History Tracking** | **(Quan trọng cho AI)** Hệ thống lưu lại toàn bộ các Scene mà người chơi đã đi qua theo thứ tự để làm "Trí nhớ" (Context) cho AI sau này. |

---

## 4. Mô Hình Dữ Liệu (Graph Database Schema - Neo4j)

**Các Nodes (Thực thể chính):**
* `(:User {id, username})`
* `(:Story {id, title, status})`
* `(:Scene {id, content, is_start, content_embedding: [Float]})` *(Thêm mảng Float chuẩn bị cho Vector Search)*
* `(:Character {id, name, description, desc_embedding: [Float]})`
* `(:Trait {id, name})`
* `(:Playthrough {id, last_played_at})`

**Các Relationships (Mối quan hệ):**
* `(Story)-[:HAS_SCENE]->(Scene)`
* `(Story)-[:HAS_CHARACTER]->(Character)`
* `(Character)-[:HAS_TRAIT]->(Trait)`
* `(Scene)-[:LEADS_TO {choice_text, effects}]->(Scene)`
* `(User)-[:HAS_PLAYTHROUGH]->(Playthrough)`
* `(Playthrough)-[:CURRENT_SCENE]->(Scene)`
* `(Playthrough)-[:AFFINITY_WITH {score}]->(Character)`
* `(Playthrough)-[:VISITED {order: Integer}]->(Scene)` *(MỚI - Lưu lại vết (footprint) của người chơi để AI đọc hiểu diễn biến câu chuyện)*

---

## 5. Yêu Cầu Phi Chức Năng (Non-Functional Requirements)
* **Transaction Safety:** Đảm bảo tính toàn vẹn dữ liệu khi vừa chuyển cảnh (Scene) vừa cập nhật điểm (Affinity) trong cùng một query.
* **Khả năng mở rộng AI (AI Scalability):** Dữ liệu text của Scene và Character phải được thiết kế độ dài phù hợp để sau này chuyển đổi thành Vector (Embeddings) không bị tràn token context window của LLM.

---

## 6. Lộ Trình Phát Triển Tiếp Theo (Phase 2 - GraphRAG & Generative AI)
* **Dynamic AI Conversations:** Người chơi có thể gõ text tự do để nói chuyện với nhân vật tại một số Cảnh (Scene) nhất định. Backend sẽ dùng LangChainGo lấy điểm hảo cảm (`AFFINITY_WITH`), tính cách (`Trait`) và lịch sử (`VISITED`) đưa vào Prompt cho Google Gemini sinh ra câu trả lời.
* **AI Lore Assistant:** Công cụ giúp tác giả tự động tạo thêm nhánh truyện hoặc sinh ra lời thoại nhân vật phụ dựa trên văn phong đã được lưu trữ trong database.