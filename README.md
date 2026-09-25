<div align="center">

# 📦 ATK Inventory Management System

### Enterprise Office Stationery & Inventory Management Platform

**A production-grade, full-stack web application for managing the complete office stationery (ATK) lifecycle — from procurement requests and multi-level approval to QR-code-based verification and real-time analytics.**

[![Go](https://img.shields.io/badge/Go-1.26.0-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.11-00ADD8?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/GORM-1.30-00ADD8?style=for-the-badge)](https://gorm.io)
[![Angular](https://img.shields.io/badge/Angular-22.1.6-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Tailwind](https://img.shields.io/badge/Tailwind-4.3.3-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com)
[![JWT](https://img.shields.io/badge/JWT-Secure-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)](https://jwt.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
![Status](https://img.shields.io/badge/status-active--development-brightgreen?style=for-the-badge)

**[Overview](#-overview)** • **[Features](#-key-features)** • **[Tech Stack](#-technology-stack)** • **[Architecture](#-architecture)** • **[Project Structure](#-project-structure)** • **[Getting Started](#-getting-started)** • **[Screenshots](#-screenshots)** • **[Author](#-author)**

</div>

---

## 📖 Overview

**ATK Inventory Management System** is an enterprise-grade platform designed to replace manual, spreadsheet-driven office stationery tracking with a secure, auditable, and fully digital workflow.

Built on a scalable RESTful architecture with a **Go (Gin) backend** and a modern **Angular 22 + Signals** frontend, the system enables organizations to manage inventory in real time — from request submission and multi-level approval, to stock-in/stock-out transactions, QR-code-based item verification, and analytical reporting.

This project was built to demonstrate a real-world, production-style implementation of:

- Scalable **RESTful API design** with role-based authorization
- **Multi-level approval workflows** modeled after real corporate processes
- **QR code generation and scanning** for tamper-resistant item tracking
- **Real-time stock analytics** with dynamic charts and trend monitoring
- **Bulk Excel import/export** for efficient data migration
- Clean separation of concerns across **frontend, backend, and data layers**

> 💡 This system reflects how modern enterprises manage shared office resources — combining governance, traceability, and operational efficiency in a single platform.

---

## ✨ Key Features

<table>
<tr>
<td width="50%" valign="top">

### 📊 Dashboard & Analytics
- Real-time monthly transaction counter with quick-glance KPIs
- 7-day transaction trend (Stock IN vs OUT) with interactive bar charts
- Category-based stock distribution with visual progress bars
- Role-aware dashboard view with responsive layout
- Recent activity feed with full audit trail

### 📦 Item & Stock Management
- Full CRUD for inventory items with category classification
- Real-time stock tracking (available, minimum threshold, unit of measure)
- Stock movement history (IN / OUT) with audit trail
- Category management (Alat Tulis, Elektronik, Kertas, Perlengkapan)
- Low-stock alerts and automated threshold monitoring
- Export stock data to Excel

### 🔄 Transaction Management
- Request submission for stationery procurement
- Stock-IN & Stock-OUT recording with mandatory documentation
- Automatic stock recalculation on every transaction
- Full transaction history with pagination and advanced filtering
- Export transaction records to Excel

### ✅ Multi-Level Approval Workflow
- Structured review process — **Staff → Manager → Admin → Super Admin**
- Approve / Reject with mandatory comments for accountability
- Real-time status tracking with approval timeline
- Centralized pending-approvals dashboard
- Notifications on status changes

</td>
<td width="50%" valign="top">

### 📱 QR Code Integration
- Unique QR code generated automatically for every item
- Live camera-based QR scanning for instant verification
- Downloadable & printable QR code labels
- Batch QR generation for multiple items
- Instant item detail retrieval from scan

### 📈 Reports & Analytics
- Transaction reports by date range, category, and status
- Stock valuation and usage statistics
- Export-ready reports (Excel) with styled headers
- Historical trend analysis with chart visualizations
- Custom date-range filtering

### 👥 User & Role Management
- Secure authentication via **JWT**
- Fine-grained **Role-Based Access Control (RBAC)**
- Self-service profile management & password changes
- Full activity/audit logging with timestamps
- Dynamic menu rendering based on user role

### 🔒 Security
- JWT authentication with password hashing (**bcrypt**)
- Middleware-enforced authorization (RBAC)
- Server-side input validation across all endpoints
- Parameterized queries via GORM (SQL-injection safe)
- Unique constraints on `username` and `email`
- **CORS** configured to whitelist trusted origins only

</td>
</tr>
</table>

### 🎨 Modern UI/UX

Responsive layout across Desktop, Tablet, and Mobile • Angular **Signals** for reactive state management • **Tailwind CSS v4.3.3** utility-first styling with CSS-first config • Toast notifications for real-time feedback • Collapsible sidebar with role-filtered navigation • Elegant modal dialogs with keyboard interactions • Installable **Progressive Web App (PWA)**

---

## 👥 Roles & Responsibilities

| Role | Level | Responsibilities |
|---|:---:|---|
| **Super Admin** | 5 | Full system administration, user & role management, system configuration, audit oversight |
| **Admin** | 4 | User management, master data, stock & transaction management |
| **Manager** | 3 | Approve transactions, monitor stock, generate reports |
| **Staff** | 2 | Create transactions, manage inventory, submit approval requests |
| **Auditor** | 1 | Read-only access for internal audit, compliance & verification |

---

## 🛠 Technology Stack

| Layer | Technology |
|---|---|
| **Backend Language** | Go 1.26.0 |
| **Backend Framework** | Gin (HTTP router) |
| **ORM** | GORM |
| **Database** | MySQL 8.4 |
| **Authentication** | JWT + bcrypt |
| **Frontend Framework** | Angular 22.1.6 (Standalone Components) |
| **Reactive State** | Angular Signals |
| **Language** | TypeScript 5.9 |
| **Styling** | Tailwind CSS v4.3.3 |
| **HTTP Client** | Angular HttpClient + RxJS |
| **Charts** | Custom SVG charts |
| **QR Code** | `qrcode` (generation) + `jsQR` (scanning) |
| **Excel Handling** | `excelize` (Go) |
| **API Style** | RESTful (JSON over HTTP) |
| **Architecture** | Clean Layered (Handler → Service → Repository) |

---

## 🏗 Architecture

The system follows a clean **layered architecture** with strict separation between HTTP handling (`handlers`), business logic (`services`), and data access (`repositories`). Role-based middleware guards every protected route, ensuring zero unauthorized access.

```
┌─────────────────────────┐
│   Client (Browser)      │
│  Angular 22 + Signals   │
└────────────┬─────────────┘
             │  RESTful API (JSON over HTTP)
┌────────────▼─────────────┐
│   Go 1.26 + Gin           │
│   Handler Layer           │
└────────────┬─────────────┘
┌────────────▼─────────────┐
│   Service Layer            │
│   (Business Logic)         │
└────────────┬─────────────┘
┌────────────▼─────────────┐
│   Repository Layer         │
│   (GORM)                   │
└────────────┬─────────────┘
┌────────────▼─────────────┐
│   MySQL 8.4                │
└─────────────────────────┘
```

---

## 📁 Project Structure

```
System_information_ATK/
│
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── config/
│   │   └── database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── item.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   └── approval.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── role_repository.go
│   │   ├── item_repository.go
│   │   └── transaction_repository.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── item_service.go
│   │   └── transaction_service.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── item_handler.go
│   │   └── transaction_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── rbac.go
│   ├── routes/
│   │   └── routes.go
│   └── utils/
│       ├── jwt.go
│       ├── hash.go
│       └── response.go
│
└── frontend/
    ├── src/
    │   ├── app/
    │   │   ├── core/
    │   │   │   ├── guards/
    │   │   │   ├── interceptors/
    │   │   │   ├── models/
    │   │   │   └── services/
    │   │   ├── shared/
    │   │   │   ├── components/
    │   │   │   │   ├── data-table/
    │   │   │   │   ├── excel-toolbar/
    │   │   │   │   ├── filter-bar/
    │   │   │   │   └── toast/
    │   │   │   └── layout/
    │   │   │       └── main-layout.component.ts
    │   │   ├── features/
    │   │   │   ├── auth/
    │   │   │   ├── dashboard/
    │   │   │   ├── users/
    │   │   │   ├── roles/
    │   │   │   ├── master/
    │   │   │   ├── stock/
    │   │   │   ├── transactions/
    │   │   │   ├── approvals/
    │   │   │   ├── reports/
    │   │   │   └── qr/
    │   │   ├── app.config.ts
    │   │   ├── app.routes.ts
    │   │   └── app.ts
    │   ├── environments/
    │   │   ├── environment.ts
    │   │   └── environment.prod.ts
    │   ├── index.html
    │   ├── main.ts
    │   └── styles.css
    ├── angular.json
    ├── package.json
    └── tsconfig.json
```

---

## 🚀 Getting Started

### Prerequisites

Make sure you have the following installed:

- **Go** 1.26.0 or higher → [Download](https://go.dev/dl/)
- **Node.js** 20+ and **npm** 10+ → [Download](https://nodejs.org)
- **Angular CLI** 22+ → `npm install -g @angular/cli@latest`
- **MySQL** 8.4 → [Download](https://dev.mysql.com/downloads/)
- **Git**

### Backend Setup

```bash
# 1. Clone the repository
git clone https://github.com/EbenEzerManurung/System_information_ATK.git
cd System_information_ATK/backend

# 2. Configure environment
cp .env.example .env
# Edit .env: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET

# 3. Install dependencies
go mod download

# 4. Run the server (auto-migrate on start)
go run main.go
```

### Frontend Setup

```bash
# 1. Navigate to frontend
cd ../frontend

# 2. Install dependencies
npm install

# 3. Configure API endpoint (if different)
# Edit src/environments/environment.ts

# 4. Start development server
ng serve --open
```

The app will be available at `http://localhost:4200`, connecting to the API at `http://localhost:8080` by default.

---

## 📸 Screenshots

### 🔐 Login
<img width="1902" height="964" alt="Login screen" src="https://github.com/user-attachments/assets/68862bba-958a-4b2e-a51b-c901342be64d" />

### 📊 Dashboard — Real-Time Analytics
<img width="1918" height="1003" alt="Dashboard 1" src="https://github.com/user-attachments/assets/8034f191-038c-4783-b93f-64888ec3dc53" />
<img width="1918" height="951" alt="Dashboard 2" src="https://github.com/user-attachments/assets/382ba873-6117-4491-b69d-a6bf21b444fc" />
<img width="1903" height="1006" alt="Dashboard 3" src="https://github.com/user-attachments/assets/abb75390-eccc-47b6-91fd-5bd95ec90dc3" />

### 🚀 Progressive Web App (PWA)
<img width="1918" height="943" alt="PWA" src="https://github.com/user-attachments/assets/38a447b9-c20f-4314-a4a0-e3e92cf757b1" />

### 👥 User Management with RBAC
<img width="1914" height="966" alt="User management 1" src="https://github.com/user-attachments/assets/d1762e70-a163-4f51-aee4-9cde97af1ecc" />
<img width="1918" height="991" alt="User management 2" src="https://github.com/user-attachments/assets/966dcf32-187d-4e47-8bd1-a4487f25d59d" />

### 📁 Master Data
<img width="1914" height="1011" alt="Master data" src="https://github.com/user-attachments/assets/b88455e0-6c07-40c7-85e3-c5a46b0401bc" />

### 📦 Stock Management
<img width="1918" height="997" alt="Stock management" src="https://github.com/user-attachments/assets/cd334269-33a6-4d76-bf75-3640d9337678" />

### 🔄 Transactions (IN vs OUT)
<img width="1891" height="997" alt="Transactions" src="https://github.com/user-attachments/assets/9c832d7c-b6c7-4153-aa71-465f106d83a9" />

### ✅ Approval Workflow
<img width="1914" height="946" alt="Approval workflow 1" src="https://github.com/user-attachments/assets/38150c76-d76c-454d-82da-86f053c4302b" />
<img width="1918" height="1008" alt="Approval workflow 2" src="https://github.com/user-attachments/assets/7e796a31-bd0f-45f1-bb93-11540f483dae" />

### 📱 QR Code Generator
<img width="1917" height="1021" alt="QR generator" src="https://github.com/user-attachments/assets/1bfa9255-72e7-4554-90d0-43853d4fbd70" />

### 📱 QR Code Scanner
<img width="1903" height="1009" alt="QR scanner" src="https://github.com/user-attachments/assets/ce64ef0d-6688-4be6-a41d-71d6030d2880" />

### 📈 Reports
<img width="1918" height="1017" alt="Reports 1" src="https://github.com/user-attachments/assets/6fc71cb2-02b9-4122-b33a-a4485634254d" />
<img width="1903" height="1039" alt="Reports 2" src="https://github.com/user-attachments/assets/c1b0f082-f656-43e7-a003-0c57c2f15680" />

<details>
<summary><b>🖥️ Local run screenshots (Backend & Frontend)</b></summary>
<br>

**Backend (Go)**
<img width="1338" height="838" alt="Backend running" src="https://github.com/user-attachments/assets/7f800daf-4b0c-4c0d-8608-e28fd9f1d4d4" />

**Frontend (Angular)**
<img width="1168" height="544" alt="Frontend running" src="https://github.com/user-attachments/assets/766a7135-0121-4f60-9c5b-10e127bbf5c0" />

</details>

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome. Feel free to check the [issues page](https://github.com/EbenEzerManurung/System_information_ATK/issues) or open a pull request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](https://tlo.mit.edu/resources/mit-github) file for details.

---

## 👨‍💻 Author

<div align="center">

**Eben Nezer Manurung**
*Backend Developer • Full Stack Developer*

[![GitHub](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/EbenEzerManurung)


⭐ If this project helped you, please consider giving it a star!

*Built with ❤️ using Go 1.26, Angular 22, Tailwind v4.3.3, and MySQL 8.4*

</div>
