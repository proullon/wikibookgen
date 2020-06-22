import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { NestedTreeControl } from '@angular/cdk/tree';
import { MatTreeModule, MatTreeNestedDataSource } from '@angular/material/tree';

import { Wikibook } from '../wikibook';
import { WikibookgenService } from '../wikibookgen.service'; 

class WikibookNode {
  title: string;
  nodes?: WikibookNode[];

  constructor(title: string) {
    this.title = title;
    this.nodes = new Array<WikibookNode>();
  }
}

const TREE_DATA: WikibookNode[] = [];

@Component({
  selector: 'app-show-wikibook',
  templateUrl: './show-wikibook.component.html',
  styleUrls: ['./show-wikibook.component.sass']
})
export class ShowWikibookComponent implements OnInit {
  
  treeControl = new NestedTreeControl<WikibookNode>(node => node.nodes);
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
        this.dataSource.data = this.wikibookToWikibookNode(wikibook);
      });
  }
 
  wikibookToWikibookNode(wikibook: Wikibook): WikibookNode[] {
    var nodes: Array<WikibookNode> = [];

    for (let ch of wikibook.volumes[0].chapters) {
      console.log('yaya ' + ch.title);
      var node = new WikibookNode(ch.title);
      for (let a of ch.articles) {
        console.log('page ' + a.title);
        var n = new WikibookNode(a.title);
        node.nodes.push(n);
      }
      console.log('chapter ' + node.title + ' has ' + node.nodes.length + ' pages');
      nodes.push(node);
    }
    console.log(nodes);
    return nodes;
  }

  hasChild = (_: number, node: WikibookNode) => !!node.nodes && node.nodes.length > 0;
}
