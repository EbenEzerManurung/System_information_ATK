import { Component, OnInit, signal, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { ToastService } from '../../core/services/toast.service';
import { ApiService } from '../../core/services/api.service';
import { ToastComponent } from '../../shared/components/toast/toast.component';

interface MenuItem {
  path: string;
  label: string;
  icon: string;
  roles: string[];
}

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterOutlet, RouterLink, RouterLinkActive, ToastComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="flex h-screen bg-primary-50">
      <aside
        [class.w-64]="!collapsed()"
        [class.w-20]="collapsed()"
        class="bg-primary-900 text-white transition-all duration-300 flex flex-col shadow-xl"
      >
        <div class="h-16 flex items-center px-4 border-b border-primary-800">
          <div class="flex items-center gap-3 overflow-hidden">
            <div class="w-9 h-9 rounded-lg bg-primary-500 flex items-center justify-center shrink-0 font-bold">A</div>
            @if (!collapsed()) {
              <div class="font-semibold text-sm whitespace-nowrap">ATK System</div>
            }
          </div>
        </div>

        <nav class="flex-1 py-4 overflow-y-auto">
          @for (item of filteredMenus(); track item.path) {
            <a
              [routerLink]="item.path"
              routerLinkActive="bg-primary-700 border-l-4 border-primary-400"
              [routerLinkActiveOptions]="{ exact: false }"
              class="flex items-center gap-3 px-4 py-3 text-sm hover:bg-primary-800 transition-all border-l-4 border-transparent"
            >
              <span class="w-5 h-5 shrink-0" [innerHTML]="item.icon"></span>
              @if (!collapsed()) {
                <span class="whitespace-nowrap">{{ item.label }}</span>
              }
            </a>
          }
        </nav>

        <button (click)="toggle()" class="p-3 border-t border-primary-800 hover:bg-primary-800 transition text-sm">
          {{ collapsed() ? '→' : '← hide' }}
        </button>
      </aside>

      <div class="flex-1 flex flex-col overflow-hidden">
        <header class="h-16 bg-white border-b border-primary-100 flex items-center justify-between px-6 shadow-sm">
          <h1 class="text-lg font-semibold text-primary-900">System Information ATK</h1>

          <div class="flex items-center gap-4">
            <div class="relative">
              <button (click)="toggleProfileMenu()" class="flex items-center gap-3 px-2 py-1 rounded-lg hover:bg-primary-50 transition">
                <div class="text-right hidden sm:block">
                  <div class="text-sm font-medium text-primary-900">{{ auth.currentUser()?.full_name }}</div>
                  <div class="text-xs text-primary-500">{{ userRole() || 'Loading...' }}</div>
                </div>
                <div class="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 font-semibold">
                  {{ initials() }}
                </div>
                <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              @if (showProfileMenu()) {
                <div class="fixed inset-0 z-40" (click)="closeProfileMenu()"></div>
                <div class="absolute right-0 mt-2 w-56 bg-white rounded-xl shadow-lg border border-primary-100 py-2 z-50 fade-in">
                  <div class="px-4 py-2 border-b border-primary-50">
                    <div class="text-sm font-semibold text-primary-900">{{ auth.currentUser()?.full_name }}</div>
                    <div class="text-xs text-slate-500">{{ auth.currentUser()?.email }}</div>
                  </div>
                  <button (click)="openEditProfile()" class="w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-primary-50 transition flex items-center gap-2">
                    <span>👤</span> Edit Profile & Password
                  </button>
                  <div class="border-t border-primary-50 my-1"></div>
                  <button (click)="auth.logout()" class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition flex items-center gap-2">
                    <span>🚪</span> Logout
                  </button>
                </div>
              }
            </div>
          </div>
        </header>

        <main class="flex-1 overflow-y-auto p-6">
          <router-outlet />
        </main>
      </div>
    </div>

    <app-toast />

    @if (showEditProfile()) {
      <div class="fixed inset-0 bg-black/40 z-[100] flex items-center justify-center p-4 overflow-y-auto" (click)="closeEditProfile()">
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 my-8 fade-in" (click)="$event.stopPropagation()">
          <h3 class="text-lg font-bold text-primary-900 mb-4">Edit Profile</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Username</label>
              <input [value]="auth.currentUser()?.username" disabled class="w-full px-3 py-2 border border-slate-200 rounded-lg bg-slate-50 text-slate-500 cursor-not-allowed" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Nama Lengkap</label>
              <input [(ngModel)]="profileForm.full_name" class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
              <input type="email" [(ngModel)]="profileForm.email" class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none" />
            </div>
          </div>
          <div class="border-t border-primary-100 my-5"></div>
          <h4 class="text-sm font-semibold text-primary-900 mb-3">Ganti Password (Opsional)</h4>
          <div class="space-y-3">

            <!-- Password Lama -->
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Password Lama</label>
              <div class="relative">
                <input [type]="showCurrentPassword() ? 'text' : 'password'" [(ngModel)]="passwordForm.currentPassword"
                       class="w-full pl-3 pr-10 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30 transition" />
                <button type="button" (click)="showCurrentPassword.set(!showCurrentPassword())"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none">
                  @if (showCurrentPassword()) {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                  } @else {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path></svg>
                  }
                </button>
              </div>
            </div>

            <!-- Password Baru -->
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Password Baru</label>
              <div class="relative">
                <input [type]="showNewPassword() ? 'text' : 'password'" [(ngModel)]="passwordForm.newPassword"
                       class="w-full pl-3 pr-10 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30 transition" />
                <button type="button" (click)="showNewPassword.set(!showNewPassword())"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none">
                  @if (showNewPassword()) {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                  } @else {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path></svg>
                  }
                </button>
              </div>
              <p class="text-xs text-slate-400 mt-1">Minimal 6 karakter</p>
            </div>

            <!-- Konfirmasi Password -->
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Konfirmasi Password</label>
              <div class="relative">
                <input [type]="showConfirmPassword() ? 'text' : 'password'" [(ngModel)]="passwordForm.confirmPassword"
                       class="w-full pl-3 pr-10 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30 transition" />
                <button type="button" (click)="showConfirmPassword.set(!showConfirmPassword())"
                        class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none">
                  @if (showConfirmPassword()) {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                  } @else {
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path></svg>
                  }
                </button>
              </div>
            </div>

          </div>
          <div class="flex justify-end gap-2 mt-6">
            <button (click)="closeEditProfile()" class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition">Batal</button>
            <button (click)="saveProfile()" [disabled]="saving()" class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition">{{ saving() ? 'Menyimpan...' : 'Simpan Perubahan' }}</button>
          </div>
        </div>
      </div>
    }
  `,
})
export class MainLayoutComponent implements OnInit {
  collapsed = signal(false);
  showProfileMenu = signal(false);
  showEditProfile = signal(false);
  saving = signal(false);

  userRole = signal<string | null>(null);

  showCurrentPassword = signal(false);
  showNewPassword = signal(false);
  showConfirmPassword = signal(false);

  profileForm = { full_name: '', email: '' };
  passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };

  menus: MenuItem[] = [
    { path: '/dashboard', label: 'Dashboard', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/></svg>', roles: ['super_admin', 'admin', 'manager', 'auditor', 'staff'] },
    { path: '/users', label: 'User', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>', roles: ['super_admin', 'admin'] },
    { path: '/roles', label: 'Role', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>', roles: ['super_admin', 'admin'] },
    { path: '/master', label: 'Master Data', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"/></svg>', roles: ['super_admin', 'admin'] },
    { path: '/stock', label: 'Stock', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>', roles: ['super_admin', 'admin', 'manager', 'staff'] },
    { path: '/transactions', label: 'Transaksi', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/></svg>', roles: ['super_admin', 'admin', 'manager', 'staff'] },
    { path: '/approvals', label: 'Approval', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>', roles: ['super_admin', 'manager'] },
    { path: '/reports', label: 'Laporan', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>', roles: ['super_admin', 'admin', 'manager', 'auditor', 'staff'] },
    { path: '/qr-generator', label: 'Generate QR', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"/></svg>', roles: ['super_admin', 'admin', 'manager', 'staff'] },
    { path: '/qr-scanner', label: 'Scan QR', icon: '<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"/></svg>', roles: ['super_admin', 'admin', 'manager', 'staff'] },
  ];

  constructor(
    public auth: AuthService,
    private toast: ToastService,
    private api: ApiService,
  ) {}

  ngOnInit() {
    this.resolveUserRole();
  }

  private resolveUserRole() {
    console.log('🔍 [RBAC] ========== MULAI RESOLVE ROLE ==========');

    let role = this.extractRole(this.auth.currentUser());
    if (role) {
      console.log('✅ [RBAC] Role dari currentUser():', role);
      this.userRole.set(role);
      return;
    }

    role = this.extractRoleFromLocalStorage();
    if (role) {
      console.log('✅ [RBAC] Role dari localStorage:', role);
      this.userRole.set(role);
      return;
    }

    console.log('🔄 [RBAC] Role belum ada, coba fetch dari API /auth/me...');
    this.api.get<any>('auth/me').subscribe({
      next: (res) => {
        const user = res?.data ?? res;
        console.log('✅ [RBAC] Response /auth/me:', user);

        const apiRole = this.extractRole(user);
        if (apiRole) {
          console.log('✅ [RBAC] Role dari API /auth/me:', apiRole);
          this.userRole.set(apiRole);
          this.saveRoleToLocalStorage(apiRole);
        } else {
          console.warn('⚠️ [RBAC] API /auth/me tidak mengembalikan role. Fallback ke auditor.');
          this.userRole.set('auditor');
        }
      },
      error: (err) => {
        console.error('❌ [RBAC] Gagal fetch /auth/me:', err);
        console.warn('⚠️ [RBAC] Fallback ke auditor (akses minimal).');
        this.userRole.set('auditor');
      }
    });
  }

  private extractRole(user: any): string {
    if (!user) return '';

    const possibleKeys = [
      'roleName', 'role_name', 'userRole', 'user_role',
      'role', 'Role', 'roles', 'RoleName', 'user_role_name'
    ];

    for (const key of possibleKeys) {
      const val = user[key];
      if (val === undefined || val === null) continue;

      if (typeof val === 'string' && val.trim() !== '') {
        return this.normalizeRole(val);
      }

      if (typeof val === 'object' && !Array.isArray(val)) {
        const nested = val.name || val.title || val.roleName || val.role_name;
        if (nested && typeof nested === 'string') return this.normalizeRole(nested);
      }

      if (Array.isArray(val) && val.length > 0) {
        const first = val[0];
        if (typeof first === 'string') return this.normalizeRole(first);
        if (typeof first === 'object' && first.name) return this.normalizeRole(first.name);
      }
    }
    return '';
  }

  private extractRoleFromLocalStorage(): string {
    const keys = ['atk_user', 'user', 'currentUser', 'auth_user', 'userData'];
    for (const key of keys) {
      try {
        const raw = localStorage.getItem(key);
        if (!raw) continue;
        const parsed = JSON.parse(raw);
        const role = this.extractRole(parsed);
        if (role) return role;
      } catch (e) {}
    }
    return '';
  }

  private saveRoleToLocalStorage(role: string) {
    try {
      const raw = localStorage.getItem('atk_user');
      if (raw) {
        const user = JSON.parse(raw);
        user.roleName = role;
        localStorage.setItem('atk_user', JSON.stringify(user));
      }
    } catch (e) {}
  }

  private normalizeRole(role: string): string {
    let normalized = String(role).toLowerCase().trim().replace(/\s+/g, '_');
    if (normalized === 'superadmin') normalized = 'super_admin';
    return normalized;
  }

  filteredMenus = computed(() => {
    const role = this.userRole();

    if (role === null) {
      console.log('⏳ [RBAC] Role belum siap, tampilkan semua menu sementara');
      return this.menus;
    }

    console.log('🎯 [RBAC] Filter menu dengan role:', role);

    if (role === 'super_admin') {
      return this.menus;
    }

    const filtered = this.menus.filter((menu) => {
      if (!menu.roles || menu.roles.length === 0) return true;
      return menu.roles.includes(role);
    });

    console.log('🎯 [RBAC] Menu setelah filter:', filtered.map(m => m.label));
    return filtered;
  });

  toggle() { this.collapsed.update((v) => !v); }

  initials(): string {
    const name = this.auth.currentUser()?.full_name ?? '';
    return name.split(' ').map((n) => n[0]).slice(0, 2).join('').toUpperCase() || 'U';
  }

  toggleProfileMenu() { this.showProfileMenu.update((v) => !v); }
  closeProfileMenu() { this.showProfileMenu.set(false); }

  openEditProfile() {
    this.closeProfileMenu();
    const user = this.auth.currentUser();
    this.profileForm = { full_name: user?.full_name ?? '', email: user?.email ?? '' };
    this.passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };

    this.showCurrentPassword.set(false);
    this.showNewPassword.set(false);
    this.showConfirmPassword.set(false);

    this.showEditProfile.set(true);
  }

  closeEditProfile() { this.showEditProfile.set(false); }

  // ============================================================
  // ✅ HELPER: Cek apakah user mengubah data profile
  // ============================================================
  private hasProfileChange(): boolean {
    const user = this.auth.currentUser();
    if (!user) return true;
    return (
      this.profileForm.full_name !== (user.full_name ?? '') ||
      this.profileForm.email !== (user.email ?? '')
    );
  }

  // ============================================================
  // ✅ HELPER: Cek apakah user mengisi form password
  // ============================================================
  isChangingPassword(): boolean {
    return !!(
      this.passwordForm.currentPassword ||
      this.passwordForm.newPassword ||
      this.passwordForm.confirmPassword
    );
  }

  // ============================================================
  // ✅ HELPER: Validasi form password
  // ============================================================
  isPasswordValid(): boolean {
    const f = this.passwordForm;
    if (!f.currentPassword || !f.newPassword || !f.confirmPassword) return false;
    if (f.newPassword.length < 6) return false;
    if (f.newPassword !== f.confirmPassword) return false;
    if (f.currentPassword === f.newPassword) return false;
    return true;
  }

  // ============================================================
  // ✅ SAVE PROFILE - DIPERBAIKI
  // ============================================================
  // Alur baru:
  // 1. Validasi keduanya (profile & password jika diisi)
  // 2. Jika ada perubahan password → Ubah password DULU
  // 3. Jika berhasil → Lanjut update profile (jika ada perubahan)
  // 4. Jika tidak ada perubahan apa pun → Tutup modal
  // ============================================================
  saveProfile() {
    const changingPassword = this.isChangingPassword();
    const changingProfile = this.hasProfileChange();

    // ---- VALIDASI PROFILE ----
    if (changingProfile) {
      if (!this.profileForm.full_name || this.profileForm.full_name.trim() === '') {
        this.toast.warning('Nama lengkap wajib diisi');
        return;
      }
      if (!this.profileForm.email || this.profileForm.email.trim() === '') {
        this.toast.warning('Email wajib diisi');
        return;
      }
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(this.profileForm.email)) {
        this.toast.warning('Format email tidak valid');
        return;
      }
    }

    // ---- VALIDASI PASSWORD ----
    if (changingPassword) {
      if (!this.passwordForm.currentPassword) {
        this.toast.warning('Password lama wajib diisi');
        return;
      }
      if (!this.passwordForm.newPassword) {
        this.toast.warning('Password baru wajib diisi');
        return;
      }
      if (this.passwordForm.newPassword.length < 6) {
        this.toast.warning('Password baru minimal 6 karakter');
        return;
      }
      if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
        this.toast.warning('Konfirmasi password tidak sama');
        return;
      }
      if (this.passwordForm.currentPassword === this.passwordForm.newPassword) {
        this.toast.warning('Password baru tidak boleh sama dengan password lama');
        return;
      }
    }

    // ---- JIKA TIDAK ADA PERUBAHAN ----
    if (!changingPassword && !changingProfile) {
      this.toast.warning('Tidak ada perubahan untuk disimpan');
      this.closeEditProfile();
      return;
    }

    this.saving.set(true);

    // ---- FUNGSI: Update Profile ----
    const doUpdateProfile = () => {
      if (!changingProfile) {
        // Tidak ada perubahan profile → selesai
        this.saving.set(false);
        this.toast.success('Password berhasil diperbarui');
        this.closeEditProfile();
        return;
      }

      console.log('📤 [PROFILE] Update profil:', this.profileForm);
      this.auth.updateProfile(this.profileForm).subscribe({
        next: (resProfile) => {
          this.saving.set(false);
          console.log('📥 [PROFILE] Response:', resProfile);
          if (resProfile.success) {
            if (changingPassword) {
              this.toast.success('Profil dan password berhasil diperbarui');
            } else {
              this.toast.success('Profil berhasil diperbarui');
            }
            this.closeEditProfile();
          } else {
            this.toast.error(resProfile.message || 'Gagal memperbarui profil');
          }
        },
        error: (err) => {
          this.saving.set(false);
          console.error('❌ [PROFILE] Error:', err);
          this.toast.error(err?.error?.message ?? 'Gagal memperbarui profil');
        }
      });
    };

    // ---- FUNGSI: Change Password ----
    const doChangePassword = () => {
      if (!changingPassword) {
        // Tidak ada perubahan password → lanjut update profile
        doUpdateProfile();
        return;
      }

      console.log('📤 [PASSWORD] Change password');
      this.auth.changePassword(this.passwordForm).subscribe({
        next: (resPass) => {
          console.log('📥 [PASSWORD] Response:', resPass);
          if (resPass.success) {
            // Reset form password
            this.passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };
            // Lanjut update profile (jika ada)
            doUpdateProfile();
          } else {
            this.saving.set(false);
            this.toast.error(resPass.message || 'Gagal mengubah password. Pastikan password lama benar.');
          }
        },
        error: (err) => {
          this.saving.set(false);
          console.error('❌ [PASSWORD] Error:', err);
          console.error('❌ [PASSWORD] Error body:', err?.error);
          const msg = err?.error?.message ?? 'Password lama salah atau terjadi kesalahan';
          this.toast.error(msg);
        }
      });
    };

    // ---- JALANKAN: Password dulu, baru Profile ----
    doChangePassword();
  }
}