import { Component, OnInit } from '@angular/core';
import {FormControl} from '@angular/forms';
import {Observable, of} from 'rxjs';
import {map, startWith, debounceTime, distinctUntilChanged, switchMap} from 'rxjs/operators';
import { Location } from '@angular/common';

import { Wikibook } from '../wikibook';
import { WikibookgenService } from '../wikibookgen.service'; 

@Component({
  selector: 'app-order-wikibook',
  templateUrl: './order-wikibook.component.html',
  styleUrls: ['./order-wikibook.component.sass']
})
export class OrderWikibookComponent implements OnInit {
  myControl = new FormControl();
  filteredOptions$: Observable<string[]>;
  language: string;
  model: string;

  constructor(
    private wikibookgenService: WikibookgenService,
    private location: Location,
  ) {
    this.language = 'en';
    this.model = 'tour';
  }
  
  ngOnInit(): void {
    this.filteredOptions$ = this.myControl.valueChanges.pipe(
      startWith(''),
      // wait 300ms after each keystroke before considering the term
      debounceTime(300),
      // ignore new term if same as previous term
      distinctUntilChanged(),
      switchMap(value => this.wikibookgenService.autocomplete(value, this.language))
    );
  }

  public orderWikibook() {
    console.log("Ordering " + this.myControl.value);
    this.wikibookgenService.order(this.myControl.value, this.language, this.model)
      .subscribe((orderID: string) => {
        console.log('wikibook order ' + orderID + ' generating, redirecting...');
        this.location.go('/wikibook/order/'+ orderID, '', null);
        this.location.forward();
      });
  }
}
