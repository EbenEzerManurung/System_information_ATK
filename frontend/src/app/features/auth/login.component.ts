import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { ToastService } from '../../core/services/toast.service';
import { ToastComponent } from '../../shared/components/toast/toast.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, ToastComponent],
  template: `
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-900 via-primary-700 to-primary-500 p-4">
      <div class="w-full max-w-md bg-white rounded-2xl shadow-2xl p-8 fade-in">
        <div class="text-center mb-8">
          <div class="w-16 h-16 rounded-2xl bg-primary-600 text-white flex items-center justify-center text-2xl font-bold mx-auto mb-4">A</div>
          <h1 class="text-2xl font-bold text-primary-900">System Information ATK</h1>
          <p class="text-sm text-slate-500 mt-1">Masuk untuk melanjutkan</p>
        </div>

        <form (ngSubmit)="onSubmit()" class="space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">Username</label>
            <input [(ngModel)]="username" name="username"
                   class="w-full px-4 py-2.5 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 outline-none transition"
                   placeholder="Masukkan username">
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">Password</label>
            <div class="relative">
              <input [(ngModel)]="password" name="password" [type]="showPassword() ? 'text' : 'password'"
                     class="w-full pl-4 pr-10 py-2.5 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 outline-none transition"
                     placeholder="Masukkan password">
              
              <!-- Tombol Toggle Password -->
              <button type="button" (click)="showPassword.set(!showPassword())"
                      class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none">
                @if (showPassword()) {
                  <!-- Ikon Mata Terbuka (Password Terlihat) -->
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                } @else {
                  <!-- Ikon Mata Tertutup (Password Tersembunyi) -->
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                  </svg>
                }
              </button>
            </div>
          </div>

          @if (serverError()) {
            <div class="p-3 rounded-lg bg-red-50 border border-red-200 text-sm text-red-700 fade-in">
              {{ serverError() }}
            </div>
          }

          <button type="submit" [disabled]="loading()"
                  class="w-full py-2.5 rounded-lg bg-primary-600 text-white font-medium hover:bg-primary-700 disabled:opacity-60 transition">
            {{ loading() ? 'Memproses...' : 'Masuk' }}
          </button>
        </form>
      </div>
    </div>

    <app-toast />
  `
})
export class LoginComponent {
  username = '';
  password = '';
  loading = signal(false);
  serverError = signal('');
  showPassword = signal(false); // <-- Tambahan state untuk toggle password

  constructor(
    private auth: AuthService,
    private router: Router,
    private toast: ToastService
  ) {}

  onSubmit() {
    this.serverError.set('');

    if (!this.username?.trim() || !this.password) {
      this.toast.warning('Username dan password wajib diisi');
      return;
    }

    this.loading.set(true);
    this.auth.login(this.username.trim(), this.password).subscribe({
      next: (res: any) => {
        this.loading.set(false);

        // ==== DETEKSI SUKSES DARI BERBAGAI FORMAT BACKEND ====
        const isSuccess =
          res?.success === true ||
          res?.status === 'success' ||
          res?.code === 200 ||
          res?.code === 201 ||
          (res?.data && (res.data.token || res.data.access_token));

        if (isSuccess) {
          this.toast.success('Login berhasil!');
          setTimeout(() => this.router.navigate(['/dashboard']), 300);
        } else {
          const msg = res?.message || 'Login gagal';
          this.serverError.set(msg);
          this.toast.error(msg);
        }
      },
      error: (err) => {
        this.loading.set(false);
        let msg = 'Username atau password salah';
        if (err.status === 0) {
          msg = 'Tidak dapat terhubung ke server';
        } else if (err.error?.message) {
          msg = err.error.message;
        }
        this.serverError.set(msg);
        this.toast.error(msg);
      }
    });
  }
}