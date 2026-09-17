import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';

// Konversi snake_case -> camelCase secara rekursif (object & array)
function toCamel(input: any): any {
  if (Array.isArray(input)) {
    return input.map(toCamel);
  }
  if (input !== null && typeof input === 'object' && !(input instanceof Date)) {
    return Object.keys(input).reduce((acc: any, key) => {
      const camelKey = key.replace(/_([a-z0-9])/g, (_, c) => c.toUpperCase());
      acc[camelKey] = toCamel(input[key]);
      return acc;
    }, {});
  }
  return input;
}

@Injectable({ providedIn: 'root' })
export class ApiService {
  private base = environment.apiUrl;

  constructor(private http: HttpClient) {}

  private buildParams(params?: Record<string, any>): HttpParams {
    let p = new HttpParams();
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== null && v !== undefined && v !== '') p = p.set(k, v);
      });
    }
    return p;
  }

  get<T>(path: string, params?: Record<string, any>): Observable<T> {
    return this.http
      .get<any>(`${this.base}/${path}`, { params: this.buildParams(params) })
      .pipe(map((res) => toCamel(res)));
  }

  post<T>(path: string, body: any): Observable<T> {
    return this.http
      .post<any>(`${this.base}/${path}`, body)
      .pipe(map((res) => toCamel(res)));
  }

  put<T>(path: string, body: any): Observable<T> {
    return this.http
      .put<any>(`${this.base}/${path}`, body)
      .pipe(map((res) => toCamel(res)));
  }

  delete<T>(path: string): Observable<T> {
    return this.http
      .delete<any>(`${this.base}/${path}`)
      .pipe(map((res) => toCamel(res)));
  }

  upload<T>(path: string, file: File): Observable<T> {
    const fd = new FormData();
    fd.append('file', file);
    return this.http
      .post<any>(`${this.base}/${path}`, fd)
      .pipe(map((res) => toCamel(res)));
  }

  exportFile(path: string, params?: Record<string, any>): Observable<Blob> {
    return this.http.get(`${this.base}/${path}`, {
      params: this.buildParams(params),
      responseType: 'blob',
    });
  }
}