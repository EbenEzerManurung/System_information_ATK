import {
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { Subject, debounceTime, distinctUntilChanged, takeUntil } from 'rxjs';

@Component({
  selector: 'app-filter-bar',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.Eager,
  template: `
    <div class="flex items-center gap-3 flex-wrap">
      <div class="relative">
        <svg
          class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-primary-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-4.35-4.35M17 11a6 6 0 11-12 0 6 6 0 0112 0z"
          />
        </svg>
        <input
          [formControl]="searchControl"
          [placeholder]="placeholder"
          class="pl-9 pr-4 py-2 text-sm w-64 border border-primary-200 rounded-lg bg-white
                      focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 outline-none transition"
        />
      </div>
      <ng-content></ng-content>
    </div>
  `,
})
export class FilterBarComponent implements OnInit, OnDestroy {
  @Input() placeholder = 'Cari...';
  @Output() search = new EventEmitter<string>();

  searchControl = new FormControl('');
  private destroy$ = new Subject<void>();

  ngOnInit() {
    this.searchControl.valueChanges
      .pipe(debounceTime(400), distinctUntilChanged(), takeUntil(this.destroy$))
      .subscribe((v) => this.search.emit(v ?? ''));
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
