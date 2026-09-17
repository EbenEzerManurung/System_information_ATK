import { Component, EventEmitter, Input, Output, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-excel-toolbar',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.Eager,
  template: `
    <div class="flex items-center gap-2">
      <label
        class="px-3 py-2 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 cursor-pointer transition flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1M12 12V4m0 0L8 8m4-4l4 4"
          />
        </svg>
        Import Excel
        <input type="file" accept=".xlsx,.xls" class="hidden" (change)="onFile($event)" />
      </label>
      <button
        (click)="exportClick.emit()"
        class="px-3 py-2 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 transition flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1M12 4v12m0 0l-4-4m4 4l4-4"
          />
        </svg>
        Export Excel
      </button>
    </div>
  `,
})
export class ExcelToolbarComponent {
  @Input() importing = false;
  @Output() import = new EventEmitter<File>();
  @Output() exportClick = new EventEmitter<void>();

  onFile(ev: Event) {
    const input = ev.target as HTMLInputElement;
    if (input.files?.length) {
      this.import.emit(input.files[0]);
      input.value = '';
    }
  }
}
