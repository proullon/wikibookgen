import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ListWikibookComponent } from './list-wikibook.component';

describe('ListWikibookComponent', () => {
  let component: ListWikibookComponent;
  let fixture: ComponentFixture<ListWikibookComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ListWikibookComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ListWikibookComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
