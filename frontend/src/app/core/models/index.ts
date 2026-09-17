// ============================================================
// USER
// ============================================================
export interface User {
  id: number;
  uuid?: string;
  username: string;
  full_name: string;
  email: string;
  phoneNumber?: string;
  roleId?: number;
  roleName?: string;
  roles?: Role[];
  status?: string;
  active?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// ============================================================
// ROLE
// ============================================================
export interface Role {
  id: number;
  name: string;
  displayName?: string;
  description: string;
  level?: number;
  permissions?: string[];
  createdAt?: string;
  updatedAt?: string;
}

// ============================================================
// MASTER
// ============================================================
export interface Master {
  id: number;
  uuid?: string;
  code: string;
  name: string;
  category: string;
  unit: string;
  description?: string;
  minStock?: number;
  maxStock?: number;
  price?: number;
  status?: string;
  stocks?: Stock[];
  createdAt?: string;
  updatedAt?: string;
}

// Alias untuk backward compat
export type MasterItem = Master;

// ============================================================
// STOCK
// ============================================================
export interface Stock {
  id: number;
  masterId: number;
  quantity: number;
  location: string;
  status: string;
  lastUpdated?: string;
  createdAt?: string;
  updatedAt?: string;
  master?: Master; // nested object dari Preload
}

// ============================================================
// TRANSACTION ITEM
// ============================================================
export interface TransactionItem {
  id?: number;
  transactionId?: number;
  masterId: number;
  quantity: number;
  price?: number;
  subTotal?: number;
  notes?: string;
  master?: Master;
}

// ============================================================
// TRANSACTION
// ============================================================
export interface Transaction {
  id: number;
  transactionCode: string;          // ← sesuai backend
  type: 'IN' | 'OUT' | 'RETURN' | 'ADJUSTMENT';
  status: string;                    // draft/pending/approved/completed/cancelled
  userId?: number;
  notes?: string;
  totalItems?: number;
  totalValue?: number;
  transactionDate?: string;
  approvedBy?: number;
  approvedAt?: string;
  createdAt?: string;
  updatedAt?: string;

  // Relasi (dari Preload)
  user?: User;
  approver?: User;
  items?: TransactionItem[];
  documents?: any[];
}

// ============================================================
// APPROVAL
// ============================================================
export interface Approval {
  id: number;
  transactionId: number;
  transactionCode?: string;
  requestedBy?: string;
  requestedAt?: string;
  status: 'PENDING' | 'APPROVED' | 'REJECTED' | 'pending' | 'approved' | 'rejected';
  signatureData?: string;
  approvedBy?: string;
  approvedAt?: string;
  notes?: string;
}

// ============================================================
// DOCUMENT
// ============================================================
export interface Document {
  id: number;
  type: string;
  transactionId?: number;
  userId?: number;
  documentCode?: string;
  qrCode?: string;
  signatureData?: string;
  metadata?: string;
  createdAt?: string;
  updatedAt?: string;
}

// ============================================================
// PAGINATION
// ============================================================
export interface PageResult<T> {
  data: T[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
}

// ============================================================
// API RESPONSE
// ============================================================
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  meta?: any;
}