import { Component, OnInit, Input } from '@angular/core';

import { Wikibook } from '../wikibook';
import { WikibookgenService } from '../wikibookgen.service'; 

@Component({
  selector: 'app-list-wikibook',
  templateUrl: './list-wikibook.component.html',
  styleUrls: ['./list-wikibook.component.sass']
})
export class ListWikibookComponent implements OnInit {

  @Input() wikibooks: Wikibook[];

  constructor(
    private wikibookgenService: WikibookgenService,
  ) { }

  ngOnInit() {
    this.getWikibooks(1, 50, '');
  }

  getWikibooks(page: number, size: number, language: string): void {
    this.wikibookgenService.listWikibook(page, size, language)
      .subscribe((wikibooks:Wikibook[]) => {
        this.wikibooks = wikibooks
      });
  }
}
