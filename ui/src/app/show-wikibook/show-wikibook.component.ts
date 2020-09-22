import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { NestedTreeControl } from '@angular/cdk/tree';
import { MatTreeModule, MatTreeNestedDataSource } from '@angular/material/tree';

import { Wikibook, GetAvailableDownloadFormatResponse } from '../wikibook';
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

  epubAvailable: boolean = false;
  epubPrintButtonText: string = 'Request print';
  pdfAvailable: boolean = false;
  pdfPrintButtonText: string = 'Request print';

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
    this.getAvailableDownloadFormat(uuid);
  }

  getAvailableDownloadFormat(uuid: string): void {
    this.wikibookgenService.getAvailableDownloadFormat(uuid)
      .subscribe((r:GetAvailableDownloadFormatResponse) => {
        if (r.epub == 'exists') {
          this.epubAvailable = true;
        }
        if (r.epub == 'printing') {
          this.epubAvailable = false;
          this.epubPrintButtonText = 'Printing';
        }

        if (r.pdf == 'exists') {
          this.pdfAvailable = true;
        }
        if (r.pdf == 'printing') {
          this.pdfAvailable = false;
          this.pdfPrintButtonText = 'Printing';
        }
      });
  }

  print(uuid: string, format: string) {
    if (format == 'epub') {
      this.epubPrintButtonText = 'Printing';
    }
    if (format == 'pdf') {
      this.pdfPrintButtonText = 'Printing';
    }

    this.wikibookgenService.print(uuid, format)
      .subscribe((r:any) => {
        this.getAvailableDownloadFormat(uuid);
      });

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
      var node = new WikibookNode(ch.title);
      for (let a of ch.articles) {
        var n = new WikibookNode(a.title);
        node.nodes.push(n);
      }
      nodes.push(node);
    }
    console.log(nodes);
    return nodes;
  }

  hasChild = (_: number, node: WikibookNode) => !!node.nodes && node.nodes.length > 0;
}
