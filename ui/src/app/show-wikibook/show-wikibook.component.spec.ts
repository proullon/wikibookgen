import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowWikibookComponent } from './show-wikibook.component';

describe('ShowWikibookComponent', () => {
  let component: ShowWikibookComponent;
  let fixture: ComponentFixture<ShowWikibookComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShowWikibookComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShowWikibookComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
