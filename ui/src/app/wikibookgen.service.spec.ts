import { TestBed } from '@angular/core/testing';

import { WikibookgenService } from './wikibookgen.service';

describe('WikibookgenService', () => {
  let service: WikibookgenService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WikibookgenService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
