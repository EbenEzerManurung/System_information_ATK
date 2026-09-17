📦 ATK Inventory Management System
Enterprise Office Stationery & Inventory Management Platform
A production-grade, full-stack web application for managing the complete office stationery (ATK) lifecycle — from procurement requests and multi-level approval to QR Code-based verification and real-time analytics.

https://img.shields.io/badge/Go-1.26.0-00ADD8?style=for-the-badge&logo=go&logoColor=white
https://img.shields.io/badge/Gin-1.11-00ADD8?style=for-the-badge&logo=gin&logoColor=white
https://img.shields.io/badge/GORM-1.30-00ADD8?style=for-the-badge
https://img.shields.io/badge/Angular-22.1.6-DD0031?style=for-the-badge&logo=angular&logoColor=white
https://img.shields.io/badge/TypeScript-5.9-3178C6?style=for-the-badge&logo=typescript&logoColor=white
https://img.shields.io/badge/Tailwind-4.3.3-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white
https://img.shields.io/badge/MySQL-8.4-4479A1?style=for-the-badge&logo=mysql&logoColor=white
https://img.shields.io/badge/JWT-Secure-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white

https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square
https://img.shields.io/badge/status-active--development-brightgreen?style=flat-square

Overview •
Features •
Tech Stack •
Architecture •
Getting Started •
Screenshots

📖 Overview
ATK Inventory Management System is an enterprise-grade platform designed to replace manual, spreadsheet-driven office stationery tracking with a secure, auditable, and fully digital workflow.

Built on a scalable RESTful architecture with a Go (Gin) backend and a modern Angular 22 + Signals frontend, the system enables organizations to manage inventory in real time — from request submission, multi-level approval, stock-in/stock-out transactions, to QR Code-based item verification and analytical reporting.

The project was built to demonstrate a real-world, production-style implementation of:

Scalable RESTful API design with role-based authorization

Multi-level approval workflows modeled after real corporate processes

QR Code generation and scanning for tamper-resistant item tracking

Real-time stock analytics with dynamic charts and trend monitoring

Bulk Excel import/export for efficient data migration

Clean separation of concerns across frontend, backend, and data layers

💡 This system reflects how modern enterprises manage shared office resources — combining governance, traceability, and operational efficiency in a single platform.

✨ Key Features
📊 Dashboard & Analytics
Real-time monthly transaction counter with quick-glance KPIs

7-day transaction trend (Stock IN vs OUT) with interactive bar charts

Category-based stock distribution with visual progress bars

Role-aware dashboard view with responsive layout

Recent activity feed with audit trail

📦 Item & Stock Management
Full CRUD for inventory items with category classification

Real-time stock tracking (available, minimum threshold, unit of measure)

Stock movement history (IN / OUT) with audit trail

Category management (Alat Tulis, Elektronik, Kertas, Perlengkapan)

Low-stock alerts and automated threshold monitoring

Export stock data to Excel

🔄 Transaction Management
Request submission for stationery procurement

Stock-IN & Stock-OUT recording with mandatory documentation

Automatic stock recalculation on every transaction

Full transaction history with pagination and advanced filtering

Export transaction records to Excel

✅ Multi-Level Approval Workflow
Structured review process — Manager → Admin → Super Admin

Approve / Reject with mandatory comments for accountability

Real-time status tracking with approval timeline

Centralized pending-approvals dashboard

Notification on status changes

📱 QR Code Integration
Unique QR Code generated automatically for every item

Live camera-based QR scanning for instant verification

Downloadable & printable QR Code labels

Batch QR generation for multiple items

Instant item detail retrieval from scan

📈 Reports & Analytics
Transaction reports by date range, category, and status

Stock valuation and usage statistics

Export-ready reports (Excel) with styled headers

Historical trend analysis with chart visualizations

Custom date-range filtering

👥 User & Role Management
Secure authentication via JWT

Fine-grained Role-Based Access Control (RBAC)

Self-service profile management & password changes

Full activity/audit logging with timestamp

Dynamic menu rendering based on user role

Supported roles:

Role	Level	Responsibilities
Super Admin	5	Full system administration, user & role management, system configuration, audit oversight
Admin	4	User management, master data, stock & transaction management
Manager	3	Approve transactions, monitor stock, generate reports
Staff	2	Create transactions, manage inventory, submit approval requests
Auditor	1	Read-only access for internal audit, compliance & verification
🎨 Modern UI/UX
Responsive layout (Desktop / Tablet / Mobile)

Angular Signals for reactive state management

Tailwind CSS v4 utility-first styling with CSS-first config

Toast notifications for real-time user feedback

Collapsible sidebar with role-filtered navigation

Elegant modal dialogs with keyboard interactions

🛠 Technology Stack
Layer	Technology
Backend Language	Go 1.26.0
Backend Framework	Gin (HTTP router)
ORM	GORM
Database	MySQL 8.4
Authentication	JWT + bcrypt
Frontend Framework	Angular 22.1.6 (Standalone Components)
Reactive State	Angular Signals
Language	TypeScript 5.9
Styling	Tailwind CSS v4.3.3
HTTP Client	Angular HttpClient + RxJS
Charts	Custom SVG Charts
QR Code	qrcode (generation) + jsQR (scanning)
Excel Handling	excelize (Go)
API Style	RESTful (JSON over HTTP)
Architecture	Clean Layered (Handler → Service → Repository)
Responsive across Desktop, Tablet, and Mobile.

