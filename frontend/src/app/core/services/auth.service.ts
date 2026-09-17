import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiResponse, User } from '../models';

interface LoginResponse {
  token: string;
  user: User;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly TOKEN_KEY = 'atk_token';
  private readonly USER_KEY = 'atk_user';

  currentUser = signal<User | null>(this.loadUser());

  constructor(private http: HttpClient, private router: Router) {}

  // ============================================================
  // LOGIN
  // ============================================================
  login(username: string, password: string): Observable<any> {
    return this.http.post<any>(
      `${environment.apiUrl}/auth/login`,
      { username, password }
    ).pipe(
      tap(res => {
        // ==== Deteksi sukses dari berbagai format backend ====
        const isSuccess =
          res?.success === true ||
          res?.status === 'success' ||
          res?.code === 200 ||
          res?.code === 201 ||
          (res?.data && (res.data.token || res.data.access_token || res.data.accessToken));

        if (!isSuccess) {
          console.warn('[AuthService] Login response tidak dikenali:', res);
          return;
        }

        // ==== Ambil token dari berbagai field ====
        const payload = res.data ?? res;
        const token =
          payload?.token ??
          payload?.access_token ??
          payload?.accessToken ??
          res?.token;

        // ==== Ambil user dari berbagai field ====
        const user =
          payload?.user ??
          res?.user ??
          payload?.profile;

        if (token) {
          localStorage.setItem(this.TOKEN_KEY, token);
        }
        if (user) {
          localStorage.setItem(this.USER_KEY, JSON.stringify(user));
          this.currentUser.set(user);
        }
      })
    );
  }

  // ============================================================
  // LOGOUT
  // ============================================================
  logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
    this.currentUser.set(null);
    this.router.navigate(['/login']);
  }

  // ============================================================
  // TOKEN / AUTH CHECK
  // ============================================================
  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  // ============================================================
  // UPDATE PROFILE
  // ============================================================
  updateProfile(payload: { full_name: string; email: string }): Observable<any> {
    return this.http.put<any>(
      `${environment.apiUrl}/profile`, payload
    ).pipe(
      tap(res => {
        const isSuccess =
          res?.success === true ||
          res?.status === 'success' ||
          res?.code === 200;

        if (isSuccess) {
          const current = this.currentUser();
          const updatedUser = { ...(current ?? {}), ...(res.data ?? res) } as User;
          localStorage.setItem(this.USER_KEY, JSON.stringify(updatedUser));
          this.currentUser.set(updatedUser);
        }
      })
    );
  }

  // ============================================================
  // CHANGE PASSWORD
  // ============================================================
  changePassword(payload: {
    currentPassword: string;
    newPassword: string;
    confirmPassword: string;
  }): Observable<any> {
    return this.http.post<any>(
      `${environment.apiUrl}/profile/change-password`, payload
    );
  }

  // ============================================================
  // HELPER: LOAD USER DARI LOCALSTORAGE
  // ============================================================
  private loadUser(): User | null {
    try {
      const raw = localStorage.getItem(this.USER_KEY);
      if (!raw) return null;
      return JSON.parse(raw) as User;
    } catch (err) {
      console.warn('[AuthService] localStorage user rusak, di-reset:', err);
      localStorage.removeItem(this.USER_KEY);
      localStorage.removeItem(this.TOKEN_KEY);
      return null;
    }
  }
}