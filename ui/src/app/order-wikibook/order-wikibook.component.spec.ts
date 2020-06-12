import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OrderWikibookComponent } from './order-wikibook.component';

describe('OrderWikibookComponent', () => {
  let component: OrderWikibookComponent;
  let fixture: ComponentFixture<OrderWikibookComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OrderWikibookComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OrderWikibookComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