🏗 Project Architecture
text
                Client (Browser)
                          │
              Angular 22 + Signals
                          │
              RESTful API (JSON over HTTP)
                          │
              Go 1.26 + Gin (Handler Layer)
                          │
                  Service Layer (Business Logic)
                          │
                  Repository Layer (GORM)
                          │
                       MySQL 8.4
The system follows a clean layered architecture with strict separation between HTTP handling (handlers), business logic (services), and data access (repositories). Role-based middleware guards every protected route, ensuring zero unauthorized access.

📁 Project Structure
text
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
    │   ├── index.html
    │   ├── main.ts
    │   └── styles.css
    ├── angular.json
    ├── package.json
    └── tsconfig.json
🚀 Getting Started
Prerequisites
Make sure you have the following installed:

Go 1.26.0 or higher → Download

Node.js 20+ and npm 10+ → Download

Angular CLI 22+ → npm install -g @angular/cli@latest

MySQL 8.4 → Download

Git

Backend Setup
bash
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
Backend runs at http://localhost:8080 by default.

Frontend Setup
bash
# 1. Navigate to frontend
cd ../frontend

# 2. Install dependencies
npm install

# 3. Configure API endpoint (if different)
# Edit src/environments/environment.ts

# 4. Start development server
ng serve
Frontend runs at http://localhost:4200.

Default Login Credentials
Username	Password	Role
super_admin	admin123	Super Admin
admin	admin123	Admin
manager	admin123	Manager
staff	admin123	Staff
auditor	admin123	Auditor (Read-Only)
⚠️ Change default passwords immediately after first login in production.

🔌 API Reference (Brief)
Authentication
Method	Endpoint	Description
POST	/api/auth/login	User login, returns JWT token
POST	/api/auth/register	Register new user (admin only)
GET	/api/auth/me	Get current authenticated user
Users
Method	Endpoint	Description
GET	/api/users	List users (paginated, searchable)
POST	/api/users	Create new user
PUT	/api/users/:id	Update user (includes username, role)
DELETE	/api/users/:id	Delete user
PUT	/api/profile	Update own profile
POST	/api/profile/change-password	Change own password
Items & Transactions
Method	Endpoint	Description
GET	/api/items	List inventory items
POST	/api/items	Create new item
GET	/api/transactions	List transactions
POST	/api/transactions	Create IN/OUT transaction
GET	/api/transactions/export	Export transactions to Excel
Approvals
Method	Endpoint	Description
GET	/api/approvals/pending	Get pending approvals
POST	/api/approvals/:id/approve	Approve transaction
POST	/api/approvals/:id/reject	Reject transaction
QR Code
Method	Endpoint	Description
GET	/api/items/:id/qr	Generate QR Code for item
POST	/api/qr/verify	Verify item from QR payload
📸 Screenshots
🔐 Login
<img width="1890" height="920" alt="Login" src="PASTE_URL_LOGIN_DISINI" />
📊 Dashboard — Real-Time Analytics
<img width="1890" height="920" alt="Dashboard" src="PASTE_URL_DASHBOARD_DISINI" />
👥 User Management with RBAC
<img width="1890" height="920" alt="Users" src="PASTE_URL_USERS_DISINI" />
📦 Stock Management
<img width="1890" height="920" alt="Stock" src="PASTE_URL_STOCK_DISINI" />
🔄 Transactions (IN vs OUT)
<img width="1890" height="920" alt="Transactions" src="PASTE_URL_TRANSACTIONS_DISINI" />
✅ Approval Workflow
<img width="1890" height="920" alt="Approvals" src="PASTE_URL_APPROVALS_DISINI" />
📱 QR Code Generator & Scanner
<img width="1890" height="920" alt="QR Code" src="PASTE_URL_QR_DISINI" />
📈 Reports & Analytics
<img width="1890" height="920" alt="Reports" src="PASTE_URL_REPORTS_DISINI" />
🔐 Security Highlights
JWT-based authentication with signed tokens and expiry

bcrypt password hashing (cost factor 10)

RBAC middleware enforcing role-based access on every protected route

Input validation on both frontend (Angular Forms) and backend (Gin binding)

SQL injection prevention via GORM's parameterized queries

Password change verification requiring current password confirmation

Unique constraints on username and email at the database level

CORS configured to whitelist trusted origins only

🗺 Roadmap
☑ Core inventory management (CRUD)
☑ Multi-level approval workflow
☑ QR Code generation & scanning
☑ Excel import/export
☑ Role-based dashboards
☑ Audit logging
☑ JWT authentication with RBAC
☑ Responsive UI with Angular Signals
□ Push notifications for approval events
□ Barcode scanner integration
□ Multi-warehouse support
□ CI/CD pipeline with GitHub Actions
□ Progressive Web App (PWA) support
□ Offline-first capability with Service Worker
□ Dark mode theme
🤝 Contributing
Contributions are welcome! Please follow these steps:

Fork the repository

Create a feature branch (git checkout -b feature/AmazingFeature)

Commit your changes (git commit -m 'Add some AmazingFeature')

Push to the branch (git push origin feature/AmazingFeature)

Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

📄 License
This project is licensed under the MIT License — see the LICENSE file for details.

👨‍💻 Author
Eben Nezer Manurung

Backend Developer • Full Stack Developer

https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white
https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white
https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white

<div align="center">
⭐ If this project helps you, please give it a star!
Built with ❤️ using Go 1.26, Angular 22, Tailwind v4, and MySQL 8.4

</div>
