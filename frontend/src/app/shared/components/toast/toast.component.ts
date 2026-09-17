import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastService, Toast } from '../../../core/services/toast.service';

@Component({
  selector: 'app-toast',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <!-- Container: fixed, kanan atas, pointer-events-none agar tidak block klik -->
    <div
      class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none"
      aria-live="polite"
      aria-atomic="true"
    >
      @for (t of toastService.toasts(); track t.id) {
        <div
          class="toast-item pointer-events-auto min-w-[280px] max-w-md px-4 py-3 rounded-xl shadow-lg text-white flex items-center gap-3"
          [ngClass]="bgClass(t.type)"
          role="alert"
        >
          <!-- Icon -->
          <span class="text-lg shrink-0" aria-hidden="true">
            {{ icon(t.type) }}
          </span>

          <!-- Message -->
          <span class="text-sm font-medium flex-1">{{ t.message }}</span>

          <!-- Close button -->
          <button
            type="button"
            (click)="toastService.remove(t.id)"
            class="opacity-80 hover:opacity-100 text-lg leading-none shrink-0 transition"
            aria-label="Tutup notifikasi"
          >
            &times;
          </button>
        </div>
      }
    </div>
  `,
  styles: [`
    .toast-item {
      animation: toastSlideIn 0.25s ease-out;
    }
    @keyframes toastSlideIn {
      from {
        opacity: 0;
        transform: translateX(24px) scale(0.95);
      }
      to {
        opacity: 1;
        transform: translateX(0) scale(1);
      }
    }
  `]
})
export class ToastComponent {
  constructor(public toastService: ToastService) {}

  /** Mapping tipe toast → class background Tailwind */
  bgClass(type: Toast['type']): Record<string, boolean> {
    return {
      'bg-green-600': type === 'success',
      'bg-red-600':   type === 'error',
      'bg-blue-600':  type === 'info',
      'bg-amber-500': type === 'warning'
    };
  }

  /** Mapping tipe toast → emoji icon */
  icon(type: Toast['type']): string {
    switch (type) {
      case 'success': return '✅';
      case 'error':   return '❌';
      case 'warning': return '⚠️';
      case 'info':    return 'ℹ️';
      default:        return '';
    }
  }
}