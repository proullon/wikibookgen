import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { NestedTreeControl } from '@angular/cdk/tree';
import { MatTreeModule, MatTreeNestedDataSource } from '@angular/material/tree';

import { Wikibook } from '../wikibook';
import { WikibookgenService } from '../wikibookgen.service'; 

interface WikibookNode {
  title: string;
  articles?: WikibookNode[];
}

const TREE_DATA: WikibookNode[] = [];

@Component({
  selector: 'app-show-wikibook',
  templateUrl: './show-wikibook.component.html',
  styleUrls: ['./show-wikibook.component.sass']
})
export class ShowWikibookComponent implements OnInit {
  
  treeControl = new NestedTreeControl<WikibookNode>(node => node.articles);
  dataSource = new MatTreeNestedDataSource<WikibookNode>();

  @Input() wikibook: Wikibook;

  constructor(
    private route: ActivatedRoute,
    private wikibookgenService: WikibookgenService,
    private location: Location
  ) {
    this.dataSource.data = TREE_DATA;
  }

  ngOnInit(): void {
    const uuid = this.route.snapshot.paramMap.get('id');
    this.getWikibook(uuid);
  }

  getWikibook(uuid: string): void {
    this.wikibookgenService.getWikibook(uuid)
      .subscribe((wikibook:Wikibook) => {
        console.log(wikibook);
        this.wikibook = wikibook;
        this.dataSource.data = wikibook.volumes[0].chapters;
      });
  }
  
  hasChild = (_: number, node: WikibookNode) => !!node.articles && node.articles.length > 0;
}
