import { Routes } from '@angular/router'; // ✅ Perbaiki typo dari '@angularrouter' menjadi '@angular/router'
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login.component').then(m => m.LoginComponent)
  },
  {
    path: '',
    loadComponent: () => import('./shared/layout/main-layout.component').then(m => m.MainLayoutComponent),
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () => import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent)
      },
      {
        path: 'users',
        loadComponent: () => import('./features/users/users.component').then(m => m.UsersComponent)
      },
      {
        path: 'roles',
        loadComponent: () => import('./features/roles/roles.component').then(m => m.RolesComponent)
      },
      {
        path: 'master',
        loadComponent: () => import('./features/master/master.component').then(m => m.MasterComponent)
      },
      {
        path: 'stock',
        loadComponent: () => import('./features/stock/stock.component').then(m => m.StockComponent)
      },
      {
        path: 'transactions',
        loadComponent: () => import('./features/transactions/transactions.component').then(m => m.TransactionsComponent)
      },
      {
        path: 'approvals',
        loadComponent: () => import('./features/approvals/approvals.component').then(m => m.ApprovalsComponent)
      },
      {
        path: 'reports',
        loadComponent: () => import('./features/reports/reports.component').then(m => m.ReportsComponent)
      },
      {
        path: 'qr-generator',
        loadComponent: () => import('./features/qr/qr-generator.component').then(m => m.QrGeneratorComponent)
      },
      {
        path: 'qr-scanner',
        loadComponent: () => import('./features/qr/qr-scanner.component').then(m => m.QrScannerComponent)
      }
    ]
  },
  { path: '**', redirectTo: 'dashboard' }
];