import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { PageResult, User } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { ExcelToolbarComponent } from '../../shared/components/excel-toolbar/excel-toolbar.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-users',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    DataTableComponent,
    ExcelToolbarComponent,
    FilterBarComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h2 class="text-2xl font-bold text-primary-900">Manajemen User</h2>
          <p class="text-sm text-slate-500">Kelola data pengguna sistem</p>
        </div>
        <div class="flex items-center gap-2">
          <app-excel-toolbar (import)="onImport($event)" (exportClick)="onExport()" />
          <button
            (click)="openForm()"
            class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition flex items-center gap-2"
          >
            <span>+</span> Tambah User
          </button>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm">
        <app-filter-bar placeholder="Cari user..." (search)="onSearch($event)" />
      </div>

      <app-data-table
        [columns]="columns"
        [page]="page()"
        [loading]="loading()"
        [sortBy]="sortBy"
        [sortDir]="sortDir"
        [trackBy]="trackById"
        (pageChange)="loadData($event)"
        (sortChange)="onSort($event)"
        (edit)="openForm($event)"
        (delete)="onDelete($event)"
      />

      @if (showForm()) {
        <div
          class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4"
          (click)="closeForm()"
        >
          <div
            class="bg-white rounded-2xl shadow-2xl w-full max-w-lg p-6 fade-in"
            (click)="$event.stopPropagation()"
          >
            <h3 class="text-lg font-bold text-primary-900 mb-4">
              {{ form.id ? 'Edit User' : 'Tambah User' }}
            </h3>
            <div class="space-y-3">
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Username</label>
                <input
                  [(ngModel)]="form.username"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              @if (!form.id) {
                <div>
                  <label class="block text-sm font-medium text-slate-700 mb-1">Password</label>
                  <div class="relative">
                    <input
                      [(ngModel)]="form.password"
                      [type]="showPassword() ? 'text' : 'password'"
                      minlength="6"
                      class="w-full pl-3 pr-10 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                    />
                    <button
                      type="button"
                      (click)="showPassword.set(!showPassword())"
                      class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none"
                    >
                      @if (showPassword()) {
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                      } @else {
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path></svg>
                      }
                    </button>
                  </div>
                  <p class="text-xs text-slate-400 mt-1">Minimal 6 karakter</p>
                </div>
              }

              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Nama Lengkap</label>
                <input
                  [(ngModel)]="form.fullName"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
                <input
                  [(ngModel)]="form.email"
                  type="email"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">No Phone</label>
                <input
                  [(ngModel)]="form.phoneNumber"
                  type="tel"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Role</label>
                <select
                  [(ngModel)]="form.roleId"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                >
                  @for (r of roles(); track r.id) {
                    <option [ngValue]="r.id">{{ r.name }}</option>
                  }
                </select>
              </div>
              <div class="flex items-center gap-2">
                <input type="checkbox" [(ngModel)]="form.active" id="active" class="rounded" />
                <label for="active" class="text-sm text-slate-700">Aktif</label>
              </div>
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <button
                (click)="closeForm()"
                class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition"
              >
                Batal
              </button>
              <button
                (click)="save()"
                [disabled]="saving()"
                class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition"
              >
                {{ saving() ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class UsersComponent implements OnInit {
  columns: TableColumn[] = [
    { key: 'username', label: 'Username', sortable: true },
    { key: 'fullName', label: 'Nama Lengkap', sortable: true },
    { key: 'email', label: 'Email' },
    { key: 'phoneNumber', label: 'No Phone' },
    { key: 'roleName', label: 'Role' },
  ];

  page = signal<PageResult<User> | null>(null);
  loading = signal(false);
  saving = signal(false);
  showForm = signal(false);
  roles = signal<any[]>([]);
  showPassword = signal(false);

  // ✅ Signal baru: track ID user yang sedang dihapus
  deletingId = signal<number | null>(null);

  form: any = {};
  search = '';
  sortBy = 'username';
  sortDir: 'asc' | 'desc' = 'asc';
  currentPage = 1;

  // ✅ TrackBy yang lebih aman
  trackById = (row: any) => row?.id ?? row?.uuid ?? Math.random();

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadData(1);
    this.api.get<any>('roles').subscribe((res) => {
      if (res.status === 'success') this.roles.set(res.data);
    });
  }

  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);
    this.api
      .get<any>('users', {
        page,
        limit: 10,
        search: this.search,
        sortBy: this.sortBy,
        sortDir: this.sortDir,
      })
      .subscribe({
        next: (res) => {
          this.loading.set(false);
          if (res.status === 'success') {
            const rows = (res.data ?? []).map((u: any) => ({
              ...u,
              username: u.username ?? u.user_name ?? '',
              fullName: u.fullName ?? u.full_name ?? '',
              phoneNumber: u.phoneNumber ?? u.phone_number ?? '',
              email: u.email ?? '',
              roleId: u.roleId ?? u.role_id ?? u.roles?.[0]?.id ?? null,
              roleName: u.roleName ?? u.role_name ?? u.roles?.[0]?.name ?? '-',
              active: (u.status ?? '').toLowerCase() === 'active',
            }));

            this.page.set({
              data: rows,
              total: res.meta?.total ?? rows.length,
              totalPages: res.meta?.totalPage ?? 1,
              page: res.meta?.page ?? page,
              size: res.meta?.limit ?? 10,
            } as PageResult<User>);
          } else {
            this.toast.error(res.message ?? 'Gagal memuat data user');
          }
        },
        error: () => {
          this.loading.set(false);
          this.toast.error('Gagal memuat data user');
        },
      });
  }

  onSearch(q: string) {
    this.search = q;
    this.loadData(1);
  }

  onSort(e: { sortBy: string; sortDir: 'asc' | 'desc' }) {
    this.sortBy = e.sortBy;
    this.sortDir = e.sortDir;
    this.loadData(this.currentPage);
  }

  openForm(row?: User) {
    if (row) {
      const r: any = row;
      this.form = {
        id: r.id,
        username: r.username ?? '',
        fullName: r.fullName ?? r.full_name ?? '',
        email: r.email ?? '',
        phoneNumber: r.phoneNumber ?? r.phone_number ?? '',
        roleId: r.roleId ?? r.role_id ?? r.roles?.[0]?.id ?? null,
        active: r.active ?? (r.status === 'active'),
      };
    } else {
      this.form = {
        active: true,
        roleId: this.roles()[0]?.id ?? null,
      };
    }
    this.showPassword.set(false);
    this.showForm.set(true);
  }

  closeForm() {
    this.showForm.set(false);
    this.form = {};
  }

  save() {
    if (!this.form.username || this.form.username.trim() === '') {
      this.toast.error('Username wajib diisi');
      return;
    }
    if (!this.form.fullName || this.form.fullName.trim() === '') {
      this.toast.error('Nama lengkap wajib diisi');
      return;
    }
    if (!this.form.email || this.form.email.trim() === '') {
      this.toast.error('Email wajib diisi');
      return;
    }
    if (!this.form.id) {
      if (!this.form.password || this.form.password.length < 6) {
        this.toast.error('Password minimal 6 karakter');
        return;
      }
    }

    this.saving.set(true);

    const payload: any = {
      username: this.form.username?.trim(),
      fullName: this.form.fullName?.trim(),
      phoneNumber: this.form.phoneNumber?.trim() ?? '',
      email: this.form.email?.trim(),
      roleId: this.form.roleId,
      status: this.form.active ? 'active' : 'inactive',
      full_name: this.form.fullName?.trim(),
      phone_number: this.form.phoneNumber?.trim() ?? '',
      role_id: this.form.roleId,
    };

    if (!this.form.id) {
      payload.password = this.form.password;
    }

    console.log('📤 [SAVE USER] Payload:', JSON.stringify(payload, null, 2));
    console.log('📤 [SAVE USER] Method:', this.form.id ? `PUT users/${this.form.id}` : 'POST users');

    const req = this.form.id
      ? this.api.put<any>(`users/${this.form.id}`, payload)
      : this.api.post<any>('users', payload);

    req.subscribe({
      next: (res) => {
        this.saving.set(false);
        console.log('📥 [SAVE USER] Response:', res);
        if (res.status === 'success') {
          this.toast.success('Data berhasil disimpan');
          this.closeForm();
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res.message ?? 'Gagal menyimpan data');
        }
      },
      error: (err) => {
        this.saving.set(false);
        console.error('❌ [SAVE USER] Error:', err);
        this.toast.error(err?.error?.message ?? 'Gagal menyimpan data');
      },
    });
  }

  // ============================================================
  // ✅ DELETE USER — DIPERBAIKI TOTAL
  // ============================================================
  onDelete(row: User) {
    // 1. Validasi ID
    const userId = (row as any)?.id ?? (row as any)?.ID;
    if (!userId) {
      console.error('❌ [DELETE] row tidak punya ID:', row);
      this.toast.error('ID user tidak valid, tidak bisa dihapus');
      return;
    }

    // 2. Cek apakah user menghapus dirinya sendiri
    const currentUser = this.getCurrentUserFromStorage();
    if (currentUser?.id === userId) {
      this.toast.error('Anda tidak bisa menghapus akun Anda sendiri');
      return;
    }

    // 3. Cek apakah sedang menghapus user ini (cegah double-click)
    if (this.deletingId() === userId) {
      return;
    }

    // 4. Konfirmasi
    const label = (row as any).username ?? (row as any).email ?? `ID ${userId}`;
    const confirmed = confirm(
      `Hapus user "${label}"?\n\nTindakan ini tidak dapat dibatalkan.`
    );
    if (!confirmed) return;

    // 5. Tandai sebagai "sedang dihapus"
    this.deletingId.set(userId);

    console.log(`🗑️ [DELETE] Menghapus user ID=${userId}`);

    // 6. Panggil API
    this.api.delete<any>(`users/${userId}`).subscribe({
      next: (res) => {
        this.deletingId.set(null);
        console.log('📥 [DELETE] Response:', res);

        if (res?.status === 'success') {
          this.toast.success(`User "${label}" berhasil dihapus`);

          // ✅ Cek jika ini user terakhir di halaman & bukan di halaman 1
          const currentData = this.page()?.data ?? [];
          const isLastItemOnPage = currentData.length === 1;

          if (isLastItemOnPage && this.currentPage > 1) {
            // Pindah ke halaman sebelumnya
            this.loadData(this.currentPage - 1);
          } else {
            // Reload halaman saat ini
            this.loadData(this.currentPage);
          }
        } else {
          this.toast.error(res?.message ?? 'Gagal menghapus user');
        }
      },
      error: (err) => {
        this.deletingId.set(null);
        console.error('❌ [DELETE] Error status:', err?.status);
        console.error('❌ [DELETE] Error body:', err?.error);

        let msg = 'Gagal menghapus user';

        if (err?.status === 0) {
          msg = 'Tidak dapat terhubung ke server';
        } else if (err?.status === 401) {
          msg = 'Sesi Anda berakhir, silakan login ulang';
        } else if (err?.status === 403) {
          msg = 'Anda tidak memiliki izin untuk menghapus user';
        } else if (err?.status === 404) {
          msg = 'User tidak ditemukan (mungkin sudah dihapus)';
          // Refresh tabel karena data mungkin sudah tidak ada
          this.loadData(this.currentPage);
        } else if (err?.status === 400) {
          msg = err?.error?.message ?? 'Tidak dapat menghapus user ini';
        } else if (err?.status >= 500) {
          msg = 'Server error, coba beberapa saat lagi';
        } else if (err?.error?.message) {
          msg = err.error.message;
        }

        this.toast.error(msg);
      },
    });
  }

  // ============================================================
  // ✅ Helper: Ambil user login dari localStorage
  // ============================================================
  private getCurrentUserFromStorage(): any {
    const keys = ['atk_user', 'user', 'currentUser'];
    for (const key of keys) {
      try {
        const raw = localStorage.getItem(key);
        if (!raw) continue;
        return JSON.parse(raw);
      } catch (e) {}
    }
    return null;
  }

  onImport(file: File) {
    this.api.upload<any>('users/import', file).subscribe({
      next: (res) => {
        if (res.status === 'success') {
          this.toast.success(`Import berhasil: ${res.data.imported} data`);
          this.loadData(1);
        } else {
          this.toast.error(res.message);
        }
      },
      error: (err) => this.toast.error(err?.error?.message ?? 'Gagal import file'),
    });
  }

  onExport() {
    this.api.exportFile('users/export', { search: this.search }).subscribe((blob) => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `users_${Date.now()}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
      this.toast.success('File berhasil diunduh');
    });
  }
}